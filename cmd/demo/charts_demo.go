package main

import "github.com/mind-vm/tabler-templ/chart"

// tc is the chart of the Tabler preview with that id in charts.json, height in rem (0 keeps the chart's
// own height).
func tc(id, label string, height float64) chart.Props {
	p := tablerCharts[id]
	p.ID, p.Label = id, label
	if height > 0 {
		p.Height = height
	}
	return p
}

// as gives a chart another id, for a second chart from the same data.
func as(p chart.Props, id string) chart.Props {
	p.ID = id
	return p
}

// small is a chart at the small height of a metric card (chart-sm), with a class.
func small(p chart.Props, class string) chart.Props {
	p.Size, p.Class = chart.SizeSmall, class
	return p
}

// noLegend hides the chart's own legend, for a page that shows its own.
func noLegend(p chart.Props) chart.Props {
	p.Legend = false
	return p
}
