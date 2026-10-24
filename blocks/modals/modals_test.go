package modals_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/modals"
	"github.com/mind-vm/tabler-templ/button"
	"github.com/mind-vm/tabler-templ/modal"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestStatus(t *testing.T) {
	h := render(t, modals.Status(modals.StatusProps{Frame: modals.Frame{ID: "paid"}, Color: "success", IconColor: "green", Title: "Paid",
		Actions: []modals.Action{{Text: "Close"}, {Text: "Invoice", Href: "/invoice", Variant: button.VariantSuccess}}}))
	assert.Contains(t, h, `<div class="modal modal-blur fade" id="paid"`)
	assert.Contains(t, h, `<div class="modal-dialog modal-sm modal-dialog-centered"><div class="modal-content"><button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button><div class="modal-status bg-success"></div><div class="modal-body text-center py-4">`)
	assert.Contains(t, h, `<i class="ti ti-circle-check icon mb-2 text-green icon-lg"`)
	assert.Contains(t, h, `<a href="/invoice" class="btn btn-success w-100">Invoice</a>`)
}

func TestConfirmDeleteAndPassword(t *testing.T) {
	h := render(t, modals.ConfirmDelete(modals.ConfirmDeleteProps{Frame: modals.Frame{ID: "del"}, Items: []string{"a.txt"}, Check: "I understand"}))
	assert.Contains(t, h, `disabled data-bs-dismiss="modal" id="del-button"`)
	assert.Contains(t, h, `data-enables="del-button"`)
	assert.Contains(t, h, `<li>• a.txt</li>`)
	h = render(t, modals.ChangePassword(modals.ChangePasswordProps{Frame: modals.Frame{ID: "pw"}, Action: "/password"}))
	assert.Contains(t, h, `<form action="/password" method="post" data-password-match="pw-new">`)
	assert.Contains(t, h, `data-bs-input="#pw-new"`)
}

func TestModalVariants(t *testing.T) {
	h := render(t, modal.Modal(modal.Props{ID: "m", NoHeader: true, NoClose: true, Sections: true, Size: modal.SizeFullWidth, Show: true}))
	assert.Contains(t, h, `class="modal show"`)
	assert.Contains(t, h, `class="modal-dialog modal-full-width"><div class="modal-content"></div>`)
	assert.NotContains(t, h, `aria-hidden`)
	h = render(t, modal.Dialog(modal.DialogProps{ID: "d", Title: "Hi"}))
	assert.Contains(t, h, `<dialog id="d" class="p-0 border-0 bg-transparent" style="width: 28rem; max-width: calc(100% - 2rem);"><form method="dialog" class="card">`)
}
