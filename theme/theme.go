// Package theme holds Tabler's theme settings: the accent color, the gray base, the font and the corner
// radius. They are data-bs-theme-* attributes on <html>, styled by tabler-themes.css (load it with
// assets.Props{Themes: true}). layout.Page and auth.Page render them from their Appearance;
// layout.ThemeSettings lets users change them in the browser.
//
//	layout.PageProps{Appearance: theme.Settings{Primary: theme.Teal, Font: theme.Serif, Radius: theme.Radius2}}
package theme

import "github.com/a-h/templ"

// Color is an accent color.
type Color string

const (
	Blue     Color = "blue" // the default
	Azure    Color = "azure"
	Indigo   Color = "indigo"
	Purple   Color = "purple"
	Pink     Color = "pink"
	Red      Color = "red"
	Orange   Color = "orange"
	Yellow   Color = "yellow"
	Lime     Color = "lime"
	Green    Color = "green"
	Teal     Color = "teal"
	Cyan     Color = "cyan"
	Inverted Color = "inverted" // black on light, white on dark
)

// Base is the gray palette behind every surface.
type Base string

const (
	Slate   Base = "slate"
	Gray    Base = "gray"
	Zinc    Base = "zinc"
	Neutral Base = "neutral" // the default
	Stone   Base = "stone"
	Mist    Base = "mist"
)

// Font is the typeface.
type Font string

const (
	SansSerif Font = "sans-serif" // the default
	Serif     Font = "serif"
	Monospace Font = "monospace"
	Comic     Font = "comic"
)

// Radius is how round corners are drawn, from square to fully round.
type Radius string

const (
	Radius0  Radius = "0"
	Radius05 Radius = "0.5"
	Radius1  Radius = "1" // the default
	Radius15 Radius = "1.5"
	Radius2  Radius = "2"
)

// The choices in Tabler's order, e.g. for a settings form.
var (
	Colors = []Color{Blue, Azure, Indigo, Purple, Pink, Red, Orange, Yellow, Lime, Green, Teal, Cyan, Inverted}
	Bases  = []Base{Slate, Gray, Zinc, Neutral, Stone, Mist}
	Fonts  = []Font{SansSerif, Serif, Monospace, Comic}
	Radii  = []Radius{Radius0, Radius05, Radius1, Radius15, Radius2}
)

// Settings is a theme; empty fields keep Tabler's defaults.
type Settings struct {
	Primary Color
	Base    Base
	Font    Font
	Radius  Radius
}

// Attrs are the settings as attributes for <html>; defaults are left out, as Tabler's theme script does.
func (s Settings) Attrs() templ.Attributes {
	a := templ.Attributes{}
	set := func(key, value, def string) {
		if value != "" && value != def {
			a["data-bs-theme-"+key] = value
		}
	}
	set("primary", string(s.Primary), string(Blue))
	set("base", string(s.Base), string(Neutral))
	set("font", string(s.Font), string(SansSerif))
	set("radius", string(s.Radius), string(Radius1))
	return a
}
