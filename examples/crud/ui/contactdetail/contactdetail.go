// Package contactdetail is a contact's page: the other contacts on the left, the contact in the middle,
// facts and its activity in the inspector (blocks/split). The activity loads lazily; the page refreshes
// its panes after the other components' events.
package contactdetail

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactdelete"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactform"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactrow"
	"github.com/mind-vm/tabler-templ/htmx"
)

// Routes: Page is a page route (no /ui), Activity a fragment.
var (
	Page     = htmx.Route("GET", "/crud/contacts/{id}")
	Activity = htmx.Route("GET", "/crud/ui/contactdetail/{id}/activity")
)

// IDs.
const (
	RootID     = "contactdetail"
	ActivityID = "contactdetail-activity"
)

var refreshOn = []string{contactform.EventSaved, contactrow.EventStatusChanged, contactdelete.EventRestored}

func trigger() string {
	parts := make([]string, len(refreshOn))
	for i, e := range refreshOn {
		parts[i] = e + " from:body"
	}
	return strings.Join(parts, ", ")
}

// Service is what the component needs from the domain.
type Service interface {
	Get(ctx context.Context, id string) (contacts.Contact, error)
	List(ctx context.Context, q contacts.Query) ([]contacts.Contact, error)
	Activity(ctx context.Context, id string) ([]contacts.Event, error)
}

type Options struct {
	Frame func(title string, body templ.Component) templ.Component
	Back  string // the list
}

type Props struct {
	Contact contacts.Contact
	Others  []contacts.Contact
	Back    string
}

// Register wires the component's endpoints.
func Register(mux *http.ServeMux, svc Service, o Options) {
	h := &handler{svc: svc, o: o}
	mux.HandleFunc(Page.Pattern(), h.page)
	mux.HandleFunc(Activity.Pattern(), h.activity)
}

type handler struct {
	svc Service
	o   Options
}

func (h *handler) contact(w http.ResponseWriter, r *http.Request) (contacts.Contact, bool) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		if htmx.IsHTMX(r) {
			htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact no longer exists."))
		} else {
			http.NotFound(w, r)
		}
		return c, false
	}
	if err != nil {
		htmx.Error(w, r, err)
		return c, false
	}
	return c, true
}

func (h *handler) page(w http.ResponseWriter, r *http.Request) {
	c, ok := h.contact(w, r)
	if !ok {
		return
	}
	others, err := h.svc.List(r.Context(), contacts.Query{})
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, View(Props{Contact: c, Others: others, Back: h.o.Back}),
		htmx.FullPage(func(b templ.Component) templ.Component { return h.o.Frame(c.Name, b) }))
}

func (h *handler) activity(w http.ResponseWriter, r *http.Request) {
	c, ok := h.contact(w, r)
	if !ok {
		return
	}
	ev, err := h.svc.Activity(r.Context(), c.ID)
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, ActivityList(ev),
		htmx.FullPage(func(b templ.Component) templ.Component { return h.o.Frame("Activity of "+c.Name, b) }))
}
