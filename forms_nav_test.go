package tablertempl_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/clipboard"
	"github.com/mind-vm/tabler-templ/form"
	"github.com/mind-vm/tabler-templ/nav"
	"github.com/mind-vm/tabler-templ/sparkline"
)

func TestClipboard(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, clipboard.Button(clipboard.Props{Target: "key", Label: "Copy"}))
	assert.Contains(t, h, `class="btn"`)
	assert.Contains(t, h, `data-bs-toggle="clipboard"`)
	assert.Contains(t, h, `data-bs-target="#key"`)
	assert.Contains(t, h, `<span class="clipboard-label">Copy</span>`)
	assert.Contains(t, h, `<span class="clipboard-feedback">Copied</span>`)

	// Icon only: the icons swap, a hidden word is announced.
	h = render(t, ctx, clipboard.Button(clipboard.Props{Style: clipboard.StyleLink, Text: "x", Delay: 500, FeedbackColor: "green", AriaLabel: "Copy it"}))
	assert.Contains(t, h, `class="link-secondary input-group-link border-0 bg-transparent p-0"`)
	assert.Contains(t, h, `data-bs-text="x"`)
	assert.Contains(t, h, `data-bs-delay="500"`)
	assert.Contains(t, h, `aria-label="Copy it"`)
	assert.Contains(t, h, `ti-clipboard icon clipboard-label`)
	assert.Contains(t, h, `ti-check icon clipboard-feedback text-green`)
	assert.Contains(t, h, `<span class="clipboard-feedback visually-hidden">Copied</span>`)
}

func TestSparklineStyle(t *testing.T) {
	h := render(t, context.Background(), sparkline.Sparkline(sparkline.Props{Values: []float64{1, -1}, NegativeColor: "orange", StrokeWidth: "1", SpotSize: "1.5"}))
	assert.Contains(t, h, `style="--tblr-sparkline-negative: var(--tblr-orange); --tblr-sparkline-stroke-width: 1; --tblr-sparkline-spot-size: 1.5"`)
	assert.NotContains(t, render(t, context.Background(), sparkline.Sparkline(sparkline.Props{Values: []float64{1}})), "style=")
}

func TestNavLinksAndOverflow(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(nav.Nav(nav.Props{Links: true, Vertical: true, Style: nav.StylePills, Attributes: templ.Attributes{"id": "pills"}}), "x"))
	assert.Equal(t, `<nav class="nav nav-pills nav-vertical" id="pills">x</nav>`, h)
	h = render(t, ctx, mit(nav.Link(nav.ItemProps{Href: "#a", Active: true}), "A"))
	assert.Equal(t, `<a href="#a" class="nav-link active" aria-current="page">A</a>`, h)

	h = render(t, ctx, nav.Nav(nav.Props{Overflow: true, NoMoreText: true}))
	assert.Contains(t, h, `<div class="nav-overflow" data-bs-toggle="nav-overflow" data-bs-more-text="false"><ul class="nav">`)
	h = render(t, ctx, mit(nav.Overflow(nav.OverflowProps{Label: "Demo"}), "<ul></ul>"))
	assert.Equal(t, `<nav class="nav-overflow" data-bs-toggle="nav-overflow" aria-label="Demo"><ul></ul></nav>`, h)
}

func TestFieldVariants(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(form.Field(form.FieldProps{Label: "E", For: "e", Required: true, RequiredLabel: " (required)", Horizontal: true}), "<input>"))
	assert.Contains(t, h, `<div class="mb-3 row"><label class="col-3 col-form-label required" for="e">`)
	assert.Contains(t, h, `<span class="visually-hidden"> (required)</span>`)
	assert.Contains(t, h, `<div class="col"><input>`)

	h = render(t, ctx, mit(form.Field(form.FieldProps{Label: "C", Fieldset: true}), "x"))
	assert.Contains(t, h, `<div class="form-label">C</div>`)
	assert.Contains(t, h, `<fieldset><legend class="visually-hidden">C</legend>x</fieldset>`)

	h = render(t, ctx, mit(form.Field(form.FieldProps{Label: "N", For: "n", LabelHidden: true}), "<input>"))
	assert.Contains(t, h, `<label class="form-label visually-hidden" for="n">N</label>`)

	h = render(t, ctx, mit(form.Field(form.FieldProps{Label: "W", Wrap: true}), "<input>"))
	assert.Contains(t, h, `<label class="d-block"><span class="form-label">W</span>`)
}

func TestFormAdditions(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, form.Check(form.CheckProps{Row: true, Switch: true, Label: "Push", Checked: true}))
	assert.Contains(t, h, `<label class="row"><span class="col">Push</span>`)
	assert.Contains(t, h, `<span class="form-check form-switch form-check-single"><input class="form-check-input" type="checkbox" checked>`)

	h = render(t, ctx, form.Check(form.CheckProps{Label: "A", Description: "More"}))
	assert.Contains(t, h, `<span class="form-check-description">More</span>`)

	assert.Equal(t, `<input class="form-check-input m-0" type="radio" aria-label="On">`,
		render(t, ctx, form.CheckInput(form.CheckProps{Class: "m-0", Radio: true, Attributes: templ.Attributes{"aria-label": "On"}})))

	assert.Contains(t, render(t, ctx, form.ColorCheck(form.ColorCheckProps{Color: "red", Round: true})), `class="form-colorinput-color bg-red rounded-circle"`)

	h = render(t, ctx, form.Select(form.SelectProps{Options: []form.Option{{Value: "a", Label: "A", Attributes: templ.Attributes{"data-custom-properties": "<b>"}}}}))
	assert.Contains(t, h, `<option value="a" data-custom-properties="&lt;b&gt;">A</option>`)

	assert.Equal(t, `<datalist id="l"><option value="a"></option></datalist>`, render(t, ctx, form.Datalist("l", []string{"a"})))
	assert.Equal(t, `<div class="valid-feedback" id="f">ok</div>`, render(t, ctx, mit(form.Feedback(form.FeedbackProps{ID: "f", Valid: true}), "ok")))

	h = render(t, ctx, form.Help(form.HelpProps{Content: "<p>x</p>", HTML: true}))
	assert.Contains(t, h, `class="form-help"`)
	assert.Contains(t, h, `data-bs-html="true"`)
	assert.Contains(t, h, `aria-label="Help"`)

	h = render(t, ctx, form.AddonButton(form.AddonButtonProps{Class: "ms-2", Icon: "ti-x", Label: "Clear", Disabled: true}))
	assert.Contains(t, h, `class="link-secondary ms-2 border-0 bg-transparent p-0"`)
	assert.Contains(t, h, `title="Clear"`)
	assert.Contains(t, h, ` disabled`)
}
