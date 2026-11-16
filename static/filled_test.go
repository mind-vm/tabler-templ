package static_test

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/icon"
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
