package legal_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/legal"
)

func TestLicenseSummary(t *testing.T) {
	var b strings.Builder
	require.NoError(t, legal.LicenseSummary(legal.LicenseProps{Name: "MIT License", Permissions: []legal.Term{{Text: "Commercial use"}}}).Render(context.Background(), &b))
	h := b.String()
	assert.Contains(t, h, `<h3 class="lh-1">MIT License</h3>`)
	assert.Contains(t, h, `<ul class="list-icon" style="--tblr-list-icon-color: var(--tblr-green);">`)
	assert.NotContains(t, h, "Limitations")
}
