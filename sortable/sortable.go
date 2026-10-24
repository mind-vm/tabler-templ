// Package sortable lets the user reorder the children of any element by drag and drop: list items, table
// rows, cards of a kanban board, tags. Tabler's Sortable component (in tabler.js) runs SortableJS with the
// Options as JSON. Load the plugin with Library in the head, before Tabler's script.
//
//	<ul class="list-group" { sortable.Attrs(sortable.Options{"animation": 150, "handle": ".sortable-handle"})... }>
package sortable

import (
	"github.com/a-h/templ"
	"github.com/mind-vm/tabler-templ/internal/plugin"
	"github.com/mind-vm/tabler-templ/versions"
)

// Version is the SortableJS version Library loads.
const Version = versions.SortableJS

// LibraryURL is SortableJS on jsDelivr (with the MultiDrag and Swap plugins).
const LibraryURL = "https://cdn.jsdelivr.net/npm/sortablejs@" + Version + "/Sortable.min.js"

// Options are SortableJS options, e.g. "group", "animation", "handle", "filter", "direction".
type Options map[string]any

// Attrs make an element's children sortable. Render Script once for elements swapped in by htmx (Library
// does).
func Attrs(o Options) templ.Attributes {
	if o == nil {
		return templ.Attributes{"data-sortable": ""}
	}
	return templ.Attributes{"data-sortable": plugin.JSON(o)}
}

// Script starts sortable elements swapped in by htmx.
func Script() templ.Component {
	return plugin.Tabler("Sortable", "[data-sortable]")
}

// Library loads SortableJS; put it in the head, before Tabler's script.
func Library() templ.Component {
	return templ.Join(plugin.Script(LibraryURL), Script())
}
