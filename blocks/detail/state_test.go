package detail_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/detail"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(ctx, &b))
	return b.String()
}

// GS-43: after a failed save the page tells the script which dialog opens again, with the values entered.
func TestGS43_StateScript(t *testing.T) {
	form := url.Values{"titel": {"Preis", "Konzept"}, "_dialog": {"d"}, "csrf": {"x"}, "password": {"p"}}
	h := render(t, templ.WithNonce(context.Background(), "n0"), detail.StateScript(detail.State{ID: "dialog-kriterien", Error: "Bitte prüfen",
		Values: detail.StateValues(form, "csrf"), Kind: "kriterium"}))
	assert.Equal(t, `<script nonce="n0">window["kitDialogState"]={"error":"Bitte prüfen","id":"dialog-kriterien","values":{"_form":["kriterium"],"titel":["Preis","Konzept"]}};</script>`, h)
	assert.Empty(t, render(t, context.Background(), detail.StateScript(detail.State{})))
	h = render(t, context.Background(), detail.StateScript(detail.State{ID: "d", Values: url.Values{"x": {"</script>"}}}, detail.ScriptConfig{StateVar: "myState"}))
	assert.Contains(t, h, `window["myState"]=`)
	assert.NotContains(t, h, "</script><")
}

// Several fields of one name get their own value back, in order.
func TestGS43_RestoreKeepsEachValue(t *testing.T) {
	h := render(t, context.Background(), detail.Script(detail.ScriptConfig{}))
	assert.Contains(t, h, "els.forEach(function(e,i)")
	assert.Contains(t, h, "(els.length>1?list[i]:list[0])")
}

// A button shared by the rows of a table sets the dialog form's action, only to a path on the same site.
func TestDialogActionAttribute(t *testing.T) {
	h := render(t, context.Background(), detail.Script(detail.ScriptConfig{}))
	assert.Contains(t, h, `"actionAttribute":"data-kit-action"`)
	assert.Contains(t, h, `f.setAttribute("action",act)`)
	assert.Contains(t, h, `/^\/(?![\/\\])/.test(act)`)
	h = render(t, context.Background(), detail.Script(detail.ScriptConfig{ActionAttribute: "data-action"}))
	assert.Contains(t, h, `"actionAttribute":"data-action"`)
}
