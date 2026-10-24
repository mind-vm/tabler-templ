// Package confetti pours Tabler's confetti over the page to celebrate a finished task or a signup:
// Attrs make any button burst it on click, Burst bursts it when it is rendered (on a page, or in content
// swapped in by htmx). Tabler's own script draws it; no plugin to load.
//
//	@button.Button(button.Props{Attributes: confetti.Attrs(confetti.Props{Count: 150})}) { Celebrate }
//	@confetti.Burst(confetti.Props{})
package confetti

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/mind-vm/tabler-templ/internal/plugin"
)

type Props struct {
	Count    int      // pieces
	Duration int      // how long they pour, in ms
	Speed    float64  // 1 is the default; lower floats, higher drops
	Colors   []string // CSS colors; default Tabler's palette
	Target   string   // a selector, e.g. "#saved": the start.bs.confetti and end.bs.confetti events fire there
}

// Attrs make an element burst confetti on click.
func Attrs(p Props) templ.Attributes {
	a := templ.Attributes{"data-bs-toggle": "confetti"}
	if p.Count > 0 {
		a["data-bs-count"] = strconv.Itoa(p.Count)
	}
	if p.Duration > 0 {
		a["data-bs-duration"] = strconv.Itoa(p.Duration)
	}
	if p.Speed > 0 {
		a["data-bs-speed"] = strconv.FormatFloat(p.Speed, 'f', -1, 64)
	}
	if p.Colors != nil {
		a["data-bs-colors"] = strings.Join(p.Colors, ", ")
	}
	if p.Target != "" {
		a["data-bs-target"] = p.Target
	}
	return a
}

func (p Props) config() map[string]any {
	c := map[string]any{}
	if p.Count > 0 {
		c["count"] = p.Count
	}
	if p.Duration > 0 {
		c["duration"] = p.Duration
	}
	if p.Speed > 0 {
		c["speed"] = p.Speed
	}
	if p.Colors != nil {
		c["colors"] = p.Colors
	}
	return c
}

// Burst bursts confetti once it is rendered: when the page is ready, or right away in content swapped in
// later. Its events fire on Target, or on the body.
func Burst(p Props) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		target := "document.body"
		if p.Target != "" {
			target = "document.querySelector(" + plugin.JSON(p.Target) + ")"
		}
		_, err := io.WriteString(w, "<script>(function () {\n\tvar burst = function () { window.tabler && tabler.Confetti.getOrCreateInstance("+target+", "+
			plugin.JSON(p.config())+").burst(); };\n\tdocument.readyState !== 'loading' ? burst() : document.addEventListener('DOMContentLoaded', burst, { once: true });\n})();</script>")
		return err
	})
}
