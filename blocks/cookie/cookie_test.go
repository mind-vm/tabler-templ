package cookie_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/cookie"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestBanner(t *testing.T) {
	h := render(t, cookie.Banner(cookie.BannerProps{Show: true, Text: templ.Raw("Cookies?"), LearnHref: "/terms", Buttons: []cookie.Button{{Text: "Allow"}}}))
	assert.Contains(t, h, `<div class="offcanvas offcanvas-bottom show h-auto" tabindex="-1" role="dialog" id="cookie-banner" aria-modal="true">`)
	assert.Contains(t, h, `<div class="col">Cookies? `)
	assert.Contains(t, h, ` <a href="/terms" target="_blank">Learn more</a></div>`)
	assert.Contains(t, h, `<button type="button" class="btn btn-primary" data-bs-dismiss="offcanvas">Allow</button>`)
}
