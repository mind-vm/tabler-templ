// Package themes holds ready-made looks for Tabler, in the spirit of shadcn/ui's themes: each Theme is an
// accent color, a gray base, a font and a corner radius, for light and dark mode. Some are combinations of
// Tabler's own theme settings (package theme); others bring their own palette from Tailwind CSS's colors,
// which shadcn builds on, and a web font. A theme only sets Tabler's CSS custom properties: the markup of
// the components stays Tabler's.
//
//	p := layout.PageProps{Title: "Orders", Head: head(), …}
//	themes.Rose.Apply(&p)            // the look, its CSS and fonts in the head
//
// Users can switch themes in layout.ThemeSettings with Section(); the choice is saved in the browser.
// An app defines its own theme the same way, as a Theme value.
package themes

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/mind-vm/tabler-templ/cdn"
	"github.com/mind-vm/tabler-templ/theme"
)

// Attribute is the attribute on <html> that names the theme.
const Attribute = "data-bs-theme-pack"

// FontsVersion is the version of the Fontsource packages the fonts load from (jsDelivr).
const FontsVersion = "5.3.0"

// Font is a web font from Fontsource (@fontsource-variable/<Package>), loaded from cdn.Base (jsDelivr or
// the app's own copy).
type Font struct {
	Family  string // the CSS family, e.g. "Inter Variable"
	Package string // e.g. "inter"
}

// URL is the font's stylesheet on cdn.Base; the browser downloads a font file only when the font is used.
func (f Font) URL() string {
	return cdn.URL("@fontsource-variable/"+f.Package, FontsVersion, "index.css")
}

// Fonts the themes use.
var (
	Inter         = Font{"Inter Variable", "inter"}
	Geist         = Font{"Geist Variable", "geist"}
	GeistMono     = Font{"Geist Mono Variable", "geist-mono"}
	DMSans        = Font{"DM Sans Variable", "dm-sans"}
	Nunito        = Font{"Nunito Variable", "nunito"}
	Manrope       = Font{"Manrope Variable", "manrope"}
	Outfit        = Font{"Outfit Variable", "outfit"}
	Figtree       = Font{"Figtree Variable", "figtree"}
	Lora          = Font{"Lora Variable", "lora"}
	JetBrainsMono = Font{"JetBrains Mono Variable", "jetbrains-mono"}
)

// Theme is a look.
type Theme struct {
	Name        string // its value of Attribute, e.g. "rose"; "" for Tabler's own look
	Label       string // e.g. "Rose"
	Description string
	// Settings are Tabler's theme settings the look uses; they are attributes on <html> like any other
	// setting, so a user's own choice in the panel still wins.
	Settings theme.Settings
	// Mode is the default color mode: "light", "dark" or "" for the browser's preference.
	Mode string
	// Primary is the accent color for light mode and PrimaryDark for dark mode (default Primary), any CSS
	// color, when Settings.Primary has no Tabler color for it. PrimaryFg is the text on it.
	Primary, PrimaryDark, PrimaryFg string
	// Grays is a gray palette (shades 50 to 950) when Settings.Base has none for it.
	Grays []string
	Font  *Font // the body font, instead of Settings.Font
	Mono  *Font // the monospace font
}

// Tailwind is a color of Tailwind CSS's palettes, e.g. Tailwind("rose", 600).
func Tailwind(name string, shade int) string {
	return palettes[name][shade]
}

// TailwindGrays is a Tailwind palette as Grays, e.g. TailwindGrays("mauve").
func TailwindGrays(name string) []string {
	var gs []string
	for _, s := range []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950} {
		gs = append(gs, palettes[name][s])
	}
	return gs
}

// Attrs are the theme's attributes for <html>: its name and its Tabler settings.
func (t Theme) Attrs() map[string]any {
	a := map[string]any{}
	for k, v := range t.Settings.Attrs() {
		a[k] = v
	}
	if t.Name != "" {
		a[Attribute] = t.Name
	}
	return a
}

func (t Theme) fonts() []Font {
	var fs []Font
	if t.Font != nil {
		fs = append(fs, *t.Font)
	}
	if t.Mono != nil {
		fs = append(fs, *t.Mono)
	}
	return fs
}

