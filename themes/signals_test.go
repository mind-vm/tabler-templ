package themes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests of GS-07 of the design concept „Gestaltung“: contrast after WCAG 2.1 AA in every pair the signal
// colours and markings are used in.

// mix is share of fg over bg, as the soft fills of pills (10 %) and alerts (16 %) lay it.
func mix(fg, bg string, share float64) [3]float64 {
	f, b := channels(fg), channels(bg)
	return [3]float64{f[0]*share + b[0]*(1-share), f[1]*share + b[1]*(1-share), f[2]*share + b[2]*(1-share)}
}

type pair struct {
	what     string
	fg, bg   [3]float64
	min      float64
	fgS, bgS string
}

func pairs() []pair {
	var ps []pair
	add := func(what, fg, bg string, min float64) {
		ps = append(ps, pair{what, channels(fg), channels(bg), min, fg, bg})
	}
	addMix := func(what, fg, dot, bg string, share, min float64) {
		ps = append(ps, pair{what, channels(fg), mix(dot, bg, share), min, fg, fmt.Sprintf("%.0f%% %s on %s", share*100, dot, bg)})
	}
	type tone struct{ name, dot, text, dark, darkDot string }
	or := func(a, b string) string {
		if a != "" {
			return a
		}
		return b
	}
	var tones []tone
	for _, s := range Signals {
		tones = append(tones, tone{s.Name, s.Dot, s.Text, s.DarkText, or(s.DarkDot, s.Dot)})
	}
	for _, m := range Markings {
		tones = append(tones, tone{m.Name, m.Dot, m.Text, m.DarkText, or(m.DarkDot, m.Dot)})
	}
	for _, t := range tones {
		for _, bg := range LightBackgrounds {
			add(t.name+" text", t.text, bg, 4.5)
			add(t.name+" dot", t.dot, bg, 3)
			addMix(t.name+" text in pill", t.text, t.dot, bg, 0.10, 4.5)
		}
		addMix(t.name+" text in alert", t.text, t.dot, LightBackgrounds[0], 0.16, 4.5)
		add(t.name+" white on filled", "#FFFFFF", t.text, 4.5)
		for _, bg := range DarkBackgrounds {
			add(t.name+" dark text", t.dark, bg, 4.5)
			add(t.name+" dot on dark", t.darkDot, bg, 3)
			addMix(t.name+" dark text in pill", t.dark, t.darkDot, bg, 0.10, 4.5)
		}
	}
	return ps
}

func TestContrast(t *testing.T) {
	assert.InDelta(t, 21, Contrast("#000000", "#FFFFFF"), 0.01)
	assert.InDelta(t, 2.27, Contrast("#E0A100", "#FFFFFF"), 0.01)
}

func TestGS07_ContrastAA(t *testing.T) {
	for _, p := range pairs() {
		c := contrastOf(p.fg, p.bg)
		assert.GreaterOrEqual(t, c, p.min, "%s: %s on %s is %.2f:1", p.what, p.fgS, p.bgS, c)
	}
}

// The signal colour itself is never a text colour: Tabler's coloured text takes the dark tone.
func TestGS07_SignalNeverText(t *testing.T) {
	css := SignalCSS()
	for _, s := range Signals {
		for _, c := range s.Classes {
			for _, sel := range []string{".bg-" + c + "-lt", ".text-" + c, ".status-" + c} {
				assert.Contains(t, css, sel, s.Name)
			}
			assert.Contains(t, css, ".alert-"+c)
			assert.Contains(t, css, "--tblr-"+c+":"+s.Dot)
		}
		assert.Contains(t, css, "--kit-signal-"+s.Name+"-text)!important")
	}
	assert.Contains(t, css, "[data-bs-theme=dark]{--kit-signal-green-text:#4ADE80")
}

// The values of the concept (Gestaltung v0.5, GS-07), so a change there is a change here on purpose.
func TestGS07_ValuesOfTheConcept(t *testing.T) {
	want := map[string][2]string{"green": {"#22A04B", "#1E7537"}, "yellow": {"#B57F00", "#8A5A00"}, "red": {"#E5484D", "#B42318"}}
	for _, s := range Signals {
		if w, ok := want[s.Name]; ok {
			assert.Equal(t, w, [2]string{s.Dot, s.Text}, s.Name)
		}
	}
}

func hexOf(rgb string) string {
	var r, g, b int
	fmt.Sscanf(rgb, "%d, %d, %d", &r, &g, &b)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// GS-07: links are text (4.5:1) and the focus ring of fields a line (3:1), in every theme, in light and dark
// mode. Both take the theme's Accent, else its primary colour (SignalCSS draws the ring in full).
func TestGS07_LinksAndFocusRing(t *testing.T) {
	assert.Contains(t, SignalCSS(), "--tblr-focus-ring-opacity:1;--tblr-focus-ring-opacity-dark:1")
	for _, th := range All {
		if th.Primary == "" {
			continue // a Tabler colour (Settings.Primary), checked by Tabler
		}
		light := th.Primary
		if th.Accent != "" {
			light = th.Accent
		}
		l := hexOf(rgbOf(light))
		// Dark mode: AccentDark, else Tabler's dark link, the dark primary colour with 40 % white.
		var d string
		if th.AccentDark != "" {
			d = hexOf(rgbOf(th.AccentDark))
		} else {
			pd := th.PrimaryDark
			if pd == "" {
				pd = th.Primary
			}
			c := mix("#FFFFFF", hexOf(rgbOf(pd)), 0.4)
			d = fmt.Sprintf("#%02X%02X%02X", int(c[0]+0.5), int(c[1]+0.5), int(c[2]+0.5))
		}
		if th.Mode != "dark" {
			for _, bg := range LightBackgrounds {
				assert.GreaterOrEqual(t, Contrast(l, bg), 4.5, "%s: link %s on %s", th.Label, l, bg)
			}
		}
		for _, bg := range DarkBackgrounds {
			assert.GreaterOrEqual(t, Contrast(d, bg), 4.5, "%s dark: link %s on %s", th.Label, d, bg)
		}
	}
}

// Fields take 16 px on touch screens and below 1200 px, so iOS does not zoom in on focus (issue #9); a
// desktop with a mouse keeps Tabler's sizes, and zooming stays allowed.
func TestTouchFields16px(t *testing.T) {
	css := SignalCSS()
	assert.Contains(t, css, "@media (max-width:1199.98px),(pointer:coarse){:is(input:not([type=checkbox],")
	assert.Contains(t, css, ",select,textarea){font-size:max(16px,1rem)}")
	assert.Contains(t, css, ".form-control-sm,.form-select-sm{min-height:36px}}")
	assert.NotContains(t, css, "maximum-scale")
}
