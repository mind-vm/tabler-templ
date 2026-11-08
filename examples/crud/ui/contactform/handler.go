package contactform

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/htmx"
	"github.com/mind-vm/tabler-templ/toast"
)

type handler struct {
	svc Service
	o   Options
}

func (h *handler) page(title string) htmx.Option {
	return htmx.FullPage(func(body templ.Component) templ.Component { return h.o.Frame(title, body) })
}

func (h *handler) props(r *http.Request, c contacts.Contact, errs contacts.FieldErrors) Props {
	return Props{Contact: c, Errors: errs, InDialog: htmx.IsHTMX(r), Cancel: h.o.Done}
}

func (h *handler) new(w http.ResponseWriter, r *http.Request) {
	htmx.Respond(w, r, View(h.props(r, contacts.Contact{Status: contacts.Lead}, nil)), h.page("New contact"))
}

func (h *handler) edit(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact no longer exists."))
		return
	}
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, View(h.props(r, c, nil)), h.page("Edit "+c.Name))
}

// read takes the contact from the form.
func read(r *http.Request) contacts.Contact {
	v, _ := strconv.Atoi(r.FormValue("version"))
	return contacts.Contact{
		ID:      r.PathValue("id"),
		Name:    r.FormValue("name"),
		Company: r.FormValue("company"),
		Email:   r.FormValue("email"),
		Status:  contacts.Status(r.FormValue("status")),
		Tags:    contacts.ParseTags(r.FormValue("tags")),
		Notes:   strings.TrimSpace(r.FormValue("notes")),
		Version: v,
	}
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	h.save(w, r, read(r), h.svc.Create, "New contact")
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	h.save(w, r, read(r), h.svc.Update, "Edit contact")
}

func (h *handler) save(w http.ResponseWriter, r *http.Request, c contacts.Contact,
	store func(context.Context, contacts.Contact) (contacts.Contact, error), title string) {
	saved, err := store(r.Context(), c)
	var fe contacts.FieldErrors
	switch {
	case errors.As(err, &fe):
		htmx.Respond(w, r, View(h.props(r, c, fe)), htmx.Status(http.StatusUnprocessableEntity), h.page(title))
		return
	case errors.Is(err, contacts.ErrConflict):
		htmx.Error(w, r, htmx.Fail(http.StatusConflict, "Someone else changed this contact in the meantime. Close the form and open it again."))
		return
	case errors.Is(err, contacts.ErrNotFound):
		htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact no longer exists."))
		return
	case err != nil:
		htmx.Error(w, r, err)
		return
	}
	if !htmx.IsHTMX(r) {
		htmx.Redirect(w, r, h.o.Done)
		return
	}
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.CloseDialog(), htmx.Trigger(EventSaved, saved.ID),
		htmx.Toast(toast.Props{Color: "success"}, saved.Name+" saved."))
}
