package faq_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/faq"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestFAQ(t *testing.T) {
	h := render(t, faq.FAQ(faq.Props{}, []faq.Section{
		{Title: "1. Introduction", Questions: []faq.Question{{Q: "Who are we?", Answer: templ.Raw("<p>Us.</p>")}, {Q: "Why?"}}},
		{Title: "2. Payments", Questions: []faq.Question{{Q: "Refund?"}}},
	}))
	assert.True(t, strings.HasPrefix(h, `<div class="card card-lg"><div class="card-body"><div class="space-y-4"><div><h2 class="mb-3">1. Introduction</h2><div class="accordion accordion-tabs" id="faq-1">`))
	assert.Contains(t, h, `<div class="accordion accordion-tabs" id="faq-2">`)
	// The first question of each section is open.
	assert.Contains(t, h, `aria-expanded="true" aria-controls="faq-1-1">Who are we?`)
	assert.Contains(t, h, `aria-expanded="false" aria-controls="faq-1-2">Why?`)
	assert.Contains(t, h, `<p>Us.</p>`)
}
