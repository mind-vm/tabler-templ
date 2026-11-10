package chart

// RawProps gives p a complete ECharts option instead of the one built from its fields (only ID, Label,
// Class, Attributes, Height, Size and Group stay in use), with formatters named by path, e.g.
// {"yAxis.0.axisLabel.formatter": "percent"}. For charts Props has no field for; Gauge builds one.
func RawProps(p Props, option map[string]any, formatters map[string]string) Props {
	p.raw, p.formatters = option, formatters
	return p
}
