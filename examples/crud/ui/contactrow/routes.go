// Package contactrow is the status cell of a contact row, edited in place: the badge with an edit
// button (View), a select that saves on change (Edit). Only the cell is swapped (item swap); the table
// around it stays.
package contactrow

import (
	"context"
	"net/http"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
)

const base = "/crud/ui/contactrow"

// Routes.
var (
	Status     = htmx.Route("GET", base+"/{id}/status")
	EditStatus = htmx.Route("GET", base+"/{id}/status/edit")
	SaveStatus = htmx.Route("POST", base+"/{id}/status")
)

// CellID is the id of a contact's status cell; the contact id is the instance.
func CellID(contactID string) string { return htmx.ID("contactrow-status", contactID) }

// EventStatusChanged fires after a status was saved; its detail is the contact id.
const EventStatusChanged = "contact:status"

// Service is what the component needs from the domain.
type Service interface {
	Get(ctx context.Context, id string) (contacts.Contact, error)
	SetStatus(ctx context.Context, id string, st contacts.Status) (contacts.Contact, error)
}

// Register wires the component's endpoints.
func Register(mux *http.ServeMux, svc Service) {
	h := &handler{svc: svc}
	mux.HandleFunc(Status.Pattern(), h.view)
	mux.HandleFunc(EditStatus.Pattern(), h.edit)
	mux.HandleFunc(SaveStatus.Pattern(), h.save)
}
