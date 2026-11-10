// Package svg draws simple charts on the server as inline SVG: bar (vertical, horizontal, stacked), line,
// area, pie, donut and radial bar (gauge). It needs no chart library and no JavaScript, so a page with only
// these charts loads no chart library (see the README, "SVG charts without JavaScript").
//
//	@svg.Chart(svg.Props{
//		Label: "Revenue per month", Type: chart.Area, Legend: true, Format: svg.German,
//		Categories: []string{"Jan", "Feb", "Mär"},
//		Series:     []chart.Serie{{Name: "Umsatz", Color: "primary", Data: []float64{1200, 1850.5, 1640}}},
//	})
//
// The fields have the names and meaning of chart.Props; FromChart converts one. The colors are Tabler's
// CSS variables (the chart palette --tblr-chart-1 … 5 by default), so the chart follows the theme and dark
// mode. The plot scales to the width of its container: a cartesian chart stretches its SVG
// (preserveAspectRatio="none", strokes that do not scale) and puts axis and value labels in HTML over it,
// so the text keeps its size; a round chart keeps its aspect ratio (viewBox) and fits the height.
//
// Every bar, slice and ring carries a <title> (the browser's tooltip), and a category of a cartesian chart
// shows all its values in one tooltip. The chart is one image for screen readers (role="img" with Label);
// Table adds the data as a visually hidden table.
//
// Limits: at most MaxSeries series and MaxCategories categories (points per serie); more are left out.
// Serie.Points (heatmap, scatter …) and the chart types other than those above are not drawn: use
// Supported, or chart.Chart for them.
package svg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/chart"
)

// Limits of a chart; series and categories beyond them are not drawn.
const (
	MaxSeries     = 12
	MaxCategories = 400
)

// Format writes a number, e.g. for a German page "1.234,5".
type Format func(float64) string

// Number formats with up to maxDecimals decimals (trailing zeros dropped), decimal as the decimal
// separator and group between groups of three digits ("" for none).
func Number(maxDecimals int, decimal, group string) Format {
	return func(v float64) string {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return ""
		}
		s := strconv.FormatFloat(v, 'f', max(maxDecimals, 0), 64)
		if strings.Contains(s, ".") {
			s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
		}
		neg := strings.HasPrefix(s, "-")
		s = strings.TrimPrefix(s, "-")
		whole, frac, _ := strings.Cut(s, ".")
		if group != "" {
			for i := len(whole) - 3; i > 0; i -= 3 {
				whole = whole[:i] + group + whole[i:]
			}
		}
		if frac != "" {
			whole += decimal + frac
		}
		if neg && strings.Trim(whole, "0"+decimal+group) != "" {
			whole = "-" + whole
		}
		return whole
	}
}

var (
	English = Number(2, ".", ",") // 12,400.5
	German  = Number(2, ",", ".") // 12.400,5
)

// Props describe a chart. The fields shared with chart.Props mean the same.
type Props struct {
	// ID makes the element chart-<ID>.
	ID         string
	Label      string // what the chart shows: the image's name for screen readers, and the table's caption
	Class      string
	Attributes templ.Attributes
	Type       chart.Type // Bar (default), Line, Area, Pie, Donut or RadialBar
	Height     float64    // in rem, default 10 (without the legend)
	Size       chart.Size
	// Series: in a pie, donut or radial bar chart each serie is one slice or ring with its value in
	// Data[0] (a percentage in a radial bar). Serie.Type ("column", "line", "area") mixes bars, lines and
	// areas in a vertical chart.
	Series     []chart.Serie
	Categories []string
	Stacked    bool // stacked bars and areas
	Horizontal bool // horizontal bars
	// Sparkline draws only the plot: no axes, grid, labels or padding.
	Sparkline   bool
	Legend      bool
	DataLabels  bool // the values on the bars, points or slices
	ShowMarkers bool // a dot on every point of a line or area
	HideGrid    bool
	Color       string   // the color of every serie without its own (a Tabler color or any CSS color)
	Colors      []string // CSS colors, one per serie; with one serie and more colors, one per bar
	// StrokeCurve: "straight" (default), "smooth" (a monotone curve: no peaks the data does not have)
	// or "stepline". StrokeWidth and StrokeDash are per serie.
	StrokeCurve string
	StrokeWidth []float64
	StrokeDash  []float64
	// YMin and YMax fix the value axis; by default it runs over the data (bars and areas from zero) to
	// round numbers.
	YMin, YMax   *float64
	YFormatter   chart.Formatter // Percent and SuffixK append "%" and "K" to the formatted number
	AnnotationsY []chart.Annotation
	// Format writes the numbers of axes, labels, tooltips and table; default English.
	Format Format
	// Table adds the data as a table that only screen readers see.
	Table bool
	// StartAngle and EndAngle (degrees, 0 at the top, clockwise) bound a radial bar's arc; default a full
	// circle. -90 and 90 make a half gauge, -120 and 120 Tabler's radial.
	StartAngle, EndAngle float64
	// Hollow is the hole of a donut or radial bar in percent of its radius, default 60 (donut) or 65.
	Hollow float64
}

