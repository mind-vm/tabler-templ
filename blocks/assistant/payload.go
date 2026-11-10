package assistant

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mind-vm/tabler-templ/chart"
)

// Normalize turns a payload into JSON's shapes — float64, []any and map[string]any — which is what the
// blocks read. A payload straight from Go (or from a JS runtime such as goja: int64, typed slices) can
// take any shape; rendering one is simpler than anticipating every one.
func Normalize(payload any) map[string]any {
	b, err := json.Marshal(payload)
	if err != nil {
		return map[string]any{}
	}
	var m map[string]any
	if json.Unmarshal(b, &m) != nil || m == nil {
		return map[string]any{}
	}
	return m
}

func str(p map[string]any, k string) string {
	s, _ := p[k].(string)
	return s
}

func strs(p map[string]any, k string) []string {
	l, _ := p[k].([]any)
	out := make([]string, 0, len(l))
	for _, e := range l {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func rows(p map[string]any, k string) []map[string]any {
	l, _ := p[k].([]any)
	out := make([]map[string]any, 0, len(l))
	for _, e := range l {
		if m, ok := e.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

// numeric reports a column of numbers, right-aligned in a table. Numeric
// strings count, as they do for charts: a model that parses CSV in a script
// usually leaves "2210.00" a string, and it is still an amount.
func numeric(rs []map[string]any, col string) bool {
	seen := false
	for _, r := range rs {
		if r[col] == nil {
			continue
		}
		if _, ok := number(r[col]); !ok {
			return false
		}
		seen = true
	}
	return seen
}

func pretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(b)
}

func cellText(v any) string {
	switch v.(type) {
	case map[string]any, []any:
		b, _ := json.Marshal(v)
		return string(b)
	}
	return cell(v)
}

// joinAnswer is an answer as one line: a multi-select's labels comma-joined.
func joinAnswer(v any, t Texts) string {
	switch a := v.(type) {
	case []any:
		parts := make([]string, 0, len(a))
		for _, e := range a {
			parts = append(parts, cell(e))
		}
		return strings.Join(parts, ", ")
	case bool:
		if a {
			return t.Yes
		}
		return t.No
	}
	return cell(v)
}

// ChartProps maps a chart block's spec onto chart.Props. The spec is
// thinking-core's emitChart: variant, title, data (rows), x, series, and for some
// variants types, size. Every chart type wants its data in its own shape; this is
// where it gets it. A key no row has gives zeros, so check the spec first (the
// pack does).
//
// A variant this table does not know is drawn as bars rather than refused:
// the spec passed the pack, so its numbers are real, and a chart in the
// wrong shape still carries them.
func ChartProps(id string, p map[string]any) chart.Props {
	data := rows(p, "data")
	x := str(p, "x")
	series := strs(p, "series")
	props := chart.Props{ID: id, Label: str(p, "title"), Size: chart.SizeLarge, Legend: len(series) > 1}
	// Straight segments, not tabler-templ's smooth default: a spline between
	// monthly totals draws dips and peaks between the points that no row
	// holds, and every number on this page has to come from the data.
	props.StrokeCurve = "straight"

	switch str(p, "variant") {
	case "line":
		props.Type = chart.Line
		props.Categories, props.Series = categories(data, x), columns(data, series)
	case "area":
		props.Type = chart.Area
		props.Categories, props.Series = categories(data, x), columns(data, series)
	case "composite":
		types, _ := p["types"].(map[string]any)
		props.Type = chart.Line
		props.Categories, props.Series = categories(data, x), columns(data, series)
		// Each serie gets its own type and line width.
		props.StrokeWidth = make([]float64, len(series))
		for i := range props.Series {
			switch t, _ := types[series[i]].(string); t {
			case "line":
				props.Series[i].Type, props.StrokeWidth[i] = "line", 2
			case "area":
				props.Series[i].Type, props.StrokeWidth[i] = "area", 2
			default:
				props.Series[i].Type = "column"
			}
		}
	case "donut", "pie", "radialBar":
		props.Type = map[string]chart.Type{"donut": chart.Donut, "pie": chart.Pie, "radialBar": chart.RadialBar}[str(p, "variant")]
		props.Legend = true
		props.Series = pieSlices(data, x, first(series))
	case "radar":
		props.Type = chart.Radar
		props.Legend = true
		props.Categories, props.Series = categories(data, x), columns(data, series)
	case "scatter":
		props.Type = chart.Scatter
		props.Series = points(data, x, series, "")
	case "bubble":
		props.Type = chart.Bubble
		props.Series = points(data, x, series, str(p, "size"))
	case "funnel":
		props.Type = chart.Bar
		props.Funnel = true
		props.DataLabels = true
		props.Categories, props.Series = categories(data, x), columns(data, series[:min(1, len(series))])
		// A funnel colours each stage, not each serie; Colors is per point.
		props.Colors = make([]string, len(data))
		for i := range data {
			props.Colors[i] = palette[i%len(palette)]
		}
	case "sparkline":
		props.Type = chart.Area
		props.Sparkline = true
		props.Size, props.Height = chart.SizeDefault, 4
		props.Series = columns(data, series[:min(1, len(series))])
	case "heatmap":
		props.Type = chart.Heatmap
		props.Legend = false
		props.Series = calendar(data, x, first(series))
		props.Color = "primary"
	default:
		props.Type = chart.Bar
		props.Categories, props.Series = categories(data, x), columns(data, series)
	}
	return props
}

// palette is Tabler's chart palette as CSS the chart script resolves, for
// the variants that colour points rather than series.
var palette = []string{"var(--tblr-primary)", "var(--tblr-green)", "var(--tblr-azure)", "var(--tblr-purple)", "var(--tblr-orange)", "var(--tblr-pink)"}

func first(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

func categories(data []map[string]any, x string) []string {
	out := make([]string, len(data))
	for i, r := range data {
		out[i] = cell(r[x])
	}
	return out
}

func columns(data []map[string]any, names []string) []chart.Serie {
	out := make([]chart.Serie, len(names))
	for i, n := range names {
		vals := make([]float64, len(data))
		for j, r := range data {
			vals[j], _ = number(r[n])
		}
		out[i] = chart.Serie{Name: n, Data: vals}
	}
	return out
}

// pieSlices is one serie per row, which is how tabler-templ spells a pie.
func pieSlices(data []map[string]any, x, value string) []chart.Serie {
	out := make([]chart.Serie, len(data))
	for i, r := range data {
		v, _ := number(r[value])
		out[i] = chart.Serie{Name: cell(r[x]), Data: []float64{v}}
	}
	return out
}

func points(data []map[string]any, x string, names []string, size string) []chart.Serie {
	out := make([]chart.Serie, 0, len(names))
	for _, n := range names {
		pts := make([]chart.Point, 0, len(data))
		for _, r := range data {
			xv, okX := number(r[x])
			yv, okY := number(r[n])
			// A point without both coordinates has nowhere to go; plotting
			// it at zero would draw a value nobody computed.
			if !okX || !okY {
				continue
			}
			pt := chart.Point{X: xv, Y: yv}
			if size != "" {
				z, _ := number(r[size])
				pt.Z = z
			}
			pts = append(pts, pt)
		}
		out = append(out, chart.Serie{Name: n, Points: pts})
	}
	if size != "" {
		// One serie per value column is the bubble shape too, but the pack
		// declares a single y column for it.
		out = out[:min(1, len(out))]
	}
	return out
}

// calendar lays daily values out as GitHub-style weeks: one serie per
// weekday, one column per week: a heatmap of raw dates on one row is
// unreadable past a month.
func calendar(data []map[string]any, x, value string) []chart.Serie {
	byDay := map[string]float64{}
	weeks := map[string]bool{}
	for _, r := range data {
		d, err := time.Parse("2006-01-02", strings.TrimSpace(cell(r[x])))
		if err != nil {
			continue
		}
		v, _ := number(r[value])
		byDay[d.Format("2006-01-02")] += v
		weeks[weekOf(d).Format("2006-01-02")] = true
	}
	order := make([]string, 0, len(weeks))
	for w := range weeks {
		order = append(order, w)
	}
	slices.Sort(order)
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	out := make([]chart.Serie, 0, 7)
	// The heatmap stacks its rows bottom up; reversed, Monday is on top.
	for i := 6; i >= 0; i-- {
		pts := make([]chart.Point, len(order))
		for j, w := range order {
			start, _ := time.Parse("2006-01-02", w)
			day := start.AddDate(0, 0, i).Format("2006-01-02")
			pt := chart.Point{X: w}
			if v, ok := byDay[day]; ok {
				pt.Y = v
			}
			pts[j] = pt
		}
		out = append(out, chart.Serie{Name: days[i], Points: pts})
	}
	return out
}

func weekOf(d time.Time) time.Time {
	offset := (int(d.Weekday()) + 6) % 7 // Monday = 0
	return d.AddDate(0, 0, -offset)
}

// number reads a cell as a number. JSON gives float64; a numeric string is
// accepted because CSV parsed in a script often stays a string, and the
// model rarely remembers to convert before charting.
func number(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int64:
		return float64(n), true
	case int:
		return float64(n), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
		return f, err == nil
	}
	return 0, false
}

// cell formats a value for a table cell or a category label.
func cell(v any) string {
	switch n := v.(type) {
	case nil:
		return ""
	case string:
		return n
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(n)
	}
	return fmt.Sprint(v)
}
