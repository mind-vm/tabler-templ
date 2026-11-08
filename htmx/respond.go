package htmx

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/toast"
)

// Ids of the page-wide elements Setup renders: the toast container and the dialog.
const (
	ToastsID        = "htmx-toasts"
	DialogID        = "htmx-dialog"
	DialogContentID = "htmx-dialog-content"
)

// EventCloseDialog is the event CloseDialog triggers; Setup's script closes the dialog on it.
const EventCloseDialog = "kit:close-dialog"

// IsHTMX reports whether htmx sent the request (HX-Request), not counting a history restore, which
// wants the whole page.
func IsHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-History-Restore-Request") != "true"
}

// An Option shapes a Respond: its status, its htmx headers, extra fragments.
type Option func(*response)

type response struct {
	status     int
	trigger    map[string]any
	order      []string
	retarget   string
	reswap     string
	pushURL    string
	replaceURL string
	oob        []templ.Component
	full       func(templ.Component) templ.Component
}

// Status sets the status code, e.g. http.StatusUnprocessableEntity for a form with errors (Setup makes
// htmx swap 422 responses in).
func Status(code int) Option { return func(r *response) { r.status = code } }

// Trigger fires an event in the browser once the response is in (HX-Trigger), e.g. to tell other
// components that something changed. detail, if given, is the event's detail (JSON).
func Trigger(event string, detail ...any) Option {
	return func(r *response) {
		if r.trigger == nil {
			r.trigger = map[string]any{}
		}
		if _, ok := r.trigger[event]; !ok {
			r.order = append(r.order, event)
		}
		if len(detail) > 0 {
			r.trigger[event] = detail[0]
		} else {
			r.trigger[event] = nil
		}
	}
}

// Retarget swaps the response into another element than the request's hx-target (HX-Retarget).
func Retarget(selector string) Option { return func(r *response) { r.retarget = selector } }

// Reswap changes how the response is swapped in (HX-Reswap), e.g. "beforeend" or "none".
func Reswap(swap string) Option { return func(r *response) { r.reswap = swap } }

// PushURL puts url into the address bar and the history (HX-Push-Url), e.g. the page's URL with the
// list's filters, so the state can be bookmarked.
func PushURL(url string) Option { return func(r *response) { r.pushURL = url } }

// ReplaceURL replaces the address bar's URL without a history entry (HX-Replace-Url).
func ReplaceURL(url string) Option { return func(r *response) { r.replaceURL = url } }

// OOB adds a fragment that htmx swaps in by its own id (out of band); c must render its root with
// hx-swap-oob, e.g. hx-swap-oob="true". Only for elements of the component's own contract.
func OOB(c templ.Component) Option { return func(r *response) { r.oob = append(r.oob, c) } }

// Toast shows a toast in Setup's container once the response is in. The text is the body; p sets the
// color, title and so on (shown and hidden after 5 seconds unless p says otherwise).
func Toast(p toast.Props, text string) Option {
	return func(r *response) { r.oob = append(r.oob, oobToast(p, text)) }
}

// ToastWith is Toast with markup as the body, e.g. a text and an Undo button.
func ToastWith(p toast.Props, body templ.Component) Option {
	return func(r *response) { r.oob = append(r.oob, oobToastWith(p, body)) }
}

// CloseDialog closes the dialog of Setup, e.g. after a form in it was saved.
func CloseDialog() Option { return Trigger(EventCloseDialog) }

// FullPage is what a request without htmx gets: the component inside the page that frame makes, e.g.
// the layout. A fragment route opened directly (bookmark, new tab, history) then shows a whole page.
// Without FullPage such a request gets the bare fragment.
func FullPage(frame func(body templ.Component) templ.Component) Option {
	return func(r *response) { r.full = frame }
}

