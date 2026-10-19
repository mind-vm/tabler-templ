package blocks_test

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tablerkit/alert"
	"github.com/mind-vm/tablerkit/blocks/detail"
	"github.com/mind-vm/tablerkit/blocks/listpage"
	"github.com/mind-vm/tablerkit/blocks/shell"
	"github.com/mind-vm/tablerkit/kit"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(ctx, &b))
	return b.String()
}

func mit(c templ.Component, kinder string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, templ.Raw(kinder)), w)
	})
}

var status = listpage.Filter{Name: "status", Label: "All states", Options: []listpage.Option{{Value: "open", Label: "Open"}, {Value: "done", Label: "Done"}}}

func TestListFromRequest(t *testing.T) {
	l := listpage.New(httptest.NewRequest("GET", "/customers?q=+berlin+&status=done&page=2", nil), listpage.Config{}, status)
	assert.Equal(t, "berlin", l.Query)
	assert.Equal(t, "done", l.Value("status"))
	assert.True(t, l.Active())
	assert.Equal(t, "/customers?page=3&q=berlin&status=done", l.Link(3))
	assert.Equal(t, "/customers?q=berlin&status=done", l.Link(1), "page 1 is the plain URL")

	// Unknown filter values count as "all"; parameter names are configurable.
	l = listpage.New(httptest.NewRequest("GET", "/x?status=hacked&seite=2", nil), listpage.Config{PageParam: "seite", PageSize: 2}, status)
	assert.Equal(t, "", l.Value("status"))
	assert.False(t, l.Active())
	assert.Equal(t, "/x?seite=2", l.Link(2))
	assert.Equal(t, "nothing here", l.Empty("nothing here", "no results"))
}

func TestPaginate(t *testing.T) {
	rows := []int{1, 2, 3, 4, 5}
	l := listpage.New(httptest.NewRequest("GET", "/x?page=2", nil), listpage.Config{PageSize: 2})
	assert.Equal(t, []int{3, 4}, listpage.Paginate(l, rows, nil))
	assert.Equal(t, 5, l.Total)
	assert.Equal(t, 3, l.Pages)
	assert.Equal(t, [2]int{3, 4}, [2]int{l.From, l.To})

	// A page beyond the last shows the last; keep filters.
	l = listpage.New(httptest.NewRequest("GET", "/x?page=99", nil), listpage.Config{PageSize: 2})
	assert.Equal(t, []int{5}, listpage.Paginate(l, rows, func(n int) bool { return n%2 == 1 }))
	assert.Equal(t, 2, l.Page)

	l = listpage.New(httptest.NewRequest("GET", "/x", nil), listpage.Config{})
	assert.Empty(t, listpage.Paginate(l, rows, func(int) bool { return false }))
	assert.Equal(t, [3]int{0, 0, 1}, [3]int{l.From, l.To, l.Pages})
}

func TestContains(t *testing.T) {
	assert.True(t, listpage.Contains("BERLIN mitte", "Mitte", "Berlin"))
	assert.False(t, listpage.Contains("berlin süd", "Berlin Mitte"))
	assert.True(t, listpage.Contains("", "x"))
}

func TestListCard(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	l := listpage.New(httptest.NewRequest("GET", "/kunden?q=a", nil), listpage.Config{Placeholder: "Name …"}, status)
	listpage.Paginate(l, make([]int, 30), nil)
	h := render(t, ctx, mit(listpage.Card(l), "<table></table>"))
	assert.Contains(t, h, `<form method="get" action="/kunden" class="kit-list-search" role="search">`)
	assert.Contains(t, h, `value="a" aria-label="Suchen" name="q" placeholder="Name …"`)
	assert.Contains(t, h, `aria-label="Suchen"`)
	assert.Contains(t, h, `>Zurücksetzen</a>`, "reset while the search is active")
	assert.Contains(t, h, "<table></table>")
	assert.Contains(t, h, "1–25 von 30")
	assert.Contains(t, h, `href="/kunden?page=2&amp;q=a"`)

	// One page: no page numbers; empty list: no footer.
	l = listpage.New(httptest.NewRequest("GET", "/kunden", nil), listpage.Config{})
	listpage.Paginate(l, make([]int, 3), nil)
	h = render(t, ctx, listpage.Footer(l))
	assert.Contains(t, h, "1–3 von 3")
	assert.NotContains(t, h, "pagination")
	listpage.Paginate(l, []int{}, nil)
	assert.Empty(t, render(t, ctx, listpage.Footer(l)))
}

