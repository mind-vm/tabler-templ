package invoice_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/invoice"
)

func TestInvoice(t *testing.T) {
	var b strings.Builder
	require.NoError(t, invoice.Invoice(invoice.Props{From: invoice.Party{Title: "Us", Lines: []string{"A", "B"}}, Title: "Invoice 1",
		Lines:  []invoice.Line{{Product: "P", Quantity: "2", Unit: "1", Amount: "2"}},
		Totals: []invoice.Total{{Label: "Total", Value: "2", Final: true}}}).Render(context.Background(), &b))
	h := b.String()
	assert.Contains(t, h, `<address> A<br> B</address>`)
	assert.Contains(t, h, `<td class="text-center">1</td>`)
	assert.Contains(t, h, `<td colspan="4" class="fw-semibold text-uppercase text-end">Total</td>`)
}
