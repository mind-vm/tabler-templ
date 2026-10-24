package users_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/users"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestCard(t *testing.T) {
	h := render(t, users.Card(users.CardProps{Avatar: "/a.jpg", Name: "Ann", Href: "/u/ann", Detail: "UI Designer", Badge: templ.Raw(`<span class="badge">Owner</span>`),
		Actions: []users.Action{{Href: "mailto:a@x.org", Icon: "ti-mail", Text: "Email"}}}))
	assert.True(t, strings.HasPrefix(h, `<div class="card"><div class="card-body p-4 text-center"><span class="avatar avatar-xl mb-3" style="background-image: url(/a.jpg);"></span><h2 class="m-0 mb-1"><a href="/u/ann">Ann</a></h2><div class="text-secondary">UI Designer</div><div class="mt-3"><span class="badge">Owner</span></div></div>`))
	assert.Contains(t, h, `<div class="d-flex"><a href="mailto:a@x.org" class="card-btn"><i class="ti ti-mail icon me-2 text-muted" aria-hidden="true"></i> Email</a></div>`)
}
