package tablertempl_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/sidebar"
)

func TestSidebarCollapsibleGroup(t *testing.T) {
	ctx := context.Background()
	item := render(t, ctx, sidebar.Item(sidebar.ItemProps{Href: "/r", Label: "Volume"}))
	closed := render(t, ctx, mit(sidebar.Group(sidebar.GroupProps{Label: "Reports & more", Collapsible: true}), item))
	assert.Contains(t, closed, "data-closed")
	assert.Contains(t, closed, `data-kit-sidebar-group aria-expanded="false" aria-controls="kit-sidebar-group-reports---more"`)
	assert.Contains(t, closed, `<ul class="kit-sidebar-menu" id="kit-sidebar-group-reports---more">`)

	open := render(t, ctx, mit(sidebar.Group(sidebar.GroupProps{Label: "Reports", Collapsible: true, Open: true, ID: "rep"}), item))
	assert.NotContains(t, open, "data-closed")
	assert.Contains(t, open, `aria-expanded="true" aria-controls="rep"`)

	plain := render(t, ctx, mit(sidebar.Group(sidebar.GroupProps{Label: "Work"}), item))
	assert.Contains(t, plain, `<div class="kit-sidebar-label">Work</div>`)
	assert.NotContains(t, plain, "data-kit-sidebar-group")
}

func TestSidebarSwitcher(t *testing.T) {
	h := render(t, context.Background(), mit(sidebar.Switcher(sidebar.SwitcherProps{Initials: "AC", Title: "Acme", Subtitle: "Enterprise", Label: "Switch workspace"}),
		`<a class="dropdown-item" href="/w/2">Globex</a>`))
	for _, want := range []string{`class="dropdown kit-sidebar-switcher"`, `data-bs-toggle="dropdown"`, `aria-label="Switch workspace"`,
		"kit-sidebar-switcher-logo", ">AC<", "Acme", "Enterprise", "Globex"} {
		assert.Contains(t, h, want)
	}
}

func TestShellTabsAndCompact(t *testing.T) {
	h := render(t, context.Background(), shell.Page(shell.Props{Title: "Ticket", Compact: true, Tabs: templ.Raw(`<ul class="nav nav-bordered"></ul>`)}))
	assert.Contains(t, h, `<body class="kit-compact">`)
	assert.Contains(t, h, `<div class="kit-shell-tabs"><ul class="nav nav-bordered"></ul></div>`)
	assert.Contains(t, h, ".kit-compact .card-body")
	plain := render(t, context.Background(), shell.Page(shell.Props{Title: "Ticket"}))
	assert.NotContains(t, plain, `class="kit-compact"`)
	assert.NotContains(t, plain, "kit-shell-tabs\">")
}
