package contactform

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/events"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/uierr"
	"github.com/mind-vm/tabler-templ/htmx"
	"github.com/mind-vm/tabler-templ/toast"
)

type handler struct {
	svc Service
	o   Options
}

func (h *handler) props(r *http.Request, c contacts.Contact, errs contacts.FieldErrors) Props {
	return Props{Contact: c, Errors: errs, InDialog: htmx.IsHTMX(r), Cancel: h.o.Done}
}

func (h *handler) new(w http.ResponseWriter, r *http.Request) {
	htmx.Respond(w, r, View(h.props(r, contacts.Contact{Status: contacts.Lead}, nil)), h.o.Frame.Page("New contact"))
}

func (h *handler) edit(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		htmx.Error(w, r, err, uierr.Rules...)
		return
	}
	htmx.Respond(w, r, View(h.props(r, c, nil)), h.o.Frame.Page("Edit "+c.Name))
}

// read takes the contact from the form. A version that is not a number is a broken form: without it the
// conflict check would compare against version 0.
func read(r *http.Request) (contacts.Contact, error) {
	v := 0
	if s := r.FormValue("version"); s != "" {
		var err error
		if v, err = strconv.Atoi(s); err != nil {
			return contacts.Contact{}, htmx.Fail(http.StatusBadRequest, "The form is out of date. Reload the page.")
		}
	}
	return contacts.Contact{
		ID:      r.PathValue("id"),
		Name:    r.FormValue("name"),
		Company: r.FormValue("company"),
		Email:   r.FormValue("email"),
		Status:  contacts.Status(r.FormValue("status")),
		Tags:    contacts.ParseTags(r.FormValue("tags")),
		Notes:   strings.TrimSpace(r.FormValue("notes")),
		Version: v,
	}, nil
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	h.save(w, r, h.svc.Create, "New contact")
}

func (h *handler) update(w http.ResponseWriter, r *http.Request) {
	h.save(w, r, h.svc.Update, "Edit contact")
}

func (h *handler) save(w http.ResponseWriter, r *http.Request,
	store func(context.Context, contacts.Contact) (contacts.Contact, error), title string) {
	c, err := read(r)
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	saved, err := store(r.Context(), c)
	var fe contacts.FieldErrors
	switch {
	case errors.As(err, &fe):
		htmx.Respond(w, r, View(h.props(r, c, fe)), htmx.Status(http.StatusUnprocessableEntity), h.o.Frame.Page(title))
		return
	case err != nil:
		htmx.Error(w, r, err, uierr.Rules...)
		return
	}
	if !htmx.IsHTMX(r) {
		htmx.Redirect(w, r, h.o.Done)
		return
	}
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.CloseDialog(), htmx.Trigger(events.Saved, saved.ID),
		htmx.Toast(toast.Props{Color: "success"}, saved.Name+" saved."))
}
