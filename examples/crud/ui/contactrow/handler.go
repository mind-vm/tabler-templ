package contactrow

import (
	"errors"
	"net/http"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
)

type handler struct{ svc Service }

func (h *handler) contact(w http.ResponseWriter, r *http.Request) (contacts.Contact, bool) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact no longer exists."))
		return c, false
	}
	if err != nil {
		htmx.Error(w, r, err)
		return c, false
	}
	return c, true
}

func (h *handler) view(w http.ResponseWriter, r *http.Request) {
	if c, ok := h.contact(w, r); ok {
		htmx.Respond(w, r, View(c))
	}
}

func (h *handler) edit(w http.ResponseWriter, r *http.Request) {
	if c, ok := h.contact(w, r); ok {
		htmx.Respond(w, r, Edit(c))
	}
}

func (h *handler) save(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.SetStatus(r.Context(), r.PathValue("id"), contacts.Status(r.FormValue("status")))
	var fe contacts.FieldErrors
	switch {
	case errors.As(err, &fe):
		htmx.Error(w, r, htmx.Fail(http.StatusUnprocessableEntity, fe["status"]))
		return
	case errors.Is(err, contacts.ErrNotFound):
		htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact no longer exists."))
		return
	case err != nil:
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, View(c), htmx.Trigger(EventStatusChanged, c.ID))
}
