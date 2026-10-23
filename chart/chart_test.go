package chart

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestChart(t *testing.T) {
	var b bytes.Buffer
	p := Props{ID: "rev", Label: "Revenue", Type: Area, Size: SizeSmall, Sparkline: true, YFormatter: Thousands,
		Series: []Serie{{Name: "Revenue", Color: "primary", Data: []float64{1, 2}}, {Name: "Other", Data: []float64{3, 4}}}}
	if err := templ.Join(Chart(p), Chart(Props{ID: "two"})).Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	html := b.String()
	for _, want := range []string{`id="chart-rev"`, `class="position-relative chart-sm"`, `aria-label="Revenue"`,
		`--chart-color-0: color-mix(in srgb, transparent, var(--tblr-primary) 100%)`, `var(--tblr-chart-2, var(--tblr-green))`,
		`&#34;height&#34;:40`, `&#34;sparkline&#34;:{&#34;enabled&#34;:true}`, `&#34;yaxis.labels.formatter&#34;:&#34;thousands&#34;`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in\n%s", want, html)
		}
	}
	if n := strings.Count(html, "window.tablerCharts"); n != 2 { // the guard and the assignment, in one script
		t.Errorf("script rendered %d times", n/2)
	}
}
