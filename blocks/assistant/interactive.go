package assistant

// The interactive blocks: each is a Go value that holds the block's state, a method that applies what
// the reader did (the form htmx posts) and returns the next state, and a templ component that draws a
// state with the hx-* attributes that post the next action. The package has no handlers: the app
// decides the URL, keeps or rebuilds the state, and answers with the component.
//
//	mux.HandleFunc("POST /plan/approval", func(w http.ResponseWriter, r *http.Request) {
//		r.ParseForm()
//		next, err := approval.Advance(r.PostForm)
//		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
//		assistant.ApprovalCard(next).Render(r.Context(), w)
//	})
//
// The page needs htmx (https://htmx.org). The forms carry no secrets and the server must check them as
// it would any form: the methods reject choices that are not on offer, but whether the reader may act
// at all is the app's call.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
)

// ErrChoice is returned for a choice the block did not offer.
var ErrChoice = errors.New("assistant: not one of the choices")

// Recommendation is a suggestion the agent wants confirmed before it acts: the question, what it
// recommends and how sure it is, Accept and the other options.
type Recommendation struct {
	ID         string // the element id; htmx swaps the card under it
	Question   string // "Want me to place this restock order?"
	Text       string // the recommendation, as Markdown
	Confidence float64
	// Accept is the label of the recommended choice, default Texts.Accept.
	Accept string
	// Options are the alternatives, shown folded under Texts.Alternatives.
	Options []string
	// Post is the URL the choice is posted to: the form field "accept" ("true") for the recommendation,
	// or "choice" with one of Options.
	Post string
	// Accepted or Chosen (an option) is the reader's decision; a decided card shows it instead of the
	// buttons.
	Accepted bool
	Chosen   string
}

// Decide applies the posted form and returns the decided card.
func (r Recommendation) Decide(form url.Values) (Recommendation, error) {
	if form.Get("accept") == "true" {
		r.Accepted, r.Chosen = true, ""
		return r, nil
	}
	c := form.Get("choice")
	if !slices.Contains(r.Options, c) {
		return r, ErrChoice
	}
	r.Accepted, r.Chosen = false, c
	return r, nil
}

// Decided reports whether the reader has decided.
func (r Recommendation) Decided() bool { return r.Accepted || r.Chosen != "" }

// Level is the confidence as a word and a color: high (≥ 0.75), medium (≥ 0.4) or low.
func Level(c float64, t Texts) (word, color string) {
	switch {
	case c >= 0.75:
		return t.High, "green"
	case c >= 0.4:
		return t.Medium, "yellow"
	default:
		return t.Low, "orange"
	}
}

// Approval asks the reader several questions one at a time before the agent acts: "1 / 3", the options
// of the current question, Back, Skip and Continue. The questions stay on the server; the answers so
// far travel in the form, so the state needs no storage.
type Approval struct {
	ID        string
	Title     string
	Questions []AskQuestion // single choice; AllowOther adds a text field
	Post      string        // the URL the form is posted to
	// Step is the current question; len(Questions) means done.
	Step int
	// Answers by question id; a skipped question has none.
	Answers map[string]string
}

// Done reports whether every question was answered or skipped.
func (a Approval) Done() bool { return a.Step >= len(a.Questions) }

// Current is the question shown, or the zero AskQuestion when done.
func (a Approval) Current() AskQuestion {
	if a.Done() || a.Step < 0 {
		return AskQuestion{}
	}
	return a.Questions[a.Step]
}

