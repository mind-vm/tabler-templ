package layout

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mind-vm/tabler-templ/theme"
)

// ThemeControl is how a setting's options are drawn.
type ThemeControl string

const (
	ControlRadio  ThemeControl = "radio"  // chips
	ControlColor  ThemeControl = "color"  // color swatches
	ControlFont   ThemeControl = "font"   // "Aa" in each font
	ControlTile   ThemeControl = "tile"   // a thumbnail per value
	ControlPreset ThemeControl = "preset" // thumbnails that set several keys at once
)

// ThemeOption is one choice of a ThemeSetting. A preset's Values are the keys and values it sets.
type ThemeOption struct {
	Value, Label string
	Values       map[string]string
}

// ThemeSetting is one setting of the panel. Key is the switcher's key (it becomes data-bs-<key> on <html>
// and tabler-<key> in localStorage), or a preset group's own name.
type ThemeSetting struct {
	Key, Legend, Hint string
	Control           ThemeControl
	Options           []ThemeOption
	SidebarOnly       bool // hidden while the navigation is horizontal
	NavbarOnly        bool // hidden while the navigation is vertical
}

// ThemeSection is a titled group of settings.
type ThemeSection struct {
	Title    string
	Settings []ThemeSetting
}

// ThemeDefaults are the switcher's settings with their default values, as in Tabler's theme script.
var ThemeDefaults = map[string]string{
	"theme": "auto", "theme-base": "neutral", "theme-font": "sans-serif", "theme-primary": "blue", "theme-radius": "1",
	"layout": "default", "navbar": "default", "navbar-position": "horizontal", "navbar-style": "default",
	"navbar-theme": "default", "offcanvas-style": "default", "sidebar": "default",
}

func titled(s string) string { return strings.ToUpper(s[:1]) + s[1:] }

func options[T ~string](values []T) []ThemeOption {
	os := make([]ThemeOption, len(values))
	for i, v := range values {
		os[i] = ThemeOption{Value: string(v), Label: titled(string(v))}
	}
	return os
}

// AppearanceSection is the panel's first section: color scheme, accent color, font, base and radius.
func AppearanceSection() ThemeSection {
	return ThemeSection{Title: "Appearance", Settings: []ThemeSetting{
		{Key: "scheme", Legend: "Color scheme", Hint: "The page and its navigation, in one pick.", Control: ControlPreset, Options: []ThemeOption{
			{Value: "auto", Label: "Auto", Values: map[string]string{"theme": "auto", "navbar-theme": "default"}},
			{Value: "light", Label: "Light", Values: map[string]string{"theme": "light", "navbar-theme": "default"}},
			{Value: "dark", Label: "Dark", Values: map[string]string{"theme": "dark", "navbar-theme": "default"}},
			{Value: "dark-nav", Label: "Dark nav", Values: map[string]string{"theme": "light", "navbar-theme": "dark"}},
			{Value: "colored", Label: "Colored", Values: map[string]string{"theme": "light", "navbar-theme": "primary"}},
		}},
		{Key: "theme-primary", Legend: "Accent color", Hint: "The accent color used across the app.", Control: ControlColor, Options: options(theme.Colors)},
		{Key: "theme-font", Legend: "Font family", Hint: "The typeface used across the app.", Control: ControlFont, Options: options(theme.Fonts)},
		{Key: "theme-base", Legend: "Theme base", Hint: "The gray palette behind every surface.", Control: ControlTile, Options: options(theme.Bases)},
		{Key: "theme-radius", Legend: "Corner radius", Hint: "How round corners are drawn.", Control: ControlTile, Options: func() []ThemeOption {
			os := make([]ThemeOption, len(theme.Radii))
			for i, r := range theme.Radii {
				os[i] = ThemeOption{Value: string(r), Label: string(r)}
			}
			return os
		}()},
	}}
}

