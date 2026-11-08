package actionbar

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/kit"
)

func TestBar(t *testing.T) {
	var b strings.Builder
	ctx := kit.WithMessages(templ.WithNonce(context.Background(), "n1"), kit.German)
	page := templ.Join(Bar(Props{Form: "f", Guard: true, Info: "Saved"}), Bar(Props{Dirty: true}))
	if err := page.Render(ctx, &b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{`data-form="f"`, "data-guard", "data-dirty", "Ungespeicherte Änderungen", "Saved", `nonce="n1"`, "beforeunload"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Count(out, "<script") != 1 || strings.Count(out, "<style") != 1 {
		t.Errorf("script or style not once per page")
	}
}
