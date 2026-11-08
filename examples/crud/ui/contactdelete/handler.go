package contactdelete

import (
	"errors"
	"net/http"
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

var errGone = htmx.Fail(http.StatusNotFound, "This contact no longer exists.")

// next is a local path to go to after deleting, or "".
func next(r *http.Request) string {
	n := r.FormValue("next")
	if strings.HasPrefix(n, "/") && !strings.HasPrefix(n, "//") {
		return n
	}
	return ""
}

func (h *handler) confirm(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		htmx.Error(w, r, errGone)
		return
	}
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, Question(c, htmx.IsHTMX(r), next(r), h.o.Done),
		htmx.FullPage(func(b templ.Component) templ.Component { return h.o.Frame("Delete "+c.Name, b) }))
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Delete(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		htmx.Error(w, r, errGone)
		return
	}
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	if !htmx.IsHTMX(r) {
		htmx.Redirect(w, r, h.o.Done)
		return
	}
	if n := next(r); n != "" {
		// The page showed the deleted contact: leave it.
		htmx.Redirect(w, r, n)
		return
	}
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.CloseDialog(), htmx.Trigger(EventDeleted, c.ID),
		htmx.ToastWith(toast.Props{Autohide: true, Delay: 10000}, Undo(c)))
}

func (h *handler) restore(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Restore(r.Context(), r.PathValue("id"))
	if errors.Is(err, contacts.ErrNotFound) {
		htmx.Error(w, r, htmx.Fail(http.StatusNotFound, "This contact cannot be restored any more."))
		return
	}
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.Trigger(EventRestored, c.ID),
		htmx.Toast(toast.Props{Color: "success"}, c.Name+" restored."))
}
