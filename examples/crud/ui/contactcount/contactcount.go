// Package contactcount is a badge with the number of contacts (of a status). It knows no other
// component's markup: it refreshes itself on the events the others fire. A page can show it more than
// once; every instance has its own id.
package contactcount

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/events"
	"github.com/mind-vm/tabler-templ/htmx"
)

// Count is the component's only route.
var Count = htmx.Route("GET", "/crud/ui/contactcount")

// RootID is the badge's id for an instance.
func RootID(inst string) string { return htmx.ID("contactcount", inst) }

// trigger: the count may change with any change of a contact.
func trigger() string { return htmx.On(events.Any...) }

// Service is what the component needs from the domain.
type Service interface {
	Count(ctx context.Context, status contacts.Status) (int, error)
}

type Props struct {
	Instance string
	Status   contacts.Status // "" counts all
	Label    string          // after the number, e.g. "active"
	Color    string          // badge color, default "blue"
	N        int
}

// Load reads the count for p.
func Load(ctx context.Context, svc Service, p Props) (Props, error) {
	n, err := svc.Count(ctx, p.Status)
	p.N = n
	return p, err
}

// Register wires the component's endpoint.
func Register(mux *http.ServeMux, svc Service) {
	mux.HandleFunc(Count.Pattern(), func(w http.ResponseWriter, r *http.Request) {
		p, err := Load(r.Context(), svc, Props{Instance: r.FormValue("instance"), Status: contacts.Status(r.FormValue("status")),
			Label: r.FormValue("label"), Color: r.FormValue("color")})
		if err != nil {
			htmx.Error(w, r, err)
			return
		}
		htmx.Respond(w, r, View(p))
	})
}

// Embed loads the count when it renders, for a page frame that has no handler of its own.
func Embed(svc Service, p Props) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		p, err := Load(ctx, svc, p)
		if err != nil {
			return err
		}
		return View(p).Render(ctx, w)
	})
}
