package csp

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

var testSheet = Sheet("test/sheet", ".t{color:red}")

func TestSheetOncePerPage(t *testing.T) {
	var b strings.Builder
	ctx := templ.WithNonce(context.Background(), "n")
	if err := templ.Join(testSheet, testSheet).Render(ctx, &b); err != nil {
		t.Fatal(err)
	}
	if got := b.String(); got != `<style nonce="n">.t{color:red}</style>` {
		t.Fatalf("got %q", got)
	}
	if !strings.Contains(All(), "/* test/sheet */\n.t{color:red}") {
		t.Fatalf("All lacks the sheet:\n%s", All())
	}
}

func TestSheetExternal(t *testing.T) {
	SetExternal(true)
	t.Cleanup(func() { SetExternal(false) })
	var b strings.Builder
	if err := testSheet.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 0 {
		t.Fatalf("external sheet rendered %q", b.String())
	}
}
