package themes

import (
	"fmt"
	"math"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
)

// Signal is a colour that means something (green good, yellow attention, red problem, blue info). It is the
// same in every theme. Dot is for dots, lines and areas (at least 3:1 against the page); Text is its dark
// tone for text, in pills, statuses, alerts and on buttons (at least 4.5:1, WCAG 2.1 AA). The signal colour
// itself is never a text colour. DarkText and DarkDot are the tones in dark mode (DarkDot empty: Dot).
type Signal struct {
	Name     string   // "green"
	Classes  []string // Tabler's colour names that mean it: "green", "success"
	Dot      string
	Text     string
	DarkText string
	DarkDot  string
}

// Signals are the signal colours of the design concept „Gestaltung“ (GS-07, v0.5).
var Signals = []Signal{
	{Name: "green", Classes: []string{"green", "success"}, Dot: "#22A04B", Text: "#1E7537", DarkText: "#4ADE80"},
	{Name: "yellow", Classes: []string{"yellow", "warning"}, Dot: "#B57F00", Text: "#8A5A00", DarkText: "#FACC15"},
	{Name: "red", Classes: []string{"red", "danger"}, Dot: "#E5484D", Text: "#B42318", DarkText: "#F87171"},
	{Name: "blue", Classes: []string{"blue", "azure", "info"}, Dot: "#2563EB", Text: "#1D4ED8", DarkText: "#60A5FA", DarkDot: "#3B82F6"},
}

// Marking is a named marking with its own colours, e.g. text from outside or an internal note (GS-07,
// GS-76). Package mark draws it; a theme may set the custom properties --kit-mark-<name>(-text|-bg).
type Marking struct {
	Name     string // "internal", "external"
	Dot      string
	Text     string
	DarkText string
	DarkDot  string
}

// Markings are the named markings: „intern, nie gesendet“ and „von außen“.
var Markings = []Marking{
	{Name: "internal", Dot: "#7C3AED", Text: "#6D28D9", DarkText: "#C4B5FD", DarkDot: "#8B5CF6"},
	{Name: "external", Dot: "#B45309", Text: "#92400E", DarkText: "#FCD34D", DarkDot: "#D97706"},
}

// Backgrounds the tones are checked on: the islands and the grey ground, and in dark mode Tabler's dark
// surfaces. Pills (10 % of the dot) are checked on both, alerts (16 %) on the island, where Tabler mixes them.
var (
	LightBackgrounds = []string{"#FFFFFF", "#F5F5F5"}
	DarkBackgrounds  = []string{"#151F2C", "#1D273B"}
)

// Contrast is the contrast ratio of two colours (#rrggbb) after WCAG 2.1, from 1 to 21. AA asks 4.5 for
// text and 3 for dots, lines and areas (GS-07).
func Contrast(fg, bg string) float64 { return contrastOf(channels(fg), channels(bg)) }

func channels(hex string) [3]float64 {
	var r, g, b int
	fmt.Sscanf(strings.TrimPrefix(hex, "#"), "%02x%02x%02x", &r, &g, &b)
	return [3]float64{float64(r), float64(g), float64(b)}
}

func luminance(c [3]float64) float64 {
	var l [3]float64
	for i, v := range c {
		v /= 255
		if v <= 0.03928 {
			l[i] = v / 12.92
		} else {
			l[i] = math.Pow((v+0.055)/1.055, 2.4)
		}
	}
	return 0.2126*l[0] + 0.7152*l[1] + 0.0722*l[2]
}

func contrastOf(a, b [3]float64) float64 {
	la, lb := luminance(a), luminance(b)
	return (math.Max(la, lb) + 0.05) / (math.Min(la, lb) + 0.05)
}

func rgbTriplet(hex string) string {
	var r, g, b int
	fmt.Sscanf(strings.TrimPrefix(hex, "#"), "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d, %d, %d", r, g, b)
}

