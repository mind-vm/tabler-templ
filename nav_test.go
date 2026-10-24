package tablertempl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/avatar"
	"github.com/mind-vm/tabler-templ/breadcrumb"
	"github.com/mind-vm/tabler-templ/button"
	"github.com/mind-vm/tabler-templ/dropdown"
	"github.com/mind-vm/tabler-templ/flag"
	"github.com/mind-vm/tabler-templ/nav"
	"github.com/mind-vm/tabler-templ/offcanvas"
	"github.com/mind-vm/tabler-templ/popover"
	"github.com/mind-vm/tabler-templ/toast"
	"github.com/mind-vm/tabler-templ/tooltip"
)

func TestDropdown(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<button type="button" class="btn dropdown-toggle" aria-expanded="false" data-bs-toggle="dropdown">Open</button>`,
		render(t, ctx, mit(dropdown.Toggle(button.Props{}), "Open")))
	assert.Equal(t, `<div class="dropdown-menu dropdown-menu-end dropdown-menu-arrow" data-bs-theme="dark">x</div>`,
		render(t, ctx, mit(dropdown.Menu(dropdown.MenuProps{End: true, Arrow: true, Dark: true}), "x")))
	assert.Equal(t, `<a class="dropdown-item active" href="/a" aria-current="true"><i class="ti ti-edit icon dropdown-item-icon" aria-hidden="true"></i>Edit</a>`,
		render(t, ctx, mit(dropdown.Item(dropdown.ItemProps{Href: "/a", Icon: "ti-edit", Active: true}), "Edit")))
	assert.Equal(t, `<button type="button" class="dropdown-item disabled" disabled>x</button>`,
		render(t, ctx, mit(dropdown.Item(dropdown.ItemProps{Disabled: true}), "x")))
	assert.Equal(t, `<label class="dropdown-item form-switch"><input class="form-check-input m-0 me-2" type="checkbox" name="n" checked> On</label>`,
		render(t, ctx, dropdown.Check(dropdown.CheckProps{Name: "n", Label: "On", Checked: true, Switch: true})))
	assert.Contains(t, render(t, ctx, dropdown.Submenu(dropdown.SubmenuProps{Label: "Move"})), `<div class="dropdown dropend"><a class="dropdown-item dropdown-toggle"`)
}

func TestAvatarAndFlag(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<span class="avatar avatar-xs avatar-rounded bg-green-lt">AB</span>`,
		render(t, ctx, avatar.Avatar(avatar.Props{Initials: "AB", Size: avatar.SizeXS, Shape: avatar.ShapeCircle, Color: "green"})))
	assert.Equal(t, `<span class="avatar" style="background-image: url(/a.jpg);"><span class="badge bg-danger" role="img" aria-label="Busy"></span> </span>`,
		render(t, ctx, avatar.Avatar(avatar.Props{Src: "/a.jpg", Status: "danger", StatusLabel: "Busy"})))
	assert.Equal(t, `<span class="avatar"><i class="ti ti-user icon" aria-hidden="true"></i></span>`, render(t, ctx, avatar.Avatar()))
	assert.Equal(t, `<div class="avatar-list avatar-list-stacked">x</div>`, render(t, ctx, mit(avatar.List(avatar.ListProps{Stacked: true}), "x")))
	assert.Contains(t, render(t, ctx, avatar.Upload(avatar.UploadProps{For: "f", Text: "Upload", Size: avatar.SizeXL})),
		`<label for="f" class="avatar avatar-xl avatar-upload cursor-pointer">`)
	assert.Equal(t, `<span class="flag flag-xs flag-country-de" role="img" aria-label="Germany"></span>`,
		render(t, ctx, flag.Flag("DE", flag.Props{Size: flag.SizeXS, Label: "Germany"})))
}

