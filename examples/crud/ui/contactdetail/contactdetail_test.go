package contactdetail

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

func TestDetail(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore(), Options{Frame: frame, Back: "/list"})
	htmx.CheckRoutes(t, mux, Page, Activity)

	w := do(t, mux, "GET", Page.URL("1"), "", false)
	b := w.Body.String()
	if !strings.HasPrefix(b, "<html><title>Ada Lovelace") || !strings.Contains(b, `hx-get="`+Activity.URL("1")+`"`) || !strings.Contains(b, `id="`+RootID+`"`) {
		t.Fatalf("page:\n%s", b)
	}
	if w = do(t, mux, "GET", Page.URL("1"), "", true); strings.Contains(w.Body.String(), "<html") {
		t.Fatal("htmx gets the fragment only")
	}
	if w = do(t, mux, "GET", Activity.URL("1"), "", true); !strings.Contains(w.Body.String(), "Contact created") {
		t.Fatalf("activity:\n%s", w.Body.String())
	}
	if w = do(t, mux, "GET", Page.URL("999"), "", false); w.Code != 404 {
		t.Fatalf("missing contact: %d", w.Code)
	}
}
