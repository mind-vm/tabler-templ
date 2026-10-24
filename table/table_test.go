package table_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/table"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestClasses(t *testing.T) {
	assert.Equal(t, "table card-table table-vcenter table-nowrap table-striped table-hover table-borderless table-center table-sm table-selectable table-mobile-md x",
		table.Classes(table.Props{Card: true, VCenter: true, Nowrap: true, Striped: true, Hover: true, Borderless: true, Center: true, Small: true,
			Selectable: true, Mobile: "md", Class: "x"}))
	assert.Equal(t, "table", table.Classes(table.Props{}))
}

func TestSortHeader(t *testing.T) {
	h := render(t, table.SortHeader(table.SortProps{Field: "sort-name", Dir: table.Asc, Class: "d-flex"}))
	assert.Equal(t, `<button type="button" class="table-sort d-flex asc" data-sort="sort-name"></button>`, h)
}

func TestSelectCheck(t *testing.T) {
	assert.Equal(t, `<input class="form-check-input m-0 align-middle table-selectable-check" type="checkbox" aria-label="Select row 1" value="true" checked>`,
		render(t, table.SelectCheck(table.CheckProps{Label: "Select row 1", Value: "true", Checked: true})))
	assert.Equal(t, `<input class="form-check-input m-0 align-middle" type="checkbox" aria-label="Select all" data-select-all>`,
		render(t, table.SelectCheck(table.CheckProps{Label: "Select all", All: true})))
}
