// Package kit holds the small helpers every tablerkit component shares: joining CSS classes and
// merging attribute maps. Components emit plain Tabler classes, so theming happens in CSS through
// Tabler's --tblr-* custom properties, never in Go.
package kit

import (
	"strings"

	"github.com/a-h/templ"
)

// Classes joins class lists, drops empty entries and duplicates, and keeps the first occurrence's order.
// Each argument may itself contain several space-separated classes.
func Classes(parts ...string) string {
	seen := map[string]bool{}
	var out []string
	for _, p := range parts {
		for _, c := range strings.Fields(p) {
			if !seen[c] {
				seen[c] = true
				out = append(out, c)
			}
		}
	}
	return strings.Join(out, " ")
}

// If returns class when cond holds, otherwise "".
func If(cond bool, class string) string {
	if cond {
		return class
	}
	return ""
}

// Attrs merges attribute maps left to right; later maps win. Nil maps are skipped.
func Attrs(maps ...templ.Attributes) templ.Attributes {
	out := templ.Attributes{}
	for _, m := range maps {
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// Default sets key to value unless the caller already set it, so callers can override every default.
func Default(a templ.Attributes, key string, value any) templ.Attributes {
	if a == nil {
		a = templ.Attributes{}
	}
	if _, ok := a[key]; !ok {
		a[key] = value
	}
	return a
}

// First returns the first props value or the zero value; components take ...Props so a bare call works.
func First[P any](props []P) P {
	var p P
	if len(props) > 0 {
		p = props[0]
	}
	return p
}
