// Package tooltip adds Tabler (Bootstrap) tooltips to any element: spread the attributes on it, e.g.
//
//	@button.Button(button.Props{Attributes: tooltip.Attrs("Save the draft", tooltip.Top)}) { Save }
//
// Tabler's script initialises every element with data-bs-toggle="tooltip".
package tooltip

import "github.com/a-h/templ"

type Placement string

const (
	Auto   Placement = ""
	Top    Placement = "top"
	Right  Placement = "right"
	Bottom Placement = "bottom"
	Left   Placement = "left"
)

// Attrs returns the attributes of a tooltip with a plain text title.
func Attrs(title string, placement Placement) templ.Attributes {
	a := templ.Attributes{"data-bs-toggle": "tooltip", "title": title}
	if placement != Auto {
		a["data-bs-placement"] = string(placement)
	}
	return a
}

// HTML returns the attributes of a tooltip whose title is HTML. The title is not escaped: never pass
// user input.
func HTML(title string, placement Placement) templ.Attributes {
	a := Attrs(title, placement)
	a["data-bs-html"] = "true"
	return a
}