// Respond writes c as the response. It sets all headers first (Vary: HX-Request, the htmx headers of
// the options, the status) and renders afterwards, so no header comes too late. A request without htmx
// gets the full page (FullPage) and no htmx headers; the extra fragments (OOB, Toast) go only to htmx.
func Respond(w http.ResponseWriter, r *http.Request, c templ.Component, opts ...Option) error {
	var res response
	for _, o := range opts {
		o(&res)
	}
	h := w.Header()
	h.Add("Vary", "HX-Request")
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", "text/html; charset=utf-8")
	}
	htmx := IsHTMX(r)
	if htmx {
		if len(res.trigger) > 0 {
			h.Set("HX-Trigger", triggerHeader(res))
		}
		set := func(k, v string) {
			if v != "" {
				h.Set(k, v)
			}
		}
		set("HX-Retarget", res.retarget)
		set("HX-Reswap", res.reswap)
		set("HX-Push-Url", res.pushURL)
		set("HX-Replace-Url", res.replaceURL)
	}
	if res.status != 0 {
		w.WriteHeader(res.status)
	}
	ctx := r.Context()
	if !htmx && res.full != nil {
		return res.full(c).Render(ctx, w)
	}
	if err := c.Render(ctx, w); err != nil {
		return err
	}
	if htmx {
		for _, o := range res.oob {
			if err := o.Render(ctx, w); err != nil {
				return err
			}
		}
	}
	return nil
}

func triggerHeader(res response) string {
	plain := true
	for _, v := range res.trigger {
		if v != nil {
			plain = false
		}
	}
	if plain {
		return strings.Join(res.order, ", ")
	}
	b, _ := json.Marshal(res.trigger)
	return string(b)
}

// HTTPError is an error with the status code and the message the user may see, e.g. a conflict.
type HTTPError struct {
	Code    int
	Message string // shown to the user; for 5xx leave it empty (a generic text is shown)
	Err     error  // the cause, logged
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error { return e.Err }

// Fail is an HTTPError with a code and a message for the user.
func Fail(code int, message string) error { return &HTTPError{Code: code, Message: message} }

// Error answers a failed request. To htmx it sends a danger toast into Setup's container (HX-Retarget,
// HX-Reswap) and leaves the page as it is; without htmx it sends a plain error page with the status.
// An *HTTPError gives the status and the user's message; any other error counts as 500, is logged and
// shown with kit.T(ctx).ErrorGeneric, so details never reach the HTML.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	code, msg := http.StatusInternalServerError, ""
	var he *HTTPError
	if errors.As(err, &he) {
		code, msg = he.Code, he.Message
	}
	if code >= 500 || msg == "" {
		slog.ErrorContext(r.Context(), "request failed", "method", r.Method, "path", r.URL.Path, "status", code, "err", err)
		msg = kit.T(r.Context()).ErrorGeneric
	}
	if !IsHTMX(r) {
		http.Error(w, msg, code)
		return
	}
	h := w.Header()
	h.Add("Vary", "HX-Request")
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("HX-Retarget", Target(ToastsID))
	h.Set("HX-Reswap", "beforeend")
	// 200, so htmx swaps the toast in; the status the request deserved goes along for scripts and logs.
	h.Set("X-Status", http.StatusText(code))
	_ = toastHTML(toast.Props{Color: "danger", Urgent: true}, msg).Render(r.Context(), w)
}

// Redirect sends the browser to url: with HX-Redirect to htmx (a full page load), with 303 See Other
// otherwise (Post/Redirect/Get).
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if IsHTMX(r) {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func withDefaults(p toast.Props) toast.Props {
	p.Show = true
	if !p.Autohide && p.Delay == 0 {
		p.Autohide, p.Delay = true, 5000
	}
	p.Class = kit.Classes(p.Class, "kit-htmx-toast")
	return p
}

func toastHTML(p toast.Props, text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return toast.Toast(withDefaults(p)).Render(templ.WithChildren(ctx, templ.Raw(templ.EscapeString(text))), w)
	})
}

// oobToast wraps a toast so htmx appends it to the container, whatever the request's target.
func oobToast(p toast.Props, text string) templ.Component {
	return oobToastWith(p, templ.Raw(templ.EscapeString(text)))
}

func oobToastWith(p toast.Props, body templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div hx-swap-oob="beforeend:`+Target(ToastsID)+`">`); err != nil {
			return err
		}
		if err := toast.Toast(withDefaults(p)).Render(templ.WithChildren(ctx, body), w); err != nil {
			return err
		}
		_, err := io.WriteString(w, `</div>`)
		return err
	})
}
