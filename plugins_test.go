package tablertempl_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/mind-vm/tabler-templ/datepicker"
	"github.com/mind-vm/tabler-templ/inputmask"
	"github.com/mind-vm/tabler-templ/rangeslider"
	"github.com/mind-vm/tabler-templ/selectbox"
	"github.com/mind-vm/tabler-templ/signature"
	"github.com/mind-vm/tabler-templ/wysiwyg"
)

func TestPlugins(t *testing.T) {
	ctx := templ.InitializeContext(context.Background())
	html := render(t, ctx, templ.Join(
		datepicker.Input(datepicker.Props{ID: "due", Name: "due", Layout: datepicker.LayoutIcon, Mode: datepicker.ModeRange, Selected: []string{"2026-06-10", "2026-06-18"}}),
		datepicker.Input(datepicker.Props{ID: "cal", Inline: true, Name: "day"}),
		inputmask.Input(inputmask.Props{ID: "zip", Mask: "00000-000", Visible: true}),
		selectbox.Select(selectbox.Props{ID: "country", Options: []selectbox.Option{{Value: "pl", Text: "Poland", Selected: true, Indicator: templ.Raw(`<span class="flag flag-xs flag-country-pl"></span>`)}}}),
		selectbox.Select(selectbox.Props{ID: "tags", Multiple: true, Create: true}),
		rangeslider.Slider(rangeslider.Props{ID: "price", Name: "price", Values: []float64{60, 90}, Step: 10, Label: "Price"}),
		signature.Pad(signature.Props{ID: "sig", Name: "sig", Clear: true}),
		wysiwyg.Editor(wysiwyg.Props{ID: "body", Value: "Hello, <b>Tabler</b>!"}),
	))
	for _, want := range []string{
		`<div class="input-icon"><input type="text" class="form-control" id="due" name="due" placeholder="Select a date" autocomplete="off"`,
		`data-bs-selection-mode="multiple-ranged"`, `data-bs-selected-dates="[&#34;2026-06-10&#34;,&#34;2026-06-18&#34;]"`,
		`<div id="cal" data-bs-inline="true" data-bs-toggle="datepicker"><input type="hidden" name="day"></div>`,
		`data-mask="00000-000" data-mask-visible="true" placeholder="00000-000"`,
		`data-custom-properties="&lt;span class=&#34;flag flag-xs flag-country-pl&#34;&gt;&lt;/span&gt;" selected>Poland</option>`,
		`&#34;create&#34;:true`, `multiple`,
		`&#34;connect&#34;:[false,true,false]`, `&#34;aria-label&#34;:&#34;Price maximum&#34;`, `<input type="hidden" name="price" value="90" data-tabler-range-value="price">`,
		`<canvas id="sig" width="400" height="400" class="signature-canvas" data-tabler-signature=""></canvas>`, `data-signature-clear`,
		`>Hello, &lt;b&gt;Tabler&lt;/b&gt;!</textarea>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s", want)
		}
	}
	// Each plugin's start script once, however many fields.
	if n := strings.Count(html, "querySelectorAll('select[data-tabler-select]')"); n != 1 {
		t.Errorf("select script %d times", n)
	}
	if n := strings.Count(html, `window.tabler.Datepicker`); n != 2 {
		t.Errorf("datepicker script %d times", n/2)
	}
	if t.Failed() {
		t.Log(html)
	}
}
