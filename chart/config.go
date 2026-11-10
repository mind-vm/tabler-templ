package chart

import (
	"fmt"
	"html"
	"math"
	"time"
)

// build turns Props into an Apache ECharts option. Colors stay CSS (var(), color-mix()); the chart
// script resolves them on the page, so they follow the theme.

type obj = map[string]any

// config is what the script gets: the ECharts option, formatters and renderers by path, and the
// relations between charts (group, brush).
type config struct {
	Option     obj               `json:"option"`
	Formatters map[string]string `json:"formatters,omitempty"`
	Group      string            `json:"group,omitempty"`
	Brush      *brush            `json:"brush,omitempty"`
}

type brush struct {
	Target string  `json:"target"`
	From   float64 `json:"from,omitempty"`
	To     float64 `json:"to,omitempty"`
}

const (
	textColor   = "var(--tblr-secondary)"
	gridColor   = "var(--tblr-border-color)"
	surface     = "var(--tblr-bg-surface)"
	bodyColor   = "var(--tblr-body-color)"
	defaultType = Bar
)

func (p Props) typ() Type {
	if p.Type == "" {
		return defaultType
	}
	return p.Type
}

func (p Props) heightPx() float64 {
	switch p.Size {
	case SizeSmall:
		return 40
	case SizeLarge:
		return 240
	}
	if p.Height > 0 {
		return p.Height * 16
	}
	return 160
}