// FromChart converts chart.Props. Fields the SVG charts do not have (toolbar, zoom, brush, y axes,
// x annotations, Options …) are dropped.
func FromChart(c chart.Props) Props {
	return Props{ID: c.ID, Label: c.Label, Class: c.Class, Attributes: c.Attributes, Type: c.Type, Height: c.Height,
		Size: c.Size, Series: c.Series, Categories: c.Categories, Stacked: c.Stacked, Horizontal: c.Horizontal || c.Funnel,
		Sparkline: c.Sparkline, Legend: c.Legend, DataLabels: c.DataLabels, ShowMarkers: c.ShowMarkers, HideGrid: c.HideGrid,
		Color: c.Color, Colors: c.Colors, StrokeCurve: c.StrokeCurve, StrokeWidth: c.StrokeWidth, StrokeDash: c.StrokeDash,
		YMax: c.YMax, YFormatter: c.YFormatter, AnnotationsY: c.AnnotationsY}
}

// Supported tells whether Chart draws charts of type t.
func Supported(t chart.Type) bool {
	switch t {
	case "", chart.Bar, chart.Line, chart.Area, chart.Pie, chart.Donut, chart.RadialBar:
		return true
	}
	return false
}

func (p Props) typ() chart.Type {
	if p.Type == "" {
		return chart.Bar
	}
	return p.Type
}

func (p Props) round() bool {
	t := p.typ()
	return t == chart.Pie || t == chart.Donut || t == chart.RadialBar
}

func (p Props) height() float64 {
	switch p.Size {
	case chart.SizeSmall:
		return 2.5
	case chart.SizeLarge:
		return 15
	}
	if p.Height > 0 {
		return p.Height
	}
	return 10
}

// plain is v in Format, without YFormatter's suffix.
func (p Props) plain(v float64) string {
	if p.Format == nil {
		return English(v)
	}
	return p.Format(v)
}

func (p Props) format(v float64) string {
	s := p.plain(v)
	switch p.YFormatter {
	case chart.Percent:
		s += "%"
	case chart.SuffixK:
		s += "K"
	}
	return s
}

func (p Props) series() []chart.Serie {
	s := p.Series
	if len(s) > MaxSeries {
		s = s[:MaxSeries]
	}
	return s
}

// palette is Tabler's chart palette, for Tabler versions without --tblr-chart-1 … 5.
var palette = [5]string{"var(--tblr-primary)", "var(--tblr-green)", "var(--tblr-azure)", "var(--tblr-purple)", "var(--tblr-orange)"}

func colorVar(c string) string {
	if c == "" || strings.ContainsAny(c, "#(") {
		return c
	}
	return "var(--tblr-" + c + ")"
}

// color is the CSS color of serie (or slice) i.
func (p Props) color(i int) string {
	if i < len(p.Colors) && p.Colors[i] != "" {
		return colorVar(p.Colors[i])
	}
	base, opacity := "", ""
	if i < len(p.Series) {
		base, opacity = colorVar(p.Series[i].Color), p.Series[i].ColorOpacity
	}
	if base == "" {
		base = colorVar(p.Color)
	}
	if base == "" {
		base = fmt.Sprintf("var(--tblr-chart-%d, %s)", i%5+1, palette[i%5])
	}
	if opacity != "" && opacity != "100%" {
		return fmt.Sprintf("color-mix(in srgb, transparent, %s %s)", base, opacity)
	}
	return base
}

