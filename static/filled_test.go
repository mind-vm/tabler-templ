package static_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/icon"
	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/static"
)

// The filled icons' webfont comes from the bundle when it has one, and the app's own address wins.
func TestWithFilledIcons(t *testing.T) {
	defer func(old string) { icon.FilledFontURL = old }(icon.FilledFontURL)
	icon.FilledFontURL = ""
	static.New("/s/", fstest.MapFS{"a.css": {Data: []byte("x")}}, static.WithFilledIcons())
	assert.Empty(t, icon.FilledFontURL, "not in the bundle: unchanged")
	f := static.New("/s/", fstest.MapFS{static.FilledIcons: {Data: []byte("woff2")}}, static.WithFilledIcons())
	assert.Equal(t, f.URL(static.FilledIcons), icon.FilledFontURL)
	icon.FilledFontURL = "/own.woff2"
	static.New("/s/", fstest.MapFS{static.FilledIcons: {Data: []byte("woff2")}}, static.WithFilledIcons())
	assert.Equal(t, "/own.woff2", icon.FilledFontURL)
}

// The order of the options does not matter: with WithFilledIcons first, kit.css is still in the bundle and
// its address carries the version (found by Vergabe).
func TestOptionOrder(t *testing.T) {
	defer func(old string) { icon.FilledFontURL = old }(icon.FilledFontURL)
	defer csp.SetExternal(false)
	fsys := fstest.MapFS{static.FilledIcons: {Data: []byte("woff2")}}
	for _, opts := range [][]static.Option{
		{static.WithFilledIcons(), static.WithKitCSS()},
		{static.WithKitCSS(), static.WithFilledIcons()},
	} {
		icon.FilledFontURL = ""
		f := static.New("/s/", fsys, opts...)
		assert.True(t, f.Has(static.KitCSS), "kit.css in the bundle")
		assert.Equal(t, f.URL(static.FilledIcons), icon.FilledFontURL)
		assert.NotEqual(t, "/s/"+static.KitCSS, f.URL(static.KitCSS), "versioned address")
	}
}
