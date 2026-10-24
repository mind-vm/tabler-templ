package cdn_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/assets"
	"github.com/mind-vm/tabler-templ/blocks/marketing"
	"github.com/mind-vm/tabler-templ/calendar"
	"github.com/mind-vm/tabler-templ/cdn"
	"github.com/mind-vm/tabler-templ/chart"
	"github.com/mind-vm/tabler-templ/colorpicker"
	"github.com/mind-vm/tabler-templ/datatable"
	"github.com/mind-vm/tabler-templ/datepicker"
	"github.com/mind-vm/tabler-templ/dropzone"
	"github.com/mind-vm/tabler-templ/icon"
	"github.com/mind-vm/tabler-templ/inputmask"
	"github.com/mind-vm/tabler-templ/lightbox"
	"github.com/mind-vm/tabler-templ/mapview"
	"github.com/mind-vm/tabler-templ/player"
	"github.com/mind-vm/tabler-templ/rangeslider"
	"github.com/mind-vm/tabler-templ/rating"
	"github.com/mind-vm/tabler-templ/selectbox"
	"github.com/mind-vm/tabler-templ/signature"
	"github.com/mind-vm/tabler-templ/sortable"
	"github.com/mind-vm/tabler-templ/themes"
	"github.com/mind-vm/tabler-templ/tour"
	"github.com/mind-vm/tabler-templ/vectormap"
	"github.com/mind-vm/tabler-templ/wysiwyg"
)

// TestSelfHosted renders everything that loads files in the browser with a local cdn.Base: nothing may
// point outside, so an app with its own copy (cmd/vendor) loads nothing from a CDN.
func TestSelfHosted(t *testing.T) {
	cdn.Base = "/static/npm/"
	defer func() { cdn.Base = cdn.JSDelivr }()

	page := templ.Join(
		assets.Head(assets.Props{Flags: true, Payments: true, Socials: true, Marketing: true, Vendors: true, Themes: true, Charts: true}),
		assets.Theme(), assets.Scripts(), themes.Head(), icon.Styles(),
		chart.Library(), calendar.Library(), colorpicker.Library(), datatable.Library(), datepicker.Library(),
		dropzone.Library(), inputmask.Library(), lightbox.Library(), mapview.Library(), player.Library(),
		rangeslider.Library(), rating.Library(), selectbox.Library(), signature.Library(), sortable.Library(),
		tour.Library(), vectormap.Library("world", "europe"), wysiwyg.Library(), marketing.TypedLibrary(),
		player.Player(player.Props{ID: "v", Src: "/movie.mp4"}),
		mapview.Map(mapview.Props{ID: "m"}),
	)
	var b bytes.Buffer
	if err := page.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	html := b.String()
	for _, bad := range []string{"jsdelivr", "unpkg", "cdnjs", "cdn.plyr.io", "fonts.googleapis"} {
		if strings.Contains(html, bad) {
			i := strings.Index(html, bad)
			t.Errorf("%s in %s", bad, html[max(0, i-120):min(len(html), i+60)])
		}
	}
	if n := strings.Count(html, "/static/npm/"); n < 30 {
		t.Errorf("only %d local URLs", n)
	}
	if marketing.MarketingCSS("1.7.1") != "/static/npm/@tabler/core@1.7.1/dist/css/tabler-marketing.min.css" {
		t.Error(marketing.MarketingCSS("1.7.1"))
	}
}

func TestRegistered(t *testing.T) {
	groups := cdn.Groups()
	for _, want := range []string{"tabler", "chart", "flags", "wysiwyg", "theme:graphite", "mapview", "player"} {
		found := false
		for _, g := range groups {
			found = found || g == want
		}
		if !found {
			t.Errorf("no group %s in %v", want, groups)
		}
	}
	for _, f := range cdn.Registered() {
		if f.Package == "" || f.Version == "" || len(f.Paths) == 0 {
			t.Errorf("%+v", f)
		}
	}
}
