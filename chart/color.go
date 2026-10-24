package chart

import "fmt"

// SerieColor is the CSS color the chart draws serie i in, e.g. for a legend outside the chart:
// --tblr-legend-color: <SerieColor(i)>.
func (p Props) SerieColor(i int) string {
	if p.Colors != nil && i < len(p.Colors) {
		return p.Colors[i]
	}
	base, opacity := "", "100%"
	if i < len(p.Series) {
		base, opacity = colorVar(firstOf(p.Series[i].Color, p.Color)), firstOf(p.Series[i].ColorOpacity, opacity)
	}
	if base == "" {
		base = fmt.Sprintf("var(--tblr-chart-%d, %s)", i%5+1, palette[i%5])
	}
	return fmt.Sprintf("color-mix(in srgb, transparent, %s %s)", base, opacity)
}
