package gallery_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/gallery"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestPhoto(t *testing.T) {
	h := render(t, gallery.Photo(gallery.PhotoProps{Src: "/a.jpg", Alt: "A", Avatar: "/p.jpg", Name: "Ann", Time: "now", Views: "467", Likes: "67", Liked: true}))
	assert.True(t, strings.HasPrefix(h, `<div class="card card-sm"><a href="#" class="d-block"><img src="/a.jpg" alt="A" class="card-img-top"></a>`))
	assert.Contains(t, h, `<span class="avatar me-3 rounded" style="background-image: url(/p.jpg);"></span><div><div>Ann</div><div class="text-secondary">now</div></div>`)
	assert.Contains(t, h, `ti-eye icon icon-inline" aria-hidden="true"></i> 467</a>`)
	assert.Contains(t, h, `text-red`)

	h = render(t, gallery.Photo(gallery.PhotoProps{Src: "/a.jpg", Initials: "SA", Name: "Sunny"}))
	assert.Contains(t, h, `>SA</span>`)
	assert.NotContains(t, h, `ti-eye`)
}

func TestTile(t *testing.T) {
	h := render(t, gallery.Tile(gallery.TileProps{Src: "/a.jpg", Title: "A", Ratio: "3x1", Class: "rounded-3 border", Lightbox: "gallery"}))
	assert.Equal(t, `<a data-fslightbox="gallery" href="/a.jpg" aria-label="A"><div class="img-responsive img-responsive-3x1 rounded-3 border" style="background-image: url(/a.jpg);"></div></a>`, h)
}