// CSS is the theme's own rules: custom properties that apply while the user has not picked that setting
// themselves (a data-bs-theme-* attribute on <html> wins).
func (t Theme) CSS() string {
	if t.Name == "" {
		return ""
	}
	sel := `[` + Attribute + `="` + t.Name + `"]`
	var b strings.Builder
	if t.Primary != "" {
		dark := t.PrimaryDark
		if dark == "" {
			dark = t.Primary
		}
		fmt.Fprintf(&b, "%s:not([data-bs-theme-primary]){--tblr-primary:light-dark(%s,%s);", sel, t.Primary, dark)
		if rgb := rgbOf(t.Primary); rgb != "" {
			fmt.Fprintf(&b, "--tblr-primary-rgb:%s;", rgb)
		}
		if t.PrimaryFg != "" {
			fmt.Fprintf(&b, "--tblr-primary-fg:%s;", t.PrimaryFg)
		}
		b.WriteString("}")
	}
	if len(t.Grays) == 11 {
		fmt.Fprintf(&b, "%s:not([data-bs-theme-base]){", sel)
		for i, s := range []int{50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950} {
			fmt.Fprintf(&b, "--tblr-gray-%d:%s;", s, t.Grays[i])
		}
		b.WriteString("}")
	}
	if t.Font != nil {
		fmt.Fprintf(&b, "%s:not([data-bs-theme-font]){--tblr-font-sans-serif:%q,system-ui,-apple-system,\"Segoe UI\",Roboto,sans-serif;--tblr-body-font-family:var(--tblr-font-sans-serif)}", sel, t.Font.Family)
	}
	if t.Mono != nil {
		fmt.Fprintf(&b, "%s{--tblr-font-monospace:%q,ui-monospace,SFMono-Regular,Menlo,monospace}", sel, t.Mono.Family)
	}
	return b.String()
}

