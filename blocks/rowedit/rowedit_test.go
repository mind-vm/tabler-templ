package rowedit_test

// The rules of the concept Gestaltung for editing a list (GS-55 to GS-59) that the block itself keeps;
// the tests carry the rule's number in their name.

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/rowedit"
	"github.com/mind-vm/tabler-templ/kit"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(templ.InitializeContext(ctx), &b))
	return b.String()
}

func with(c templ.Component, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, templ.Join(children...)), w)
	})
}

func row(id string) templ.Component {
	return with(rowedit.Row(rowedit.RowProps{Action: "/offers/1/positions/" + id, Kind: "position"}),
		templ.Raw(`<input type="hidden" name="id" value="`+id+`"><input name="qty" value="1">`),
		with(rowedit.Actions(),
			rowedit.Save(rowedit.SaveProps{Label: "Save position " + id}),
			rowedit.Delete(rowedit.DeleteProps{Action: "/offers/1/positions/" + id + "/delete", Confirm: "Delete position " + id + "?"})))
}

func dialog() templ.Component {
	return with(rowedit.Dialog(rowedit.DialogProps{ID: "dialog-positions", Title: "Edit positions", Context: "A-1 · Acme"}),
		with(rowedit.List(rowedit.ListProps{Columns: "8rem 1fr 3rem"}),
			with(rowedit.Head(), templ.Raw("<span>No.</span><span>Qty</span><span></span>")),
			row("p1"), row("p2"),
			with(rowedit.Row(rowedit.RowProps{Action: "/offers/1/positions", Kind: "position-new", New: true}),
				templ.Raw(`<input name="qty" value="">`), with(rowedit.Actions(), rowedit.Add(rowedit.AddProps{})))),
		rowedit.Total("Total", "1.234,50 €"))
}

// GS-55: the dialog covers (almost) the whole page; the detail script reopens it after each save.
func TestGS55_FullPageDialog(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Contains(t, html, "kit-dialog-large")
	assert.Contains(t, html, "modal-fullscreen")
	assert.Contains(t, html, "kit-rowedit-dialog")
	assert.Contains(t, html, `data-kit-dialog`)
	assert.Contains(t, html, `<div class="kit-dialog-context">A-1 · Acme</div>`)
}

// GS-56: every row is a form of its own with Save before Delete (Enter saves), and the hint on top.
func TestGS56_EveryRowIsAForm(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Equal(t, 3, strings.Count(html, `<form method="post"`))
	assert.Regexp(t, `<form method="post" action="/offers/1/positions/p1" class="kit-rowedit-row"[^>]*data-form="position"`, html)
	assert.Contains(t, html, "Each row is saved on its own.")
	r := html[strings.Index(html, `/positions/p1"`):]
	assert.Less(t, strings.Index(r, "kit-rowedit-save"), strings.Index(r, "kit-rowedit-delete"), "Save is the form's first button, so Enter saves")
	assert.Regexp(t, `style="--kit-rowedit-columns:8rem 1fr 3rem;?"`, html)
	assert.Contains(t, html, `aria-label="Save position p1"`)

	de := render(t, kit.WithMessages(context.Background(), kit.German), dialog())
	assert.Contains(t, de, "Jede Zeile wird für sich gespeichert.")
	assert.Contains(t, de, "Hinzufügen")
}

// GS-57: the new row is the last one.
func TestGS57_NewRowAtTheBottom(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Greater(t, strings.Index(html, `class="kit-rowedit-row kit-rowedit-new"`), strings.LastIndex(html, `data-form="position"`))
	assert.Contains(t, html, "kit-rowedit-add")
}

// GS-58: Delete asks first and works on an invalid row.
func TestGS58_DeleteAsksFirst(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Contains(t, html, `data-kit-confirm="Delete position p1?"`)
	assert.Contains(t, html, `formaction="/offers/1/positions/p1/delete"`)
	assert.Contains(t, html, "formnovalidate")
	s := render(t, context.Background(), templ.Join(rowedit.Script(), rowedit.Script()))
	assert.Equal(t, 1, strings.Count(s, "<script"), "the script once per page")
	assert.Contains(t, s, `data-kit-confirm`)
	assert.Contains(t, s, "window.confirm")
}