// style is the colors as custom properties --s0, --s1 … on the chart element.
func (p Props) style(n int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "--chart-svg-h: %srem;", num(p.height()))
	for i := range n {
		fmt.Fprintf(&b, "--s%d: %s;", i, p.color(i))
	}
	return b.String()
}

// num writes a coordinate with at most two decimals.
func num(f float64) string {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	if s == "-0" {
		return "0"
	}
	return s
}

func value(s chart.Serie, i int) float64 {
	if i < len(s.Data) {
		return s.Data[i]
	}
	return math.NaN()
}

// The view is what the template draws, all positions computed.

type tick struct {
	Pos  float64 // percent along the axis
	Text string
	Alt  bool // hidden on narrow charts
}

type shape struct {
	D, Class, Style, Title string
}

type label struct {
	Style, Text, Class string
}

type text struct {
	X, Y        float64
	Text, Class string
}

type legendItem struct{ Name, Color, Value string }

type view struct {
	Round   bool
	ViewBox string
	Inset   string // the plot's position in the box
	YTicks  []tick // labels left of the plot (values, or categories of horizontal bars)
	XTicks  []tick // labels below the plot
	Shapes  []shape
	Labels  []label
	Texts   []text
	Legend  []legendItem
	Head    []string
	Rows    [][]string
	Colors  int
}

func (p Props) view() view {
	if p.round() {
		return p.roundView()
	}
	return p.cartesianView()
}

// scale is the value axis: lo to hi in steps.
type scale struct{ lo, hi, step float64 }

func niceNum(x float64, round bool) float64 {
	exp := math.Floor(math.Log10(x))
	f := x / math.Pow(10, exp)
	var nf float64
	switch {
	case round && f < 1.5, !round && f <= 1:
		nf = 1
	case round && f < 3, !round && f <= 2:
		nf = 2
	case round && f < 7, !round && f <= 5:
		nf = 5
	default:
		nf = 10
	}
	return nf * math.Pow(10, exp)
}

func niceScale(lo, hi float64, ticks int, fixLo, fixHi *float64) scale {
	if fixLo != nil {
		lo = *fixLo
	}
	if fixHi != nil {
		hi = *fixHi
	}
	if hi <= lo {
		if lo == 0 {
			hi = 1
		} else {
			lo, hi = lo-math.Abs(lo)/2, hi+math.Abs(hi)/2
		}
	}
	step := niceNum(niceNum(hi-lo, false)/float64(ticks-1), true)
	s := scale{math.Floor(lo/step) * step, math.Ceil(hi/step) * step, step}
	if fixLo != nil {
		s.lo = *fixLo
	}
	if fixHi != nil {
		s.hi = *fixHi
	}
	return s
}

func (s scale) pos(v float64) float64 { return (v - s.lo) / (s.hi - s.lo) * 100 }

func (s scale) ticks() []float64 {
	var t []float64
	for v := s.lo; v <= s.hi+s.step/1e6; v += s.step {
		t = append(t, math.Round(v/s.step)*s.step)
	}
	return t
}

type kind int

const (
	kBar kind = iota
	kLine
	kArea
)

func (p Props) kind(s chart.Serie) kind {
	if p.Horizontal {
		return kBar
	}
	switch s.Type {
	case "column", "bar":
		return kBar
	case "line":
		return kLine
	case "area":
		return kArea
	}
	switch p.typ() {
	case chart.Line:
		return kLine
	case chart.Area:
		return kArea
	}
	return kBar
}

func runes(s string) int { return utf8.RuneCountInString(s) }

