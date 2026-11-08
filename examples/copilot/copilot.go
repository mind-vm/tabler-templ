// Package copilot is a demo app with an AI copilot beside its pages: orders and customers, and a
// panel at the right where the reader asks questions. The copilot reads through tools, shows what it
// read, fills in the app's forms and points at elements; it never saves. Everything it asks to do goes
// through the app's Catalog (catalog.go), which refuses what is not on the list: bank details, e-mail
// addresses, foreign pages. Answers are rendered without links and images (assistant.AnswerStrict).
//
// The parts, all from blocks/assistant:
//
//   - Copilot and CopilotButton: the panel, open from page to page.
//   - Stream with Turn, ToolChips, BlockView, AnswerStrict: the answer, streamed as HTML fragments.
//   - Catalog and Stream.Action: the checked actions, run on the page by ActionScript, also after
//     going to another page.
//   - Transcript: the server keeps the conversation, so the panel shows it again on the next page.
//
// The agent (agent.go) is a stand-in for a language model, so the demo runs offline.
//
// cmd/demo mounts it at /copilot. To run it alone:
//
//	mux := http.NewServeMux()
//	copilot.Register(mux)
//	http.ListenAndServe("localhost:8080", http.NewCrossOriginProtection().Handler(mux))
package copilot

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/blocks/assistant"
	"github.com/mind-vm/tabler-templ/blocks/shell"
)

// App is the demo with its data and the conversations of its sessions.
type App struct {
	store *Store
	agent *agent

	mu       sync.Mutex
	sessions map[string]*session
}

type session struct {
	transcript *assistant.Transcript
	busy       bool
	turns      int // numbers the turns, so element ids stay unique on the page
}

const cookieName = "copilot-session"

// Register mounts the demo with fresh data.
func Register(mux *http.ServeMux) *App { return RegisterWith(mux, NewStore(), 350*time.Millisecond) }

// RegisterWith mounts the demo on store; pause is the agent's pause between steps.
func RegisterWith(mux *http.ServeMux, store *Store, pause time.Duration) *App {
	a := &App{store: store, agent: &agent{store: store, role: "office", pause: pause}, sessions: map[string]*session{}}
	mux.HandleFunc("GET "+Home, a.list)
	mux.HandleFunc("POST "+orders, a.create)
	mux.HandleFunc("GET "+orders+"/{id}", a.order)
	mux.HandleFunc("POST "+orders+"/{id}", a.update)
	mux.HandleFunc("GET "+customer+"/{id}", a.customer)
	mux.HandleFunc("POST "+customer+"/{id}", a.updateCustomer)
	mux.HandleFunc("POST "+Home+"/assistant/run", a.run)
	return a
}

// session finds the reader's conversation by cookie, or starts one.
func (a *App) session(w http.ResponseWriter, r *http.Request) (string, *session) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if c, err := r.Cookie(cookieName); err == nil {
		if s, ok := a.sessions[c.Value]; ok {
			return c.Value, s
		}
	}
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	id := hex.EncodeToString(b)
	s := &session{transcript: &assistant.Transcript{Max: 30}}
	a.sessions[id] = s
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: id, Path: Home, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return id, s
}

func (a *App) render(w http.ResponseWriter, r *http.Request, p page, body templ.Component) {
	id, s := a.session(w, r)
	p.Session, p.Transcript = id, s.transcript
	w.Header().Set("Cache-Control", "no-store")
	_ = a.frame(p).Render(templ.WithChildren(r.Context(), body), w)
}

func (a *App) customerName(id string) string {
	c, _ := a.store.Customer(id)
	return c.Name
}

func (a *App) list(w http.ResponseWriter, r *http.Request) {
	overdue := r.URL.Query().Get("overdue") == "1"
	a.listWith(w, r, overdue, OrderInput{Due: a.store.Today().AddDate(0, 0, 14).Format(dateLayout)}, nil)
}