// SignalCSS sets the signal colours and makes Tabler's coloured text use their dark tones: pills
// (bg-*-lt), statuses, text-*, alerts, filled badges and buttons. The focus ring of fields takes the
// primary colour in full, so it reaches 3:1 in every theme. On touch screens and below 1200 px fields
// take 16 px. It applies in every theme; load it after Tabler and the theme.
func SignalCSS() string {
	var b strings.Builder
	// The focus ring is a line (3:1): the primary colour in full, not at Tabler's 50 % over the ground.
	b.WriteString(":root{--tblr-focus-ring-opacity:1;--tblr-focus-ring-opacity-dark:1;")
	for _, s := range Signals {
		fmt.Fprintf(&b, "--kit-signal-%s:%s;--kit-signal-%s-text:%s;", s.Name, s.Dot, s.Name, s.Text)
		for _, c := range s.Classes {
			fmt.Fprintf(&b, "--tblr-%s:%s;--tblr-%s-rgb:%s;", c, s.Dot, c, rgbTriplet(s.Dot))
		}
	}
	for _, m := range Markings {
		fmt.Fprintf(&b, "--kit-mark-%s:%s;--kit-mark-%s-text:%s;--kit-mark-%s-bg:color-mix(in oklab,%s 10%%,transparent);", m.Name, m.Dot, m.Name, m.Text, m.Name, m.Dot)
	}
	b.WriteString("}[data-bs-theme=dark]{")
	for _, s := range Signals {
		fmt.Fprintf(&b, "--kit-signal-%s-text:%s;", s.Name, s.DarkText)
		if s.DarkDot != "" {
			fmt.Fprintf(&b, "--kit-signal-%s:%s;", s.Name, s.DarkDot)
			for _, c := range s.Classes {
				fmt.Fprintf(&b, "--tblr-%s:%s;--tblr-%s-rgb:%s;", c, s.DarkDot, c, rgbTriplet(s.DarkDot))
			}
		}
	}
	for _, m := range Markings {
		fmt.Fprintf(&b, "--kit-mark-%s-text:%s;", m.Name, m.DarkText)
		if m.DarkDot != "" {
			fmt.Fprintf(&b, "--kit-mark-%s:%s;", m.Name, m.DarkDot)
		}
	}
	b.WriteString("}")
	for _, s := range Signals {
		t := "var(--kit-signal-" + s.Name + "-text)"
		var text, fill, alert, btn, outline []string
		for _, c := range s.Classes {
			text = append(text, ".bg-"+c+"-lt", ".text-"+c, ".status-"+c, ".badge-outline.text-"+c)
			fill = append(fill, ".badge.bg-"+c)
			alert = append(alert, ".alert-"+c)
			btn = append(btn, ".btn-"+c)
			outline = append(outline, ".btn-outline-"+c, ".btn-ghost-"+c)
		}
		fmt.Fprintf(&b, "%s{color:%s!important}", strings.Join(text, ","), t)
		fmt.Fprintf(&b, "%s{background-color:var(--kit-signal-%s-text)!important;color:#fff!important}", strings.Join(fill, ","), s.Name)
		fmt.Fprintf(&b, "%s{--tblr-alert-color:%s}", strings.Join(alert, ","), t)
		fmt.Fprintf(&b, "%s{--tblr-btn-color:#fff;--tblr-btn-bg:var(--kit-signal-%[2]s-text);--tblr-btn-border-color:var(--kit-signal-%[2]s-text);--tblr-btn-hover-color:#fff;--tblr-btn-hover-bg:color-mix(in oklab,var(--kit-signal-%[2]s-text) 85%%,#000);--tblr-btn-hover-border-color:color-mix(in oklab,var(--kit-signal-%[2]s-text) 85%%,#000);--tblr-btn-active-bg:color-mix(in oklab,var(--kit-signal-%[2]s-text) 75%%,#000)}", strings.Join(btn, ","), s.Name)
		fmt.Fprintf(&b, "%s{--tblr-btn-color:%s;--tblr-btn-hover-color:#fff;--tblr-btn-hover-bg:var(--kit-signal-%s-text);--tblr-btn-hover-border-color:var(--kit-signal-%s-text)}", strings.Join(outline, ","), t, s.Name, s.Name)
	}
	b.WriteString(touchFields)
	return b.String()
}

// touchFields gives fields 16 px on touch screens and below 1200 px: iOS Safari zooms in on a field
// below 16 px when it gets focus, and the page stays zoomed. :is() lifts the selector to (0,1,1), so it
// beats .form-control-sm. Desktop with a mouse keeps Tabler's sizes; zooming stays allowed.
const touchFields = `@media (max-width:1199.98px),(pointer:coarse){` +
	`:is(input:not([type=checkbox],[type=radio],[type=range],[type=button],[type=submit],[type=reset],[type=image],[type=color],[type=hidden]),select,textarea){font-size:max(16px,1rem)}` +
	`.form-control-sm,.form-select-sm{min-height:36px}}`

var signalSheet = csp.SheetFunc("themes/signals", SignalCSS)

// SignalStyles renders SignalCSS once per page (with the nonce of templ.WithNonce). Put it in the head
// after Tabler's stylesheets and the theme.
func SignalStyles() templ.Component { return signalSheet }