func (p Props) cartesianView() view {
	series := p.series()
	n := 0
	for _, s := range series {
		n = max(n, len(s.Data))
	}
	n = max(n, len(p.Categories))
	n = min(n, MaxCategories)
	cat := func(i int) string {
		if i < len(p.Categories) {
			return p.Categories[i]
		}
		return strconv.Itoa(i + 1)
	}
	v := view{Colors: len(series)}

	// The range of the values: stacked series add up, bars and areas start at zero.
	kinds := make([]kind, len(series))
	bars, anyBar, zero := 0, false, false
	for i, s := range series {
		kinds[i] = p.kind(s)
		if kinds[i] == kBar {
			bars++
			anyBar = true
		}
		zero = zero || kinds[i] != kLine
	}
	lo, hi := math.Inf(1), math.Inf(-1)
	pos, neg := make([]float64, n), make([]float64, n)
	for i, s := range series {
		for j := range n {
			x := value(s, j)
			if math.IsNaN(x) {
				continue
			}
			if p.Stacked && kinds[i] != kLine {
				if x >= 0 {
					pos[j] += x
					x = pos[j]
				} else {
					neg[j] += x
					x = neg[j]
				}
			}
			lo, hi = min(lo, x), max(hi, x)
		}
	}
	if math.IsInf(lo, 1) {
		lo, hi = 0, 1
	}
	if zero {
		lo, hi = min(lo, 0), max(hi, 0)
	}
	nt := 5
	if p.height() < 5 {
		nt = 3
	}
	sc := niceScale(lo, hi, nt, p.YMin, p.YMax)
	base := sc.pos(math.Max(sc.lo, math.Min(0, sc.hi)))
	if sc.lo > 0 {
		base = 0
	}

	// Category positions: band centers when there are bars, else from edge to edge.
	center := func(j int) float64 {
		if anyBar || p.Horizontal {
			return (float64(j) + .5) * 100 / float64(n)
		}
		if n == 1 {
			return 50
		}
		return float64(j) * 100 / float64(n-1)
	}
	band := 100 / float64(n)
	if !anyBar && n > 1 {
		band = 100 / float64(n-1)
	}
	// along maps (category position, value position) to SVG coordinates; y grows downwards.
	at := func(c, val float64) (float64, float64) {
		if p.Horizontal {
			return val, c
		}
		return c, 100 - val
	}
	rect := func(c0, c1, v0, v1 float64) string {
		x0, y0 := at(c0, v0)
		x1, y1 := at(c1, v1)
		return fmt.Sprintf("M%s %sH%sV%sH%sZ", num(x0), num(y0), num(x1), num(y1), num(x0))
	}

	// Insets and axis labels.
	f := p.format
	if !p.Sparkline {
		values := sc.ticks()
		width := 0
		for _, t := range values {
			width = max(width, runes(f(t)))
		}
		catWidth := 0
		for j := range n {
			catWidth = max(catWidth, runes(cat(j)))
		}
		step := 1
		if !p.Horizontal {
			target := max(2, min(24, 72/(catWidth+2)))
			step = (n + target - 1) / target
		}
		// Every other label is dropped on a narrow chart (below 24rem) when the labels would not fit there.
		crowded := (n+step-1)/step*(catWidth+2) > 36
		k := 0
		for j := 0; j < n; j += step {
			t := tick{Pos: center(j), Text: cat(j), Alt: crowded && k%2 == 1}
			k++
			if p.Horizontal {
				v.YTicks = append(v.YTicks, t)
			} else {
				v.XTicks = append(v.XTicks, t)
			}
		}
		for _, t := range values {
			if p.Horizontal {
				v.XTicks = append(v.XTicks, tick{Pos: sc.pos(t), Text: f(t)})
			} else {
				v.YTicks = append(v.YTicks, tick{Pos: 100 - sc.pos(t), Text: f(t)})
			}
		}
		left, right, top := width, max(1, width/2), 0.5
		if p.Horizontal {
			left = min(catWidth, 16)
		} else if !anyBar {
			right = min(max(1, catWidth/2), 6)
		}
		if p.DataLabels && !p.Horizontal {
			top = 1.5
		}
		if p.DataLabels && p.Horizontal {
			right = max(right, width+1)
		}
		v.Inset = fmt.Sprintf("left:calc(%dch + .5rem);right:%dch;top:%srem;bottom:1.5rem", left, right, num(top))
	} else {
		v.Inset = "inset:0"
	}

	// Grid.
	if !p.Sparkline && !p.HideGrid {
		for _, t := range sc.ticks() {
			x0, y0 := at(0, sc.pos(t))
			x1, y1 := at(100, sc.pos(t))
			v.Shapes = append(v.Shapes, shape{D: fmt.Sprintf("M%s %sL%s %s", num(x0), num(y0), num(x1), num(y1)), Class: "chart-svg-grid"})
		}
	}

	// A tooltip per category with all its values, on a band behind the marks.
	{
		for j := range n {
			lines := []string{cat(j)}
			for _, s := range series {
				if x := value(s, j); !math.IsNaN(x) {
					lines = append(lines, s.Name+": "+f(x))
				}
			}
			c0, c1 := max(0, center(j)-band/2), min(100, center(j)+band/2)
			v.Shapes = append(v.Shapes, shape{D: rect(c0, c1, 0, 100), Class: "chart-svg-hit", Title: strings.Join(lines, "\n")})
		}
	}

	// Areas first, then bars, then lines on top.
	label := func(c, val float64, s string, inside bool) {
		x, y := at(c, val)
		cls, st := "chart-svg-v", fmt.Sprintf("left:%s%%;top:%s%%", num(x), num(y))
		switch {
		case inside:
			cls += " chart-svg-in"
		case p.Horizontal:
			cls += " chart-svg-end"
		}
		v.Labels = append(v.Labels, mkLabel(st, s, cls))
	}
	curve := p.StrokeCurve
	lower := make([]float64, n)
	upPos, upNeg := make([]float64, n), make([]float64, n)
	var lines []shape
	for _, want := range []kind{kArea, kBar} {
		bar := 0
		for i, s := range series {
			if kinds[i] != want {
				continue
			}
			color := fmt.Sprintf("var(--s%d)", i)
			switch want {
			case kArea:
				var pts, low []pt
				for j := range n {
					x := value(s, j)
					if math.IsNaN(x) {
						pts, low = append(pts, gap), append(low, gap)
						continue
					}
					b := 0.0
					if p.Stacked {
						b = lower[j]
						x += b
						lower[j] = x
					}
					cx, cy := at(center(j), sc.pos(x))
					bx, by := at(center(j), sc.pos(b))
					if !p.Stacked {
						bx, by = at(center(j), base)
					}
					pts, low = append(pts, pt{cx, cy}), append(low, pt{bx, by})
				}
				v.Shapes = append(v.Shapes, shape{D: areaPath(pts, low, curve, p.Stacked), Class: "chart-svg-area", Style: "fill:color-mix(in srgb, transparent, " + color + " 16%)"})
				lines = append(lines, p.line(i, pts, curve)...)
				if p.DataLabels {
					for j, q := range pts {
						if !q.gap() {
							label(center(j), 100-q.y, f(value(s, j)), false)
						}
					}
				}
			case kBar:
				group := band * .6
				if bars == 1 || p.Stacked {
					group = band * .5
				}
				w := group
				if !p.Stacked {
					w = group / float64(bars)
				}
				for j := range n {
					x := value(s, j)
					if math.IsNaN(x) {
						continue
					}
					c0 := center(j) - group/2
					if !p.Stacked {
						c0 += float64(bar) * w
					}
					v0, v1 := base, sc.pos(x)
					if p.Stacked {
						if x >= 0 {
							v0, v1 = sc.pos(upPos[j]), sc.pos(upPos[j]+x)
							upPos[j] += x
						} else {
							v0, v1 = sc.pos(upNeg[j]), sc.pos(upNeg[j]+x)
							upNeg[j] += x
						}
					}
					v0, v1 = clamp(v0), clamp(v1)
					st := "fill:" + color
					if len(series) == 1 && len(p.Colors) > 1 {
						st = "fill:" + p.color(j%len(p.Colors))
					}
					v.Shapes = append(v.Shapes, shape{D: rect(c0, c0+w, v0, v1), Class: "chart-svg-bar", Style: st, Title: cat(j) + "\n" + s.Name + ": " + f(x)})
					if p.DataLabels {
						if p.Stacked {
							if math.Abs(v1-v0) >= 8 { // no label on a sliver
								label(c0+w/2, (v0+v1)/2, f(x), true)
							}
						} else {
							label(c0+w/2, v1, f(x), false)
						}
					}
				}
				bar++
			}
		}
	}
	for i, s := range series {
		if kinds[i] != kLine {
			continue
		}
		var pts []pt
		for j := range n {
			x := value(s, j)
			if math.IsNaN(x) {
				pts = append(pts, gap)
				continue
			}
			cx, cy := at(center(j), sc.pos(x))
			pts = append(pts, pt{cx, cy})
			if p.DataLabels {
				label(center(j), sc.pos(x), f(x), false)
			}
		}
		lines = append(lines, p.line(i, pts, curve)...)
	}
	v.Shapes = append(v.Shapes, lines...)

	for _, a := range p.AnnotationsY {
		x, ok := toFloat(a.Value)
		if !ok || x < sc.lo || x > sc.hi {
			continue
		}
		x0, y0 := at(0, sc.pos(x))
		x1, y1 := at(100, sc.pos(x))
		color := colorVar(a.Color)
		if color == "" {
			color = "var(--tblr-red)"
		}
		dash := 4
		if a.Dash != nil {
			dash = *a.Dash
		}
		v.Shapes = append(v.Shapes, shape{D: fmt.Sprintf("M%s %sL%s %s", num(x0), num(y0), num(x1), num(y1)), Class: "chart-svg-line",
			Style: fmt.Sprintf("stroke:%s;stroke-width:1;stroke-dasharray:%d", color, dash), Title: strings.TrimSpace(a.Label + " " + f(x))})
		if a.Label != "" {
			v.Labels = append(v.Labels, mkLabel(fmt.Sprintf("left:%s%%;top:%s%%;background:%s", num(x1), num(y1), color), a.Label, "chart-svg-note"))
		}
	}

	// Legend and table.
	if p.Legend {
		for i, s := range series {
			v.Legend = append(v.Legend, legendItem{Name: s.Name, Color: fmt.Sprintf("var(--s%d)", i)})
		}
	}
	if p.Table {
		v.Head = []string{""}
		for _, s := range series {
			v.Head = append(v.Head, s.Name)
		}
		for j := range n {
			row := []string{cat(j)}
			for _, s := range series {
				row = append(row, f(value(s, j)))
			}
			v.Rows = append(v.Rows, row)
		}
	}
	return v
}

