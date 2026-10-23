package tablertempl_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/alert"
	"github.com/mind-vm/tabler-templ/badge"
	"github.com/mind-vm/tabler-templ/button"
	"github.com/mind-vm/tabler-templ/card"
	"github.com/mind-vm/tabler-templ/form"
	"github.com/mind-vm/tabler-templ/icon"
	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/modal"
	"github.com/mind-vm/tabler-templ/pagination"
	"github.com/mind-vm/tabler-templ/sidebar"
	"github.com/mind-vm/tabler-templ/table"
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

func TestButton(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(button.Button(button.Props{Variant: button.VariantPrimary, Size: button.SizeSmall, Icon: "ti-plus"}), "Add"))
	assert.Equal(t, `<button class="btn btn-primary btn-sm"><i class="ti ti-plus" aria-hidden="true"></i>Add</button>`, h)

	h = render(t, ctx, mit(button.Button(button.Props{Href: "/x", Variant: button.VariantDanger, Outline: true}), "Go"))
	assert.Equal(t, `<a href="/x" class="btn btn-outline-danger">Go</a>`, h)

	h = render(t, ctx, button.Button(button.Props{IconOnly: true, Icon: "ti-x", Label: "Remove", Type: "button"}))
	assert.Contains(t, h, `class="btn btn-icon"`)
	assert.Contains(t, h, `aria-label="Remove"`)

	// Attributes override defaults.
	h = render(t, ctx, button.Button(button.Props{Label: "A", Attributes: templ.Attributes{"title": "B"}}))
	assert.Contains(t, h, `title="B"`)
	assert.NotContains(t, h, `title="A"`)
}

// Links are sanitized: a javascript: URL from user input must never become a clickable href.
func TestUnsafeHrefIsSanitized(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(button.Button(button.Props{Href: "javascript:alert(1)"}), "x"))
	assert.NotContains(t, h, "javascript:")
	h = render(t, ctx, pagination.Pagination(pagination.Props{Page: 1, Pages: 2, Href: func(int) string { return "javascript:alert(1)" }}))
	assert.NotContains(t, h, "javascript:")
	assert.Contains(t, render(t, ctx, mit(button.Button(button.Props{Href: "/kunden?q=a&b=1"}), "x")), `href="/kunden?q=a&amp;b=1"`)
}

func TestCardAndIcon(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(card.Card(card.Props{Section: true, Class: "h-100 card"}), "x"))
	assert.Equal(t, `<section class="card h-100">x</section>`, h)
	assert.Equal(t, `<i class="ti ti-alert-triangle" role="img" aria-label="Warning"></i>`,
		render(t, ctx, icon.Icon("ti-alert-triangle", icon.Props{Label: "Warning"})))
	assert.Equal(t, `<span class="badge bg-green-lt">ok</span>`, render(t, ctx, mit(badge.Badge(badge.Props{Color: "green", Light: true}), "ok")))
}

func TestModal(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(modal.Modal(modal.Props{ID: "m1", Title: "Edit", Centered: true, Fade: true}), "<p>body</p>"))
	assert.Contains(t, h, `class="modal fade" id="m1"`)
	assert.Contains(t, h, `aria-labelledby="m1-title"`)
	assert.Contains(t, h, `class="modal-dialog modal-dialog-centered"`)
	assert.Contains(t, h, `<div class="modal-body"><p>body</p></div>`)
	assert.Contains(t, h, `aria-label="Close"`)
	assert.Equal(t, "#m1", modal.Trigger("m1")["data-bs-target"])

	f := render(t, ctx, mit(modal.Modal(modal.Props{ID: "m3", DialogClass: "x", Form: templ.Attributes{"method": "post"}}), "b"))
	assert.Contains(t, f, `<div class="modal-dialog x"><form class="modal-content" method="post">`)
	assert.Contains(t, f, `<div class="modal-body">b</div></form>`)

	full := render(t, ctx, modal.Modal(modal.Props{ID: "m2", Size: modal.SizeFullscreen, Centered: true}))
	assert.Contains(t, full, `class="modal-dialog modal-fullscreen"`)
}

// Messages come from the context; props still win.
func TestMessages(t *testing.T) {
	de := kit.WithMessages(context.Background(), kit.German)
	assert.Contains(t, render(t, de, modal.Modal(modal.Props{ID: "m"})), `aria-label="Schließen"`)
	assert.Contains(t, render(t, de, modal.Modal(modal.Props{ID: "m", CloseLabel: "Zu"})), `aria-label="Zu"`)
	assert.Contains(t, render(t, de, alert.Alert(alert.Props{Dismissible: true})), `aria-label="Schließen"`)
	href := func(n int) string { return "?p=" + string(rune('0'+n)) }
	assert.Contains(t, render(t, de, pagination.Pagination(pagination.Props{Page: 1, Pages: 3, Href: href})), `aria-label="Nächste Seite"`)
}

