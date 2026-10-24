package layout_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/errorpage"
	"github.com/mind-vm/tabler-templ/blocks/layout"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

var menu = []layout.Entry{
	{Title: "Home", Href: "/", Icon: "ti-home", Section: "Overview"},
	{Title: "Layout", Icon: "ti-layout-2", Active: true, Columns: 2, Children: []layout.Item{
		{Title: "Boxed", Href: "/boxed"}, {Title: "Fluid", Href: "/fluid", Active: true, Badge: "New"},
		{Title: "More", Children: []layout.Item{{Title: "RTL", Href: "/rtl"}}},
	}},
}

func TestPage(t *testing.T) {
	h := render(t, layout.Page(layout.PageProps{Title: "T", Position: layout.Vertical, Width: layout.Fluid, Sticky: true, NavbarTheme: "dark",
		SidebarMode: layout.Folded, Header: layout.PageHeader(layout.PageHeaderProps{Pretitle: "Overview", Title: "Orders"})}))
	assert.Contains(t, h, `data-bs-layout="fluid"`)
	assert.Contains(t, h, `data-bs-navbar="sticky"`)
	assert.Contains(t, h, `data-bs-navbar-position="vertical"`)
	assert.Contains(t, h, `data-bs-navbar-theme="dark"`)
	assert.Contains(t, h, `data-bs-sidebar="folded"`)
	assert.Contains(t, h, `<div class="page-pretitle">Overview</div><h1 class="page-title">Orders</h1>`)
	assert.Contains(t, h, `<div class="page-body"><div class="container-xl">`)
}

func TestMenu(t *testing.T) {
	h := render(t, layout.Sidebar(layout.SidebarProps{}, menu))
	assert.Contains(t, h, `<aside class="navbar navbar-vertical navbar-expand-lg">`)
	assert.Contains(t, h, `<li class="nav-section-title">Overview</li>`)
	assert.Contains(t, h, `<a class="nav-link" href="/">`)
	// The active dropdown is open in a sidebar and keeps open.
	assert.Contains(t, h, `data-bs-auto-close="false"`)
	assert.Contains(t, h, `<div class="dropdown-menu show"><div class="dropdown-menu-columns">`)
	assert.Contains(t, h, `<a class="dropdown-item active" href="/fluid" aria-current="page">`)
	assert.Contains(t, h, `<span class="badge badge-sm bg-green-lt text-uppercase ms-auto">New</span>`)
	assert.Contains(t, h, `<div class="dropend"><a class="dropdown-item dropdown-toggle"`)

	h = render(t, layout.Navbar(layout.NavbarProps{Dark: true}, menu))
	assert.Contains(t, h, `<header class="navbar navbar-expand-md d-print-none" data-bs-theme="dark">`)
	assert.Contains(t, h, `<nav aria-label="Primary" class="nav-overflow" data-bs-toggle="nav-overflow">`)
	assert.Contains(t, h, `data-bs-auto-close="outside"`)
	assert.NotContains(t, h, "nav-section-title")
}

func TestErrorPage(t *testing.T) {
	h := render(t, errorpage.Page(errorpage.Props{Title: "404"}, errorpage.Error{Code: "404", Text: "Not found"}))
	assert.Contains(t, h, `<body class="border-top-wide border-primary">`)
	assert.Contains(t, h, `<div class="empty-header">404</div><h1 class="empty-title">Oops… You just found an error page</h1>`)
	assert.Contains(t, h, `Take me home`)
}
