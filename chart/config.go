package chart

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// config is what the chart script gets: the ApexCharts options and the formatters to put in by path. The
// script resolves the colors itself.
type config struct {
	Config     map[string]any    `json:"config"`
	Formatters map[string]string `json:"formatters,omitempty"`
}

func (p Props) typ() Type {
	if p.Type == "" {
		return Bar
	}
	return p.Type
}

func (p Props) height() float64 {
	switch p.Size {
	case SizeSmall:
		return 2.5
	case SizeLarge:
		return 15
	}
	if p.Height > 0 {
		return p.Height
	}
	return 10
}

func (p Props) round() bool {
	t := p.typ()
	return t == Pie || t == Donut || t == RadialBar || t == PolarArea
}

// style is the series colors as custom properties on the chart element: a serie's own color, the
// chart's color, or the next slot of Tabler's chart palette (--tblr-chart-1 … 5).
func (p Props) style() string {
	var b strings.Builder
	for i, s := range p.Series {
		base := colorVar(firstOf(s.Color, p.Color))
		if base == "" {
			base = fmt.Sprintf("var(--tblr-chart-%d, %s)", i%5+1, palette[i%5])
		}
		fmt.Fprintf(&b, "--chart-color-%d: color-mix(in srgb, transparent, %s %s);", i, base, firstOf(s.ColorOpacity, "100%"))
		if p.typ() == Area {
			fmt.Fprintf(&b, "--chart-fill-%d: color-mix(in srgb, transparent, %s 16%%);", i, base)
		}
	}
	if p.typ() == BoxPlot {
		c := "primary"
		if len(p.Series) > 0 {
			c = firstOf(p.Series[0].Color, p.Color, c)
		}
		fmt.Fprintf(&b, "--chart-box-upper: %s;--chart-box-lower: color-mix(in srgb, transparent, %s 40%%);", colorVar(c), colorVar(c))
	}
	return b.String()
}

// palette is Tabler's chart palette, for Tabler versions that do not define --tblr-chart-1 … 5 yet.
var palette = [5]string{"var(--tblr-primary)", "var(--tblr-green)", "var(--tblr-azure)", "var(--tblr-purple)", "var(--tblr-orange)"}

func colorVar(c string) string {
	if c == "" || strings.ContainsAny(c, "#(") {
		return c
	}
	return "var(--tblr-" + c + ")"
}

