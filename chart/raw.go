package chart

// RawProps gives p a complete ApexCharts config instead of the one built from its fields (only ID,
// Label, Class, Attributes and Size stay in use), with formatters named by path, e.g.
// {"plotOptions.radialBar.dataLabels.value.formatter": "percent"}. For charts Props has no field for,
// such as a gauge.
func RawProps(p Props, config map[string]any, formatters map[string]string) Props {
	p.raw, p.formatters = config, formatters
	return p
}
