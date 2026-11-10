// Package assistant draws an AI assistant's conversation on the server: the composer, one turn per
// message (the user's message, the agent's steps, its questions, the result blocks and the answer), and
// the blocks an agent emits — markdown, a table, a chart, a panel of figures.
//
// Everything is rendered as HTML fragments and streamed to the page over server-sent events; Script is
// the page's whole client (about 150 lines): it posts a message, reads the stream and puts each
// fragment in its slot. The package knows no model and no agent framework. The block payloads follow
// the contract of thinking-core's ext.BlocksPack (github.com/mind-vm/thinking-core), so an agent built on
// it can be drawn as it is, but any agent can fill the same structs.
//
//	@assistant.Conversation(assistant.ConversationProps{Run: "/api/run", Reply: "/api/reply", Examples: examples})
//	@assistant.Script()
//
// A run handler opens a Stream and sends fragments:
//
//	s := assistant.Open(w)
//	s.Send(ctx, assistant.SlotTurn, "", "", assistant.Turn(msg))
//	s.Send(ctx, assistant.SlotSteps, "step-1-1", "", assistant.Step(assistant.StepProps{ID: "step-1-1", N: 1, Code: code}))
//	s.Send(ctx, assistant.SlotBlocks, b.DOMID, "", assistant.BlockView(b))
//	s.Send(ctx, assistant.SlotAnswer, "", assistant.ModeSet, assistant.Answer(text))
//	s.Send(ctx, assistant.SlotStatus, "", assistant.ModeSet, assistant.Usage(assistant.UsageProps{Steps: 3}))
package assistant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/chart"
	"github.com/mind-vm/tabler-templ/chart/svg"
)

// Slots of a turn, where Stream.Send puts a fragment.
const (
	SlotTurn   = "turn"   // a new Turn, appended to the conversation
	SlotSteps  = "steps"  // the agent's steps (scripts, tool calls)
	SlotAsks   = "asks"   // questions the agent asks the user
	SlotBlocks = "blocks" // result blocks
	SlotAnswer = "answer" // the closing text
	SlotStatus = "status" // "working" while running, then Usage or Failure
)

// ModeSet replaces the slot's content instead of adding to it.
const ModeSet = "set"

// Block is one result block as the page draws it.
type Block struct {
	// Type is the block type: "markdown", "table", "chart" or "stats". Another type is shown as JSON.
	Type string
	// Payload is the block's spec in JSON's shapes (float64, []any, map[string]any); Normalize turns
	// any Go value into them.
	Payload map[string]any
	// DOMID is the element the block occupies. A block sent again with the same DOMID replaces the
	// earlier copy (a revised chart, a table that grows) instead of adding one.
	DOMID string
	// ChartID names the chart element of a chart block. Give every send a fresh one, unlike DOMID: the
	// replaced copy's chart is still being torn down when the new one starts.
	ChartID string
	// Revised marks a block that replaces one from an earlier turn, so the reader knows the figure up
	// the page they remember is gone.
	Revised bool
}

// AskQuestion is one question of an "questions" ask; the json tags are thinking-core's sandbox.Question.
type AskQuestion struct {
	ID          string      `json:"id"`
	Header      string      `json:"header,omitempty"`
	Question    string      `json:"question"`
	Options     []AskOption `json:"options"`
	MultiSelect bool        `json:"multiSelect,omitempty"`
	AllowOther  bool        `json:"allowOther,omitempty"`
}

// AskOption is one choice of an AskQuestion.
type AskOption struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// AskRequest is a question the agent waits on.
type AskRequest struct {
	RequestID string `json:"requestId"`
	// InputType is "text" (default), "confirm" (yes or no), "choice" (one of Options) or "questions"
	// (a form of Questions).
	InputType   string        `json:"inputType"`
	Message     string        `json:"message,omitempty"`
	Options     []string      `json:"options"`
	Placeholder string        `json:"placeholder"`
	Questions   []AskQuestion `json:"questions"`
}

// AskFromArgs decodes an ask from an event's arguments (thinking-core's input_request event: the
// message is the event's content). It goes through JSON, so typed values and decoded JSON both work.
func AskFromArgs(message string, args map[string]any) (AskRequest, error) {
	var r AskRequest
	b, err := json.Marshal(args)
	if err != nil {
		return r, err
	}
	if err := json.Unmarshal(b, &r); err != nil {
		return r, err
	}
	if message != "" {
		r.Message = message
	}
	return r, nil
}

