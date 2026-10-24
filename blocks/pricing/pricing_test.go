package pricing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/pricing"
)

func TestTable(t *testing.T) {
	var b strings.Builder
	require.NoError(t, pricing.Table(pricing.TableProps{Title: "Plans"}, []pricing.Plan{{Name: "A"}, {Name: "B", Featured: "green"}},
		[]pricing.Group{{Title: "Features", Rows: []pricing.Row{{Feature: "X", Included: []bool{false, true}}}}}).Render(context.Background(), &b))
	h := b.String()
	assert.Contains(t, h, `<th colspan="3" class="subheader">Features</th>`)
	assert.Contains(t, h, `ti-x icon text-red`)
	assert.Contains(t, h, `<a href="#" class="btn btn-green w-100">Choose plan</a>`)
}
