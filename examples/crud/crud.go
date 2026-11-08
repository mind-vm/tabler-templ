// Package crud is the CRUD demo app for the co-located HTMX components of the concept: contacts with a
// list, a form, inline editing, deleting with undo, a live count and a detail page. Each component is
// one package under ui/ with its contract (routes, ids, events), handlers, markup and tests; this file
// only composes them.
//
// cmd/demo mounts it at /crud. To run it alone:
//
//	mux := http.NewServeMux()
//	crud.Register(mux)
//	http.ListenAndServe("localhost:8080", http.NewCrossOriginProtection().Handler(mux))
package crud

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactcount"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactdelete"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactdetail"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactform"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contactrow"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/contacttable"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/layout"
	"github.com/mind-vm/tabler-templ/htmx"
)

// Home is the list page.
const Home = "/crud"

// Register mounts the demo with a fresh in-memory store.
func Register(mux *http.ServeMux) { RegisterWith(mux, contacts.NewStore()) }

// RegisterWith mounts the demo on store.
func RegisterWith(mux *http.ServeMux, store *contacts.Store) {
	lay := layout.Layout{Counts: store, Home: Home}

	contactcount.Register(mux, store)
	contactrow.Register(mux, store)
	contactform.Register(mux, store, contactform.Options{Frame: lay.Frame, Done: Home})
	contactdelete.Register(mux, store, contactdelete.Options{Frame: lay.Frame, Done: Home})
	contactdetail.Register(mux, store, contactdetail.Options{Frame: lay.Frame, Back: Home})
	table := contacttable.Register(mux, store, contacttable.Options{Page: Home, Detail: func(id string) string { return contactdetail.Page.URL(id) }})

	mux.HandleFunc("GET "+Home, func(w http.ResponseWriter, r *http.Request) {
		p, err := table.Props(r)
		if err != nil {
			htmx.Error(w, r, err)
			return
		}
		page := layout.Page{Title: "Contacts", List: true, Status: p.List.Value("status"), Actions: contactform.NewButton(),
			Crumbs: []shell.Crumb{{Text: "Contacts"}}}
		w.Header().Add("Vary", "HX-Request")
		_ = lay.Page(page).Render(templ.WithChildren(r.Context(), contacttable.View(p)), w)
	})
}
