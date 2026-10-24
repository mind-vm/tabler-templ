package layout

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/theme"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestThemeSettings(t *testing.T) {
	h := render(t, templ.Join(ThemeSettings(ThemeSettingsProps{}), ThemeSettings(ThemeSettingsProps{ID: "two", NoButton: true, Sections: []ThemeSection{AppearanceSection()}})))
	for _, want := range []string{`id="offcanvas-settings"`, `data-bs-target="#offcanvas-settings"`, `name="theme-primary" value="teal"`,
		`name="navbar-position" value="vertical"`, `data-bs-theme-base="zinc"`, `name="scheme" value="auto" class="form-imagecheck-input" checked`} {
		if !strings.Contains(h, want) {
			t.Errorf("missing %s", want)
		}
	}
	two := h[strings.Index(h, `id="two"`):]
	if strings.Contains(two, `name="navbar-position"`) {
		t.Error("only the appearance section was asked for")
	}
	if n := strings.Count(h, "form[data-theme-settings]"); n != 1 {
		t.Errorf("script rendered %d times", n)
	}
}

func TestPageAppearance(t *testing.T) {
	h := render(t, Page(PageProps{Theme: "dark", Appearance: theme.Settings{Primary: theme.Teal, Radius: theme.Radius2}, OffcanvasFloating: true}))
	for _, want := range []string{`data-bs-theme="dark"`, `data-bs-theme-primary="teal"`, `data-bs-theme-radius="2"`, `data-bs-offcanvas-style="floating"`} {
		if !strings.Contains(h, want) {
			t.Errorf("missing %s", want)
		}
	}
}
