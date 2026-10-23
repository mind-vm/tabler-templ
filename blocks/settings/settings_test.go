package settings_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/plans"
	"github.com/mind-vm/tabler-templ/blocks/settings"
)

func TestSettingsAndPlans(t *testing.T) {
	var b strings.Builder
	require.NoError(t, settings.Frame(settings.FrameProps{Groups: []settings.Group{{Title: "Account", Links: []settings.Link{{Text: "Me", Href: "/me", Active: true}}}}}).Render(context.Background(), &b))
	h := b.String()
	assert.Contains(t, h, `<h2 class="subheader">Account</h2><nav aria-label="Account" class="list-group list-group-transparent">`)
	assert.Contains(t, h, `class="list-group-item list-group-item-action active d-flex align-items-center" aria-current="true"`)
	b.Reset()
	require.NoError(t, plans.Plan(plans.PlanProps{Name: "Pro", Price: "$9", Featured: "green", Features: []plans.Feature{{Text: "API", Missing: true}}}).Render(context.Background(), &b))
	h = b.String()
	assert.Contains(t, h, `class="ribbon ribbon-top ribbon-bookmark bg-green"`)
	assert.Contains(t, h, `ti-x icon text-danger`)
	assert.Contains(t, h, `class="btn btn-green w-100"`)
}
