package crud

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactcount"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactdelete"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactdetail"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactform"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactrow"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contacttable"
	"github.com/mind-vm/tabler-templ/htmx"
)

// Every route of every component is registered under its own pattern.
func TestContract(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)
	htmx.CheckRoutes(t, mux,
		contactcount.Count,
		contactrow.Status, contactrow.EditStatus, contactrow.SaveStatus,
		contactform.New, contactform.Edit, contactform.Create, contactform.Update,
		contactdelete.Confirm, contactdelete.Delete, contactdelete.Restore,
		contacttable.List,
		contactdetail.Page, contactdetail.Activity,
	)
}

// The list page renders with the table, both count instances and the htmx client.
func TestHome(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/crud?status=lead", nil))
	body := w.Body.String()
	for _, want := range []string{`id="contacttable"`, `id="contactcount-header"`, `id="contactcount-active"`, `id="htmx-toasts"`,
		`id="htmx-dialog"`, "htmx.min.js", "Linus Torvalds"} {
		if !strings.Contains(body, want) {
			t.Errorf("home lacks %q", want)
		}
	}
	if strings.Contains(body, "Ada Lovelace") {
		t.Error("the status filter from the URL is not applied")
	}
}

// No URL in the markup is a string literal: hx-get, hx-post, href and action come from the contract.
func TestNoLiteralURLs(t *testing.T) {
	literal := regexp.MustCompile(`\b(hx-get|hx-post|hx-put|hx-delete|action)="[^"]`)
	files, _ := filepath.Glob("ui/*/*.templ")
	if len(files) == 0 {
		t.Fatal("no templates found")
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for i, line := range strings.Split(string(b), "\n") {
			if literal.MatchString(line) {
				t.Errorf("%s:%d: literal URL: %s", f, i+1, strings.TrimSpace(line))
			}
		}
	}
}
