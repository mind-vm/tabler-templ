// Package contacttable is the list of contacts: search, a status filter, the order and pages, all in the
// URL (blocks/listpage), so a list can be bookmarked. Without JavaScript the search is a plain GET form
// and the pages are links; with htmx the same form and links are boosted into the table (hx-boost) and
// the page's address follows (HX-Push-Url). The table refreshes itself on the other components' events.
package contacttable

import (
	"context"
	"net/http"
	"strings"

	"github.com/mind-vm/tabler-templ/blocks/listpage"
	"github.com/mind-vm/tabler-templ/examples/crud/contacts"
	"github.com/mind-vm/tabler-templ/examples/crud/ui/events"
	"github.com/mind-vm/tabler-templ/htmx"
)

// List is the component's route: the table for the query in the URL.
var List = htmx.Route("GET", "/crud/ui/contacttable")

// RootID is the table's id.
const RootID = "contacttable"

// trigger: the rows may change after these events (a status change only touches its own cell).
func trigger() string { return htmx.On(events.ListChanged...) }

// Service is what the component needs from the domain.
type Service interface {
	List(ctx context.Context, q contacts.Query) ([]contacts.Contact, error)
}

type Options struct {
	// Page is the page that shows the table, e.g. "/crud": its address takes the list's query, and a
	// request to List without htmx is sent there.
	Page string
	// Detail is the link of a contact's page.
	Detail func(id string) string
	// PageSize, default 10.
	PageSize int
}

type Props struct {
	List   *listpage.List
	Rows   []contacts.Contact
	Detail func(id string) string
}

// Table builds the table's props from the request; the page that embeds the table calls it too.
type Table struct {
	svc Service
	o   Options
}

// Register wires the component's endpoint and returns the Table for the page.
func Register(mux *http.ServeMux, svc Service, o Options) *Table {
	if o.PageSize == 0 {
		o.PageSize = 10
	}
	t := &Table{svc: svc, o: o}
	mux.HandleFunc(List.Pattern(), t.serve)
	return t
}

var statusFilter = listpage.Filter{Name: "status", Label: "All states", Options: func() []listpage.Option {
	out := make([]listpage.Option, len(contacts.Statuses))
	for i, s := range contacts.Statuses {
		out[i] = listpage.Option{Value: string(s), Label: s.Label()}
	}
	return out
}()}

var sortFilter = listpage.Filter{Name: "sort", Label: "By name", Options: []listpage.Option{
	{Value: "company", Label: "By company"}, {Value: "updated", Label: "Recently updated"}}}

// Props reads the query from r and loads the page of contacts.
func (t *Table) Props(r *http.Request) (Props, error) {
	l := listpage.New(r, listpage.Config{Placeholder: "Name, company, e-mail, tag …", PageSize: t.o.PageSize}, statusFilter, sortFilter)
	l.Path = List.URL() // the form and the page links ask the component, not the page
	rows, err := t.svc.List(r.Context(), contacts.Query{Search: l.Query, Status: contacts.Status(l.Value("status")), Sort: l.Value("sort")})
	if err != nil {
		return Props{}, err
	}
	return Props{List: l, Rows: listpage.Paginate(l, rows, nil), Detail: t.o.Detail}, nil
}

// pageURL is the page's address with the list's query, without empty parameters.
func (t *Table) pageURL(r *http.Request, l *listpage.List) string {
	if l != nil {
		if _, q, ok := strings.Cut(l.Link(l.Page), "?"); ok {
			return t.o.Page + "?" + q
		}
		return t.o.Page
	}
	if r.URL.RawQuery == "" {
		return t.o.Page
	}
	return t.o.Page + "?" + r.URL.RawQuery
}

func (t *Table) serve(w http.ResponseWriter, r *http.Request) {
	if !htmx.IsHTMX(r) {
		// Opened directly (bookmark, new tab): the page shows the table with the same query.
		htmx.Redirect(w, r, t.pageURL(r, nil))
		return
	}
	p, err := t.Props(r)
	if err != nil {
		htmx.Error(w, r, err)
		return
	}
	opts := []htmx.Option{}
	if r.Header.Get("HX-Boosted") == "true" { // a search or a page link, not a refresh after an event
		opts = append(opts, htmx.PushURL(t.pageURL(r, p.List)))
	}
	htmx.Respond(w, r, View(p), opts...)
}
