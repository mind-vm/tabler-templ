// Package uierr says once how the contacts domain's errors read in the UI, for htmx.Error.
package uierr

import (
	"net/http"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
)

// Rules map the domain's errors to a status and a message.
var Rules = []htmx.Rule{
	htmx.Is(contacts.ErrNotFound, http.StatusNotFound, "This contact no longer exists."),
	htmx.Is(contacts.ErrConflict, http.StatusConflict, "Someone else changed this contact in the meantime. Close the form and open it again."),
}