// LayoutSection is the panel's second section: navigation, width, navbar, offcanvas and sidebar.
func LayoutSection() ThemeSection {
	return ThemeSection{Title: "Layout", Settings: []ThemeSetting{
		{Key: "navbar-position", Legend: "Navigation", Hint: "Where the main menu lives.", Control: ControlTile,
			Options: []ThemeOption{{Value: "horizontal", Label: "Navbar"}, {Value: "vertical", Label: "Sidebar"}}},
		{Key: "layout", Legend: "Container width", Hint: "How wide the page content runs.", Control: ControlTile,
			Options: options([]string{"default", "fluid", "boxed"})},
		{Key: "navbar", Legend: "Navbar behavior", Hint: "Keep the navbar visible while scrolling.", Control: ControlTile, NavbarOnly: true,
			Options: []ThemeOption{{Value: "default", Label: "Default"}, {Value: "sticky", Label: "Sticky"}}},
		{Key: "navbar-style", Legend: "Navigation style", Hint: "Attach the navigation to the edges, or let it float.", Control: ControlTile,
			Options: []ThemeOption{{Value: "default", Label: "Attached"}, {Value: "floating", Label: "Floating"}}},
		{Key: "offcanvas-style", Legend: "Offcanvas style", Hint: "Attach side panels to the edges, or let them float.", Control: ControlTile,
			Options: []ThemeOption{{Value: "default", Label: "Attached"}, {Value: "floating", Label: "Floating"}}},
		{Key: "sidebar", Legend: "Sidebar", Hint: "How the sidebar behaves.", Control: ControlTile, SidebarOnly: true,
			Options: []ThemeOption{{Value: "default", Label: "Expanded"}, {Value: "folded", Label: "Folded"}, {Value: "folded-hover", Label: "On hover"}}},
	}}
}

// themeChecked tells whether an option is the default (presets: every key it sets is at its default).
func themeChecked(s ThemeSetting, o ThemeOption) bool {
	if s.Control == ControlPreset {
		for k, v := range o.Values {
			if ThemeDefaults[k] != v {
				return false
			}
		}
		return true
	}
	return ThemeDefaults[s.Key] == o.Value
}

// themeConfig is what the panel's script needs: defaults, presets and the settings shown for one
// navigation only.
func themeConfig(sections []ThemeSection) string {
	presets := map[string]map[string]map[string]string{}
	var sidebarOnly, navbarOnly []string
	for _, sec := range sections {
		for _, s := range sec.Settings {
			if s.Control == ControlPreset {
				presets[s.Key] = map[string]map[string]string{}
				for _, o := range s.Options {
					presets[s.Key][o.Value] = o.Values
				}
			}
			if s.SidebarOnly {
				sidebarOnly = append(sidebarOnly, s.Key)
			}
			if s.NavbarOnly {
				navbarOnly = append(navbarOnly, s.Key)
			}
		}
	}
	b, _ := json.Marshal(map[string]any{"defaults": ThemeDefaults, "presets": presets, "sidebarOnly": sidebarOnly, "navbarOnly": navbarOnly})
	return string(b)
}

// The thumbnails, drawn as inline SVG on the theme's own custom properties so they follow the color mode,
// the base and the accent color (a port of the preview's ThemeSettings.astro).

const (
	tPage    = "var(--tblr-bg-surface-secondary)"
	tSurface = "var(--tblr-bg-surface)"
	tBorder  = "var(--tblr-border-color)"
	tPrimary = "var(--tblr-primary)"
	tOnAcc   = "rgba(255, 255, 255, 0.7)"
)

var (
	tShade  = map[string]string{"light": "var(--tblr-gray-50)", "dark": "var(--tblr-gray-900)"}
	tBar    = map[string]string{"light": "var(--tblr-white)", "dark": "var(--tblr-gray-800)", "primary": "var(--tblr-primary)"}
	tLineOn = map[string]string{"light": "var(--tblr-gray-300)", "dark": "var(--tblr-gray-700)"}
)

