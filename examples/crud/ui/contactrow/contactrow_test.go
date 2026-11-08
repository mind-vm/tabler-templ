package contactrow

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestStatusCell(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, contacts.NewStore())
	htmx.CheckRoutes(t, mux, Status, EditStatus, SaveStatus)

	w := do(t, mux, "GET", EditStatus.URL("1"), "", true)
	if !strings.Contains(w.Body.String(), `<td id="`+CellID("1")+`"`) || !strings.Contains(w.Body.String(), `<select name="status"`) {
		t.Fatalf("edit cell:\n%s", w.Body.String())
	}
	w = do(t, mux, "POST", SaveStatus.URL("1"), "status=inactive", true)
	if w.Code != 200 || w.Header().Get("HX-Trigger") != `{"`+events.StatusChanged+`":"1"}` || !strings.Contains(w.Body.String(), "Inactive") {
		t.Fatalf("save: %d %v\n%s", w.Code, w.Header(), w.Body.String())
	}
	w = do(t, mux, "POST", SaveStatus.URL("1"), "status=bogus", true)
	if w.Header().Get("HX-Retarget") != htmx.Target(htmx.ToastsID) || !strings.Contains(w.Body.String(), "Pick a status.") {
		t.Fatalf("invalid status: %v\n%s", w.Header(), w.Body.String())
	}
}
