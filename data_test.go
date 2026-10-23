package tablertempl_test

import (
	"context"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/datagrid"
	"github.com/mind-vm/tabler-templ/divider"
	"github.com/mind-vm/tabler-templ/iconlist"
	"github.com/mind-vm/tabler-templ/empty"
	"github.com/mind-vm/tabler-templ/listgroup"
	"github.com/mind-vm/tabler-templ/progress"
	"github.com/mind-vm/tabler-templ/ribbon"
	"github.com/mind-vm/tabler-templ/status"
	"github.com/mind-vm/tabler-templ/steps"
	"github.com/mind-vm/tabler-templ/timeline"
	"github.com/mind-vm/tabler-templ/tracking"
)

func TestProgress(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, progress.Progress(progress.Props{Value: 38, Color: "red", Size: progress.SizeSmall, Striped: true, Animated: true}))
	assert.Equal(t, `<div class="progress progress-sm"><div class="progress-bar bg-red progress-bar-striped progress-bar-animated" style="width: 38%;" role="progressbar" aria-valuenow="38" aria-valuemin="0" aria-valuemax="100" aria-label="38% Complete"><span class="visually-hidden">38% Complete</span></div></div>`, h)
	assert.Contains(t, render(t, ctx, progress.Progress(progress.Props{Indeterminate: true})), `progress-bar-indeterminate`)
	assert.Contains(t, render(t, ctx, progress.Stacked(progress.Props{}, []progress.Bar{{Value: 10, Color: "blue"}})), `<div class="progress progress-separated" style="width: 10%;"`)
	assert.Contains(t, render(t, ctx, progress.Background(progress.BackgroundProps{Value: 65, Text: "Poland", ShowValue: true})),
		`<div class="progress-bar bg-primary-lt"`)
	h = render(t, ctx, progress.Steps(progress.StepsProps{Labels: []string{"A", "B", "C"}, Active: 2}))
	assert.Equal(t, `<ol class="progress-steps"><li class="progress-steps-item bg-primary"><span class="visually-hidden">A</span></li><li class="progress-steps-item bg-primary" aria-current="step"><span class="visually-hidden">B</span></li><li class="progress-steps-item"><span class="visually-hidden">C</span></li></ol>`, h)
}

func TestSteps(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, steps.Steps(steps.Props{Active: 2, Color: "green", Counter: true}, []steps.Step{{Title: "Cart"}, {Title: "Billing"}}))
	assert.Equal(t, `<ul class="steps steps-green steps-counter"><li class="step-item">Cart</li><li class="step-item active">Billing</li></ul>`, h)
	h = render(t, ctx, steps.Steps(steps.Props{Active: 1}, []steps.Step{{}, {}}))
	assert.Contains(t, h, `<ol class="steps" aria-label="Progress">`)
	assert.Contains(t, h, `<a href="#" aria-current="step"><span class="visually-hidden">Step 1 (current)</span></a>`)
	assert.Contains(t, h, `<span><span class="visually-hidden">Step 2</span></span>`)
}

func TestStatusRibbonEmpty(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<span class="status status-green status-lite"><span class="status-dot status-dot-animated"></span>Up</span>`,
		render(t, ctx, mit(status.Status(status.Props{Color: "green", Lite: true, Animated: true}), "Up")))
	assert.Contains(t, render(t, ctx, status.Indicator(status.DotProps{Color: "red", Animated: true})), `class="status-indicator status-red status-indicator-animated"`)
	assert.Equal(t, `<div class="ribbon ribbon-top ribbon-start bg-red">NEW</div>`, render(t, ctx, ribbon.Ribbon(ribbon.Props{Top: true, Start: true, Color: "red", Text: "NEW"})))
	h := render(t, ctx, empty.Empty(empty.Props{Header: "404", Title: "Oops", Subtitle: "Not found", TitleLevel: 1}))
	assert.Equal(t, `<div class="empty"><div class="empty-header">404</div><h1 class="empty-title">Oops</h1><p class="empty-subtitle text-secondary">Not found</p></div>`, h)
}

func TestDatagridTrackingTimeline(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(datagrid.Grid(datagrid.Props{Horizontal: true, Columns: 2, TitleWidth: "8rem"}), "x"))
	assert.Equal(t, `<div class="datagrid datagrid-horizontal datagrid-cols-2" style="--tblr-datagrid-title-width: 8rem;">x</div>`, h)
	assert.Equal(t, `<div class="datagrid-item"><dt class="datagrid-title">Name</dt><dd class="datagrid-content">Ann</dd></div>`,
		render(t, ctx, mit(datagrid.Item(datagrid.ItemProps{Title: "Name", List: true}), "Ann")))
	assert.Contains(t, render(t, ctx, tracking.Tracking(tracking.Props{}, []tracking.Block{{Color: "success", Title: "Up"}, {}})),
		`<div class="tracking-block bg-success" data-bs-placement="top" data-bs-toggle="tooltip" title="Up"></div><div class="tracking-block"></div>`)
	h = render(t, ctx, mit(timeline.Event(timeline.EventProps{Icon: "ti-check", IconColor: "green", Time: "1 day ago", Title: "Backup"}), "<p>ok</p>"))
	assert.Equal(t, `<li class="timeline-event"><div class="timeline-event-icon bg-green-lt"><i class="ti ti-check icon" aria-hidden="true"></i></div><div class="card timeline-event-card"><div class="card-body"><div class="text-secondary float-end">1 day ago</div><h4>Backup</h4><p>ok</p></div></div></li>`, h)
}

func TestListGroup(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="list-group list-group-sm list-group-flush list-group-hoverable">x</div>`,
		render(t, ctx, mit(listgroup.List(listgroup.Props{Size: listgroup.SizeSmall, Flush: true, Hoverable: true}), "x")))
	h := render(t, ctx, listgroup.Item(listgroup.ItemProps{Href: "/p", Active: true, Title: "Ann", Description: "Designer", Meta: templ.Raw("2 h")}))
	// Whitespace between the slots is irrelevant: the item is a flex row.
	assert.Contains(t, h, `<a href="/p" class="list-group-item list-group-item-action active" aria-current="true">`)
	assert.Contains(t, h, `<span class="list-group-item-content"><span class="list-group-item-title">Ann</span> <span class="list-group-item-description">Designer</span>`)
	assert.Contains(t, h, `<span class="list-group-item-meta">2 h</span>`)
	assert.Equal(t, `<div class="list-group-item list-group-item-danger">x</div>`, render(t, ctx, mit(listgroup.Item(listgroup.ItemProps{Color: "danger"}), "x")))
	assert.Equal(t, "true", listgroup.CollapseToggle("n", true)["aria-expanded"])
}

func TestDividerIconList(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="hr-text hr-text-start text-red">or</div>`, render(t, ctx, mit(divider.Text(divider.Props{Position: divider.Start, Class: "text-red"}), "or")))
	h := render(t, ctx, mit(iconlist.List(iconlist.Props{Color: "success", RowGap: "1rem"}), "x"))
	assert.Equal(t, `<ul class="list-icon" style="--tblr-list-icon-row-gap: 1rem; --tblr-list-icon-color: var(--tblr-success);">x</ul>`, h)
	assert.Equal(t, `<li><i class="ti ti-check icon" aria-hidden="true"></i>Done</li>`, render(t, ctx, mit(iconlist.Item(), "Done")))
}
