package marketing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/marketing"
)

func render(t *testing.T, ctx context.Context, cs ...templ.Component) string {
	t.Helper()
	var b strings.Builder
	for _, c := range cs {
		require.NoError(t, c.Render(ctx, &b))
	}
	return b.String()
}

func TestLandingStylesOncePerPage(t *testing.T) {
	ctx := templ.InitializeContext(context.Background())
	h := render(t, ctx,
		marketing.Steps(marketing.HeaderProps{}, []marketing.Step{{Title: "One"}}),
		marketing.Bento(marketing.HeaderProps{}, []marketing.Tile{{Title: "A", Span: 4, Color: "dark"}, {Title: "B", Span: 5}}))
	assert.Equal(t, 1, strings.Count(h, "<style"))
	assert.Contains(t, h, `<ol class="mkt-steps">`)
	assert.Contains(t, h, "mkt-bento-span-4")
	assert.Contains(t, h, "bg-dark text-dark-fg")
	assert.Contains(t, h, "mkt-bento-span-2", "an unsupported span falls back to 2")
}

func TestPlans(t *testing.T) {
	plans := []marketing.Plan{
		{Name: "Free", Monthly: marketing.PlanPrice{Amount: "$0"}},
		{Name: "Pro", Featured: true, Monthly: marketing.PlanPrice{Amount: "$10", Period: "/mo"}, Yearly: marketing.PlanPrice{Amount: "$8", Period: "/mo"}},
	}
	h := render(t, context.Background(), marketing.Plans(marketing.PlansProps{ID: "p", Save: "-20%"}, plans))
	assert.Contains(t, h, `id="p-yearly"`)
	assert.Contains(t, h, "-20%")
	assert.Contains(t, h, "Most popular", "a featured plan gets a default badge")
	assert.Equal(t, 2, strings.Count(h, "$0"), "without yearly prices the monthly ones show in both")
	assert.Contains(t, h, "$8")

	h = render(t, context.Background(), marketing.Plans(marketing.PlansProps{NoSwitch: true}, plans))
	assert.NotContains(t, h, "btn-check")
}

func TestHeroAnnouncement(t *testing.T) {
	h := render(t, context.Background(), marketing.Hero(marketing.HeroProps{Title: "T",
		Announcement: marketing.Announcement(marketing.AnnouncementProps{Badge: "New", Text: "v2", Href: "/news"})}))
	assert.Contains(t, h, `<a href="/news" class="mkt-announcement">`)
	assert.Less(t, strings.Index(h, "mkt-announcement"), strings.Index(h, "hero-title"))
}

func TestCTABoxAndQuote(t *testing.T) {
	h := render(t, context.Background(), marketing.CTABox(marketing.CTABoxProps{Title: "Go", Color: "dark"}))
	assert.Contains(t, h, "mkt-cta-box bg-dark text-dark-fg")
	assert.Contains(t, h, "justify-content-center text-center", "without media the text is centered")

	h = render(t, context.Background(), marketing.Quote(marketing.QuoteProps{Text: "Great", Name: "Ann", Stars: 3}))
	assert.Equal(t, 3, strings.Count(h, "ti-filled"))
	assert.Contains(t, h, `aria-label="3 of 5 stars"`)
}

func TestMarquee(t *testing.T) {
	items := []templ.Component{templ.Raw("<b>A</b>"), templ.Raw("<b>B</b>")}
	h := render(t, context.Background(), marketing.Marquee(marketing.MarqueeProps{Label: "Customers", Speed: "slow", Reverse: true}, items))
	assert.Contains(t, h, `class="mkt-marquee mkt-marquee-slow mkt-marquee-reverse"`)
	assert.Equal(t, 2, strings.Count(h, "<b>A</b>"), "the items are rendered twice for the loop")
	assert.Contains(t, h, `aria-label="Customers"`)
	assert.Contains(t, h, `aria-hidden="true" inert`, "the copy is hidden from assistive technology")
}

func TestCTABoxInline(t *testing.T) {
	h := render(t, context.Background(), marketing.CTABox(marketing.CTABoxProps{Title: "Go", Color: "light", Inline: true,
		Actions: []templ.Component{templ.Raw(`<a class="btn">Buy</a>`)}}))
	assert.Contains(t, h, `class="mkt-cta-inline"`)
	assert.NotContains(t, h, "justify-content-center")
	assert.Contains(t, h, "Buy")
}

func TestPageTheme(t *testing.T) {
	h := render(t, context.Background(), marketing.Page(marketing.PageProps{Title: "T", Theme: "dark",
		HTMLAttributes: templ.Attributes{"data-bs-theme-pack": "graphite"}}))
	assert.Contains(t, h, `data-bs-theme="dark"`)
	assert.Contains(t, h, `data-bs-theme-pack="graphite"`)
	assert.Contains(t, h, `lang="en"`)
}
