package themes

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/blocks/layout"
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
