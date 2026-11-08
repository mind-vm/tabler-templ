package contactdelete

import (
	"net/http"

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

// next is a local path to go to after deleting, or "".
func next(r *http.Request) string {
	if n := r.FormValue("next"); htmx.LocalPath(n) {
		return n
	}
	return ""
}

func (h *handler) confirm(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		htmx.Error(w, r, err, uierr.Rules...)
		return
	}
	htmx.Respond(w, r, Question(c, htmx.IsHTMX(r), next(r), h.o.Done), h.o.Frame.Page("Delete "+c.Name))
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Delete(r.Context(), r.PathValue("id"))
	if err != nil {
		htmx.Error(w, r, err, uierr.Rules...)
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
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.CloseDialog(), htmx.Trigger(events.Deleted, c.ID),
		htmx.ToastWith(toast.Props{Autohide: true, Delay: 10000}, Undo(c)))
}

func (h *handler) restore(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.Restore(r.Context(), r.PathValue("id"))
	if err != nil {
		htmx.Error(w, r, err, htmx.Is(contacts.ErrNotFound, http.StatusNotFound, "This contact cannot be restored any more."))
		return
	}
	htmx.Respond(w, r, templ.NopComponent, htmx.Reswap("none"), htmx.Trigger(events.Restored, c.ID),
		htmx.Toast(toast.Props{Color: "success"}, c.Name+" restored."))
}