// Texts are the texts the package renders. Set them once per request with WithTexts; empty fields fall
// back to English.
type Texts struct {
	Placeholder   string // the composer's placeholder
	Send          string
	Working       string
	Revised       string // badge on a replaced block
	NoRows        string
	RowsShown     string // fmt with the shown and the total count: "Showing %d of %s rows."
	Unknown       string // title of a block of unknown type, before the type
	Script        string // a step's title, before its number
	Ran           string
	Retried       string // a step whose result is an error
	Failed        string // title of Failure
	Asking        string
	SendAnswers   string
	Yes           string
	No            string
	Other         string
	OtherHint     string
	YourAnswer    string
	ChooseSeveral string
	YouAnswered   string
	MissingAnswer string // fmt with the question: "Answer “%s” first."
	Usage         string // fmt with steps, prompt and completion tokens, blocks
	Thinking      string // a running Thinking's title
	ThoughtFor    string // fmt with the duration: "Thought for %s"
	Sources       string // fmt with the count: "%d sources"
	FollowUps     string
	ToolCalls     string // fmt with the count: "%d tool calls"
	Copy          string
	Copied        string
	Chunks        string // fmt with the count: "%d chunks"
	Characters    string // fmt with the count: "%d characters"
	Insights      string
	Of            string // fmt with the position and the count: "%d / %d"
	Previous      string
	Next          string
	Accept        string
	Alternatives  string
	Confidence    string // fmt with Texts.High, Medium or Low: "%s confidence"
	High          string
	Medium        string
	Low           string
	Chose         string // a decided Recommendation, before the choice
	Skip          string
	Continue      string
	Back          string
	Skipped       string
	Reject        string
	AcceptAll     string
	RejectAll     string
	Accepted      string
	Rejected      string
	Pending       string // fmt with the count: "%d pending"
	Model         string // aria-label of the model picker
	Selection     string // aria-label of the selection toolbar
	ActionOpen    string // an action card's title, before the action's: "Open page"
	ActionFill    string
	ActionShow    string
	RunAgain      string // the card's button
	Filled        string // banner over a form the assistant filled in
	NotOnPage     string // fmt with the action's title: "%s is not on this page."
	NotFilled     string // fmt with the fields: "Not filled in: %s. Please enter by hand."
	Unchanged     string // the form already held the values
	Refused       string // ActionRefused's title, before the action's name
}

// English are the default texts.
var English = Texts{
	Placeholder:   "Ask something…",
	Send:          "Send",
	Working:       "working",
	Revised:       "revised",
	NoRows:        "No rows.",
	RowsShown:     "Showing %d of %s rows.",
	Unknown:       "Unrecognised block: ",
	Script:        "Script",
	Ran:           "ran",
	Retried:       "error → retried",
	Failed:        "The run failed",
	Asking:        "The agent is asking",
	SendAnswers:   "Send answers",
	Yes:           "Yes",
	No:            "No",
	Other:         "Other…",
	OtherHint:     "Type your own answer",
	YourAnswer:    "Your answer",
	ChooseSeveral: " (choose several)",
	YouAnswered:   "You answered",
	MissingAnswer: "Answer “%s” first.",
	Usage:         "%d steps · %d+%d tokens · %d blocks",
	Thinking:      "Thinking",
	ThoughtFor:    "Thought for %s",
	Sources:       "%d sources",
	FollowUps:     "Follow-ups",
	ToolCalls:     "%d tool calls",
	Copy:          "Copy",
	Copied:        "Copied",
	Chunks:        "%d chunks",
	Characters:    "%d characters",
	Insights:      "Insights",
	Of:            "%d / %d",
	Previous:      "Previous",
	Next:          "Next",
	Accept:        "Accept",
	Alternatives:  "Alternatives",
	Confidence:    "%s confidence",
	High:          "High",
	Medium:        "Medium",
	Low:           "Low",
	Chose:         "You chose",
	Skip:          "Skip",
	Continue:      "Continue",
	Back:          "Back",
	Skipped:       "skipped",
	Reject:        "Reject",
	AcceptAll:     "Accept all",
	RejectAll:     "Reject all",
	Accepted:      "Accepted",
	Rejected:      "Rejected",
	Pending:       "%d pending",
	Model:         "Model",
	Selection:     "Ask about the selection",
	ActionOpen:    "Open page",
	ActionFill:    "Fill in",
	ActionShow:    "Show",
	RunAgain:      "Run again",
	Filled:        "Filled in by the assistant: check the values and save yourself.",
	NotOnPage:     "%s is not on this page.",
	NotFilled:     "Not filled in: %s. Please enter by hand.",
	Unchanged:     "The form already holds these values; nothing changed.",
	Refused:       "Not allowed",
}