func firstOf(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

func annotation(a Annotation, axis string) map[string]any {
	color := colorVar(firstOf(a.Color, "red"))
	m := map[string]any{axis: a.Value, "strokeDashArray": 4, "borderColor": color}
	if a.Dash != nil {
		m["strokeDashArray"] = *a.Dash
	}
	if a.Label != "" {
		pos, anchor, off := "top", "middle", 0
		if axis == "y" {
			pos, anchor, off = "left", "start", 8
		}
		m["label"] = map[string]any{
			"text": a.Label, "borderWidth": 0, "position": pos, "textAnchor": anchor, "offsetX": off,
			"style": map[string]any{"background": color, "color": "var(--tblr-white)", "fontSize": "11px", "fontFamily": "inherit"},
		}
	}
	return m
}

func (s Serie) data() any {
	if s.Points != nil {
		return s.Points
	}
	return s.Data
}

// build ports Tabler's chartConfig (shared/lib/chart-script.ts).
func (p Props) build() config {
	f := map[string]string{}
	for k, v := range p.formatters {
		f[k] = v
	}
	if p.raw != nil {
		return config{Config: p.raw, Formatters: f}
	}
	t := p.typ()
	c := map[string]any{}

	ch := map[string]any{"id": p.ID, "type": string(t), "fontFamily": "inherit", "height": p.height() * 16}
	if p.Sparkline {
		ch["sparkline"] = map[string]any{"enabled": true}
	} else {
		ch["parentHeightOffset"] = 0
		ch["toolbar"] = map[string]any{"show": p.Toolbar}
	}
	if !p.Animations {
		ch["animations"] = map[string]any{"enabled": false}
	}
	if p.Stacked {
		ch["stacked"] = true
	}
	if p.Group != "" {
		ch["group"] = p.Group
	}
	if p.BrushTarget != "" {
		ch["brush"] = map[string]any{"enabled": true, "target": p.BrushTarget}
		ch["selection"] = map[string]any{
			"enabled": true,
			"xaxis":   map[string]any{"min": p.SelectionFrom, "max": p.SelectionTo},
			"fill":    map[string]any{"color": "var(--tblr-primary)", "opacity": 0.16},
			"stroke":  map[string]any{"color": "var(--tblr-primary)", "width": 1, "opacity": 1, "dashArray": 0},
		}
	}
	c["chart"] = ch

	switch t {
	case Bar, RangeBar:
		bar := map[string]any{"columnWidth": "50%"}
		if p.Funnel {
			bar = map[string]any{"horizontal": true, "barHeight": "80%", "isFunnel": true, "isFunnel3d": false, "distributed": true}
		} else if p.Horizontal {
			bar = map[string]any{"barHeight": "50%", "horizontal": true}
		}
		c["plotOptions"] = map[string]any{"bar": bar}
	case RadialBar:
		c["plotOptions"] = map[string]any{"radialBar": map[string]any{
			"startAngle": -120, "endAngle": 120,
			"hollow": map[string]any{"margin": 16, "size": "50%"},
			"dataLabels": map[string]any{
				"show":  true,
				"name":  map[string]any{"show": false, "offsetY": 0, "fontSize": "16px", "color": "inherit"},
				"value": map[string]any{"offsetY": -8, "fontSize": "24px", "fontWeight": 600},
			},
		}}
	case Heatmap:
		c["plotOptions"] = map[string]any{"heatmap": map[string]any{"radius": 3}}
		c["dataLabels"] = map[string]any{"enabled": false}
	case Treemap:
		c["plotOptions"] = map[string]any{"treemap": map[string]any{"distributed": false, "enableShades": true, "shadeIntensity": 0.4}}
		c["dataLabels"] = map[string]any{"enabled": true}
	case Radar:
		radar := map[string]any{"polygons": map[string]any{
			"strokeColors": "var(--apx-grid)", "connectorColors": "var(--apx-grid)",
			"fill": map[string]any{"colors": []string{"transparent", "transparent"}},
		}}
		if p.RadarSize > 0 {
			radar["size"] = p.RadarSize
		}
		c["plotOptions"] = map[string]any{"radar": radar}
		c["markers"] = map[string]any{"size": 3}
		c["fill"] = map[string]any{"opacity": 0.15}
	case PolarArea:
		c["plotOptions"] = map[string]any{"polarArea": map[string]any{
			"rings": map[string]any{"strokeColor": "var(--apx-grid)"}, "spokes": map[string]any{"connectorColors": "var(--apx-grid)"},
		}}
	case BoxPlot:
		c["plotOptions"] = map[string]any{"boxPlot": map[string]any{"colors": map[string]any{"upper": "var(--chart-box-upper)", "lower": "var(--chart-box-lower)"}}}
		c["stroke"] = map[string]any{"colors": []string{"var(--tblr-secondary)"}}
	case Bubble:
		c["fill"] = map[string]any{"opacity": 0.8}
		c["dataLabels"] = map[string]any{"enabled": false}
	case Area:
		fills := make([]string, len(p.Series))
		for i := range p.Series {
			fills[i] = fmt.Sprintf("var(--chart-fill-%d)", i)
		}
		c["fill"] = map[string]any{"colors": fills, "type": "solid"}
	}
	if t == Bar || t == Area {
		c["dataLabels"] = map[string]any{"enabled": p.DataLabels}
	}

	if p.Title != "" {
		c["title"] = map[string]any{"text": html.EscapeString(p.Title), "margin": 0, "floating": true, "offsetX": 10, "style": map[string]any{"fontSize": "18px"}}
	}

	if t == RadialBar && p.StrokeLinecap != "" {
		c["stroke"] = map[string]any{"lineCap": p.StrokeLinecap}
	} else if t == Area || t == Line {
		s := map[string]any{"width": 2, "lineCap": "round", "curve": firstOf(p.StrokeCurve, "smooth")}
		if p.StrokeWidth != nil {
			s["width"] = p.StrokeWidth
		}
		if p.StrokeDash != nil {
			s["dashArray"] = p.StrokeDash
		}
		c["stroke"] = s
	}

	if len(p.Series) > 0 {
		if p.round() {
			vals, labels := make([]float64, len(p.Series)), make([]string, len(p.Series))
			for i, s := range p.Series {
				if len(s.Data) > 0 {
					vals[i] = s.Data[0]
				}
				labels[i] = s.Name
			}
			c["series"], c["labels"] = vals, labels
		} else {
			series := make([]map[string]any, len(p.Series))
			for i, s := range p.Series {
				series[i] = map[string]any{"name": s.Name, "data": s.data()}
				if s.Type != "" {
					series[i]["type"] = s.Type
				}
			}
			c["series"] = series
		}
	}

	if p.AnnotationsX != nil || p.AnnotationsY != nil {
		a := map[string]any{}
		for _, x := range p.AnnotationsX {
			a["xaxis"] = append(asSlice(a["xaxis"]), annotation(x, "x"))
		}
		for _, y := range p.AnnotationsY {
			a["yaxis"] = append(asSlice(a["yaxis"]), annotation(y, "y"))
		}
		c["annotations"] = a
	}

	// The negative padding pulls a cartesian plot flush to the edges of its card; radial and block
	// layouts draw their labels outside the plot.
	loose := t == Radar || t == PolarArea || t == Treemap || p.Funnel || p.YAxes != nil
	grid := map[string]any{}
	if !p.Sparkline && !loose {
		grid["padding"] = map[string]any{"top": -20, "right": 0, "left": -4, "bottom": -4}
	}
	if p.HideGrid {
		grid["show"] = false
	} else {
		grid["strokeDashArray"] = 4
		if p.ShowX {
			grid["xaxis"] = map[string]any{"lines": map[string]any{"show": true}}
		}
	}
	c["grid"] = grid

	if p.DataLabels {
		c["dataLabels"] = map[string]any{"enabled": true}
	}

	if p.Categories != nil || p.Datetime || t == Heatmap || t == Bubble || t == BoxPlot || t == RangeBar {
		x := map[string]any{"labels": map[string]any{"padding": 0}, "tooltip": map[string]any{"enabled": false}}
		if t == Area || t == Bar {
			x["axisBorder"] = map[string]any{"show": false}
		}
		if p.Categories != nil {
			x["categories"] = p.Categories
		}
		if p.Datetime {
			x["type"] = "datetime"
		}
		c["xaxis"] = x
		if p.XFormatter != "" {
			f["xaxis.labels.formatter"] = string(p.XFormatter)
		}
	}

	axisTitle := func(text string, offset int) map[string]any {
		m := map[string]any{"text": html.EscapeString(text), "style": map[string]any{"fontSize": "11px", "fontWeight": 500, "fontFamily": "inherit"}}
		if offset != 0 {
			m["offsetX"] = offset
		}
		return m
	}
	switch {
	case t == Radar:
		c["yaxis"] = map[string]any{"show": false}
	case p.YAxes != nil:
		axes := make([]map[string]any, len(p.YAxes))
		for i, a := range p.YAxes {
			m := map[string]any{"labels": map[string]any{"padding": 4}}
			if a.SeriesName != "" {
				m["seriesName"] = a.SeriesName
			}
			if a.Opposite {
				m["opposite"] = true
			}
			if a.Min != nil {
				m["min"] = *a.Min
			}
			if a.Max != nil {
				m["max"] = *a.Max
			}
			if a.Title != "" {
				off := -8
				if a.Opposite {
					off = 8
				}
				m["title"] = axisTitle(a.Title, off)
			}
			axes[i] = m
		}
		c["yaxis"] = axes
	case !p.round():
		y := map[string]any{"labels": map[string]any{"padding": 4}}
		if p.YMax != nil {
			y["max"] = *p.YMax
		}
		if p.YTitle != "" {
			y["title"] = axisTitle(p.YTitle, 0)
		}
		if p.YTooltip {
			y["tooltip"] = map[string]any{"enabled": true}
		}
		c["yaxis"] = y
		if p.YFormatter != "" {
			f["yaxis.labels.formatter"] = string(p.YFormatter)
		}
	}

	if p.Datetime {
		n := 0
		if len(p.Series) > 0 {
			n = len(p.Series[0].Data)
		}
		c["labels"] = days(firstOf(p.StartDate, "2020-06-20"), n)
	}

	if p.Colors != nil {
		c["colors"] = p.Colors
	} else if len(p.Series) > 0 {
		colors := make([]string, len(p.Series))
		for i := range p.Series {
			colors[i] = fmt.Sprintf("var(--chart-color-%d)", i)
		}
		c["colors"] = colors
	}

	if p.Legend {
		c["legend"] = map[string]any{
			"show": true, "position": "bottom", "offsetY": 12,
			"markers":    map[string]any{"size": 5, "strokeWidth": 0, "shape": "circle"},
			"itemMargin": map[string]any{"horizontal": 8, "vertical": 8},
		}
	} else {
		c["legend"] = map[string]any{"show": false}
	}

	if p.HideTooltip || t == Pie || t == Donut {
		tt := map[string]any{}
		if p.HideTooltip {
			tt["enabled"] = false
		}
		if t == Pie || t == Donut {
			tt["fillSeriesColor"] = false
		}
		c["tooltip"] = tt
	}
	if p.HidePoints {
		c["markers"] = map[string]any{"size": 0, "hover": map[string]any{"sizeOffset": 0}}
	}
	if p.ShowMarkers {
		c["markers"] = map[string]any{"size": 2}
	}

	merge(c, p.Options)
	return config{Config: c, Formatters: f}
}

func asSlice(v any) []any {
	s, _ := v.([]any)
	return s
}

// merge copies src into dst, descending into maps present on both sides.
func merge(dst, src map[string]any) {
	for k, v := range src {
		if sm, ok := v.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				merge(dm, sm)
				continue
			}
		}
		dst[k] = v
	}
}

// days are consecutive dates from start (YYYY-MM-DD).
func days(start string, n int) []string {
	d, err := time.Parse(time.DateOnly, start)
	if err != nil {
		return nil
	}
	out := make([]string, n)
	for i := range out {
		out[i] = d.AddDate(0, 0, i).Format(time.DateOnly)
	}
	return out
}
