package contactform

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/events"
	"github.com/mind-vm/tabler-templ/htmx"
)

func do(t *testing.T, mux *http.ServeMux, method, url, form string, htmxReq bool) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(method, url, strings.NewReader(form))
	if form != "" {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if htmxReq {
		r.Header.Set("HX-Request", "true")
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

func frame(title string, body templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		io.WriteString(w, "<html><title>"+title+"</title>")
		body.Render(ctx, w)
		_, err := io.WriteString(w, "</html>")
		return err
	})
}

func setup(t *testing.T) *http.ServeMux {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore(), Options{Frame: frame, Done: "/list"})
	htmx.CheckRoutes(t, mux, New, Edit, Create, Update)
	return mux
}

func TestCreate(t *testing.T) {
	mux := setup(t)
	// Invalid: 422 with the form and the errors at the fields, in the dialog.
	w := do(t, mux, "POST", Create.URL(), "name=&email=x&status=lead", true)
	b := w.Body.String()
	if w.Code != 422 || !strings.Contains(b, `id="`+FormID+`"`) || !strings.Contains(b, "Enter a name.") || !strings.Contains(b, "is-invalid") || strings.Contains(b, "<html") {
		t.Fatalf("invalid: %d\n%s", w.Code, b)
	}
	// Valid: nothing swapped, dialog closed, event and toast.
	w = do(t, mux, "POST", Create.URL(), "name=Zoe&status=lead", true)
	trig := w.Header().Get("HX-Trigger")
	if w.Code != 200 || w.Header().Get("HX-Reswap") != "none" || !strings.Contains(trig, events.Saved) || !strings.Contains(trig, htmx.EventCloseDialog) ||
		!strings.Contains(w.Body.String(), "Zoe saved.") {
		t.Fatalf("valid: %d %v\n%s", w.Code, w.Header(), w.Body.String())
	}
	// Without JavaScript: Post/Redirect/Get, and the form as a page.
	if w = do(t, mux, "POST", Create.URL(), "name=Max&status=lead", false); w.Code != 303 || w.Header().Get("Location") != "/list" {
		t.Fatalf("no js: %d %v", w.Code, w.Header())
	}
	if w = do(t, mux, "GET", New.URL(), "", false); !strings.Contains(w.Body.String(), "<html><title>New contact") || strings.Contains(w.Body.String(), "hx-post") {
		t.Fatalf("page:\n%s", w.Body.String())
	}
}

func TestConflict(t *testing.T) {
	mux := setup(t)
	form := "name=Ada&status=active&version=1"
	if w := do(t, mux, "POST", Update.URL("1"), form, true); w.Code != 200 || w.Header().Get("HX-Retarget") != "" {
		t.Fatalf("first update: %d %v", w.Code, w.Header())
	}
	w := do(t, mux, "POST", Update.URL("1"), form, true)
	if w.Header().Get("HX-Retarget") != htmx.Target(htmx.ToastsID) || !strings.Contains(w.Body.String(), "Someone else changed") {
		t.Fatalf("conflict: %v\n%s", w.Header(), w.Body.String())
	}
}

// A broken version field is a 400, not an update against version 0.
func TestBadVersion(t *testing.T) {
	mux := setup(t)
	w := do(t, mux, "POST", Update.URL("1"), "name=Ada&status=active&version=x", true)
	if w.Code != 400 || !strings.Contains(w.Body.String(), "out of date") {
		t.Fatalf("bad version: %d\n%s", w.Code, w.Body.String())
	}
}
