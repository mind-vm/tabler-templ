package htmx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/toast"
)

func TestRoute(t *testing.T) {
	toggle := Route("post", "/ui/todo/{id}/toggle")
	if toggle.Pattern() != "POST /ui/todo/{id}/toggle" || toggle.Method() != "POST" {
		t.Fatalf("pattern %q", toggle.Pattern())
	}
	if got := toggle.URL("a b/c"); got != "/ui/todo/a%20b%2Fc/toggle" {
		t.Fatalf("url %q", got)
	}
	files := Route("GET", "/files/{path...}")
	if got := files.URL("a/b c"); got != "/files/a/b%20c" {
		t.Fatalf("rest url %q", got)
	}
	if got := Route("GET", "/crud/{$}").URL(); got != "/crud/" {
		t.Fatalf("end url %q", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("URL with a missing value must panic")
		}
	}()
	toggle.URL()
}

func TestHelpers(t *testing.T) {
	if ID("todo", "") != "todo" || ID("todo", "side bar") != "todo-side-bar" {
		t.Fatal("ID")
	}
	if Vals("instance", "b") != `{"instance":"b"}` || Target("x") != "#x" {
		t.Fatal("Vals or Target")
	}
	if Query("q", "a b", "page", "") != "?q=a+b" || Query("q", "") != "" {
		t.Fatalf("Query %q", Query("q", "a b", "page", ""))
	}
}

type fakeT struct{ errs []string }

func (f *fakeT) Helper()                        {}
func (f *fakeT) Errorf(format string, a ...any) { f.errs = append(f.errs, fmt.Sprintf(format, a...)) }

func TestCheckRoutes(t *testing.T) {
	list, item, missing := Route("GET", "/ui/x"), Route("POST", "/ui/x/{id}"), Route("DELETE", "/ui/x/{id}")
	mux := http.NewServeMux()
	nop := func(http.ResponseWriter, *http.Request) {}
	mux.HandleFunc(list.Pattern(), nop)
	mux.HandleFunc("POST /ui/x/{name}", nop) // registered under another pattern
	f := &fakeT{}
	CheckRoutes(f, mux, list, item, missing)
	if len(f.errs) != 2 || !strings.Contains(f.errs[0], `want "POST /ui/x/{id}"`) || !strings.Contains(f.errs[1], "DELETE /ui/x/{id}") {
		t.Fatalf("errors: %q", f.errs)
	}
}

func text(s string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error { _, err := io.WriteString(w, s); return err })
}

func frame(body templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		io.WriteString(w, "<html>")
		body.Render(ctx, w)
		_, err := io.WriteString(w, "</html>")
		return err
	})
}

func req(method string, htmx bool) *http.Request {
	r := httptest.NewRequest(method, "/ui/x", nil)
	if htmx {
		r.Header.Set("HX-Request", "true")
	}
	return r
}

func TestRespondHTMX(t *testing.T) {
	w := httptest.NewRecorder()
	err := Respond(w, req("POST", true), text("<div id=x></div>"), Status(422), Trigger("todo:changed"), Trigger("todo:count", 3),
		PushURL("/todos?q=a"), Toast(toast.Props{Color: "success"}, "Saved <b>"), FullPage(frame))
	if err != nil {
		t.Fatal(err)
	}
	h := w.Header()
	if w.Code != 422 || h.Get("Vary") != "HX-Request" || h.Get("HX-Push-Url") != "/todos?q=a" {
		t.Fatalf("status %d headers %v", w.Code, h)
	}
	if h.Get("HX-Trigger") != `{"todo:changed":null,"todo:count":3}` {
		t.Fatalf("trigger %q", h.Get("HX-Trigger"))
	}
	body := w.Body.String()
	if !strings.HasPrefix(body, "<div id=x></div>") || strings.Contains(body, "<html>") {
		t.Fatalf("fragment expected:\n%s", body)
	}
	if !strings.Contains(body, `hx-swap-oob="beforeend:#htmx-toasts"`) || !strings.Contains(body, "Saved &lt;b&gt;") || !strings.Contains(body, "toast-success") {
		t.Fatalf("toast:\n%s", body)
	}
}

func TestRespondWithoutHTMX(t *testing.T) {
	w := httptest.NewRecorder()
	Respond(w, req("GET", false), text("<div id=x></div>"), Trigger("a"), Toast(toast.Props{}, "hi"), FullPage(frame))
	if w.Header().Get("HX-Trigger") != "" || w.Header().Get("Vary") != "HX-Request" {
		t.Fatalf("headers %v", w.Header())
	}
	if w.Body.String() != "<html><div id=x></div></html>" {
		t.Fatalf("full page expected: %s", w.Body.String())
	}
	// A history restore wants the page too.
	r := req("GET", true)
	r.Header.Set("HX-History-Restore-Request", "true")
	w = httptest.NewRecorder()
	Respond(w, r, text("x"), FullPage(frame))
	if w.Body.String() != "<html>x</html>" {
		t.Fatalf("history restore: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	Respond(w, req("POST", true), text("x"), Trigger("a"), Trigger("b"))
	if w.Header().Get("HX-Trigger") != "a, b" {
		t.Fatalf("plain trigger %q", w.Header().Get("HX-Trigger"))
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, req("POST", true), Fail(http.StatusConflict, "Changed by someone else."))
	if w.Code != 200 || w.Header().Get("HX-Retarget") != "#htmx-toasts" || w.Header().Get("HX-Reswap") != "beforeend" ||
		!strings.Contains(w.Body.String(), "Changed by someone else.") || !strings.Contains(w.Body.String(), "toast-danger") {
		t.Fatalf("htmx error: %d %v %s", w.Code, w.Header(), w.Body.String())
	}
	w = httptest.NewRecorder()
	Error(w, req("POST", true), errors.New("db: password=secret"))
	if strings.Contains(w.Body.String(), "secret") || !strings.Contains(w.Body.String(), "Something went wrong") {
		t.Fatalf("internal error leaked: %s", w.Body.String())
	}
	w = httptest.NewRecorder()
	Error(w, req("POST", false), Fail(http.StatusConflict, "Conflict."))
	if w.Code != 409 || !strings.Contains(w.Body.String(), "Conflict.") {
		t.Fatalf("plain error: %d %s", w.Code, w.Body.String())
	}
}

func TestRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	Redirect(w, req("POST", true), "/todos")
	if w.Code != 204 || w.Header().Get("HX-Redirect") != "/todos" {
		t.Fatalf("htmx redirect %d %v", w.Code, w.Header())
	}
	w = httptest.NewRecorder()
	Redirect(w, req("POST", false), "/todos")
	if w.Code != 303 || w.Header().Get("Location") != "/todos" {
		t.Fatalf("redirect %d %v", w.Code, w.Header())
	}
}

func TestSetup(t *testing.T) {
	var b strings.Builder
	if err := Setup().Render(templ.WithNonce(context.Background(), "n"), &b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{`id="htmx-toasts"`, `id="htmx-dialog"`, `id="htmx-dialog-content"`, `nonce="n"`,
		"responseHandling", "data-offline=\"No connection to the server.\"", "kit:close-dialog", "tablerCharts.destroy"} {
		if !strings.Contains(out, want) {
			t.Errorf("Setup lacks %q", want)
		}
	}
	if strings.Count(out, "data-allow-eval") != 1 { // only the script's own check
		t.Error("eval must be off by default")
	}
	b.Reset()
	Setup(SetupProps{AllowEval: true}).Render(context.Background(), &b)
	if strings.Count(b.String(), "data-allow-eval") != 2 {
		t.Error("AllowEval must set the attribute")
	}
}
