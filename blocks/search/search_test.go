package search_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/search"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestFilters(t *testing.T) {
	h := render(t, search.Filters(search.FiltersProps{}))
	assert.Contains(t, h, `<form action="./" method="get" autocomplete="off" novalidate>`)

	h = render(t, search.Categories("Category", []search.Category{{Name: "Games", Count: "24", Active: true}, {Name: "Toys", Count: "54"}}))
	assert.Contains(t, h, `<div class="subheader mb-2">Category</div><nav class="list-group list-group-transparent mb-3">`)
	assert.Contains(t, h, `active d-flex align-items-center" aria-current="true">Games<small class="text-secondary ms-auto">24</small></a>`)

	h = render(t, search.Choices(search.ChoicesProps{Title: "Rating", Name: "stars", Radio: true, Options: []search.Option{{Label: "5 stars", Checked: true}, {Label: "4 stars"}}}))
	assert.Contains(t, h, `<fieldset class="border-0 p-0 mb-3"><legend class="subheader mb-2 float-none">Rating</legend>`)
	assert.Contains(t, h, `type="radio"`)
	assert.Equal(t, 1, strings.Count(h, "checked"))

	h = render(t, search.Range(search.RangeProps{Title: "Price", Prefix: "$", From: "3"}))
	assert.Contains(t, h, `value="3"`)
	assert.Contains(t, h, `—`)

	h = render(t, search.Select(search.SelectProps{Title: "Shipping", Name: "country", Options: []string{"UK", "USA"}, Selected: "USA"}))
	assert.Contains(t, h, `name="country"`)
	assert.Contains(t, h, `selected`)

	h = render(t, search.Buttons(search.ButtonsProps{ResetHref: "#"}))
	assert.Contains(t, h, `Confirm changes`)
	assert.Contains(t, h, `<a href="#"`)
}
