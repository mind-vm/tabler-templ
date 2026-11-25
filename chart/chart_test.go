package chart

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, cs ...templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := templ.Join(cs...).Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestChart(t *testing.T) {
	p := Props{ID: "rev", Label: "Revenue", Type: Area, Size: SizeSmall, Sparkline: true, YFormatter: Thousands,
		Series: []Serie{{Name: "Revenue", Color: "primary", Data: []float64{1, 2}}, {Name: "Other", Data: []float64{3, 4}}}}
	out := render(t, Chart(p), Chart(Props{ID: "two"}))
	for _, want := range []string{`id="chart-rev"`, `class="position-relative chart-sm"`, `aria-label="Revenue"`, `style="height:40px;"`,
		`color-mix(in srgb, transparent, var(--tblr-primary) 100%)`, `var(--tblr-chart-2, var(--tblr-green))`,
		`&#34;yAxis.0.axisLabel.formatter&#34;:&#34;thousands&#34;`, `echarts.init`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s in\n%s", want, out)
		}
	}
	if n := strings.Count(out, "window.tablerCharts"); n != 2 { // the guard and the assignment, in one script
		t.Errorf("script rendered %d times", n/2)
	}
}

// option is the ECharts option of a rendered chart.
func option(t *testing.T, p Props) map[string]any {
	t.Helper()
	m := regexp.MustCompile(`data-tabler-chart="([^"]*)"`).FindStringSubmatch(render(t, Chart(p)))
	var c struct {
		Option map[string]any `json:"option"`
	}
	if err := json.Unmarshal([]byte(html.UnescapeString(m[1])), &c); err != nil {
		t.Fatal(err)
	}
	return c.Option
}

func TestStyleAttributeMerges(t *testing.T) {
	out := render(t, Chart(Props{ID: "g", Height: 5, Attributes: templ.Attributes{"style": "width: 10rem", "data-x": "1"}}))
	if strings.Count(out, "style=") != 1 || !strings.Contains(out, `style="height:80px;width: 10rem;"`) || !strings.Contains(out, `data-x="1"`) {
		t.Errorf("attributes: %s", out)
	}
}

func TestOptionsMergeOverBuilt(t *testing.T) {
	o := option(t, Props{ID: "m", Type: Line, Series: []Serie{{Name: "a", Data: []float64{1}}},
		Options: map[string]any{"grid": map[string]any{"top": 40}, "legend": map[string]any{"show": true}}})
	if g := o["grid"].(map[string]any); g["top"] != 40.0 || g["containLabel"] != true {
		t.Errorf("grid = %v", g)
	}
}

func TestDatetimeLabelsCoverLongestSerie(t *testing.T) {
	o := option(t, Props{ID: "d", Type: Line, Datetime: true, StartDate: "2026-01-01",
		Series: []Serie{{Name: "a", Data: []float64{1}}, {Name: "b", Data: []float64{1, 2, 3}}}})
	b := o["series"].([]any)[1].(map[string]any)["data"].([]any)
	if last := b[2].([]any); last[0] != "2026-01-03" {
		t.Errorf("b = %v", b)
	}
}

func TestPieIsOneSlicePerSerie(t *testing.T) {
	o := option(t, Props{ID: "p", Type: Donut, Series: []Serie{{Name: "A", Data: []float64{3}}, {Name: "B", Data: []float64{7}}}})
	s := o["series"].([]any)[0].(map[string]any)
	if s["type"] != "pie" || s["radius"].([]any)[0] != "55%" {
		t.Fatalf("series = %v", s)
	}
	if d := s["data"].([]any)[1].(map[string]any); d["name"] != "B" || d["value"] != 7.0 {
		t.Errorf("slice = %v", d)
	}
}