func TestNav(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(nav.Nav(nav.Props{Style: nav.StyleTabs, CardHeader: true, Fill: true}), "x"))
	assert.Equal(t, `<ul class="nav nav-tabs card-header-tabs nav-fill">x</ul>`, h)
	h = render(t, ctx, mit(nav.Item(nav.ItemProps{Target: "home", Active: true, Icon: "ti-home"}), "Home"))
	assert.Contains(t, h, `<li class="nav-item" role="presentation"><a href="#home" class="nav-link active"`)
	assert.Contains(t, h, `data-bs-toggle="tab"`)
	assert.Contains(t, h, `aria-selected="true"`)
	assert.Contains(t, h, `<i class="ti ti-home icon me-2" aria-hidden="true"></i>Home`)
	assert.Contains(t, render(t, ctx, mit(nav.Item(nav.ItemProps{Href: "/a", Active: true}), "A")), `aria-current="page"`)
	assert.Equal(t, `<div class="tab-pane fade active show" id="home" role="tabpanel" aria-labelledby="home-tab">x</div>`,
		render(t, ctx, mit(nav.Pane(nav.PaneProps{ID: "home", Active: true, Fade: true}), "x")))
	assert.Contains(t, render(t, ctx, mit(nav.Nav(nav.Props{Overflow: true}), "")), `<div class="nav-overflow" data-bs-toggle="nav-overflow"><ul class="nav">`)
	h = render(t, ctx, mit(nav.Segment(nav.SegmentProps{Name: "v", ID: "v1", Value: "list", Active: true}), "List"))
	assert.Equal(t, `<input type="radio" class="nav-link-input" name="v" id="v1" value="list" checked> <label for="v1" class="nav-link">List</label>`, h)
}

func TestBreadcrumb(t *testing.T) {
	h := render(t, context.Background(), breadcrumb.Breadcrumb(breadcrumb.Props{Separator: breadcrumb.SeparatorArrows, Muted: true},
		[]breadcrumb.Crumb{{Text: "Home", Href: "/", Icon: "ti-home", IconOnly: true}, {Text: "Data"}}))
	assert.Equal(t, `<nav aria-label="Breadcrumbs"><ol class="breadcrumb breadcrumb-arrows breadcrumb-muted"><li class="breadcrumb-item"><a href="/" aria-label="Home"><i class="ti ti-home icon" aria-hidden="true"></i></a></li><li class="breadcrumb-item active" aria-current="page">Data</li></ol></nav>`, h)
}

func TestOffcanvasAndToast(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(offcanvas.Offcanvas(offcanvas.Props{ID: "f", Title: "Filters", Placement: offcanvas.End, Floating: true}), "x"))
	assert.Contains(t, h, `<div class="offcanvas offcanvas-end offcanvas-floating" tabindex="-1" role="dialog" id="f" aria-labelledby="f-title">`)
	assert.Contains(t, h, `<h2 class="offcanvas-title" id="f-title">Filters</h2><button type="button" class="btn-close" data-bs-dismiss="offcanvas" aria-label="Close">`)
	assert.Contains(t, render(t, ctx, offcanvas.Offcanvas(offcanvas.Props{ID: "s", Responsive: "lg"})), `class="offcanvas-lg offcanvas-start"`)
	assert.Equal(t, "#f", offcanvas.Trigger("f")["data-bs-target"])

	h = render(t, ctx, mit(toast.Toast(toast.Props{ID: "t", Title: "Ann", Time: "now", Show: true, Color: "success"}), "Hi"))
	assert.Contains(t, h, `<div class="toast toast-success show"`)
	assert.Contains(t, h, `data-bs-autohide="false"`)
	assert.Contains(t, h, `<strong class="me-auto">Ann</strong> <small>now</small> <button type="button" class="ms-2 btn-close" data-bs-dismiss="toast"`)
	assert.Contains(t, render(t, ctx, toast.Toast(toast.Props{Autohide: true, Delay: 3000})), `data-bs-delay="3000"`)
	assert.Contains(t, render(t, ctx, toast.Container(toast.ContainerProps{Position: toast.BottomEnd})), `class="toast-container position-fixed p-3 bottom-0 end-0"`)
}

func TestTooltipAndPopover(t *testing.T) {
	assert.Equal(t, "top", tooltip.Attrs("x", tooltip.Top)["data-bs-placement"])
	assert.Equal(t, "true", tooltip.HTML("<b>x</b>", tooltip.Auto)["data-bs-html"])
	a := popover.Attrs(popover.Props{Title: "T", Content: "C", Trigger: popover.Hover, Body: true})
	assert.Equal(t, "hover", a["data-bs-trigger"])
	assert.Equal(t, "body", a["data-bs-container"])
	assert.Equal(t, "C", a["data-bs-content"])
}