func mkLabel(style, s, class string) label { return label{Style: style, Text: s, Class: class} }

func clamp(x float64) float64 { return math.Max(0, math.Min(100, x)) }

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}

// line is the stroke of serie i and its markers.
func (p Props) line(i int, pts []pt, curve string) []shape {
	st := fmt.Sprintf("stroke:var(--s%d)", i)
	if i < len(p.StrokeWidth) && p.StrokeWidth[i] > 0 {
		st += ";stroke-width:" + num(p.StrokeWidth[i])
	}
	if i < len(p.StrokeDash) && p.StrokeDash[i] > 0 {
		st += ";stroke-dasharray:" + num(p.StrokeDash[i])
	}
	out := []shape{{D: linePath(pts, curve), Class: "chart-svg-line", Style: st}}
	if p.ShowMarkers {
		var b strings.Builder
		for _, q := range pts {
			if !q.gap() {
				fmt.Fprintf(&b, "M%s %sh0", num(q.x), num(q.y))
			}
		}
		out = append(out, shape{D: b.String(), Class: "chart-svg-dot", Style: fmt.Sprintf("stroke:var(--s%d)", i)})
	}
	return out
}

type pt struct{ x, y float64 }

var gap = pt{math.NaN(), math.NaN()}

func (q pt) gap() bool { return math.IsNaN(q.x) }