// GS-59: the total is shown only; a locked list says why instead of its buttons.
func TestGS59_TotalAndLocked(t *testing.T) {
	html := render(t, context.Background(), dialog())
	i := strings.Index(html, `<div class="kit-rowedit-total">`)
	require.Positive(t, i)
	assert.NotContains(t, html[i:], "<input")
	locked := render(t, context.Background(), rowedit.Locked("Submitted: read-only."))
	assert.Contains(t, locked, "Submitted: read-only.")
	assert.NotContains(t, locked, "<button")
}

// GS-54: the list scrolls in the dialog on narrow screens, never the page.
func TestGS54_ListScrollsInItsBox(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Contains(t, html, `<div class="kit-rowedit-scroll">`)
	assert.Contains(t, html, ".kit-rowedit-scroll{max-width:100%;overflow-x:auto}")
}

// GS-56: Save and Delete stand at the right, also while the row scrolls sideways.
func TestGS56_ActionsStayAtTheRight(t *testing.T) {
	html := render(t, context.Background(), dialog())
	assert.Contains(t, html, ".kit-rowedit-actions{position:sticky;right:0")
}

// GS-56: long texts in an optional second line of the row, inside its form; closed, open after an error.
func TestGS56_DetailsLine(t *testing.T) {
	row := with(rowedit.Row(rowedit.RowProps{Action: "/p/1"}),
		templ.Raw(`<input name="qty">`), with(rowedit.Actions(), rowedit.Save(rowedit.SaveProps{})),
		with(rowedit.Details(rowedit.DetailsProps{Summary: "Langtext"}), templ.Raw(`<textarea name="long"></textarea>`)))
	html := render(t, context.Background(), row)
	form := html[strings.Index(html, "<form"):strings.Index(html, "</form>")]
	assert.Contains(t, form, `<details class="kit-rowedit-details"><summary>Langtext</summary>`)
	assert.Contains(t, form, `<textarea name="long">`)
	assert.Less(t, strings.Index(form, "kit-rowedit-save"), strings.Index(form, "<details"), "Save stays in the first line")
	open := render(t, context.Background(), rowedit.Details(rowedit.DetailsProps{Summary: "x", Open: true}))
	assert.Contains(t, open, `<details class="kit-rowedit-details" open>`)
}

// GS-54 on a list to edit: below the breakpoint every row is a card; each field carries its label.
func TestGS54_RowsBecomeCards(t *testing.T) {
	list := with(rowedit.List(rowedit.ListProps{Columns: "8rem 1fr 3rem", Mobile: "md"}),
		with(rowedit.Row(rowedit.RowProps{Action: "/p/1"}),
			with(rowedit.Cell("Menge"), templ.Raw(`<input name="qty" aria-label="Menge">`)),
			with(rowedit.Actions(), rowedit.Save(rowedit.SaveProps{}))))
	html := render(t, context.Background(), list)
	assert.Contains(t, html, `class="kit-rowedit kit-rowedit-mobile-md"`)
	assert.Contains(t, html, `<div class="kit-rowedit-cell"><span class="kit-rowedit-label" aria-hidden="true">Menge</span><input name="qty"`)
	assert.Contains(t, html, "@media (max-width:767.98px){.kit-rowedit-mobile-md{min-width:0}")
	assert.Contains(t, html, ".kit-rowedit-mobile-md .kit-rowedit-actions{order:2;position:static")
	css := html[strings.Index(html, "@media (max-width:767.98px){.kit-rowedit-mobile-md"):]
	assert.Less(t, strings.Index(css, " .kit-rowedit-row{display:flex"), strings.Index(css, " .kit-rowedit-head{display:none}"), "the head is hidden after the rows become flex (it has both classes)")
}
