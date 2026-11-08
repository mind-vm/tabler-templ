package contacttable

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestTable(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore(), Options{Page: "/list", Detail: func(id string) string { return "/c/" + id }})
	htmx.CheckRoutes(t, mux, List)

	// A search (boosted form): the table only, and the page's address follows without empty parameters.
	r := httptest.NewRequest("GET", List.URL()+"?q=bell&status=&sort=", nil)
	r.Header.Set("HX-Request", "true")
	r.Header.Set("HX-Boosted", "true")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	b := w.Body.String()
	if w.Header().Get("HX-Push-Url") != "/list?q=bell" || !strings.Contains(b, "Ken Thompson") || strings.Contains(b, "Ada Lovelace") || !strings.Contains(b, `href="/c/9"`) {
		t.Fatalf("search: %v\n%s", w.Header(), b)
	}
	// A refresh after an event does not touch the address.
	if w = do(t, mux, "GET", List.URL(), "", true); w.Header().Get("HX-Push-Url") != "" {
		t.Fatalf("refresh pushed %q", w.Header().Get("HX-Push-Url"))
	}
	// Opened directly: to the page, with the same query.
	if w = do(t, mux, "GET", List.URL()+"?q=bell", "", false); w.Code != 303 || w.Header().Get("Location") != "/list?q=bell" {
		t.Fatalf("direct: %d %v", w.Code, w.Header())
	}
}