// runs splits points at gaps.
func runs(pts []pt) [][]pt {
	var out [][]pt
	var cur []pt
	for _, q := range pts {
		if q.gap() {
			if len(cur) > 0 {
				out = append(out, cur)
			}
			cur = nil
			continue
		}
		cur = append(cur, q)
	}
	if len(cur) > 0 {
		out = append(out, cur)
	}
	return out
}

func linePath(pts []pt, curve string) string {
	var b strings.Builder
	for _, r := range runs(pts) {
		fmt.Fprintf(&b, "M%s %s", num(r[0].x), num(r[0].y))
		segment(&b, r, curve)
	}
	return b.String()
}

// segment writes the path from r[0] through the rest of r.
func segment(b *strings.Builder, r []pt, curve string) {
	switch {
	case curve == "stepline":
		for _, q := range r[1:] {
			fmt.Fprintf(b, "H%sV%s", num(q.x), num(q.y))
		}
	case curve == "smooth" && len(r) > 2:
		t := monotone(r)
		for i := 0; i < len(r)-1; i++ {
			dx := (r[i+1].x - r[i].x) / 3
			fmt.Fprintf(b, "C%s %s %s %s %s %s", num(r[i].x+dx), num(r[i].y+t[i]*dx), num(r[i+1].x-dx), num(r[i+1].y-t[i+1]*dx), num(r[i+1].x), num(r[i+1].y))
		}
	default:
		for _, q := range r[1:] {
			fmt.Fprintf(b, "L%s %s", num(q.x), num(q.y))
		}
	}
}

