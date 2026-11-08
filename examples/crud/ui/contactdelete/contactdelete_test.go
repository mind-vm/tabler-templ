package contactdelete

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
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

func TestDeleteAndUndo(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore(), Options{Frame: frame, Done: "/list"})
	htmx.CheckRoutes(t, mux, Confirm, Delete, Restore)

	if w := do(t, mux, "GET", Confirm.URL("1"), "", true); !strings.Contains(w.Body.String(), "Delete Ada Lovelace?") || strings.Contains(w.Body.String(), "<html") {
		t.Fatalf("question:\n%s", w.Body.String())
	}
	w := do(t, mux, "POST", Delete.URL("1"), "", true)
	if !strings.Contains(w.Header().Get("HX-Trigger"), EventDeleted) || !strings.Contains(w.Body.String(), Restore.URL("1")) {
		t.Fatalf("delete: %v\n%s", w.Header(), w.Body.String())
	}
	w = do(t, mux, "POST", Restore.URL("1"), "", true)
	if !strings.Contains(w.Header().Get("HX-Trigger"), EventRestored) || !strings.Contains(w.Body.String(), "Ada Lovelace restored.") {
		t.Fatalf("restore: %v\n%s", w.Header(), w.Body.String())
	}
	// From the contact's own page: leave it after deleting. An outside "next" is ignored.
	if w = do(t, mux, "POST", Delete.URL("2"), "next=/list", true); w.Header().Get("HX-Redirect") != "/list" {
		t.Fatalf("next: %v", w.Header())
	}
	if w = do(t, mux, "POST", Delete.URL("3"), "next=//evil.example", true); w.Header().Get("HX-Redirect") != "" {
		t.Fatalf("open redirect: %v", w.Header())
	}
	if w = do(t, mux, "POST", Delete.URL("4"), "", false); w.Code != 303 || w.Header().Get("Location") != "/list" {
		t.Fatalf("no js: %d %v", w.Code, w.Header())
	}
}
