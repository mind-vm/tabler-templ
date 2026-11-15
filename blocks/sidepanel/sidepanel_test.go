package sidepanel_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/sidepanel"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(templ.InitializeContext(context.Background()), &b))
	return b.String()
}

func panel(open bool, active string) *sidepanel.Panel {
	return &sidepanel.Panel{Open: open, Active: active, Tabs: []sidepanel.Tab{
		{Name: "notes", Icon: "ti-notes", Title: "Notes", Content: templ.Raw("<p>n</p>")},
		{Name: "assistant", Icon: "ti-sparkles", Title: "Assistant", Content: templ.Raw("<p>a</p>"), Forms: true},
	}}
}

func TestReadAndRemember(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	open, tab := sidepanel.Read(r, "notes", "assistant")
	assert.False(t, open)
	assert.Empty(t, tab)

	w := httptest.NewRecorder()
	sidepanel.Remember(w, "assistant")
	r.AddCookie(w.Result().Cookies()[0])
	open, tab = sidepanel.Read(r, "notes", "assistant")
	assert.True(t, open)
	assert.Equal(t, "assistant", tab)

	r = httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: sidepanel.Cookie, Value: "other"})
	open, _ = sidepanel.Read(r, "notes", "assistant")
	assert.False(t, open, "an unknown tab counts as closed")
}

// GS-13: closed, the panel is in the page but hidden; open from the cookie, it is rendered open (no flicker)
// on its tab, and the page makes room for it.
func TestGS13_PanelFromTheServer(t *testing.T) {
	closed := render(t, sidepanel.Aside(panel(false, "")))
	assert.Contains(t, closed, `<aside id="kit-sidepanel" class="kit-sidepanel" data-active="notes" aria-label="Notes, Assistant" hidden>`)
	assert.Contains(t, closed, `role="tablist"`)

	open := render(t, sidepanel.Aside(panel(true, "assistant")))
	assert.Contains(t, open, `data-active="assistant"`)
	assert.Contains(t, open, "data-kit-from-server")
	assert.NotContains(t, open, `aria-label="Notes, Assistant" hidden`)
	assert.Contains(t, open, `id="kit-sidepanel-notes" class="kit-sidepanel-content" data-kit-sidepanel-content="notes" aria-label="Notes" hidden`)
	assert.Contains(t, open, `data-kit-sidepanel-forms`)

	assert.Equal(t, "kit-sidepanel-open", sidepanel.BodyAttributes(panel(true, ""), nil)["class"])
	assert.Equal(t, "x kit-sidepanel-open", sidepanel.BodyAttributes(panel(true, ""), templ.Attributes{"class": "x"})["class"])
	assert.Nil(t, sidepanel.BodyAttributes(panel(false, ""), nil))
}

func TestButtons(t *testing.T) {
	html := render(t, sidepanel.Buttons(panel(true, "assistant")))
	assert.Contains(t, html, `data-kit-sidepanel-tab="notes"`)
	assert.Contains(t, html, `aria-controls="kit-sidepanel"`)
	assert.Contains(t, html, `aria-expanded="true"`, "the open tab's button")
	assert.Equal(t, 1, strings.Count(html, `aria-expanded="true"`))
}

func TestNoTabsNoPanel(t *testing.T) {
	assert.Empty(t, render(t, sidepanel.Aside(&sidepanel.Panel{})))
	assert.Empty(t, render(t, sidepanel.Buttons(nil)))
}

func TestOneTabHasNoTabRow(t *testing.T) {
	p := &sidepanel.Panel{Tabs: []sidepanel.Tab{{Name: "notes", Icon: "ti-notes", Title: "Notes"}}}
	html := render(t, sidepanel.Aside(p))
	assert.NotContains(t, html, `role="tablist"`)
	assert.Contains(t, html, "Notes")
}

func TestScriptOncePerPage(t *testing.T) {
	html := render(t, templ.Join(sidepanel.Script(), sidepanel.Script()))
	assert.Equal(t, 1, strings.Count(html, "<script"))
	assert.Contains(t, html, "kit-sidepanel=")
	assert.Contains(t, html, "window.kitSidepanel")
}
