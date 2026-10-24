package profile_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/profile"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestHeader(t *testing.T) {
	h := render(t, profile.Header(profile.HeaderProps{Avatar: "/a.jpg", Name: "Dunn", Bio: "Bio.",
		Facts: []profile.Fact{{Icon: "ti-map", Text: "Harbin"}, {Icon: "ti-mail", Text: "d@x.org", Href: "mailto:d@x.org"}}}))
	assert.True(t, strings.HasPrefix(h, `<div class="page-header"><div class="container"><div class="row align-items-center"><div class="col-auto"><span class="avatar avatar-lg" style="background-image: url(/a.jpg);"></span>`))
	assert.Contains(t, h, `<h1 class="fw-bold m-0">Dunn</h1><div class="my-2">Bio.</div>`)
	assert.Contains(t, h, `Harbin</div> <div class="list-inline-item">`)
	assert.Contains(t, h, `<a href="mailto:d@x.org" class="text-reset">d@x.org</a>`)
	assert.NotContains(t, h, `ms-auto`)
}

func TestInfo(t *testing.T) {
	h := render(t, profile.Info("", []profile.Item{{Icon: "ti-book", Label: "Went to:", Value: "Ljubljana"}}))
	assert.Contains(t, h, `<h2 class="card-title">Basic info</h2>`)
	assert.Contains(t, h, `<span>Went to: <strong>Ljubljana</strong></span>`)
}
