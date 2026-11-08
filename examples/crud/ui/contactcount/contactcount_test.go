package contactcount

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactform"
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

func TestCount(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore())
	htmx.CheckRoutes(t, mux, Count)
	w := do(t, mux, "GET", Count.URL()+htmx.Query("instance", "side", "status", "active", "label", "active"), "", true)
	b := w.Body.String()
	if !strings.Contains(b, `id="`+RootID("side")+`"`) || !strings.Contains(b, "12") || !strings.Contains(b, " active") ||
		!strings.Contains(b, contactform.EventSaved+" from:body") {
		t.Fatalf("count:\n%s", b)
	}
	if RootID("") == RootID("side") {
		t.Fatal("instances need their own ids")
	}
}
