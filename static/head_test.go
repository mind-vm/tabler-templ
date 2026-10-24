package static

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/a-h/templ"
	"github.com/mind-vm/tabler-templ/assets"
)

func tabler() *Files {
	fsys := fstest.MapFS{IconsCSS: {Data: []byte(".ti{}")}, "icons/fonts/tabler-icons.woff2": {Data: []byte("wOF2")}}
	for _, n := range []string{"css/tabler.min.css", "css/tabler.rtl.min.css", "css/tabler-vendors.min.css",
		"css/tabler-flags.min.css", "js/tabler.min.js", "js/tabler-theme.min.js"} {
		fsys[CoreDir+n] = &fstest.MapFile{Data: []byte(n)}
	}
	return New("/static/", fsys)
}

func render(t *testing.T, c templ.Component, ctx ...context.Context) string {
	t.Helper()
	var b strings.Builder
	cx := context.Background()
	if len(ctx) > 0 {
		cx = ctx[0]
	}
	if err := c.Render(cx, &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// Head links the stylesheets assets.Head links, in its order, at versioned addresses of the bundle.
func TestHead(t *testing.T) {
	f := tabler()
	v := "/static/v" + f.Version() + "/"
	got := render(t, Head(f, assets.Props{Flags: true, Vendors: true}))
	want := `<link rel="stylesheet" href="` + v + `core/css/tabler.min.css">` +
		`<link rel="stylesheet" href="` + v + `core/css/tabler-flags.min.css">` +
		`<link rel="stylesheet" href="` + v + `core/css/tabler-vendors.min.css">` +
		`<link rel="stylesheet" href="` + v + `icons/tabler-icons.min.css">`
	if got != want {
		t.Errorf("Head:\n%s\nwant\n%s", got, want)
	}
	if got := render(t, Head(f, assets.Props{RTL: true})); !strings.Contains(got, v+"core/css/tabler.rtl.min.css") {
		t.Errorf("RTL: %s", got)
	}
	if got := render(t, Head(f)); strings.Count(got, "<link") != 2 {
		t.Errorf("default: %s", got)
	}
	for _, n := range Stylesheets(assets.Props{Flags: true, Charts: true}) {
		if !f.Has(n) {
			t.Errorf("missing %s", n)
		}
	}
	if got := render(t, Head(f, assets.Props{Charts: true})); !strings.Contains(got, "apexcharts") {
		t.Errorf("Charts: %s", got)
	}
	// Every address Head links is served with a year of cache.
	for _, n := range Stylesheets(assets.Props{Flags: true, Vendors: true}) {
		if w := get(t, f, f.URL(n)); w.Code != http.StatusOK || w.Header().Get("Cache-Control") != CacheVersioned {
			t.Errorf("%s: %d %v", n, w.Code, w.Header())
		}
	}
}

func TestScripts(t *testing.T) {
	f := tabler()
	v := "/static/v" + f.Version() + "/core/js/"
	if got := render(t, Scripts(f)); got != `<script src="`+v+`tabler.min.js"></script>` {
		t.Errorf("Scripts: %s", got)
	}
	if got := render(t, Theme(f)); got != `<script src="`+v+`tabler-theme.min.js"></script>` {
		t.Errorf("Theme: %s", got)
	}
	got := render(t, Scripts(f), templ.WithNonce(context.Background(), "n0nce"))
	if got != `<script nonce="n0nce" src="`+v+`tabler.min.js"></script>` {
		t.Errorf("nonce: %s", got)
	}
}
