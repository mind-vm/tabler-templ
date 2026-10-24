package billing_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/billing"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestBilling(t *testing.T) {
	h := render(t, billing.Plan(billing.PlanProps{Name: "Pro", Status: "Active", Price: "$49", Period: "/month", Usage: &billing.Usage{Label: "Requests", Value: 74}}))
	assert.Contains(t, h, `<span class="badge bg-success-lt ms-1">Active</span>`)
	assert.Contains(t, h, `style="width: 74%;"`)
	h = render(t, billing.Invoices(billing.CardProps{Title: "Invoices"}, []billing.Invoice{{Number: "INV-1", Status: "Open", StatusColor: "warning", PayHref: "#"}}))
	assert.Contains(t, h, `<button type="button" class="btn btn-primary btn-sm">Pay now</button>`)
	assert.Contains(t, h, `aria-label="Download invoice INV-1"`)
	h = render(t, billing.Methods(billing.MethodsProps{}, []billing.Method{{ID: "m1", Provider: "visa", Title: "Visa", Default: true, Checked: true}}))
	assert.Contains(t, h, `<input type="radio" id="m1" name="payment-method" value="" class="form-selectgroup-input" checked>`)
	assert.Contains(t, h, `payment payment-sm payment-provider-visa me-3`)
	h = render(t, billing.CardForm(billing.CardFormProps{Action: "/pay", Submit: "Pay"}))
	assert.Contains(t, h, `autocomplete="cc-number"`)
	assert.Contains(t, h, `<button type="submit" class="btn btn-primary w-100">Pay</button>`)
}
