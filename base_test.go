package tablertempl_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/colorswatch"
	"github.com/mind-vm/tabler-templ/pattern"
	"github.com/mind-vm/tabler-templ/prose"
	"github.com/mind-vm/tabler-templ/social"
)

func TestProse(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="prose text-secondary"><p>Hi</p></div>`, render(t, ctx, mit(prose.Prose(prose.Props{Class: "text-secondary"}), "<p>Hi</p>")))
	assert.Equal(t, `<div class="prose"><h1>A</h1></div>`, render(t, ctx, prose.HTML("<h1>A</h1>")))
	assert.Equal(t, `<div class="callout">Tip</div>`, render(t, ctx, mit(prose.Callout(), "Tip")))
}

func TestPattern(t *testing.T) {
	p := pattern.Props{Type: pattern.Grid, Color: "primary", Size: "lg", Opacity: "light", Class: "border"}
	assert.Equal(t, "bg-pattern-grid bg-pattern-primary bg-pattern-lg border bg-pattern-opacity-light", pattern.Classes(p))
	assert.Equal(t, `<div class="bg-pattern-dots"></div>`, render(t, context.Background(), pattern.Pattern(pattern.Props{Type: pattern.Dots})))
}

func TestSocial(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<span class="social social-md social-app-github"></span>`, render(t, ctx, social.Social("github", social.Props{Size: "md"})))
	assert.Equal(t, `<span class="icon social social-gray social-app-x" role="img" aria-label="X"></span>`,
		render(t, ctx, social.Social("x", social.Props{Class: "icon", Gray: true, Label: "X"})))
}

func TestColorSwatch(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, colorswatch.Swatch(colorswatch.Props{Color: "blue", Name: "Blue", Value: "#066fd1", Initials: "bl"}))
	assert.Contains(t, h, `<span class="avatar bg-blue text-blue-fg avatar-square">bl</span>`)
	assert.Contains(t, h, `Blue<br><code>#066fd1</code>`)
	h = render(t, ctx, colorswatch.Swatch(colorswatch.Props{Color: "cyan", Light: true, Name: "Cyan Light", Value: "#17a2b8", Initials: "cy"}))
	assert.Contains(t, h, `class="avatar bg-cyan-lt avatar-square"`)
	assert.Equal(t, "bg-gradient bg-gradient-from-blue bg-gradient-to-transparent bg-gradient-via-white bg-gradient-to-ts",
		colorswatch.GradientClasses(colorswatch.GradientProps{From: "blue", Via: "white", To: "transparent", Direction: colorswatch.ToTopStart}))
	assert.Equal(t, `<div class="border bg-gradient bg-gradient-from-red" id="g"></div>`,
		render(t, ctx, colorswatch.Gradient(colorswatch.GradientProps{From: "red", Class: "border", Attributes: templ.Attributes{"id": "g"}})))
}