// rgbOf converts an oklch() color to "r, g, b" for Tabler's --tblr-primary-rgb; "" for other colors.
func rgbOf(c string) string {
	if !strings.HasPrefix(c, "oklch(") {
		return ""
	}
	f := strings.Fields(strings.NewReplacer("oklch(", "", ")", "", "%", "", "deg", "").Replace(c))
	if len(f) < 3 {
		return ""
	}
	num := func(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
	l, ch, h := num(f[0])/100, num(f[1]), num(f[2])*math.Pi/180
	a, bb := ch*math.Cos(h), ch*math.Sin(h)
	l_ := math.Pow(l+0.3963377774*a+0.2158037573*bb, 3)
	m_ := math.Pow(l-0.1055613458*a-0.0638541728*bb, 3)
	s_ := math.Pow(l-0.0894841775*a-1.2914855480*bb, 3)
	lin := [3]float64{
		4.0767416621*l_ - 3.3077115913*m_ + 0.2309699292*s_,
		-1.2684380046*l_ + 2.6097574011*m_ - 0.3413193965*s_,
		-0.0041960863*l_ - 0.7034186147*m_ + 1.7076147010*s_,
	}
	var out []string
	for _, x := range lin {
		if x <= 0.0031308 {
			x *= 12.92
		} else {
			x = 1.055*math.Pow(x, 1/2.4) - 0.055
		}
		out = append(out, strconv.Itoa(int(math.Round(math.Max(0, math.Min(1, x))*255))))
	}
	return strings.Join(out, ", ")
}

func ptr(f Font) *Font { return &f }

// The themes. Tabler, Zinc, Slate, Ocean, Forest and Paper are combinations of Tabler's own settings; the
// others bring a Tailwind palette and a web font.
var (
	Tabler = Theme{Label: "Tabler", Description: "Tabler's own look: blue on neutral gray."}
	Zinc   = Theme{Name: "zinc", Label: "Zinc", Description: "Black and white on cool gray, like shadcn/ui's default.",
		Settings: theme.Settings{Primary: theme.Inverted, Base: theme.Zinc, Radius: theme.Radius1}}
	Slate = Theme{Name: "slate", Label: "Slate", Description: "Indigo on blue-gray.",
		Settings: theme.Settings{Primary: theme.Indigo, Base: theme.Slate}}
	Ocean = Theme{Name: "ocean", Label: "Ocean", Description: "Azure on blue-gray, with rounder corners.",
		Settings: theme.Settings{Primary: theme.Azure, Base: theme.Slate, Radius: theme.Radius15}}
	Forest = Theme{Name: "forest", Label: "Forest", Description: "Green on warm gray, with tighter corners.",
		Settings: theme.Settings{Primary: theme.Green, Base: theme.Stone, Radius: theme.Radius05}}
	Paper = Theme{Name: "paper", Label: "Paper", Description: "Black and white, serif type and square corners, like print.",
		Settings: theme.Settings{Primary: theme.Inverted, Base: theme.Stone, Font: theme.Serif, Radius: theme.Radius0}}

	Rose = Theme{Name: "rose", Label: "Rose", Description: "Rose on zinc, set in Inter.",
		Settings: theme.Settings{Base: theme.Zinc}, Primary: Tailwind("rose", 600), PrimaryDark: Tailwind("rose", 500), Font: ptr(Inter)}
	Violet = Theme{Name: "violet", Label: "Violet", Description: "Violet on mauve, set in Geist, with round corners.",
		Settings: theme.Settings{Radius: theme.Radius15}, Primary: Tailwind("violet", 600), PrimaryDark: Tailwind("violet", 500),
		Grays: TailwindGrays("mauve"), Font: ptr(Geist), Mono: ptr(GeistMono)}
	Amber = Theme{Name: "amber", Label: "Amber", Description: "Amber with dark text on warm gray, set in Nunito.",
		Settings: theme.Settings{Base: theme.Stone, Radius: theme.Radius15}, Primary: Tailwind("amber", 500), PrimaryDark: Tailwind("amber", 400),
		PrimaryFg: Tailwind("amber", 950), Font: ptr(Nunito)}
	Emerald = Theme{Name: "emerald", Label: "Emerald", Description: "Emerald on olive gray, set in DM Sans.",
		Primary: Tailwind("emerald", 600), PrimaryDark: Tailwind("emerald", 500), Grays: TailwindGrays("olive"), Font: ptr(DMSans)}
	Graphite = Theme{Name: "graphite", Label: "Graphite", Description: "Monochrome with Geist and Geist Mono, small corners.",
		Settings: theme.Settings{Base: theme.Neutral, Radius: theme.Radius05}, Primary: Tailwind("neutral", 900), PrimaryDark: Tailwind("neutral", 50),
		PrimaryFg: "light-dark(var(--tblr-light), var(--tblr-dark))", Font: ptr(Geist), Mono: ptr(GeistMono)}
	Sunset = Theme{Name: "sunset", Label: "Sunset", Description: "Orange on taupe, set in Outfit, with the roundest corners.",
		Settings: theme.Settings{Radius: theme.Radius2}, Primary: Tailwind("orange", 600), PrimaryDark: Tailwind("orange", 500),
		Grays: TailwindGrays("taupe"), Font: ptr(Outfit)}
	Midnight = Theme{Name: "midnight", Label: "Midnight", Description: "Dark by default: indigo on blue-gray, set in Manrope.",
		Mode: "dark", Settings: theme.Settings{Base: theme.Slate}, Primary: Tailwind("indigo", 600), PrimaryDark: Tailwind("indigo", 400),
		Font: ptr(Manrope), Mono: ptr(JetBrainsMono)}
)

func init() {
	for _, t := range All {
		for _, f := range t.fonts() {
			cdn.Register("theme:"+strings.ToLower(t.Label), "@fontsource-variable/"+f.Package, FontsVersion, "index.css")
		}
	}
}

// All are the themes in the order the panel shows them.
var All = []Theme{Tabler, Zinc, Slate, Ocean, Forest, Paper, Rose, Violet, Amber, Emerald, Graphite, Sunset, Midnight}

// Find returns the theme with that name ("" or "tabler" is Tabler's own look).
func Find(name string) (Theme, bool) {
	if name == "tabler" {
		name = ""
	}
	for _, t := range All {
		if t.Name == name {
			return t, true
		}
	}
	return Theme{}, false
}