// monotone returns the tangents of a monotone cubic through r (Fritsch–Carlson): the curve has no peak
// or dip between two points that the data does not have.
func monotone(r []pt) []float64 {
	n := len(r)
	m := make([]float64, n-1)
	for i := range m {
		m[i] = (r[i+1].y - r[i].y) / (r[i+1].x - r[i].x)
	}
	t := make([]float64, n)
	t[0], t[n-1] = m[0], m[n-2]
	for i := 1; i < n-1; i++ {
		if m[i-1]*m[i] > 0 {
			t[i] = (m[i-1] + m[i]) / 2
		}
	}
	for i := range m {
		if m[i] == 0 {
			t[i], t[i+1] = 0, 0
			continue
		}
		a, b := t[i]/m[i], t[i+1]/m[i]
		if s := a*a + b*b; s > 9 {
			k := 3 / math.Sqrt(s)
			t[i], t[i+1] = k*a*m[i], k*b*m[i]
		}
	}
	return t
}

// areaPath fills between pts and low (the baseline, or the serie below in a stack).
func areaPath(pts, low []pt, curve string, stacked bool) string {
	var b strings.Builder
	start := 0
	for i := 0; i <= len(pts); i++ {
		if i < len(pts) && !pts[i].gap() {
			continue
		}
		if i > start {
			up, down := pts[start:i], low[start:i]
			fmt.Fprintf(&b, "M%s %s", num(up[0].x), num(up[0].y))
			segment(&b, up, curve)
			rev := make([]pt, len(down))
			for k := range down {
				rev[k] = down[len(down)-1-k]
			}
			fmt.Fprintf(&b, "L%s %s", num(rev[0].x), num(rev[0].y))
			if stacked {
				segment(&b, rev, curve)
			} else {
				segment(&b, rev, "")
			}
			b.WriteString("Z")
		}
		start = i + 1
	}
	return b.String()
}

// Round charts are drawn in a 100×100 viewBox around (50, 50); angles in degrees, 0 at the top.

func polar(r, a float64) string {
	rad := a * math.Pi / 180
	return num(50+r*math.Sin(rad)) + " " + num(50-r*math.Cos(rad))
}

// arc continues a path along the circle of radius r from a0 to a1, in two halves above 180°.
func arc(r, a0, a1 float64) string {
	sweep := 1
	if a1 < a0 {
		sweep = 0
	}
	if math.Abs(a1-a0) > 180 {
		mid := (a0 + a1) / 2
		return fmt.Sprintf("A%s %s 0 0 %d %sA%s %s 0 0 %d %s", num(r), num(r), sweep, polar(r, mid), num(r), num(r), sweep, polar(r, a1))
	}
	return fmt.Sprintf("A%s %s 0 0 %d %s", num(r), num(r), sweep, polar(r, a1))
}

func sector(R, r, a0, a1 float64) string {
	d := "M" + polar(R, a0) + arc(R, a0, a1)
	if r <= 0 {
		return d + "L50 50Z"
	}
	return d + "L" + polar(r, a1) + arc(r, a1, a0) + "Z"
}

