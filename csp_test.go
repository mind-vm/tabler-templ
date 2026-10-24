package tablertempl_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/assets"
	"github.com/mind-vm/tabler-templ/blocks/detail"
	"github.com/mind-vm/tabler-templ/blocks/invoice"
	"github.com/mind-vm/tabler-templ/blocks/layout"
	"github.com/mind-vm/tabler-templ/blocks/maps"
	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/chart"
	"github.com/mind-vm/tabler-templ/confetti"
	"github.com/mind-vm/tabler-templ/datepicker"
	"github.com/mind-vm/tabler-templ/form"
	"github.com/mind-vm/tabler-templ/icon"
	"github.com/mind-vm/tabler-templ/lightbox"
	"github.com/mind-vm/tabler-templ/mapview"
	"github.com/mind-vm/tabler-templ/modal"
	"github.com/mind-vm/tabler-templ/player"
	"github.com/mind-vm/tabler-templ/sidebar"
	"github.com/mind-vm/tabler-templ/themes"
	"github.com/mind-vm/tabler-templ/vectormap"
)

var openTag = regexp.MustCompile(`<(script|style)\b[^>]*>`)

// TestNonceRendered renders the packages' scripts and styles with templ's nonce: every <script> and
// <style> carries it.
func TestNonceRendered(t *testing.T) {
	ctx := templ.WithNonce(context.Background(), "n0nce")
	for name, c := range map[string]templ.Component{
		"assets":     templ.Join(assets.Head(), assets.Theme(), assets.Scripts()),
		"chart":      templ.Join(chart.Library(), chart.Script()),
		"confetti":   confetti.Burst(confetti.Props{}),
		"datepicker": templ.Join(datepicker.Library(), datepicker.Script()),
		"detail":     templ.Join(detail.Styles(), detail.Script(detail.ScriptConfig{})),
		"form":       form.PasswordToggleScript(),
		"icon":       icon.Styles(),
		"invoice":    invoice.PrintButton(),
		"layout":     layout.ThemeSettings(layout.ThemeSettingsProps{}),
		"lightbox":   lightbox.Library(),
		"maps":       maps.Script(),
		"mapview":    templ.Join(mapview.Library(), mapview.Script()),
		"modal":      modal.DialogScript(),
		"player":     templ.Join(player.Library(), player.Script()),
		"shell":      shell.Styles(),
		"sidebar":    templ.Join(sidebar.Styles(), sidebar.Script()),
		"themes":     themes.Head(themes.All...),
		"vectormap":  vectormap.Library("world"),
	} {
		html := render(t, ctx, c)
		tags := openTag.FindAllString(html, -1)
		if len(tags) == 0 {
			t.Errorf("%s: no <script> or <style>", name)
		}
		for _, tag := range tags {
			if !strings.Contains(tag, `nonce="n0nce"`) {
				t.Errorf("%s: %s without the nonce", name, tag)
			}
		}
	}
}

// TestNonceSource checks the packages' sources (not the demo): every line that writes a <script> or <style>
// goes through package internal/csp, which adds templ's nonce.
func TestNonceSource(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch path {
			case ".git", "cmd", "docs", "internal/csp", "internal/tools":
				return filepath.SkipDir
			}
			return nil
		}
		if !(strings.HasSuffix(path, ".templ") || strings.HasSuffix(path, ".go")) ||
			strings.HasSuffix(path, "_templ.go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if (strings.Contains(line, "<script") || strings.Contains(line, "<style")) && !strings.Contains(line, "csp.") {
				t.Errorf("%s:%d: <script> or <style> without csp: %s", path, i+1, strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
