package chart

import (
	"math"
	"strconv"
	"strings"
)

// Format is how numbers are written: up to Decimals decimals (trailing zeros dropped), Decimal as the
// decimal separator and Group between groups of three digits ("" for none). The zero value is the
// default: ECharts' own English numbers (12,400.5). Props.Format sets it for an ECharts chart; the SVG
// charts (package chart/svg) use the same type.
type Format struct {
	Decimals int    `json:"decimals"`
	Decimal  string `json:"decimal"`
	Group    string `json:"group"`
}

// Number is a Format with up to maxDecimals decimals, decimal as the decimal separator and group
// between groups of three digits ("" for none).
func Number(maxDecimals int, decimal, group string) Format {
	if decimal == "" {
		decimal = "."
	}
	return Format{Decimals: max(maxDecimals, 0), Decimal: decimal, Group: group}
}

var (
	English = Number(2, ".", ",") // 12,400.5
	German  = Number(2, ",", ".") // 12.400,5
)

// IsZero reports whether f is unset, the default format.
func (f Format) IsZero() bool { return f == Format{} }

// Format writes v; the zero Format writes it in English.
func (f Format) Format(v float64) string {
	if f.IsZero() {
		f = English
	}
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ""
	}
	s := strconv.FormatFloat(v, 'f', max(f.Decimals, 0), 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	whole, frac, _ := strings.Cut(s, ".")
	if f.Group != "" {
		for i := len(whole) - 3; i > 0; i -= 3 {
			whole = whole[:i] + f.Group + whole[i:]
		}
	}
	if frac != "" {
		whole += f.Decimal + frac
	}
	if neg && strings.Trim(whole, "0"+f.Decimal+f.Group) != "" {
		whole = "-" + whole
	}
	return whole
}