func colorVar(c string) string {
	if c == "" {
		return ""
	}
	for _, r := range c {
		if r == '#' || r == '(' {
			return c
		}
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

// mix is color c at opacity pct (0–100).
func mix(c string, pct int) string {
	return fmt.Sprintf("color-mix(in srgb, transparent, %s %d%%)", c, pct)
}

func serieBase(p Props, i int) string {
	if i < len(p.Series) {
		if c := colorVar(firstOf(p.Series[i].Color, p.Color)); c != "" {
			return c
		}
	}
	return fmt.Sprintf("var(--tblr-chart-%d, %s)", i%5+1, palette[i%5])
}

// palette is Tabler's chart palette, for Tabler versions that do not define --tblr-chart-1 … 5 yet.
var palette = [5]string{"var(--tblr-primary)", "var(--tblr-green)", "var(--tblr-azure)", "var(--tblr-purple)", "var(--tblr-orange)"}

func days(start string, n int) []string {
	d, err := time.Parse(time.DateOnly, firstOf(start, "2020-06-20"))
	if err != nil {
		return nil
	}
	out := make([]string, n)
	for i := range out {
		out[i] = d.AddDate(0, 0, i).Format(time.DateOnly)
	}
	return out
}

func ms(day string) float64 {
	d, _ := time.Parse(time.DateOnly, day)
	return float64(d.UnixMilli())
}

func num(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	}
	return 0
}

func floats(v any) []float64 {
	f, _ := v.([]float64)
	return f
}

func tooltip(trigger string) obj {
	return obj{
		"trigger": trigger, "confine": true, "backgroundColor": surface, "borderColor": gridColor,
		"textStyle": obj{"color": bodyColor, "fontSize": 12}, "extraCssText": "box-shadow: var(--tblr-shadow-dropdown); border-radius: var(--tblr-border-radius);",
	}
}

func axis(kind string, show bool) obj {
	return obj{
		"type": kind, "show": show,
		"axisLine":  obj{"show": false, "lineStyle": obj{"color": gridColor}},
		"axisTick":  obj{"show": false},
		"axisLabel": obj{"color": textColor, "fontSize": 11, "hideOverlap": true},
		"splitLine": obj{"show": false, "lineStyle": obj{"color": gridColor, "type": []int{4, 4}}},
	}
}

// build ports Props.build to ECharts.
func (p Props) build() config {
	f := map[string]string{}
	for k, v := range p.formatters {
		f[k] = v
	}
	if p.raw != nil {
		return config{Option: p.raw, Formatters: f, Group: p.Group}
	}
	c := p.option(f)
	merge(c.Option, p.Options)
	return c
}

func (p Props) option(f map[string]string) config {
	t := p.typ()
	o := obj{
		"animation":       p.Animations,
		"backgroundColor": "transparent",
		"textStyle":       obj{"color": bodyColor},
	}
	names := make([]string, len(p.Series))
	for i, s := range p.Series {
		names[i] = s.Name
	}
	legend := obj{"show": p.Legend, "bottom": 0, "icon": "circle", "itemWidth": 8, "itemHeight": 8, "itemGap": 16,
		"textStyle": obj{"color": textColor, "fontSize": 12}, "data": names}
	o["legend"] = legend
	if p.Title != "" {
		o["title"] = obj{"text": html.EscapeString(p.Title), "left": 10, "textStyle": obj{"fontSize": 18, "color": bodyColor}}
	}
	if p.Toolbar {
		o["toolbox"] = obj{"right": 8, "top": 0, "feature": obj{"dataZoom": obj{"yAxisIndex": "none"}, "restore": obj{}, "saveAsImage": obj{}},
			"iconStyle": obj{"borderColor": textColor}}
	}
	grid := obj{"left": 0, "right": 8, "top": 12, "bottom": 4, "containLabel": true}
	if p.Legend {
		grid["bottom"] = 32
	}
	if p.Toolbar {
		grid["top"] = 32
	}
	if p.Sparkline {
		grid = obj{"left": 0, "right": 0, "top": 2, "bottom": 2, "containLabel": false}
	}
	tt := tooltip("axis")
	tt["show"] = !p.HideTooltip
	o["tooltip"] = tt

	c := config{Option: o, Formatters: f, Group: p.Group}
	if p.BrushTarget != "" {
		c.Brush = &brush{Target: p.BrushTarget, From: p.SelectionFrom, To: p.SelectionTo}
	}

	switch t {
	case Pie, Donut:
		pie(p, t, o)
		return c
	case PolarArea:
		polarArea(p, o)
		return c
	case RadialBar:
		gauges(p, o)
		return c
	case Radar:
		radar(p, o)
		return c
	case Treemap:
		treemap(p, o)
		return c
	}
	if p.Funnel {
		funnel(p, o)
		return c
	}

	o["grid"] = grid
	show := !p.Sparkline
	x, y := axis("category", show), axis("value", show)
	y["splitLine"].(obj)["show"] = show && !p.HideGrid
	x["splitLine"].(obj)["show"] = show && p.ShowX
	if p.YMax != nil {
		y["max"] = *p.YMax
	}
	if p.YTitle != "" {
		y["name"], y["nameLocation"], y["nameGap"] = p.YTitle, "middle", 40
	}
	if t == Line || t == Area {
		x["boundaryGap"] = false
	}
	if t == Line || t == Area || t == Scatter || t == Candlestick || t == Bubble {
		y["scale"] = true // a nice minimum, not zero
	}
	var labels []string
	switch {
	case p.Datetime && t != Candlestick && t != RangeBar:
		x["type"] = "time"
		n := 0
		for _, s := range p.Series {
			n = max(n, len(s.Data))
		}
		labels = days(p.StartDate, n)
	case p.Categories != nil:
		x["data"] = p.Categories
	}
	if p.YFormatter != "" {
		f["yAxis.0.axisLabel.formatter"] = string(p.YFormatter)
	}
	if p.XFormatter != "" {
		f["xAxis.0.axisLabel.formatter"] = string(p.XFormatter)
	}

	yAxes := []obj{y}
	axisOf := map[string]int{}
	if p.YAxes != nil {
		yAxes = nil
		for i, a := range p.YAxes {
			ya := axis("value", show)
			ya["splitLine"].(obj)["show"] = i == 0 && !p.HideGrid
			if a.Opposite {
				ya["position"] = "right"
			}
			if a.Min != nil {
				ya["min"] = *a.Min
			}
			if a.Max != nil {
				ya["max"] = *a.Max
			}
			if a.Title != "" {
				ya["name"], ya["nameLocation"], ya["nameGap"], ya["nameRotate"] = a.Title, "middle", 48, 90
				ya["nameTextStyle"] = obj{"color": textColor, "fontSize": 11, "fontWeight": 500}
			}
			if a.SeriesName != "" {
				axisOf[a.SeriesName] = i
			}
			yAxes = append(yAxes, ya)
		}
	}

	var series []obj
	switch t {
	case Heatmap:
		series = heatmap(p, o, x, y)
	case BoxPlot:
		series = boxplot(p, x)
	case Candlestick:
		series = candlestick(p, x)
		x["boundaryGap"] = true
	case RangeBar:
		series = timeline(p, x, y, f)
		x, y = y, x // categories on the left
		yAxes = []obj{y}
		tt["trigger"] = "item"
	case Bubble:
		series = bubble(p, x)
		tt["trigger"] = "item"
	default:
		for i, s := range p.Series {
			series = append(series, cartesian(p, t, i, s, labels, axisOf))
		}
		if t == Scatter {
			tt["trigger"] = "item"
		}
	}

	if t == Bar && p.Horizontal {
		x, yAxes[0] = yAxes[0], x
		if p.YFormatter != "" {
			delete(f, "yAxis.0.axisLabel.formatter")
			f["xAxis.0.axisLabel.formatter"] = string(p.YFormatter)
		}
	}
	if len(series) > 0 && (p.AnnotationsX != nil || p.AnnotationsY != nil) {
		series[0]["markLine"] = markLines(p)
	}

	xAxes := []obj{x}
	if p.BrushTarget != "" && len(labels) > 0 {
		// The slider drives a hidden second time axis, so the brush chart itself stays unzoomed.
		hidden := obj{"type": "time", "show": false, "min": ms(labels[0]), "max": ms(labels[len(labels)-1])}
		xAxes = append(xAxes, hidden)
		o["dataZoom"] = []obj{{
			"type": "slider", "xAxisIndex": 1, "filterMode": "none", "startValue": p.SelectionFrom, "endValue": p.SelectionTo,
			"top": 0, "bottom": 0, "left": 0, "right": 0, "height": "auto", "showDetail": false, "showDataShadow": false,
			"backgroundColor": "transparent", "borderColor": "transparent", "fillerColor": mix("var(--tblr-primary)", 16),
			"handleStyle": obj{"color": "var(--tblr-primary)", "borderColor": "var(--tblr-primary)"}, "moveHandleSize": 0,
			"brushSelect": false,
		}}
	}
	o["xAxis"], o["yAxis"], o["series"] = xAxes, yAxes, series
	return c
}

func cartesian(p Props, t Type, i int, s Serie, labels []string, axisOf map[string]int) obj {
	kind := map[Type]string{Bar: "bar", Line: "line", Area: "line", Scatter: "scatter"}[t]
	switch s.Type {
	case "column", "bar":
		kind = "bar"
	case "line":
		kind = "line"
	case "area":
		kind = "line"
	}
	color := p.SerieColor(i)
	se := obj{"name": s.Name, "type": kind, "itemStyle": obj{"color": color}, "emphasis": obj{"focus": "none"}}
	var data any = s.Data
	if labels != nil {
		pairs := make([][2]any, len(s.Data))
		for j, v := range s.Data {
			pairs[j] = [2]any{labels[j], v}
		}
		data = pairs
	}
	if s.Points != nil {
		pts := make([][2]any, len(s.Points))
		for j, pt := range s.Points {
			pts[j] = [2]any{pt.X, pt.Y}
		}
		data = pts
	}
	se["data"] = data
	if n, ok := axisOf[s.Name]; ok {
		se["yAxisIndex"] = n
	}
	if p.Stacked {
		se["stack"] = "total"
	}
	switch kind {
	case "line":
		w := 2.0
		if i < len(p.StrokeWidth) {
			w = p.StrokeWidth[i]
		}
		ls := obj{"width": w, "cap": "round", "color": color}
		if i < len(p.StrokeDash) && p.StrokeDash[i] > 0 {
			ls["type"] = []float64{p.StrokeDash[i], p.StrokeDash[i]}
		}
		se["lineStyle"] = ls
		switch firstOf(p.StrokeCurve, "smooth") {
		case "smooth":
			se["smooth"] = 0.35
		case "stepline":
			se["step"] = "end"
		}
		se["showSymbol"] = p.ShowMarkers
		se["symbolSize"] = 5
		if t == Area || s.Type == "area" {
			se["areaStyle"] = obj{"color": mix(serieBase(p, i), 16), "opacity": 1}
		}
	case "bar":
		se["barCategoryGap"] = "50%"
		se["barMaxWidth"] = 48
	case "scatter":
		se["symbolSize"] = 9
	}
	if p.DataLabels {
		pos := "top"
		if kind == "bar" {
			pos = "inside"
		}
		se["label"] = obj{"show": true, "position": pos, "fontSize": 11, "color": map[bool]string{true: "#fff", false: textColor}[kind == "bar"]}
	}
	return se
}

func markLines(p Props) obj {
	var data []obj
	add := func(a Annotation, key string) {
		color := colorVar(firstOf(a.Color, "red"))
		dash := 4
		if a.Dash != nil {
			dash = *a.Dash
		}
		pos := "insideEndTop"
		if key == "yAxis" {
			pos = "insideStartTop"
		}
		data = append(data, obj{key: a.Value,
			"lineStyle": obj{"color": color, "type": []int{dash, dash}, "width": 1},
			"label": obj{"show": a.Label != "", "formatter": a.Label, "position": pos, "color": "#fff", "backgroundColor": color,
				"padding": []int{2, 6}, "borderRadius": 3, "fontSize": 11}})
	}
	for _, a := range p.AnnotationsX {
		add(a, "xAxis")
	}
	for _, a := range p.AnnotationsY {
		add(a, "yAxis")
	}
	return obj{"silent": true, "symbol": "none", "data": data}
}

func pie(p Props, t Type, o obj) {
	data := make([]obj, len(p.Series))
	for i, s := range p.Series {
		v := 0.0
		if len(s.Data) > 0 {
			v = s.Data[0]
		}
		data[i] = obj{"name": s.Name, "value": v, "itemStyle": obj{"color": p.SerieColor(i)}}
	}
	radius := []string{"0%", "80%"}
	if t == Donut {
		radius = []string{"55%", "80%"}
	}
	se := obj{"type": "pie", "radius": radius, "data": data, "center": []string{"50%", "45%"},
		"itemStyle": obj{"borderColor": surface, "borderWidth": 2},
		"label":     obj{"show": p.DataLabels, "position": "inside", "formatter": "{d}%", "color": "#fff", "fontSize": 11},
		"labelLine": obj{"show": false}}
	if !p.Legend {
		se["center"] = []string{"50%", "50%"}
	}
	tt := tooltip("item")
	tt["show"] = !p.HideTooltip
	o["tooltip"], o["series"] = tt, []obj{se}
}

// gauges draw a radial bar: one ring per serie, from -120° to 120°.
func gauges(p Props, o obj) {
	var series []obj
	n := len(p.Series)
	for i, s := range p.Series {
		v := 0.0
		if len(s.Data) > 0 {
			v = s.Data[0]
		}
		w := 14.0
		if n > 1 {
			w = 8
		}
		r := 90 - float64(i)*(w+6)/1.6
		series = append(series, obj{
			"type": "gauge", "name": s.Name, "startAngle": 210, "endAngle": -30, "min": 0, "max": 100, "radius": fmt.Sprintf("%.0f%%", r),
			"center":   []string{"50%", "55%"},
			"progress": obj{"show": true, "roundCap": p.StrokeLinecap == "round" || n == 1, "width": w, "itemStyle": obj{"color": p.SerieColor(i)}},
			"axisLine": obj{"roundCap": p.StrokeLinecap == "round" || n == 1, "lineStyle": obj{"width": w, "color": [][2]any{{1, "var(--tblr-bg-surface-secondary)"}}}},
			"pointer":  obj{"show": false}, "axisTick": obj{"show": false}, "splitLine": obj{"show": false}, "axisLabel": obj{"show": false},
			"anchor": obj{"show": false}, "title": obj{"show": false},
			"detail": obj{"show": n == 1, "valueAnimation": p.Animations, "formatter": "{value}%", "fontSize": 24, "fontWeight": 600, "color": bodyColor, "offsetCenter": []any{0, 0}},
			"data":   []obj{{"value": v, "name": s.Name}},
		})
	}
	tt := tooltip("item")
	tt["show"] = !p.HideTooltip
	o["tooltip"], o["series"] = tt, series
}

func radar(p Props, o obj) {
	maxV := 0.0
	for _, s := range p.Series {
		for _, v := range s.Data {
			maxV = math.Max(maxV, v)
		}
	}
	ind := make([]obj, len(p.Categories))
	for i, c := range p.Categories {
		ind[i] = obj{"name": c, "max": math.Ceil(maxV * 1.1)}
	}
	data := make([]obj, len(p.Series))
	for i, s := range p.Series {
		data[i] = obj{"name": s.Name, "value": s.Data, "itemStyle": obj{"color": p.SerieColor(i)},
			"lineStyle": obj{"color": p.SerieColor(i), "width": 2}, "areaStyle": obj{"color": mix(serieBase(p, i), 15)}, "symbolSize": 5}
	}
	o["radar"] = obj{"indicator": ind, "radius": "65%", "center": []string{"50%", "48%"},
		"axisName": obj{"color": textColor, "fontSize": 11}, "splitArea": obj{"show": false},
		"splitLine": obj{"lineStyle": obj{"color": gridColor}}, "axisLine": obj{"lineStyle": obj{"color": gridColor}}}
	tt := tooltip("item")
	tt["show"] = !p.HideTooltip
	o["tooltip"], o["series"] = tt, []obj{{"type": "radar", "data": data}}
}

func treemap(p Props, o obj) {
	var data []obj
	color := serieBase(p, 0)
	if len(p.Series) > 0 {
		lo, hi := math.Inf(1), math.Inf(-1)
		for _, pt := range p.Series[0].Points {
			lo, hi = math.Min(lo, num(pt.Y)), math.Max(hi, num(pt.Y))
		}
		for _, pt := range p.Series[0].Points {
			shade := 100
			if hi > lo {
				shade = 45 + int(55*(num(pt.Y)-lo)/(hi-lo)) // lighter for smaller values
			}
			data = append(data, obj{"name": pt.X, "value": num(pt.Y), "itemStyle": obj{"color": mix(color, shade)}})
		}
	}
	o["series"] = []obj{{"type": "treemap", "data": data, "roam": false, "nodeClick": false, "breadcrumb": obj{"show": false},
		"width": "100%", "height": "100%", "top": 0, "left": 0,
		"itemStyle": obj{"borderColor": surface, "borderWidth": 2, "gapWidth": 2, "borderRadius": 3},
		"label":     obj{"show": true, "color": "#fff", "fontSize": 12}}}
	tt := tooltip("item")
	o["tooltip"] = tt
}

func funnel(p Props, o obj) {
	var data []obj
	if len(p.Series) > 0 {
		for i, v := range p.Series[0].Data {
			d := obj{"value": v}
			if i < len(p.Categories) {
				d["name"] = p.Categories[i]
			}
			if i < len(p.Colors) {
				d["itemStyle"] = obj{"color": p.Colors[i]}
			}
			data = append(data, d)
		}
	}
	o["legend"].(obj)["data"] = p.Categories
	o["series"] = []obj{{"type": "funnel", "sort": "none", "gap": 3, "left": "5%", "width": "90%", "top": 4, "bottom": map[bool]int{true: 32, false: 4}[p.Legend], "data": data,
		"minSize": "10%", "label": obj{"show": true, "position": "inside", "formatter": "{b}: {c}", "color": "#fff", "fontSize": 11},
		"itemStyle": obj{"borderWidth": 0}}}
	tt := tooltip("item")
	o["tooltip"] = tt
}

func heatmap(p Props, o, x, y obj) []obj {
	var cats []any
	if len(p.Series) > 0 {
		for _, pt := range p.Series[0].Points {
			cats = append(cats, pt.X)
		}
	}
	rows := make([]string, len(p.Series))
	var data [][3]any
	maxV := 0.0
	for r, s := range p.Series {
		rows[r] = s.Name
		for c, pt := range s.Points {
			if pt.Y == nil {
				data = append(data, [3]any{c, r, "-"}) // no value: an empty cell
				continue
			}
			v := num(pt.Y)
			maxV = math.Max(maxV, v)
			data = append(data, [3]any{c, r, v})
		}
	}
	x["data"], y["type"], y["data"] = cats, "category", rows
	y["splitLine"].(obj)["show"] = false
	base := serieBase(p, 0)
	o["visualMap"] = obj{"show": false, "min": 0, "max": maxV, "inRange": obj{"color": []string{mix(base, 8), base}}}
	o["tooltip"].(obj)["trigger"] = "item"
	return []obj{{"type": "heatmap", "data": data, "itemStyle": obj{"borderColor": surface, "borderWidth": 3, "borderRadius": 4}}}
}

func boxplot(p Props, x obj) []obj {
	var out []obj
	for i, s := range p.Series {
		var cats []any
		var data [][]float64
		for _, pt := range s.Points {
			cats = append(cats, pt.X)
			data = append(data, floats(pt.Y))
		}
		if i == 0 {
			x["data"] = cats
		}
		base := serieBase(p, i)
		out = append(out, obj{"type": "boxplot", "name": s.Name, "data": data, "boxWidth": []string{"20%", "50%"},
			"itemStyle": obj{"color": mix(base, 30), "borderColor": base, "borderWidth": 1.5}})
	}
	return out
}

func candlestick(p Props, x obj) []obj {
	var out []obj
	for i, s := range p.Series {
		var cats []string
		var data [][4]float64
		for _, pt := range s.Points {
			cats = append(cats, time.UnixMilli(int64(num(pt.X))).UTC().Format("Jan 2006"))
			y := floats(pt.Y) // Point: open, high, low, close; ECharts: open, close, low, high
			if len(y) == 4 {
				data = append(data, [4]float64{y[0], y[3], y[2], y[1]})
			}
		}
		if i == 0 {
			x["data"] = cats
		}
		out = append(out, obj{"type": "candlestick", "name": s.Name, "data": data, "barMaxWidth": 14,
			"itemStyle": obj{"color": "var(--tblr-green)", "color0": "var(--tblr-red)", "borderColor": "var(--tblr-green)", "borderColor0": "var(--tblr-red)"}})
	}
	return out
}

// timeline draws ranges with a custom series; the script puts in the "timeline" renderer.
func timeline(p Props, x, y obj, f map[string]string) []obj {
	var out []obj
	var rows []any
	for i, s := range p.Series {
		var data []obj
		for _, pt := range s.Points {
			r := len(rows)
			rows = append(rows, pt.X)
			if y := floats(pt.Y); len(y) == 2 {
				data = append(data, obj{"name": pt.X, "value": []any{r, y[0], y[1]}})
			}
		}
		out = append(out, obj{"type": "custom", "name": s.Name, "renderItem": "timeline", "encode": obj{"x": []int{1, 2}, "y": 0},
			"itemStyle": obj{"color": p.SerieColor(i)}, "data": data})
		f[fmt.Sprintf("series.%d.renderItem", i)] = "timeline"
	}
	x["data"] = rows
	x["inverse"] = true
	y["type"] = "time"
	y["splitLine"].(obj)["show"] = true
	return out
}

func bubble(p Props, x obj) []obj {
	maxZ := 0.0
	for _, s := range p.Series {
		for _, pt := range s.Points {
			maxZ = math.Max(maxZ, num(pt.Z))
		}
	}
	x["type"] = "value"
	x["splitLine"].(obj)["show"] = p.ShowX
	var out []obj
	for i, s := range p.Series {
		var data []obj
		for _, pt := range s.Points {
			z := num(pt.Z)
			size := 10.0
			if maxZ > 0 {
				size = 6 + 24*z/maxZ
			}
			data = append(data, obj{"value": []any{pt.X, pt.Y, pt.Z}, "symbolSize": size})
		}
		out = append(out, obj{"type": "scatter", "name": s.Name, "data": data, "itemStyle": obj{"color": mix(serieBase(p, i), 80)}})
	}
	return out
}

// polarArea draws a polar area chart as bars in polar coordinates: equal angles, radius by value,
// rings for the scale. One serie per slice, so the legend lists them.
func polarArea(p Props, o obj) {
	n := len(p.Series)
	names := make([]string, n)
	for i, s := range p.Series {
		names[i] = s.Name
	}
	center := []string{"50%", "50%"}
	if p.Legend {
		center = []string{"50%", "44%"}
	}
	o["polar"] = obj{"radius": []string{"0%", "78%"}, "center": center}
	o["angleAxis"] = obj{"type": "category", "data": names, "startAngle": 90, "axisLine": obj{"show": false}, "axisTick": obj{"show": false},
		"axisLabel": obj{"show": false}, "splitLine": obj{"show": true, "lineStyle": obj{"color": gridColor}}}
	o["radiusAxis"] = obj{"type": "value", "axisLine": obj{"show": false}, "axisTick": obj{"show": false},
		"axisLabel": obj{"color": textColor, "fontSize": 10}, "splitLine": obj{"show": true, "lineStyle": obj{"color": gridColor}}}
	var series []obj
	for i, s := range p.Series {
		data := make([]any, n)
		for j := range data {
			data[j] = 0
		}
		if len(s.Data) > 0 {
			data[i] = s.Data[0]
		}
		series = append(series, obj{"type": "bar", "coordinateSystem": "polar", "name": s.Name, "stack": "slices", "data": data,
			"barCategoryGap": 0, "itemStyle": obj{"color": p.SerieColor(i), "opacity": 0.85, "borderColor": surface, "borderWidth": 1}})
	}
	tt := tooltip("item")
	tt["show"] = !p.HideTooltip
	o["tooltip"], o["series"] = tt, series
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
