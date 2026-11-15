package mark_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/mark"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(templ.InitializeContext(ctx), &b))
	return b.String()
}

// GS-76: text from outside is marked, escaped and never a link.
func TestGS76_TextFromOutside(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, mark.Text(mark.External, "Siehe https://example.org <a href=\"x\">hier</a>\nZeile 2"))
	assert.Contains(t, h, `data-mark="external"`)
	assert.Contains(t, h, ">von außen</span>")
	assert.Contains(t, h, "&lt;a href=")
	assert.NotContains(t, h, "<a ")
	assert.Contains(t, h, "Zeile 2")
}

// GS-07, GS-76: an internal note says it is never sent.
func TestGS76_InternalNote(t *testing.T) {
	h := render(t, kit.WithMessages(context.Background(), kit.German), mark.Text(mark.Internal, "Rückruf"))
	assert.Contains(t, h, "intern, nie gesendet")
	assert.Contains(t, render(t, context.Background(), mark.Badge(mark.Internal)), "internal, never sent")
	assert.Contains(t, render(t, context.Background(), mark.Text(mark.External, "x", mark.TextProps{Label: "von außen · E-Mail"})), "von außen · E-Mail")
}
