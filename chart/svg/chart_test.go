package svg

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/chart"
)

var update = flag.Bool("update", false, "rewrite the golden files in testdata")

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

var (
	months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	solved = []float64{12, 19, 8, 15, 22, 18}
	opened = []float64{5, 7, 4, 9, 6, -3}
	two    = []chart.Serie{{Name: "Solved", Data: solved}, {Name: "Opened", Data: opened}}
	slices = []chart.Serie{{Name: "Direct", Data: []float64{40}}, {Name: "Search", Data: []float64{25}}, {Name: "Social", Data: []float64{35}}}
)

// goldens are the charts whose SVG is checked byte for byte against testdata/<name>.golden; run
// go test ./chart/svg -update after a deliberate change and look at the diff.
var goldens = map[string]Props{
	"bar":        {ID: "bar", Label: "Tickets", Categories: months, Series: two, Legend: true, Table: true},
	"bar-stack":  {Label: "Tickets", Categories: months, Series: two, Stacked: true, DataLabels: true},
	"bar-horiz":  {Label: "Tickets", Categories: months, Series: two[:1], Horizontal: true, DataLabels: true, Colors: []string{"green", "red", "#123456"}},
	"line":       {Label: "Revenue", Type: chart.Line, Categories: months, Series: []chart.Serie{{Name: "Revenue", Data: []float64{1200.5, 1900, 800, 1500, 2200, 1800}}}, StrokeCurve: "smooth", ShowMarkers: true, Format: German, AnnotationsY: []chart.Annotation{{Value: 1000.0, Label: "Target"}}},
	"area-stack": {Label: "Storage", Type: chart.Area, Categories: months, Series: []chart.Serie{{Name: "Media", Data: solved}, {Name: "Backups", Data: solved[:4]}}, Stacked: true, YFormatter: chart.SuffixK},
	"sparkline":  {Label: "Trend", Type: chart.Area, Sparkline: true, Size: chart.SizeSmall, Series: two[:1]},
	"donut":      {Label: "Traffic", Type: chart.Donut, Series: slices, Legend: true, DataLabels: true, Table: true},
	"pie":        {Label: "Traffic", Type: chart.Pie, Series: slices[:2]},
	"radial":     {Label: "Done", Type: chart.RadialBar, Series: []chart.Serie{{Name: "Done", Data: []float64{72.5}}}, Format: German},
	"gauge":      {Label: "Load", Type: chart.RadialBar, StartAngle: -90, EndAngle: 90, Series: []chart.Serie{{Name: "Load", Color: "orange", Data: []float64{42}}}},
	"rings":      {Label: "Goals", Type: chart.RadialBar, Series: slices, Legend: true},
}

func TestGolden(t *testing.T) {
	for name, p := range goldens {
		t.Run(name, func(t *testing.T) {
			got := render(t, Chart(p))
			// The sheet is the same for every chart; the goldens hold the markup.
			got = strings.Replace(got, "<style>"+CSS+"</style>", "", 1)
			path := filepath.Join("testdata", name+".golden")
			if *update {
				if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				t.Errorf("%s differs from %s (go test ./chart/svg -update):\n%s", name, path, got)
			}
		})
	}
}

func TestNoScriptOneSheet(t *testing.T) {
	html := render(t, templ.Join(Chart(goldens["bar"]), Chart(goldens["donut"])))
	if strings.Contains(html, "<script") {
		t.Error("an SVG chart renders a script")
	}
	if n := strings.Count(html, "<style>"); n != 1 {
		t.Errorf("the sheet is rendered %d times", n)
	}
	for _, want := range []string{`id="chart-bar"`, `role="img" aria-label="Tickets"`, `<title>Jan` + "\n" + `Solved: 12` + "\n" + `Opened: 5</title>`,
		`--s0: var(--tblr-chart-1, var(--tblr-primary));`, `<table class="visually-hidden"><caption>Tickets</caption>`, `<td>-3</td>`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
}

func TestNumber(t *testing.T) {
	for _, c := range []struct {
		f    Format
		in   float64
		want string
	}{
		{English, 12400.5, "12,400.5"}, {German, 12400.5, "12.400,5"}, {German, -1234567.126, "-1.234.567,13"},
		{English, 0.1 + 0.2, "0.3"}, {English, -0.001, "0"}, {Number(0, ",", " "), 999999.6, "1 000 000"}, {English, 100, "100"},
	} {
		if got := c.f(c.in); got != c.want {
			t.Errorf("%v: got %q, want %q", c.in, got, c.want)
		}
	}
}

func TestScale(t *testing.T) {
	s := niceScale(0, 22, 5, nil, nil)
	if s.lo != 0 || s.hi != 30 || s.step != 10 {
		t.Errorf("0..22 → %+v", s)
	}
	s = niceScale(-3, 41, 5, nil, nil)
	if s.lo != -10 || s.hi != 50 || s.step != 10 {
		t.Errorf("-3..41 → %+v", s)
	}
	m := 100.0
	if s = niceScale(0, 0, 5, nil, &m); s.hi != 100 || s.lo != 0 {
		t.Errorf("YMax → %+v", s)
	}
}

// The limits hold: series and categories past them are not drawn.
func TestLimits(t *testing.T) {
	var series []chart.Serie
	data := make([]float64, MaxCategories+50)
	for range MaxSeries + 3 {
		series = append(series, chart.Serie{Name: "s", Data: data})
	}
	v := Props{Type: chart.Line, Series: series, Table: true}.view()
	if v.Colors != MaxSeries || len(v.Rows) != MaxCategories {
		t.Errorf("%d series, %d rows", v.Colors, len(v.Rows))
	}
}

func TestFromChartAndSupported(t *testing.T) {
	p := FromChart(chart.Props{ID: "x", Type: chart.Bar, Funnel: true, Legend: true, Series: two})
	if !p.Horizontal || !p.Legend || p.ID != "x" || len(p.Series) != 2 {
		t.Errorf("FromChart: %+v", p)
	}
	if !Supported(chart.Donut) || Supported(chart.Heatmap) {
		t.Error("Supported")
	}
}
