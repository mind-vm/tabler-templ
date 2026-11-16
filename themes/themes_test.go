package themes

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/auth"
	"github.com/mind-vm/tabler-templ/blocks/layout"
	"github.com/mind-vm/tabler-templ/blocks/marketing"
	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/theme"
)

func TestThemes(t *testing.T) {
	seen := map[string]bool{}
	for _, th := range All {
		if seen[th.Name] {
			t.Errorf("%s twice", th.Name)
		}
		seen[th.Name] = true
		if got, ok := Find(strings.ToLower(th.Label)); !ok || got.Name != th.Name {
			t.Errorf("Find(%q)", th.Label)
		}
		if th.Primary != "" && !strings.HasPrefix(th.Primary, "oklch(") {
			t.Errorf("%s: primary %q", th.Name, th.Primary)
		}
		if th.Grays != nil && len(th.Grays) != 11 {
			t.Errorf("%s: %d grays", th.Name, len(th.Grays))
		}
		for _, g := range th.Grays {
			if g == "" {
				t.Errorf("%s: a gray is missing from the palette", th.Name)
			}
		}
	}
	if len(palettes) < 20 {
		t.Errorf("%d palettes", len(palettes))
	}
}

func TestCSS(t *testing.T) {
	css := Violet.CSS()
	for _, want := range []string{`[data-bs-theme-pack="violet"]:not([data-bs-theme-primary]){--tblr-primary:light-dark(oklch(`,
		`:not([data-bs-theme-base]){--tblr-gray-50:`, `"Geist Variable"`, `--tblr-font-monospace:"Geist Mono Variable"`} {
		if !strings.Contains(css, want) {
			t.Errorf("missing %s in %s", want, css)
		}
	}
	if Zinc.CSS() != "" || Tabler.CSS() != "" {
		t.Error("themes made of Tabler's settings need no CSS")
	}
	if rgb := rgbOf("oklch(54.6% 0.1724 254.2deg)"); rgb != "6, 111, 209" { // Tabler's blue
		t.Errorf("rgb %s", rgb)
	}
}

func TestApply(t *testing.T) {
	p := layout.PageProps{Appearance: theme.Settings{Radius: theme.Radius0}}
	Midnight.Apply(&p)
	if p.Theme != "dark" || p.HTMLAttributes[Attribute] != "midnight" || p.Appearance.Base != theme.Slate || p.Appearance.Radius != theme.Radius0 {
		t.Errorf("%+v", p)
	}
	var b bytes.Buffer
	if err := p.Head.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"@fontsource-variable/manrope@", "<style>", "tabler-theme-pack"} {
		if !strings.Contains(b.String(), want) {
			t.Errorf("head misses %s", want)
		}
	}
	if s := Section(); len(s.Settings[0].Options) != len(All) || s.Settings[0].Options[0].Value != "default" {
		t.Errorf("section %+v", s.Settings[0].Options[0])
	}
}

func TestApplyMarketing(t *testing.T) {
	p := marketing.PageProps{Appearance: theme.Settings{Radius: theme.Radius2}}
	Graphite.ApplyMarketing(&p)
	if p.HTMLAttributes[Attribute] != "graphite" || p.Appearance.Base != theme.Neutral || p.Appearance.Radius != theme.Radius2 || p.Head == nil {
		t.Errorf("%+v", p)
	}
}

// Head without arguments renders the CSS of all themes as a sheet, so static.WithKitCSS can put it into the
// bundle instead of every page.
func TestStylesIsASheet(t *testing.T) {
	var b strings.Builder
	require.NoError(t, Styles().Render(templ.InitializeContext(context.Background()), &b))
	assert.Contains(t, b.String(), `[data-bs-theme-pack="graphite"]`)
	assert.Contains(t, csp.All(), `[data-bs-theme-pack="graphite"]`)
}

// A theme's corners hold with only its attribute on <html> (data-bs-theme-pack, e.g. through kit.css),
// unless the page sets data-bs-theme-radius itself (found by Vergabe).
func TestCSSRadius(t *testing.T) {
	assert.Contains(t, Violet.CSS(), `[data-bs-theme-pack="violet"]:not([data-bs-theme-radius]){--tblr-border-radius-scale:1.5}`)
	assert.Contains(t, Graphite.CSS(), `:not([data-bs-theme-radius]){--tblr-border-radius-scale:0.5}`)
	assert.Empty(t, Zinc.CSS(), "Tabler's own corners need no rule")
}

// An app page (blocks/shell) takes a theme like the other pages: attributes on <html> and the head.
func TestApplyShell(t *testing.T) {
	p := shell.Props{Title: "x", Appearance: theme.Settings{Radius: theme.Radius0}}
	Midnight.ApplyShell(&p)
	var b bytes.Buffer
	require.NoError(t, shell.Page(p).Render(templ.InitializeContext(context.Background()), &b))
	h := b.String()
	for _, want := range []string{`data-bs-theme-pack="midnight"`, `data-bs-theme="dark"`, `data-bs-theme-radius="0"`, `data-bs-theme-base="slate"`, `lang="en"`, "tabler-theme-pack"} {
		assert.Contains(t, h, want)
	}
}

// Fixed links only the theme's own fonts and CSS and leaves out the panel's script (found by Vergabe: the
// sign-in pages linked the fonts of every theme).
func TestFixed(t *testing.T) {
	p := auth.PageProps{}
	Violet.Fixed().ApplyAuth(&p)
	var b bytes.Buffer
	require.NoError(t, p.Head.Render(templ.InitializeContext(context.Background()), &b))
	h := b.String()
	assert.Equal(t, 2, strings.Count(h, `rel="stylesheet"`), "Geist and Geist Mono: %s", h)
	assert.Contains(t, h, `[data-bs-theme-pack="violet"]`)
	assert.NotContains(t, h, `[data-bs-theme-pack="graphite"]`)
	assert.NotContains(t, h, "tabler-theme-pack")
	assert.Equal(t, "violet", p.HTMLAttributes[Attribute])

	csp.SetExternal(true)
	defer csp.SetExternal(false)
	b.Reset()
	require.NoError(t, p.Head.Render(templ.InitializeContext(context.Background()), &b))
	assert.NotContains(t, b.String(), "<style", "with kit.css the CSS is in the bundle")
}
