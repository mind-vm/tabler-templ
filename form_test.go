package tablertempl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/button"
	"github.com/mind-vm/tabler-templ/form"
)

func TestFormControls(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<input class="form-control form-control-lg form-control-rounded is-valid is-valid-lite" value="" list="l" autocomplete="off" name="q">`,
		render(t, ctx, form.Input(form.InputProps{Name: "q", Size: form.SizeLarge, Rounded: true, State: form.StateValid, Lite: true, List: "l", Autocomplete: "off"})))
	assert.Contains(t, render(t, ctx, form.Input(form.InputProps{Type: "color", Value: "#066fd1"})), `class="form-control form-control-color"`)
	h := render(t, ctx, form.Select(form.SelectProps{Multiple: true, SelectedValues: []string{"b"}, Groups: []form.OptionGroup{{Label: "G", Options: []form.Option{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}}}}))
	assert.Contains(t, h, `<select class="form-select" multiple>`)
	assert.Contains(t, h, `<optgroup label="G"><option value="a">A</option><option value="b" selected>B</option></optgroup>`)
	assert.Contains(t, render(t, ctx, form.Check(form.CheckProps{Switch: true, Large: true, LabelOn: "On", LabelOff: "Off"})),
		`<label class="form-check form-switch form-switch-lg">`)
	h = render(t, ctx, mit(form.Field(form.FieldProps{Label: "Bio", Description: "56/100", Help: "Tell us"}), "<textarea></textarea>"))
	assert.Contains(t, h, `<span class="form-help" data-bs-content="Tell us" data-bs-placement="top" data-bs-toggle="popover">?</span>`)
	assert.Contains(t, h, `<span class="form-label-description">56/100</span>`)
}

func TestFormGroups(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="input-group input-group-flat">x</div>`, render(t, ctx, mit(form.Group(form.GroupProps{Flat: true}), "x")))
	h := render(t, ctx, mit(form.InputIcon(form.InputIconProps{Icon: "ti-search", Start: true}), "<input>"))
	assert.Equal(t, `<div class="input-icon"><span class="input-icon-addon"><i class="ti ti-search icon" aria-hidden="true"></i></span><input></div>`, h)
	h = render(t, ctx, form.PasswordToggle(form.PasswordToggleProps{Target: "pw"}))
	assert.Contains(t, h, `data-password-toggle="pw" data-label-show="Show password" data-label-hide="Hide password"`)
	h = render(t, ctx, form.SelectItem(form.SelectItemProps{Name: "lang", Value: "go", Label: "Go", Radio: true, Checked: true}))
	assert.Equal(t, `<label class="form-selectgroup-item"><input type="radio" name="lang" value="go" class="form-selectgroup-input" checked> <span class="form-selectgroup-label">Go</span></label>`, h)
	assert.Contains(t, render(t, ctx, form.ColorCheck(form.ColorCheckProps{Name: "c", Value: "lime", Color: "lime", Radio: true})),
		`<span class="form-colorinput-color bg-lime"></span>`)
	assert.Contains(t, render(t, ctx, form.ImageCheck(form.ImageCheckProps{Name: "i", Value: "1", Src: "/a.jpg"})), `<img src="/a.jpg" alt="" class="form-imagecheck-image">`)
}

func TestFormWidgets(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, form.OTP(form.OTPProps{Name: "code", Length: 6, Groups: []int{3, 3}, Separator: "-", Connected: true}))
	assert.Equal(t, `<div class="otp otp-connected" data-bs-toggle="otp" data-bs-length="6" data-bs-groups="[3,3]" data-bs-separator="-"><input type="text" name="code" maxlength="6"></div>`, h)
	h = render(t, ctx, form.NumberInput(form.NumberInputProps{Value: "2", Min: "1", Max: "5", Width: "7rem", AriaLabel: "Qty"}))
	assert.Equal(t, `<input type="text" class="form-control" data-bs-toggle="number-input" value="2" min="1" max="5" style="--tblr-number-input-input-width: 7rem;" aria-label="Qty">`, h)
	assert.Contains(t, render(t, ctx, form.NumberInput(form.NumberInputProps{DecrementIcon: "ti-chevron-down"})), `data-bs-number-input-action="decrement"`)
	h = render(t, ctx, form.Strength(form.StrengthProps{Segments: 3, Messages: map[string]string{"weak": "Too easy"}}))
	assert.Contains(t, h, `<div class="strength" data-bs-strength data-bs-messages="{&#34;weak&#34;:&#34;Too easy&#34;}" aria-label="Password strength"><span class="strength-segment"></span><span class="strength-segment"></span><span class="strength-segment"></span></div><div class="strength-text"></div>`)
	assert.Equal(t, `<input type="range" class="form-range" min="0" max="100" value="40" name="r">`, render(t, ctx, form.Range(form.RangeProps{Name: "r", Value: 40})))
	assert.Contains(t, render(t, ctx, mit(form.Floating(form.FloatingProps{Label: "Email", For: "e"}), "<input>")), `<div class="form-floating"><input><label for="e">Email</label></div>`)
	assert.Contains(t, render(t, ctx, mit(form.Fieldset(form.FieldsetProps{Legend: "Contact"}), "x")), `<fieldset class="form-fieldset"><legend class="form-label">Contact</legend>x</fieldset>`)
}

func TestButtonGroup(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, `<div class="btn-group-vertical" role="group">x</div>`, render(t, ctx, mit(button.Group(button.GroupProps{Vertical: true}), "x")))
	assert.Equal(t, `<input type="radio" class="btn-check" name="r" id="r1" autocomplete="off" checked> <label for="r1" class="btn">One</label>`,
		render(t, ctx, button.Check(button.CheckProps{ID: "r1", Name: "r", Label: "One", Radio: true, Checked: true})))
}