// A scatter with numbers as x draws them on a value axis (#89); with names as x it keeps the category axis.
func TestScatterXAxis(t *testing.T) {
	pts := func(xs ...any) []Serie {
		s := Serie{Name: "b"}
		for i, x := range xs {
			s.Points = append(s.Points, Point{X: x, Y: float64(i)})
		}
		return []Serie{s}
	}
	for _, c := range []struct {
		name   string
		series []Serie
		want   string
	}{
		{"numbers", pts(1.0, 4.0), "value"},
		{"ints", pts(1, 4), "value"},
		{"names", pts("Jan", "Feb"), "category"},
		{"mixed", pts(1.0, "Feb"), "category"},
	} {
		x := option(t, Props{ID: "s", Type: Scatter, Series: c.series})["xAxis"].([]any)[0].(map[string]any)
		if x["type"] != c.want {
			t.Errorf("%s: x axis = %v, want %s", c.name, x["type"], c.want)
		}
	}
	x := option(t, Props{ID: "s", Type: Scatter, Categories: []string{"a", "b"}, Series: pts(0, 1)})["xAxis"].([]any)[0].(map[string]any)
	if x["type"] != "category" {
		t.Errorf("with categories: x axis = %v", x["type"])
	}
}

func TestRawProps(t *testing.T) {
	o := option(t, RawProps(Props{ID: "r"}, Gauge(GaugeOptions{Value: 58, Color: "red"}), nil))
	s := o["series"].([]any)[0].(map[string]any)
	if s["type"] != "gauge" || s["data"].([]any)[0].(map[string]any)["value"] != 58.0 {
		t.Errorf("gauge = %v", s)
	}
}

// whole is the whole config of a rendered chart: option, formatters and unit.
func whole(t *testing.T, p Props) map[string]any {
	t.Helper()
	m := regexp.MustCompile(`data-tabler-chart="([^"]*)"`).FindStringSubmatch(render(t, Chart(p)))
	var c map[string]any
	if err := json.Unmarshal([]byte(html.UnescapeString(m[1])), &c); err != nil {
		t.Fatal(err)
	}
	return c
}

// Unit goes to the value axis, the tooltip and the data labels; on horizontal bars the value axis is x.
func TestUnit(t *testing.T) {
	bars := []Serie{{Name: "a", Data: []float64{1, 2}}, {Name: "b", Data: []float64{3, 4}}}
	paths := func(p Props) string {
		u, _ := whole(t, p)["unit"].(map[string]any)
		if u == nil {
			return ""
		}
		var s []string
		for _, x := range u["paths"].([]any) {
			s = append(s, x.(string))
		}
		return u["text"].(string) + " " + strings.Join(s, ",")
	}
	for _, c := range []struct {
		p    Props
		want string
	}{
		{Props{ID: "v", Series: bars, Unit: "€", DataLabels: true},
			"€ tooltip.valueFormatter,yAxis.0.axisLabel.formatter,series.0.label.formatter,series.1.label.formatter"},
		{Props{ID: "h", Series: bars, Unit: " kWh ", Horizontal: true, YFormatter: Thousands}, "kWh tooltip.valueFormatter,xAxis.0.axisLabel.formatter"},
		{Props{ID: "s", Type: Area, Series: bars, Unit: "h", Sparkline: true, HideTooltip: true}, ""},
		{Props{ID: "d", Type: Donut, Series: bars, Unit: "€"}, "€ tooltip.valueFormatter"},
		{Props{ID: "p", Series: bars, Unit: "€", YFormatter: Percent}, ""},
		{Props{ID: "r", Type: RadialBar, Series: bars, Unit: "€"}, ""},
		{Props{ID: "x", Series: bars, Unit: `<img src=x onerror="alert(1)">{a}`}, "img src=x on tooltip.valueFormatter,yAxis.0.axisLabel.formatter"},
		{Props{ID: "n", Series: bars}, ""},
	} {
		if got := paths(c.p); got != c.want {
			t.Errorf("%s: got %q, want %q", c.p.ID, got, c.want)
		}
	}
	if f := whole(t, Props{ID: "h", Series: bars, Unit: "€", Horizontal: true, YFormatter: Thousands})["formatters"].(map[string]any); f["xAxis.0.axisLabel.formatter"] != "thousands" {
		t.Errorf("formatters = %v", f)
	}
}