func (p Props) roundView() view {
	series := p.series()
	v := view{Round: true, ViewBox: "0 0 100 100", Colors: len(series)}
	f := p.format
	val := func(s chart.Serie) float64 {
		x := value(s, 0)
		if math.IsNaN(x) || x < 0 {
			return 0
		}
		return x
	}
	if p.typ() == chart.RadialBar {
		a0, a1 := p.StartAngle, p.EndAngle
		if a0 == 0 && a1 == 0 {
			a1 = 360
		}
		if a1-a0 <= 180 && a0 >= -90 && a1 <= 90 {
			v.ViewBox = "0 0 100 62"
		}
		hollow := p.Hollow
		if hollow <= 0 {
			hollow = 65
		}
		k := float64(max(len(series), 1))
		ring := 50 * (1 - hollow/100)
		gapW := math.Min(2, ring/k/4)
		w := ring/k - gapW
		if k == 1 {
			w = ring - 2
		}
		pct := func(x float64) string { return p.plain(x) + "%" }
		for i, s := range series {
			r := 50 - 1 - w/2 - float64(i)*(w+gapW)
			x := math.Min(val(s), 100)
			title := strings.TrimPrefix(s.Name+": "+pct(value(s, 0)), ": ")
			v.Shapes = append(v.Shapes, shape{D: "M" + polar(r, a0) + arc(r, a0, a1), Class: "chart-svg-track", Style: "stroke-width:" + num(w), Title: title})
			if x > 0 {
				v.Shapes = append(v.Shapes, shape{D: "M" + polar(r, a0) + arc(r, a0, a0+(a1-a0)*x/100), Class: "chart-svg-ring",
					Style: fmt.Sprintf("stroke:var(--s%d);stroke-width:%s", i, num(w)), Title: title})
			}
			if p.Legend {
				v.Legend = append(v.Legend, legendItem{Name: s.Name, Color: fmt.Sprintf("var(--s%d)", i), Value: pct(value(s, 0))})
			}
			if p.Table {
				v.Rows = append(v.Rows, []string{s.Name, pct(value(s, 0))})
			}
		}
		if len(series) == 1 {
			y := 50.0
			if v.ViewBox != "0 0 100 100" {
				y = 41
			}
			s := series[0]
			if s.Name != "" {
				v.Texts = append(v.Texts, text{X: 50, Y: y + 5, Text: pct(value(s, 0)), Class: "chart-svg-big"}, text{X: 50, Y: y + 14, Text: s.Name, Class: "chart-svg-name"})
			} else {
				v.Texts = append(v.Texts, text{X: 50, Y: y + 6, Text: pct(value(s, 0)), Class: "chart-svg-big"})
			}
		}
	} else {
		total := 0.0
		for _, s := range series {
			total += val(s)
		}
		hollow := 0.0
		if p.typ() == chart.Donut {
			hollow = p.Hollow
			if hollow <= 0 {
				hollow = 60
			}
		}
		R, r := 49.0, 49*hollow/100
		a := 0.0
		for i, s := range series {
			x := val(s)
			share := 0.0
			if total > 0 {
				share = x / total
			}
			title := s.Name + ": " + f(value(s, 0)) + " (" + p.plain(math.Round(share*1000)/10) + "%)"
			if share > 0 {
				v.Shapes = append(v.Shapes, shape{D: sector(R, r, a, a+share*360), Class: "chart-svg-slice", Style: fmt.Sprintf("fill:var(--s%d)", i), Title: title})
				if p.DataLabels && share >= .05 {
					mid := a + share*180
					rad := mid * math.Pi / 180
					lr := (R + r) / 2
					if r == 0 {
						lr = R * .62
					}
					v.Texts = append(v.Texts, text{X: 50 + lr*math.Sin(rad), Y: 50 - lr*math.Cos(rad) + 2, Text: p.plain(math.Round(share*100)) + "%", Class: "chart-svg-pl"})
				}
			}
			a += share * 360
			if p.Legend {
				v.Legend = append(v.Legend, legendItem{Name: s.Name, Color: fmt.Sprintf("var(--s%d)", i), Value: f(value(s, 0))})
			}
			if p.Table {
				v.Rows = append(v.Rows, []string{s.Name, f(value(s, 0))})
			}
		}
	}
	if p.Table {
		v.Head = []string{"", ""}
	}
	return v
}
