package assistant

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func frag(t *testing.T, c templ.Component) string {
	t.Helper()
	out, err := Fragment(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestRecommendationDecide(t *testing.T) {
	r := Recommendation{ID: "rec", Question: "Order?", Confidence: 0.8, Options: []string{"Half", "Later"}, Post: "/rec"}
	if _, err := r.Decide(url.Values{"choice": {"Everything"}}); !errors.Is(err, ErrChoice) {
		t.Fatalf("an option not offered was taken: %v", err)
	}
	got, err := r.Decide(url.Values{"accept": {"true"}})
	if err != nil || !got.Accepted || !got.Decided() {
		t.Fatalf("accept: %+v %v", got, err)
	}
	got, err = r.Decide(url.Values{"choice": {"Later"}})
	if err != nil || got.Accepted || got.Chosen != "Later" {
		t.Fatalf("alternative: %+v %v", got, err)
	}
	out := frag(t, RecommendationCard(got))
	if strings.Contains(out, `name="accept"`) || !strings.Contains(out, "Later") {
		t.Fatalf("a decided card still offers the buttons or hides the choice:\n%s", out)
	}
	open := frag(t, RecommendationCard(r))
	for _, want := range []string{`hx-post="/rec"`, `hx-swap="outerHTML"`, `name="accept" value="true"`, `value="Half"`, "High confidence"} {
		if !strings.Contains(open, want) {
			t.Errorf("open card lacks %s", want)
		}
	}
}

func approval() Approval {
	return Approval{ID: "ap", Post: "/ap", Questions: []AskQuestion{
		{ID: "size", Question: "How many?", Options: []AskOption{{Label: "One"}, {Label: "Three"}}},
		{ID: "when", Question: "When?", Options: []AskOption{{Label: "Now"}}, AllowOther: true},
		{ID: "note", Question: "Note?", Options: []AskOption{{Label: "None"}}},
	}}
}

func TestApprovalAdvance(t *testing.T) {
	a := approval()
	a, err := a.Advance(url.Values{"step": {"0"}, "answers": {"{}"}, "action": {"continue"}, "answer": {"Three"}})
	if err != nil || a.Step != 1 || a.Answers["size"] != "Three" {
		t.Fatalf("continue: %+v %v", a, err)
	}
	// The next request carries the state the card rendered.
	out := frag(t, ApprovalCard(a))
	if !strings.Contains(out, `name="step" value="1"`) || !strings.Contains(out, "2 / 3") {
		t.Fatalf("state not in the form:\n%s", out)
	}
	form := url.Values{"step": {"1"}, "answers": {a.answersJSON()}, "action": {"continue"}, "answer": {otherValue}, "other": {"Friday"}}
	if a, err = approval().Advance(form); err != nil || a.Answers["when"] != "Friday" || a.Answers["size"] != "Three" {
		t.Fatalf("other: %+v %v", a, err)
	}
	if a, err = approval().Advance(url.Values{"step": {"2"}, "answers": {a.answersJSON()}, "action": {"skip"}}); err != nil || !a.Done() {
		t.Fatalf("skip: %+v %v", a, err)
	}
	if _, ok := a.Answers["note"]; ok {
		t.Fatal("a skipped question has an answer")
	}
	if out := frag(t, ApprovalCard(a)); !strings.Contains(out, "skipped") || strings.Contains(out, `name="action"`) {
		t.Fatalf("done card:\n%s", out)
	}
	back, _ := approval().Advance(url.Values{"step": {"2"}, "answers": {"{}"}, "action": {"back"}})
	if back.Step != 1 {
		t.Fatalf("back: step %d", back.Step)
	}
}

func TestApprovalRejectsForgedState(t *testing.T) {
	for name, f := range map[string]url.Values{
		"option":  {"step": {"0"}, "action": {"continue"}, "answer": {"Ninety"}},
		"step":    {"step": {"7"}, "action": {"continue"}},
		"answers": {"step": {"0"}, "answers": {"not json"}},
		"action":  {"step": {"0"}, "action": {"drop table"}},
		// "Other…" without AllowOther is not an answer.
		"other": {"step": {"0"}, "action": {"continue"}, "answer": {otherValue}, "other": {"Seven"}},
	} {
		if _, err := approval().Advance(f); err == nil {
			t.Errorf("%s: forged form accepted", name)
		}
	}
	a, err := approval().Advance(url.Values{"step": {"0"}, "answers": {`{"ghost":"x"}`}, "action": {"skip"}})
	if err != nil || len(a.Answers) != 0 {
		t.Fatalf("an answer to a question that does not exist was kept: %+v", a.Answers)
	}
}

func TestDiffResolve(t *testing.T) {
	d := Diff{ID: "d", Post: "/d", Columns: []string{"Flavor", "Supplier"}, Rows: []DiffRow{
		{ID: "r1", Cells: []DiffCell{{Old: "Mint"}, {Old: "Acme", New: "Dairy Co"}}},
		{ID: "r2", Cells: []DiffCell{{Old: "Plum", New: "Plum sorbet"}, {Old: "Acme"}}},
		{ID: "r3", Cells: []DiffCell{{Old: "Fig", New: ""}, {Old: "Acme", New: "Acme"}}},
	}}
	got, err := d.Resolve(url.Values{"row": {"r1"}, "decision": {"reject"}})
	if err != nil || got.Rows[0].State != "rejected" || d.Rows[0].State != "" {
		t.Fatalf("resolve changed the original or missed the row: %+v %v", got.Rows, err)
	}
	all, _ := got.Resolve(url.Values{"row": {"*"}, "decision": {"accept"}})
	if all.Rows[0].State != "rejected" || all.Rows[1].State != "accepted" || all.Pending() != 0 {
		t.Fatalf("accept all overrode a decision or left rows: %+v", all.Rows)
	}
	if _, err := d.Resolve(url.Values{"row": {"r9"}, "decision": {"accept"}}); err == nil {
		t.Fatal("unknown row accepted")
	}
	if _, err := d.Resolve(url.Values{"row": {"r1"}, "decision": {"maybe"}}); err == nil {
		t.Fatal("unknown decision accepted")
	}
	out := frag(t, DiffCard(d))
	if !strings.Contains(out, `<del class="assistant-old">Acme</del>`) || !strings.Contains(out, "Dairy Co") {
		t.Fatalf("pending change not shown:\n%s", out)
	}
	if strings.Contains(out, `<del class="assistant-old">Fig</del>`) {
		t.Fatal("an unchanged cell shown as an edit")
	}
	if !strings.Contains(out, `hx-vals="{&#34;decision&#34;:&#34;accept&#34;,&#34;row&#34;:&#34;r1&#34;}"`) {
		t.Fatalf("row action not posted:\n%s", out)
	}
	if out := frag(t, DiffCard(all)); strings.Contains(out, "hx-post") {
		t.Fatal("a fully decided diff still offers actions")
	}
}

func TestTaskListPollsWhileRunning(t *testing.T) {
	l := TaskList{ID: "t", Poll: "/tasks", Tasks: []Task{{Title: "Verify", Status: "done"}, {Title: "Draft", Status: "running", Progress: 0.5}}}
	if out := frag(t, TaskListCard(l)); !strings.Contains(out, `hx-get="/tasks"`) || !strings.Contains(out, `hx-trigger="every 1000ms"`) {
		t.Fatalf("running list does not poll:\n%s", out)
	}
	l.Tasks[1].Status = "failed"
	if out := frag(t, TaskListCard(l)); strings.Contains(out, "hx-get") {
		t.Fatal("finished list still polls")
	}
}

func TestNewBlockTypes(t *testing.T) {
	for typ, js := range map[string]string{
		"code":     `{"title":"churn.ts","diff":true,"code":"@@ -1 +1 @@\n-old\n+new"}`,
		"context":  `{"chunks":[{"title":"SOP","source":"Dairy.pdf","kind":"pdf","chars":290,"text":"Keep cold."}]}`,
		"screen":   `{"title":"Agent's screen","url":"https://example.com","status":"working","image":"/shot.png"}`,
		"tasks":    `{"title":"Tasks","tasks":[{"title":"Draft emails","status":"running","progress":0.68}]}`,
		"insights": `{"items":[{"text":"Mint fell","metric":"-4.41%","delta":-0.0441,"values":[3,2,1]},{"text":"Pistachio rose"}]}`,
	} {
		out := render(t, Block{Type: typ, DOMID: "b-" + typ, Payload: payload(t, js)})
		if strings.Contains(out, "Unrecognised block") {
			t.Errorf("%s: drawn as unknown", typ)
		}
	}
	code := render(t, Block{Type: "code", DOMID: "c", Payload: payload(t, `{"diff":true,"code":"+a\n-b\n c"}`)})
	if !strings.Contains(code, "assistant-line-add") || !strings.Contains(code, "assistant-line-del") || !strings.Contains(code, `data-bs-target="#c-code"`) {
		t.Fatalf("diff lines or copy target missing:\n%s", code)
	}
	ins := render(t, Block{Type: "insights", DOMID: "i", Payload: payload(t, `{"items":[{"text":"a"},{"text":"b"},{"text":"c"}]}`)})
	if strings.Count(ins, `type="radio"`) != 3 || strings.Count(ins, "checked") != 1 || !strings.Contains(ins, `for="i-insight-2"`) {
		t.Fatalf("insight pager:\n%s", ins)
	}
}

func TestParts(t *testing.T) {
	th := frag(t, Thinking(ThinkingProps{ID: "th", Duration: "4s", Thoughts: []Thought{{Kind: "search", Title: "Searched trends"}}}))
	if !strings.Contains(th, "Thought for 4s") || !strings.Contains(th, "ti-world-search") {
		t.Fatalf("thinking:\n%s", th)
	}
	src := frag(t, Sources([]Source{{Title: "Trends", URL: "https://www.trends.example/x"}}))
	if !strings.Contains(src, "1 sources") || !strings.Contains(src, "trends.example") {
		t.Fatalf("sources:\n%s", src)
	}
	if frag(t, Sources(nil)) != "" || frag(t, FollowUps(nil)) != "" {
		t.Fatal("empty lists render something")
	}
	if out := frag(t, FollowUps([]string{"Why?"})); !strings.Contains(out, `data-example="Why?"`) {
		t.Fatalf("follow-up does not send:\n%s", out)
	}
	tc := frag(t, ToolChips("tc", []ToolCall{{Name: "search", Detail: "q"}, {Name: "edit", Running: true}}))
	if !strings.Contains(tc, "2 tool calls") || !strings.Contains(tc, "spinner") {
		t.Fatalf("tool chips:\n%s", tc)
	}
}

func TestConversationExtras(t *testing.T) {
	out := frag(t, Conversation(ConversationProps{Run: "/r", Models: []string{"Fast", "Deep"}, Model: "Deep",
		SelectionActions: []SelectionAction{{Label: "Explain", Prompt: "Explain %s"}}}))
	if !strings.Contains(out, `<option value="Deep" selected>`) || !strings.Contains(out, `data-selection-prompt="Explain %s"`) {
		t.Fatalf("model picker or selection toolbar missing:\n%s", out)
	}
}
