package tables_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/tables"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestRanking(t *testing.T) {
	h := render(t, tables.Ranking(tables.RankingProps{Title: "Top Pages", Label: "Page", ValueLabel: "Visitors"}, []tables.RankingRow{{Text: "/", Percent: 82, Value: "4896"}}))
	assert.Contains(t, h, `<table class="table table-borderless table-sm">`)
	assert.Contains(t, h, `<div class="progressbg-text">/</div>`)
	assert.Contains(t, h, `<td class="w-1 fw-medium text-end">4896</td>`)
}

func TestInvoices(t *testing.T) {
	h := render(t, tables.Invoices(tables.InvoicesProps{Count: 8, Total: 16, Page: 3, Pages: 5},
		[]tables.Invoice{{No: "001401", Subject: "Design Works", Country: "us", CountryName: "United States", Client: "Carlson", Status: "Paid", StatusColor: "success"}}))
	assert.Contains(t, h, `<table class="table card-table table-vcenter table-selectable text-nowrap datatable" id="invoices-table">`)
	assert.Contains(t, h, `data-select-all`)
	assert.Contains(t, h, `aria-label="Select invoice 001401"`)
	assert.Contains(t, h, `Showing <strong>1 to 1</strong> of <strong>16 entries</strong>`)
	assert.Contains(t, h, `aria-label="Invoices pagination"`)
}

func TestAdvanced(t *testing.T) {
	p := tables.AdvancedProps{ID: "emp", Title: "Employee", Search: true, Page: 1, Pages: 2,
		Columns: []tables.Column{{Field: "sort-name", Label: "Name"}, {Field: "sort-category", Label: "Category", Class: "py-0"}},
		Rows:    []tables.Row{{Name: "Ann", Cells: []templ.Component{templ.Raw("Ann"), tables.OnChecked("Agencies", nil)}}}}
	h := render(t, tables.Advanced(p))
	assert.Contains(t, h, `aria-label="Search" id="emp-search">`)
	assert.Contains(t, h, `<button type="button" class="table-sort d-flex justify-content-between" data-sort="sort-name">Name</button>`)
	assert.Contains(t, h, `<tbody class="table-tbody">`)
	assert.Contains(t, h, `<td class="sort-category py-0"><span class="on-unchecked">Agencies</span>`)
	assert.Contains(t, h, `data-table-id="emp"`)
	assert.Contains(t, h, `<span id="emp-page-count" class="me-1">20</span>`)
	s := render(t, tables.AdvancedScript(p))
	assert.Contains(t, s, `"keys":["sort-name","sort-category"]`)
	assert.Contains(t, s, `new List(`)
}
