package payment_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/payment"
)

func TestPayment(t *testing.T) {
	var b strings.Builder
	require.NoError(t, payment.Payment("visa", payment.Props{Size: payment.SizeXS, Dark: true, Label: "Visa"}).Render(context.Background(), &b))
	assert.Equal(t, `<span class="payment payment-xs payment-provider-visa-dark" role="img" aria-label="Visa"></span>`, b.String())
	assert.Equal(t, payment.Provider{Code: "2c2p", Name: "2c2p"}, payment.Providers[0])
	assert.Len(t, payment.Providers, 169)
}
