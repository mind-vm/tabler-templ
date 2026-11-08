// Package contactform creates and edits a contact. With htmx the form opens in the page's dialog
// (htmx.Setup) and closes after saving; without JavaScript the same routes show it as a page and
// redirect after saving (Post/Redirect/Get). Invalid input comes back as 422 with the errors at the
// fields; an edit of a contact someone else changed meanwhile is a 409 toast.
package contactform

import (
	"context"
	"net/http"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
)

const base = "/crud/ui/contactform"

// Routes.
var (
	New    = htmx.Route("GET", base+"/new")
	Edit   = htmx.Route("GET", base+"/{id}/edit")
	Create = htmx.Route("POST", base)
	Update = htmx.Route("POST", base+"/{id}")
)

// FormID is the form's id; a 422 response replaces the form under it.
const FormID = "contactform"

// Service is what the component needs from the domain.
type Service interface {
	Get(ctx context.Context, id string) (contacts.Contact, error)
	Create(ctx context.Context, c contacts.Contact) (contacts.Contact, error)
	Update(ctx context.Context, c contacts.Contact) (contacts.Contact, error)
}

type Options struct {
	// Frame puts the form into a page for requests without htmx.
	Frame htmx.Frame
	// Done is where the browser goes after saving without htmx, e.g. the list.
	Done string
}

// Register wires the component's endpoints.
func Register(mux *http.ServeMux, svc Service, o Options) {
	h := &handler{svc: svc, o: o}
	mux.HandleFunc(New.Pattern(), h.new)
	mux.HandleFunc(Edit.Pattern(), h.edit)
	mux.HandleFunc(Create.Pattern(), h.create)
	mux.HandleFunc(Update.Pattern(), h.update)
}
