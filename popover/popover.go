// Package popover adds Tabler (Bootstrap) popovers to any element: spread the attributes on it, e.g.
//
//	@button.Button(button.Props{Type: "button", Attributes: popover.Attrs(popover.Props{Title: "Hint", Content: "…"})}) { Help }
//
// Tabler's script initialises every element with data-bs-toggle="popover".
package popover

import "github.com/a-h/templ"

type Placement string

const (
	Auto   Placement = ""
	Top    Placement = "top"
	Right  Placement = "right"
	Bottom Placement = "bottom"
	Left   Placement = "left"
)

// Trigger is the event that opens the popover; the default is a click.
type Trigger string

const (
	Click  Trigger = ""
	Hover  Trigger = "hover"
	Focus  Trigger = "focus"
	Manual Trigger = "manual"
)

type Props struct {
	Title     string
	Content   string
	Placement Placement
	Trigger   Trigger
	// Body attaches the popover to <body> (data-bs-container="body"), so overflow does not clip it.
	Body bool
	// HTML treats Title and Content as HTML; they are not escaped, never pass user input.
	HTML bool
}

// Attrs returns the attributes of a popover.
func Attrs(p Props) templ.Attributes {
	a := templ.Attributes{"data-bs-toggle": "popover", "data-bs-content": p.Content}
	if p.Title != "" {
		a["title"] = p.Title
	}
	if p.Placement != Auto {
		a["data-bs-placement"] = string(p.Placement)
	}
	if p.Trigger != Click {
		a["data-bs-trigger"] = string(p.Trigger)
	}
	if p.Body {
		a["data-bs-container"] = "body"
	}
	if p.HTML {
		a["data-bs-html"] = "true"
	}
	return a
}
