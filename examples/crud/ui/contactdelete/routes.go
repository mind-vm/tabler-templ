// Package contactdelete deletes a contact after asking: the question opens in the page's dialog
// (a page without JavaScript), the deletion fires EventDeleted and shows a toast with Undo, which
// restores the contact and fires EventRestored.
package contactdelete

import (
	"context"
	"net/http"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
)

const base = "/crud/ui/contactdelete"

// Routes. Delete is a POST, not a DELETE, so the form works without JavaScript.
var (
	Confirm = htmx.Route("GET", base+"/{id}")
	Delete  = htmx.Route("POST", base+"/{id}")
	Restore = htmx.Route("POST", base+"/{id}/restore")
)

// Events; the detail is the contact id.
const (
	EventDeleted  = "contact:deleted"
	EventRestored = "contact:restored"
)

// Service is what the component needs from the domain.
type Service interface {
	Get(ctx context.Context, id string) (contacts.Contact, error)
	Delete(ctx context.Context, id string) (contacts.Contact, error)
	Restore(ctx context.Context, id string) (contacts.Contact, error)
}

type Options struct {
	Frame func(title string, body templ.Component) templ.Component
	// Done is where the browser goes after deleting from a page that shows the contact itself (the
	// question's "next"), and after deleting without htmx.
	Done string
}

// Register wires the component's endpoints.
func Register(mux *http.ServeMux, svc Service, o Options) {
	h := &handler{svc: svc, o: o}
	mux.HandleFunc(Confirm.Pattern(), h.confirm)
	mux.HandleFunc(Delete.Pattern(), h.delete)
	mux.HandleFunc(Restore.Pattern(), h.restore)
}