func (a *App) listWith(w http.ResponseWriter, r *http.Request, overdue bool, in OrderInput, errs FieldErrors) {
	active, title := "orders", "Orders"
	if overdue {
		active, title = "overdue", "Overdue orders"
	}
	if len(errs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	a.render(w, r, page{Title: title, Active: active, Actions: newOrderButton(), Flash: r.URL.Query().Get("saved"),
		Crumbs: []shell.Crumb{{Text: "Orders"}}}, a.orderList(a.store.Orders("", overdue), overdue, in, errs))
}

func input(r *http.Request) OrderInput {
	return OrderInput{Customer: r.PostFormValue("customer"), Title: r.PostFormValue("title"), Amount: r.PostFormValue("amount"),
		Status: r.PostFormValue("status"), Due: r.PostFormValue("due"), Note: r.PostFormValue("note")}
}

func (a *App) create(w http.ResponseWriter, r *http.Request) {
	in := input(r)
	o, errs := a.store.Create(in)
	if len(errs) > 0 {
		a.listWith(w, r, false, in, errs)
		return
	}
	http.Redirect(w, r, orders+"/"+o.ID+"?saved=Order+"+o.ID+"+created.", http.StatusSeeOther)
}

func (a *App) order(w http.ResponseWriter, r *http.Request) {
	o, ok := a.store.Order(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	in := OrderInput{Status: string(o.Status), Due: o.Due.Format(dateLayout), Note: o.Note}
	a.orderWith(w, r, o, in, nil)
}

func (a *App) orderWith(w http.ResponseWriter, r *http.Request, o Order, in OrderInput, errs FieldErrors) {
	if len(errs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	a.render(w, r, page{Title: o.ID + " · " + o.Title, Active: "orders", Flash: r.URL.Query().Get("saved"),
		Crumbs: []shell.Crumb{{Text: "Orders", Href: Home}, {Text: o.ID}}}, a.orderPage(o, in, errs))
}

func (a *App) update(w http.ResponseWriter, r *http.Request) {
	in := input(r)
	o, errs, err := a.store.Update(r.PathValue("id"), in)
	switch {
	case errors.Is(err, ErrNotFound):
		http.NotFound(w, r)
	case len(errs) > 0:
		a.orderWith(w, r, o, in, errs)
	default:
		http.Redirect(w, r, orders+"/"+o.ID+"?saved=Saved.", http.StatusSeeOther)
	}
}

func (a *App) customer(w http.ResponseWriter, r *http.Request) {
	c, ok := a.store.Customer(r.PathValue("id"))
	if !ok {
		http.NotFound(w, r)
		return
	}
	a.customerWith(w, r, c, nil)
}

func (a *App) customerWith(w http.ResponseWriter, r *http.Request, c Customer, errs FieldErrors) {
	if len(errs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	a.render(w, r, page{Title: c.Name, Active: c.ID, Flash: r.URL.Query().Get("saved"),
		Crumbs: []shell.Crumb{{Text: "Customers"}, {Text: c.Name}}}, a.customerPage(c, errs))
}

func (a *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
	c, errs, err := a.store.UpdateCustomer(r.PathValue("id"), r.PostFormValue("name"), r.PostFormValue("city"),
		r.PostFormValue("email"), r.PostFormValue("iban"))
	switch {
	case errors.Is(err, ErrNotFound):
		http.NotFound(w, r)
	case len(errs) > 0:
		c.Name = r.PostFormValue("name")
		a.customerWith(w, r, c, errs)
	default:
		http.Redirect(w, r, customer+"/"+c.ID+"?saved=Saved.", http.StatusSeeOther)
	}
}

// run answers one message as a stream of fragments and actions, recorded in the session's transcript.
func (a *App) run(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&in); err != nil || in.Message == "" || len(in.Message) > 4000 {
		http.Error(w, "message missing or too long", http.StatusBadRequest)
		return
	}
	_, sess := a.session(w, r)
	a.mu.Lock()
	if sess.busy {
		a.mu.Unlock()
		http.Error(w, "the copilot is still answering", http.StatusConflict)
		return
	}
	sess.busy = true
	sess.turns++
	n := sess.turns
	a.mu.Unlock()
	defer func() { a.mu.Lock(); sess.busy = false; a.mu.Unlock() }()

	s := assistant.Open(w)
	if s == nil {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	s.Record(sess.transcript)
	ctx := r.Context()
	_ = s.Turn(ctx, in.Message)
	t := &turn{ctx: ctx, s: s, n: n, a: a.agent}
	start := time.Now()
	a.agent.plan(t, in.Message)
	_ = s.Send(ctx, assistant.SlotStatus, "", assistant.ModeSet, assistant.Usage(assistant.UsageProps{
		Text: time.Since(start).Round(100*time.Millisecond).String() + " · no language model, demo agent"}))
}
