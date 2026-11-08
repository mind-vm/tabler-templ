package static

import (
	"context"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/internal/csp"
)

func TestWithKitCSS(t *testing.T) {
	t.Cleanup(func() { csp.SetExternal(false) })
	d := New("/x/", fstest.MapFS{"app.css": {Data: []byte(".app{}")}}, WithKitCSS())
	if !d.Has(KitCSS) {
		t.Fatal("bundle lacks the kit's CSS")
	}
	w := get(t, d, d.URL(KitCSS), "Accept-Encoding", "gzip")
	if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/css") || w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("kit.css: %d %v", w.Code, w.Header())
	}
	css := csp.All()
	for _, name := range []string{"/* sidebar */", "/* blocks/shell */", "/* icon */"} {
		if !strings.Contains(css, name) {
			t.Errorf("kit.css lacks %s", name)
		}
	}

	var b strings.Builder
	page := shell.Page(shell.Props{Title: "x", Head: Head(d)})
	if err := page.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	h := b.String()
	if strings.Contains(h, "<style") {
		t.Errorf("page still has inline CSS")
	}
	kit, tabler := strings.Index(h, d.URL(KitCSS)), strings.Index(h, "tabler.min.css")
	if kit < 0 || tabler < 0 || kit > tabler {
		t.Errorf("kit.css must be linked before Tabler's CSS:\n%s", h)
	}
}

func TestKitWithoutOption(t *testing.T) {
	d := bundle()
	var b strings.Builder
	if err := Kit(d).Render(context.Background(), &b); err != nil || b.Len() != 0 || d.Has(KitCSS) {
		t.Fatalf("Kit without WithKitCSS: %q %v", b.String(), err)
	}
}