// ChartOptions choose how chart blocks are drawn.
type ChartOptions struct {
	// SVG draws chart blocks on the server with package chart/svg instead of with ApexCharts, so the page
	// needs no chart library. The variants the SVG charts do not have (radar, scatter, bubble, heatmap)
	// are shown as a table of their data; composite, funnel and sparkline are drawn.
	SVG bool
	// Format writes the numbers of SVG charts, e.g. svg.German; default English.
	Format svg.Format
}

type chartOptionsKey struct{}

// WithChartOptions sets how the request's chart blocks are drawn; without it they use ApexCharts.
func WithChartOptions(ctx context.Context, o ChartOptions) context.Context {
	return context.WithValue(ctx, chartOptionsKey{}, o)
}

// ChartOptionsFrom returns the request's chart options, the zero value (ApexCharts) if none were set.
func ChartOptionsFrom(ctx context.Context) ChartOptions {
	o, _ := ctx.Value(chartOptionsKey{}).(ChartOptions)
	return o
}

type textsKey struct{}

// WithTexts sets the texts for the request.
func WithTexts(ctx context.Context, t Texts) context.Context {
	return context.WithValue(ctx, textsKey{}, t)
}

// T returns the request's texts, empty fields filled from English.
func T(ctx context.Context) Texts {
	t, ok := ctx.Value(textsKey{}).(Texts)
	if !ok {
		return English
	}
	v, en := reflect.ValueOf(&t).Elem(), reflect.ValueOf(English)
	for i := range v.NumField() {
		if f := v.Field(i); f.Kind() == reflect.String && f.String() == "" {
			f.SetString(en.Field(i).String())
		}
	}
	return t
}

// Fragment renders a component as a fragment of a page that is already loaded.
//
// Components render their page-wide scripts once per render (chart.Chart its start-up script), and a
// fragment is a render of its own, so every streamed chart would carry the whole script again. Fragment
// counts those scripts as rendered: the page has them from Script or Conversation.
func Fragment(ctx context.Context, c templ.Component) (string, error) {
	ctx = templ.InitializeContext(ctx)
	if err := chart.Script().Render(ctx, io.Discard); err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := c.Render(ctx, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Stream writes fragments as server-sent events in the shape Script reads:
//
//	data: {"type":"html","args":{"slot":"blocks","id":"block-1","mode":"","html":"…"}}
//	data: {"type":"action","args":{"id":"action-1","name":"order-new","kind":"fill",…}}
//
// It is safe for use from several goroutines.
type Stream struct {
	mu  sync.Mutex
	w   http.ResponseWriter
	f   http.Flusher
	rec *Transcript
}

// Open starts the event stream. It returns nil when w cannot flush.
func Open(w http.ResponseWriter) *Stream {
	f, ok := w.(http.Flusher)
	if !ok {
		return nil
	}
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	f.Flush()
	return &Stream{w: w, f: f}
}

// Send renders c and sends it into slot. id is the element the fragment replaces if the page already
// shows it (empty: always add); mode ModeSet replaces the slot's content. A component that fails to
// render is sent as a Failure.
func (s *Stream) Send(ctx context.Context, slot, id, mode string, c templ.Component) error {
	html, err := Fragment(ctx, c)
	if err != nil {
		html, _ = Fragment(ctx, Failure(err.Error()))
	}
	return s.SendHTML(slot, id, mode, html)
}

// SendHTML sends an already rendered fragment.
func (s *Stream) SendHTML(slot, id, mode, html string) error {
	if s.rec != nil && slot != SlotTurn {
		s.rec.add(slot, id, mode, html)
	}
	b, err := json.Marshal(map[string]any{"type": "html", "args": map[string]string{"slot": slot, "id": id, "mode": mode, "html": html}})
	if err != nil {
		return err
	}
	return s.write(b)
}

func (s *Stream) write(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, err := fmt.Fprintf(s.w, "data: %s\n\n", b); err != nil {
		return err
	}
	s.f.Flush()
	return nil
}