// Advance applies the posted form — the state (fields "step" and "answers") and the action ("action":
// "continue", "skip" or "back", with "answer" or "other") — to a, which supplies the questions. It
// returns the next state.
func (a Approval) Advance(form url.Values) (Approval, error) {
	step, err := strconv.Atoi(form.Get("step"))
	if err != nil || step < 0 || step > len(a.Questions) {
		return a, fmt.Errorf("assistant: bad step %q", form.Get("step"))
	}
	answers := map[string]string{}
	if s := form.Get("answers"); s != "" {
		if err := json.Unmarshal([]byte(s), &answers); err != nil {
			return a, fmt.Errorf("assistant: bad answers: %w", err)
		}
	}
	// Answers only to questions that exist: the field came from the page.
	for id := range answers {
		if !slices.ContainsFunc(a.Questions, func(q AskQuestion) bool { return q.ID == id }) {
			delete(answers, id)
		}
	}
	a.Step, a.Answers = step, answers
	if a.Done() {
		return a, nil
	}
	q := a.Current()
	switch form.Get("action") {
	case "back":
		a.Step = max(0, a.Step-1)
	case "skip":
		delete(a.Answers, q.ID)
		a.Step++
	case "continue", "":
		v := form.Get("answer")
		if other := form.Get("other"); q.AllowOther && other != "" && (v == "" || v == otherValue) {
			v = other
		} else if !slices.ContainsFunc(q.Options, func(o AskOption) bool { return o.Label == v }) {
			return a, ErrChoice
		}
		a.Answers[q.ID] = v
		a.Step++
	default:
		return a, fmt.Errorf("assistant: unknown action %q", form.Get("action"))
	}
	return a, nil
}

// otherValue marks the "Other…" radio; its text comes from the field "other".
const otherValue = "\x00other"

func (a Approval) answersJSON() string {
	if a.Answers == nil {
		return "{}"
	}
	b, _ := json.Marshal(a.Answers)
	return string(b)
}

// Diff is a table of edits the agent proposes: each row shows the old and the new value of the cells it
// changes, and the reader accepts or rejects it, row by row or all at once.
type Diff struct {
	ID      string
	Title   string
	Columns []string
	Rows    []DiffRow
	Post    string // the URL the decisions are posted to
}

// DiffRow is one row; Cells follow Columns.
type DiffRow struct {
	ID    string
	Cells []DiffCell
	// State is "" (pending), "accepted" or "rejected".
	State string
}

// DiffCell is a cell's current value and, if the agent changes it, the proposed one.
type DiffCell struct {
	Old string
	New string // empty: unchanged
}

// Changed reports whether the cell has a proposed value different from the current one.
func (c DiffCell) Changed() bool { return c.New != "" && c.New != c.Old }

// Pending counts the rows not decided yet.
func (d Diff) Pending() int {
	n := 0
	for _, r := range d.Rows {
		if r.State == "" {
			n++
		}
	}
	return n
}

// Resolve applies a decision from the posted form: "row" (a row id, or "*" for every pending row) and
// "decision" ("accept" or "reject"). It returns a copy; d is left as it was.
func (d Diff) Resolve(form url.Values) (Diff, error) {
	state := map[string]string{"accept": "accepted", "reject": "rejected"}[form.Get("decision")]
	if state == "" {
		return d, fmt.Errorf("assistant: bad decision %q", form.Get("decision"))
	}
	id := form.Get("row")
	d.Rows = slices.Clone(d.Rows)
	found := false
	for i, r := range d.Rows {
		if (id == "*" && r.State == "") || r.ID == id {
			d.Rows[i].State = state
			found = true
		}
	}
	if !found && id != "*" {
		return d, fmt.Errorf("assistant: no row %q", id)
	}
	return d, nil
}

// TaskList is the agent's tasks with their live status. While a task runs and Poll is set, the list
// asks Poll for itself again every Every (default one second) and stops once all are finished.
type TaskList struct {
	ID    string
	Title string
	Tasks []Task
	Poll  string
	Every int // milliseconds
}

// Task is one task. Status is "pending", "running", "done" or "failed".
type Task struct {
	Title    string
	Status   string
	Detail   string   // "12/12", "7 SKUs"
	Progress float64  // 0–1 for a bar while running; 0 shows none
	Steps    []string // sub-steps, folded under the task
}

// Running reports whether a task is pending or running.
func (l TaskList) Running() bool {
	return slices.ContainsFunc(l.Tasks, func(t Task) bool { return t.Status == "pending" || t.Status == "running" })
}

func (l TaskList) trigger() string {
	ms := l.Every
	if ms <= 0 {
		ms = 1000
	}
	return "every " + strconv.Itoa(ms) + "ms"
}

func tasksFrom(p map[string]any) TaskList {
	l := TaskList{Title: str(p, "title")}
	for _, t := range rows(p, "tasks") {
		pr, _ := number(t["progress"])
		l.Tasks = append(l.Tasks, Task{Title: str(t, "title"), Status: str(t, "status"), Detail: cellText(t["detail"]),
			Progress: pr, Steps: strs(t, "steps")})
	}
	return l
}