func TestPagination(t *testing.T) {
	assert.Equal(t, []int{1, 0, 4, 5, 6, 0, 10}, pagination.Numbers(5, 10))
	assert.Equal(t, []int{1, 2, 0, 10}, pagination.Numbers(1, 10))
	href := func(n int) string { return "/l?seite=" + strings.Repeat("I", n) }
	assert.Empty(t, render(t, context.Background(), pagination.Pagination(pagination.Props{Page: 1, Pages: 1, Href: href})))
	h := render(t, context.Background(), pagination.Pagination(pagination.Props{Page: 2, Pages: 3, Href: href}))
	assert.Contains(t, h, `<li class="page-item active"><a class="page-link" href="/l?seite=II" aria-current="page">2</a></li>`)
}

func TestForm(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(form.Field(form.FieldProps{Label: "Name", For: "n", Required: true, Hint: "h", Error: "e"}),
		"<input>"))
	assert.Equal(t, `<div class="mb-3"><label class="form-label required" for="n">Name</label><input><div class="invalid-feedback d-block">e</div><div class="form-hint">h</div></div>`, h)
	assert.Equal(t, `<input class="form-control is-invalid" value="a&amp;b" required id="n" name="name">`,
		render(t, ctx, form.Input(form.InputProps{ID: "n", Name: "name", Value: "a&b", Required: true, Invalid: true})))
	h = render(t, ctx, form.Select(form.SelectProps{Name: "s", Empty: "All", Selected: "b",
		Options: []form.Option{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}}))
	assert.Contains(t, h, `<option value="">All</option>`)
	assert.Contains(t, h, `<option value="b" selected>B</option>`)
}

func TestFieldOhneWrapperKlasse(t *testing.T) {
	h := render(t, context.Background(), mit(form.Field(form.FieldProps{Class: "-", Label: "L"}), "<input>"))
	assert.Equal(t, `<div><label class="form-label">L</label><input></div>`, h)
}

func TestCheck(t *testing.T) {
	h := render(t, context.Background(), form.Check(form.CheckProps{Name: "q", Value: "1", Label: "Q", Checked: true}))
	assert.Equal(t, `<label class="form-check"><input class="form-check-input" type="checkbox" value="1" checked name="q"> <span class="form-check-label">Q</span></label>`, h)
}

func TestTable(t *testing.T) {
	h := render(t, context.Background(), mit(table.Table(table.Props{Responsive: true, Class: "card-table"}), "<tbody></tbody>"))
	assert.Equal(t, `<div class="table-responsive"><table class="table card-table"><tbody></tbody></table></div>`, h)
}

func TestClasses(t *testing.T) {
	assert.Equal(t, "btn btn-sm x", kit.Classes("btn", "", "btn-sm btn", " x "))
}

func TestSidebar(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	item := sidebar.Item(sidebar.ItemProps{Href: "/kunden", Icon: "ti-building", Label: "Kunden", Active: true, Count: 3, CountTitle: "offen"})
	nav := mit(sidebar.Sidebar(sidebar.Props{ID: "leiste", Label: "Navigation"}), render(t, ctx, item))
	h := render(t, ctx, nav)
	assert.Contains(t, h, `<aside id="leiste" class="kit-sidebar offcanvas-lg offcanvas-start d-print-none" tabindex="-1" aria-label="Navigation">`)
	assert.Contains(t, h, `class="kit-sidebar-link active" aria-current="page" data-kit-tooltip="Kunden">`)
	assert.Contains(t, h, `<span class="badge bg-red-lt kit-sidebar-badge" title="offen">3</span>`)
	assert.NotContains(t, h, "data-collapsed")

	// Collapsed: attribute on the sidebar, tooltips on the items (they read the state from the context).
	collapsed := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		inner := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { return item.Render(ctx, w) })
		return sidebar.Sidebar(sidebar.Props{ID: "leiste", Collapsed: true}).Render(templ.WithChildren(ctx, inner), w)
	})
	h = render(t, ctx, collapsed)
	assert.Contains(t, h, `data-collapsed`)
	assert.Contains(t, h, `title="Kunden"`)

	// Trigger: offcanvas on small screens, the script's toggle on large ones; German labels.
	h = render(t, ctx, sidebar.Trigger(sidebar.TriggerProps{Target: "leiste"}))
	assert.Contains(t, h, `data-bs-toggle="offcanvas" data-bs-target="#leiste"`)
	assert.Contains(t, h, `aria-label="Menü öffnen"`)
	assert.Contains(t, h, `data-kit-sidebar-toggle="leiste"`)
	assert.Contains(t, h, `aria-expanded="true"`)

	// User: the children become the dropdown's items.
	h = render(t, ctx, mit(sidebar.User(sidebar.UserProps{Initials: "BB", Name: "Birgit Büro", Detail: "Büro"}), `<a class="dropdown-item" href="/x">Abmelden</a>`))
	assert.Contains(t, h, `data-bs-toggle="dropdown"`)
	assert.Contains(t, h, `<div class="dropdown-menu"><a class="dropdown-item" href="/x">Abmelden</a></div>`)

	assert.Contains(t, render(t, ctx, sidebar.Styles()), ".kit-sidebar[data-collapsed]")
	assert.Contains(t, render(t, ctx, sidebar.Script()), sidebar.CookieName)
}

func TestSidebarCollapsedCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	assert.False(t, sidebar.Collapsed(r))
	r.AddCookie(&http.Cookie{Name: sidebar.CookieName, Value: "collapsed"})
	assert.True(t, sidebar.Collapsed(r))
}