func TestShellPage(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	p := shell.Props{Lang: "de", Title: "Kunden", AppName: "Acme",
		Head:        templ.Raw(`<link rel="stylesheet" href="/tabler.css">`),
		Scripts:     templ.Raw(`<script src="/tabler.js"></script>`),
		Sidebar:     templ.Raw(`<aside id="sidebar" class="kit-sidebar"></aside>`),
		Breadcrumbs: []shell.Crumb{{Text: "Ablauf"}, {Text: "Kunden", Href: "/kunden"}, {Text: "Stadt"}},
		Subtitle:    templ.Raw("6 Kunden"),
		Actions:     templ.Raw(`<a class="btn">Neu</a>`),
		Flash:       &shell.Flash{Variant: alert.VariantSuccess, Text: "Gespeichert"},
	}
	h := render(t, ctx, mit(shell.Page(p), "<p>Inhalt</p>"))
	assert.True(t, strings.HasPrefix(h, "<!doctype html><html lang=\"de\">"), h[:40])
	assert.Contains(t, h, "<title>Kunden · Acme</title>")
	assert.Contains(t, h, `</style><link rel="stylesheet" href="/tabler.css">`, "the kit's CSS first, so the app's theme overrides it")
	assert.Contains(t, h, `data-bs-target="#sidebar"`, "sidebar buttons in the header")
	assert.Contains(t, h, `aria-label="Brotkrumen"`)
	assert.Contains(t, h, `<li class="breadcrumb-item">Ablauf</li><li class="breadcrumb-item"><a href="/kunden">Kunden</a></li><li class="breadcrumb-item active" aria-current="page">Stadt</li>`)
	assert.Contains(t, h, `<h2 class="page-title">Kunden</h2>`)
	assert.Contains(t, h, `<div class="kit-shell-subtitle">6 Kunden</div>`)
	assert.Contains(t, h, `<div class="btn-list"><a class="btn">Neu</a></div>`)
	assert.Contains(t, h, "alert-success")
	assert.Contains(t, h, "<p>Inhalt</p>")
	assert.Contains(t, h, `<script src="/tabler.js"></script><script>`, "the sidebar script after the app's")

	// Without sidebar (sign-in): the brand instead of the buttons, no sidebar script.
	h = render(t, ctx, shell.Page(shell.Props{Title: "Anmelden", Brand: templ.Raw(`<a class="brand">Acme</a>`)}))
	assert.Contains(t, h, `<a class="brand">Acme</a>`)
	assert.NotContains(t, h, "data-kit-sidebar-toggle")
	assert.NotContains(t, h, "kit_sidebar=")
}

func TestDetailSection(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	edit := detail.EditButton(detail.EditButtonProps{Dialog: "dialog-kunde"})
	h := render(t, ctx, mit(detail.Section(detail.SectionProps{Title: "Kunde", Extra: "3", Action: edit, Class: "h-100"}), "<p>x</p>"))
	assert.Contains(t, h, `<section class="card kit-section h-100">`)
	assert.Contains(t, h, `<h3 class="kit-section-title">Kunde <small>3</small></h3>`)
	assert.Contains(t, h, `data-bs-target="#dialog-kunde"`)
	assert.Contains(t, h, `title="Bearbeiten"`)
	assert.Contains(t, h, `<span class="kit-section-action-text">Bearbeiten</span>`)
	assert.Contains(t, h, "<p>x</p>")
	assert.NotContains(t, render(t, ctx, detail.Section(detail.SectionProps{})), "kit-section-head", "no head without title and action")

	h = render(t, ctx, mit(detail.Fields(), render(t, ctx, detail.Field("Telefon", "", "–"))+render(t, ctx, detail.Field("Ort", "Bonn", "–"))))
	assert.Equal(t, `<dl class="kit-fields"><dt>Telefon</dt><dd class="kit-field-empty">–</dd><dt>Ort</dt><dd>Bonn</dd></dl>`, h)
}

func TestDetailDialog(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, mit(detail.Dialog(detail.DialogProps{ID: "dialog-kunde", Title: "Kunde", Tooltip: "Kunde (E-32)", Context: "10001 · Stadt",
		Section: true, Error: "Kundennummer ist fest"}), "<form></form>"))
	assert.Contains(t, h, `class="modal fade kit-dialog" id="dialog-kunde"`)
	assert.Contains(t, h, `data-kit-dialog`)
	assert.Contains(t, h, `data-kit-section`)
	assert.Contains(t, h, `<div class="kit-dialog-context">10001 · Stadt</div>`)
	assert.Contains(t, h, `<h2 class="modal-title" id="dialog-kunde-title" title="Kunde (E-32)">Kunde</h2>`)
	assert.Contains(t, h, `<div class="alert kit-dialog-error alert-danger" role="alert">`)
	assert.Contains(t, h, "<form></form>")

	h = render(t, ctx, detail.Dialog(detail.DialogProps{ID: "dialog-liste", Title: "Liste", Large: true}))
	assert.Contains(t, h, "kit-dialog-large")
	assert.Contains(t, h, "modal-fullscreen")
	assert.NotContains(t, h, "data-kit-section")
	assert.NotContains(t, h, "kit-dialog-error")
}

func TestDetailScript(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, detail.Script(detail.ScriptConfig{DialogField: "_drawer", StateVar: "drawerZustand"}))
	assert.True(t, strings.HasPrefix(h, "<script>(function(cfg){"))
	assert.Contains(t, h, `"dialogField":"_drawer"`)
	assert.Contains(t, h, `"returnField":"_return"`, "defaults for what is not set")
	assert.Contains(t, h, `"stateVar":"drawerZustand"`)
	assert.Contains(t, h, `"cancel":"Abbrechen"`)
	assert.Contains(t, h, `"sectionHint":"Gespeichert wird nur dieser Abschnitt."`)
	assert.Equal(t, 1, strings.Count(h, "</script>"))

	// A configured name cannot end the script.
	h = render(t, ctx, detail.Script(detail.ScriptConfig{StateVar: "</script><script>alert(1)"}))
	assert.Equal(t, 1, strings.Count(h, "</script>"))
}
