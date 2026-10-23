package tablertempl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/accordion"
	"github.com/mind-vm/tabler-templ/card"
	"github.com/mind-vm/tabler-templ/carousel"
	"github.com/mind-vm/tabler-templ/chat"
	"github.com/mind-vm/tabler-templ/hovergallery"
	"github.com/mind-vm/tabler-templ/imagelabel"
	"github.com/mind-vm/tabler-templ/legend"
	"github.com/mind-vm/tabler-templ/mention"
	"github.com/mind-vm/tabler-templ/mockup"
	"github.com/mind-vm/tabler-templ/placeholder"
	"github.com/mind-vm/tabler-templ/pricing"
	"github.com/mind-vm/tabler-templ/shape"
	"github.com/mind-vm/tabler-templ/shimmer"
	"github.com/mind-vm/tabler-templ/signal"
	"github.com/mind-vm/tabler-templ/spinner"
	"github.com/mind-vm/tabler-templ/switchicon"
	"github.com/mind-vm/tabler-templ/tag"
	"github.com/mind-vm/tabler-templ/trending"
)

func TestAccordion(t *testing.T) {
	h := render(t, context.Background(), accordion.Accordion(accordion.Props{ID: "faq", Style: accordion.Flush}, []accordion.Item{{Key: "a", Title: "A", Open: true, Icon: "ti-user"}, {Key: "b", Title: "B", Plus: true}}))
	assert.Contains(t, h, `<div class="accordion accordion-flush" id="faq">`)
	assert.Contains(t, h, `<h3 class="accordion-header"><button id="faq-a-button" class="accordion-button" type="button" data-bs-toggle="collapse" data-bs-target="#faq-a" aria-expanded="true" aria-controls="faq-a"><span class="accordion-button-icon">`)
	assert.Contains(t, h, `<div id="faq-a" class="accordion-collapse collapse show" data-bs-parent="#faq" role="region" aria-labelledby="faq-a-button">`)
	assert.Contains(t, render(t, context.Background(), accordion.Accordion(accordion.Props{ID: "x", HeadingLevel: 2}, []accordion.Item{{Key: "a"}})), `<h2 class="accordion-header">`)
	assert.Contains(t, h, `class="accordion-button-toggle accordion-button-toggle-plus"`)
}

func TestSmallComponents(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="spinner-grow text-blue spinner-grow-sm" role="status"></div>`, render(t, ctx, spinner.Spinner(spinner.Props{Grow: true, Color: "blue", Small: true})))
	assert.Contains(t, render(t, ctx, tag.Tag(tag.Props{Text: "Poland", Flag: "pl", Removable: true})), `<button type="button" class="btn-close" aria-label="Remove Poland"></button>`)
	assert.Equal(t, `<span class="signal signal-green" role="img" aria-label="High"><span class="signal-bar active"></span><span class="signal-bar active"></span><span class="signal-bar"></span></span>`,
		render(t, ctx, signal.Signal(signal.Props{Value: 2, Color: "green", Label: "High"})))
	h := render(t, ctx, switchicon.SwitchIcon(switchicon.Props{Variant: switchicon.Fade, Active: true}))
	assert.Contains(t, h, `class="switch-icon switch-icon-fade active" aria-label="Toggle" aria-pressed="true" data-bs-toggle="switch-icon"`)
	assert.Contains(t, h, `class="ti ti-filled icon" data-icon="ti-heart"`)
	assert.Contains(t, render(t, ctx, trending.Trending(trending.Props{Value: -4})), `<span class="text-red d-inline-flex align-items-center lh-1"><span class="visually-hidden">Decreased by </span>-4%`)
	assert.Contains(t, render(t, ctx, mit(shimmer.Shimmer(shimmer.Props{Loader: true}), "Thinking")), `<span class="text-shimmer"><span class="text-shimmer-loader" aria-hidden="true"></span>Thinking</span>`)
	assert.Contains(t, render(t, ctx, legend.Legend(legend.Props{Text: "Web", Color: "blue", Value: "42", Unit: "%"})), `<span class="legend legend-blue"><span class="legend-dot"></span> Web<span class="legend-value">42 <span class="legend-unit">%</span></span></span>`)
	assert.Contains(t, render(t, ctx, mit(mention.Mention(mention.Props{Href: "#", Avatar: "/a.jpg", At: true}), "Ann")), `<span class="visually-hidden">@</span>Ann</a>`)
	assert.Contains(t, render(t, ctx, shape.Shape("ti-star", shape.Props{Size: "lg", Color: "blue"})), `<span class="shape shape-lg shape-blue">`)
}

func TestMediaComponents(t *testing.T) {
	ctx := context.Background()
	assert.Contains(t, render(t, ctx, imagelabel.Label(imagelabel.LabelProps{Corner: imagelabel.TopEnd, Icon: "ti-crown", IconOnly: true, Text: "Premium"})), `<span class="img-label img-label-top-end img-label-icon">`)
	assert.Contains(t, render(t, ctx, hovergallery.Gallery(hovergallery.Props{Indicators: true, AspectRatio: "1 / 1"}, []hovergallery.Image{{Src: "/a.jpg"}})),
		`<div class="hover-gallery hover-gallery-indicators" style="--tblr-hover-gallery-aspect-ratio: 1 / 1;"><div class="hover-gallery-item"><img src="/a.jpg" alt=""></div></div>`)
	assert.Contains(t, render(t, ctx, mockup.Browser(mockup.Props{URL: "example.com", Secure: true})), `<div class="mockup-browser-url"><i class="ti ti-lock icon"`)
	h := render(t, ctx, carousel.Carousel(carousel.Props{ID: "c", Indicators: carousel.Thumbs, Fade: true}, []carousel.Slide{{Src: "/1.jpg", Title: "One"}, {Src: "/2.jpg"}}))
	assert.Contains(t, h, `class="carousel slide carousel-fade"`)
	assert.Contains(t, h, `class="ratio ratio-4x3 active" style="background-image: url(/1.jpg);"`)
	assert.Contains(t, h, `<div class="carousel-caption d-none d-md-block"><h3>One</h3></div>`)
	assert.Equal(t, `<div class="placeholder col-9 placeholder-lg"></div>`, render(t, ctx, placeholder.Placeholder(placeholder.Props{Class: "col-9", Size: "lg"})))
	h = render(t, ctx, pricing.Card(pricing.CardProps{Title: "Pro", Currency: "$", Price: "49", Featured: true, Features: []pricing.Feature{{Text: "Support", Excluded: true}}}))
	assert.Contains(t, h, `<div class="pricing-card featured"><h3 class="pricing-title">Pro</h3>`)
	assert.Contains(t, h, `ti-x icon icon-inline me-1 text-red`)
	assert.Equal(t, `<div class="card card-gradient card-gradient-rainbow card-gradient-start card-gradient-animated">x</div>`,
		render(t, ctx, mit(card.Card(card.Props{Gradient: "rainbow", GradientEdge: "start", GradientAnimated: true}), "x")))
	assert.Contains(t, render(t, ctx, mit(chat.Bubble(chat.MessageProps{Author: "You", Time: "10:26", Me: true}), "<p>Hi</p>")), `<div class="chat-bubble chat-bubble-me">`)
}
