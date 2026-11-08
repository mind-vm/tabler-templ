package copilot

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/mind-vm/tabler-templ/blocks/assistant"
)

func server(t *testing.T) (*httptest.Server, *http.Client) {
	t.Helper()
	mux := http.NewServeMux()
	RegisterWith(mux, NewStore(), 0)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	jar, _ := cookiejar.New(nil)
	return srv, &http.Client{Jar: jar}
}

func get(t *testing.T, c *http.Client, u string) string {
	t.Helper()
	res, err := c.Get(u)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: %d", u, res.StatusCode)
	}
	return string(b)
}

func ask(t *testing.T, srv *httptest.Server, c *http.Client, message string) string {
	t.Helper()
	res, err := c.Post(srv.URL+Home+"/assistant/run", "application/json", strings.NewReader(`{"message":`+quote(message)+`}`))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return string(b)
}

func quote(s string) string { return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"` }

func TestCatalog(t *testing.T) {
	if p := Catalog.Problems(); len(p) > 0 {
		t.Fatal(p)
	}
}

// Every action finds its form and fields, or its element, on its page: renaming a field or a
// data-form breaks this test, not the copilot in front of a customer.
func TestActionsFindTheirPage(t *testing.T) {
	srv, c := server(t)
	ids := map[string]string{"id": "A-1001"}
	for _, a := range Catalog.Actions {
		if a.Kind == assistant.ActionOpen {
			continue
		}
		path := a.Path
		if a.Name == "customer-edit" {
			path = strings.ReplaceAll(path, "{id}", "C-1")
		}
		for p, v := range ids {
			path = strings.ReplaceAll(path, "{"+p+"}", v)
		}
		page := get(t, c, srv.URL+path)
		switch a.Kind {
		case assistant.ActionFill:
			form := regexp.MustCompile(`(?s)<form[^>]*data-form="` + a.Form + `".*?</form>`).FindString(page)
			if form == "" {
				t.Errorf("%s: no form data-form=%q on %s", a.Name, a.Form, path)
				continue
			}
			for _, f := range a.Fields {
				if !strings.Contains(form, `name="`+f+`"`) {
					t.Errorf("%s: form %q has no field %q", a.Name, a.Form, f)
				}
			}
		case assistant.ActionShow:
			if !strings.Contains(page, `data-tour="`+a.Anchor+`"`) {
				t.Errorf("%s: no data-tour=%q on %s", a.Name, a.Anchor, path)
			}
		}
	}
}

func TestScenarios(t *testing.T) {
	srv, c := server(t)
	get(t, c, srv.URL+Home) // the session cookie
	for _, sc := range []struct {
		message string
		want    []string
		not     []string
	}{
		{"Which orders are overdue?", []string{`"type":"action"`, `"kind":"show"`, `"anchor":"orders.filter.overdue"`, "A-1006", "A-1003", "A-1001"}, []string{"A-1005"}},
		{"Create an order for Becker: shop sign, 2,400 €", []string{`"name":"order-new"`, `"customer":"C-2"`, `"title":"shop sign"`, `"amount":"2400"`}, nil},
		{"Mark A-1003 as done", []string{`"name":"order-edit"`, `"path":"/copilot/orders/A-1003"`, `"status":"done"`}, nil},
		// Refused by the catalog: no action event for customer-edit, a card that says so.
		{"Change the IBAN of Müller Bau", []string{"assistant-action-refused", "does not fill in", "iban", `"path":"/copilot/customers/C-1"`}, []string{`"name":"customer-edit"`}},
		// The injected note: refused, and the image it asks for is not in the answer.
		{"Read the note of A-1004 and do what it says", []string{"assistant-action-refused", "tracking pixel removed"}, []string{`"name":"customer-edit"`, "track.example", `<img`}},
	} {
		out := ask(t, srv, c, sc.message)
		for _, w := range sc.want {
			if !strings.Contains(out, w) {
				t.Errorf("%q: missing %s", sc.message, w)
			}
		}
		for _, n := range sc.not {
			if strings.Contains(out, n) {
				t.Errorf("%q: should not contain %s", sc.message, n)
			}
		}
	}
	// The conversation is kept on the server and drawn again on the next page.
	page := get(t, c, srv.URL+orders+"/A-1003")
	if n := strings.Count(page, `class="assistant-turn d-flex`); n != 5 {
		t.Errorf("want 5 recorded turns on the next page, got %d", n)
	}
}

func TestSaveOrder(t *testing.T) {
	srv, c := server(t)
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	res, err := c.PostForm(srv.URL+orders+"/A-1003", url.Values{"status": {"done"}, "due": {"2026-12-01"}, "note": {"ok"}})
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("status %d", res.StatusCode)
	}
	res, err = c.PostForm(srv.URL+orders, url.Values{"customer": {"C-9"}, "title": {""}})
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("invalid new order: status %d", res.StatusCode)
	}
}