func num(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

func tSVG(inner, attrs string) string {
	return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 96 64" class="form-imagecheck-image w-100" aria-hidden="true" focusable="false"` + attrs + `>` + inner + `</svg>`
}

func tFrame(fill string) string {
	return `<rect x="0.5" y="0.5" width="95" height="63" rx="6" fill="` + fill + `" stroke="` + tBorder + `" />`
}

const tOutline = `<rect x="0.5" y="0.5" width="95" height="63" rx="6" fill="none" stroke="` + tBorder + `" />`

func tTopBar(fill string) string {
	return `<path d="M6.5 0.5h83a6 6 0 0 1 6 6v13.5h-95v-13.5a6 6 0 0 1 6 -6z" fill="` + fill + `" /><rect x="0.5" y="19.5" width="95" height="0.5" fill="` + tBorder + `" />`
}

func tRail(w float64, fill string) string {
	return fmt.Sprintf(`<path d="M6.5 0.5h%sv63h-%sa6 6 0 0 1 -6 -6v-51a6 6 0 0 1 6 -6z" fill="%s" /><rect x="%s" y="0.5" width="0.5" height="63" fill="%s" />`,
		num(w-6), num(w-6), fill, num(w), tBorder)
}

func tLines(x, y float64, widths []float64, fill string, step float64) string {
	var b strings.Builder
	for i, w := range widths {
		fmt.Fprintf(&b, `<rect x="%s" y="%s" width="%s" height="5" rx="2.5" fill="%s" />`, num(x), num(y+float64(i)*step), num(w), fill)
	}
	return b.String()
}

func tDots(x, y float64, n int, fill string) string {
	var b strings.Builder
	for i := range n {
		fmt.Fprintf(&b, `<circle cx="%s" cy="%s" r="2.5" fill="%s" />`, num(x), num(y+float64(i)*11), fill)
	}
	return b.String()
}

func tPageHalf(mode string, right bool) string {
	if right {
		return `<path d="M48 0.5h41.5a6 6 0 0 1 6 6v51a6 6 0 0 1 -6 6h-41.5z" fill="` + tShade[mode] + `" />`
	}
	return `<path d="M6.5 0.5h41.5v63h-41.5a6 6 0 0 1 -6 -6v-51a6 6 0 0 1 6 -6z" fill="` + tShade[mode] + `" />`
}

func tBarHalf(fill string, right bool) string {
	if right {
		return `<path d="M48 0.5h41.5a6 6 0 0 1 6 6v13.5h-47.5z" fill="` + fill + `" />`
	}
	return `<path d="M6.5 0.5h41.5v19.5h-47.5v-13.5a6 6 0 0 1 6 -6z" fill="` + fill + `" />`
}

// presetSVG draws a preset tile: page halves and navigation bar halves in the preset's modes.
func presetSVG(values map[string]string) string {
	mode := values["theme"]
	if mode == "" {
		mode = "light"
	}
	nav := values["navbar-theme"]
	first, second := mode, mode
	if mode == "auto" {
		first, second = "light", "dark"
	}
	barFill := func(m string) string {
		switch nav {
		case "primary":
			return tBar["primary"]
		case "dark":
			return tBar["dark"]
		}
		return tBar[m]
	}
	return tSVG(tPageHalf(first, false)+tPageHalf(second, true)+tBarHalf(barFill(first), false)+tBarHalf(barFill(second), true)+
		tLines(12, 30, []float64{46, 72, 30}, tLineOn[first], 11)+tOutline, "")
}

func tCard(x, y, w, h float64) string {
	return fmt.Sprintf(`<rect x="%s" y="%s" width="%s" height="%s" rx="3" fill="%s" stroke="%s" />`, num(x), num(y), num(w), num(h), tSurface, tBorder) +
		tLines(x+6, y+7, []float64{float64(int(w*0.6 + 0.5)), float64(int(w*0.8 + 0.5)), float64(int(w*0.4 + 0.5))}, tBorder, 9)
}

func tThumb(y float64) string {
	return `<rect x="90" y="` + num(y) + `" width="3" height="14" rx="1.5" fill="` + tBorder + `" />`
}

func tNavItems(xs []float64, y float64) string {
	var b strings.Builder
	for _, x := range xs {
		fmt.Fprintf(&b, `<rect x="%s" y="%s" width="14" height="4" rx="2" fill="%s" />`, num(x), num(y), tOnAcc)
	}
	return b.String()
}

// tileSVG draws the thumbnail of one value of one setting.
func tileSVG(key, value string) string {
	content := tLines(12, 30, []float64{46, 72, 30}, tBorder, 11)
	switch key {
	case "navbar-position":
		if value == "vertical" {
			return tSVG(tFrame(tPage)+tRail(30, tPrimary)+tLines(7, 10, []float64{16, 16, 16, 16}, tOnAcc, 10)+tLines(40, 12, []float64{40, 48, 28}, tBorder, 11)+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tTopBar(tPrimary)+tNavItems([]float64{8, 26, 44}, 8)+content+tOutline, "")
	case "layout":
		switch value {
		case "boxed":
			return tSVG(`<rect x="0.5" y="0.5" width="95" height="63" rx="6" fill="var(--tblr-bg-surface-inverted)" /><rect x="12" y="8" width="72" height="48" rx="4" fill="`+tPage+`" />`+tCard(18, 16, 60, 32)+tOutline, "")
		case "fluid":
			return tSVG(tFrame(tPage)+tCard(4, 8, 88, 48)+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tCard(22, 8, 52, 48)+tOutline, "")
	case "navbar":
		if value == "sticky" {
			return tSVG(tFrame(tPage)+tLines(12, 4, []float64{46, 72, 30, 60}, tBorder, 11)+tThumb(44)+`<rect x="0.5" y="19.5" width="95" height="5" fill="var(--tblr-gray-900)" opacity="0.18" />`+tTopBar(tPrimary)+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tTopBar(tPrimary)+content+tThumb(24)+tOutline, "")
	case "navbar-style":
		if value == "floating" {
			return tSVG(tFrame(tPage)+`<rect x="5" y="5" width="86" height="14" rx="4" fill="`+tPrimary+`" />`+tNavItems([]float64{11, 29, 47}, 10)+content+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tTopBar(tPrimary)+tNavItems([]float64{8, 26, 44}, 8)+content+tOutline, "")
	case "offcanvas-style":
		dim := `<rect x="0.5" y="0.5" width="95" height="63" rx="6" fill="var(--tblr-gray-900)" opacity="0.18" />`
		if value == "floating" {
			return tSVG(tFrame(tPage)+tLines(10, 12, []float64{40, 48, 28}, tBorder, 11)+dim+`<rect x="53" y="5" width="38" height="54" rx="4" fill="`+tSurface+`" stroke="`+tBorder+`" />`+tLines(59, 12, []float64{20, 26, 14}, tBorder, 11)+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tLines(10, 12, []float64{40, 48, 28}, tBorder, 11)+dim+`<path d="M58 0.5h31.5a6 6 0 0 1 6 6v51a6 6 0 0 1 -6 6h-31.5z" fill="`+tSurface+`" /><rect x="57.5" y="0.5" width="0.5" height="63" fill="`+tBorder+`" />`+tLines(64, 12, []float64{20, 26, 14}, tBorder, 11)+tOutline, "")
	case "sidebar":
		folded := tFrame(tPage) + tRail(16, tPrimary) + tDots(8, 12, 4, tOnAcc) + tLines(26, 12, []float64{50, 60, 34}, tBorder, 11)
		switch value {
		case "folded":
			return tSVG(folded+tOutline, "")
		case "folded-hover":
			return tSVG(folded+`<rect x="16" y="0.5" width="4" height="63" fill="var(--tblr-gray-900)" opacity="0.12" /><rect x="16.5" y="0.5" width="28" height="63" fill="`+tSurface+`" stroke="`+tBorder+`" />`+tLines(21, 9.5, []float64{18, 18, 18, 18}, tBorder, 11)+tOutline, "")
		}
		return tSVG(tFrame(tPage)+tRail(30, tPrimary)+tLines(7, 10, []float64{16, 16, 16, 16}, tOnAcc, 10)+tLines(40, 12, []float64{40, 48, 28}, tBorder, 11)+tOutline, "")
	case "theme-base":
		var b strings.Builder
		for i, shade := range []int{100, 300, 500, 700, 900} {
			fmt.Fprintf(&b, `<rect x="%d" y="24" width="12" height="16" rx="3" fill="var(--tblr-gray-%d)" />`, 10+i*16, shade)
		}
		return tSVG(tFrame(tSurface)+b.String(), ` data-bs-theme-base="`+value+`"`)
	case "theme-radius":
		r, _ := strconv.ParseFloat(value, 64)
		return tSVG(tFrame(tPage)+`<rect x="28" y="12" width="40" height="40" rx="`+num(r*10)+`" fill="`+tPrimary+`" />`, "")
	}
	return tSVG(tFrame(tPage)+tOutline, "")
}
