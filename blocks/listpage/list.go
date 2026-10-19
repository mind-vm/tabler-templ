// Package listpage is a block for list pages: a card with search and filters above a table, and the
// range ("1–25 of 80") with page numbers below it. Search, filters and page live in the URL
// (?q=…&status=…&page=…), so a list can be bookmarked and shared; the app filters its own rows and
// Paginate cuts out the page.
//
//	l := listpage.New(r, listpage.Config{Placeholder: "Name, city …"},
//		listpage.Filter{Name: "status", Label: "All states", Options: []listpage.Option{{"open", "Open"}, {"done", "Done"}}})
//	rows = listpage.Paginate(l, rows, func(c Customer) bool {
//		return (l.Value("status") == "" || c.Status == l.Value("status")) && listpage.Contains(l.Query, c.Name, c.City)
//	})
//
//	@listpage.Card(l) { @table.Table() { … rows … } }
package listpage

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DefaultPageSize is the number of rows per page when Config.PageSize is 0.
const DefaultPageSize = 25

type Option struct{ Value, Label string }

// Filter is a select above the list; the empty value means "all" and shows Label.
type Filter struct {
	Name, Label string
	Options     []Option
	Value       string // set by New from the URL
}

type Config struct {
	Placeholder string // of the search field
	PageParam   string // URL parameter of the page, default "page"
	QueryParam  string // URL parameter of the search, default "q"
	PageSize    int    // default DefaultPageSize
}

// List is the state of one list page.
type List struct {
	Config
	Path    string
	Query   string
	Filters []Filter
	// Set by Paginate:
	Page, Pages     int
	Total, From, To int
}

// New reads search, filters and page from the request. Unknown filter values count as "all".
func New(r *http.Request, c Config, filters ...Filter) *List {
	if c.PageParam == "" {
		c.PageParam = "page"
	}
	if c.QueryParam == "" {
		c.QueryParam = "q"
	}
	if c.PageSize <= 0 {
		c.PageSize = DefaultPageSize
	}
	q := r.URL.Query()
	l := &List{Config: c, Path: r.URL.Path, Query: strings.TrimSpace(q.Get(c.QueryParam)), Filters: filters}
	for i, f := range l.Filters {
		l.Filters[i].Value = ""
		for _, o := range f.Options {
			if o.Value == q.Get(f.Name) {
				l.Filters[i].Value = o.Value
			}
		}
	}
	l.Page, _ = strconv.Atoi(q.Get(c.PageParam))
	return l
}

// Value returns the chosen value of a filter, "" for all.
func (l *List) Value(name string) string {
	for _, f := range l.Filters {
		if f.Name == name {
			return f.Value
		}
	}
	return ""
}

// Active reports whether search or a filter narrows the list.
func (l *List) Active() bool {
	if l.Query != "" {
		return true
	}
	for _, f := range l.Filters {
		if f.Value != "" {
			return true
		}
	}
	return false
}

// Link leads to a page and keeps search and filters.
func (l *List) Link(page int) string {
	v := url.Values{}
	if l.Query != "" {
		v.Set(l.QueryParam, l.Query)
	}
	for _, f := range l.Filters {
		if f.Value != "" {
			v.Set(f.Name, f.Value)
		}
	}
	if page > 1 {
		v.Set(l.PageParam, strconv.Itoa(page))
	}
	if len(v) == 0 {
		return l.Path
	}
	return l.Path + "?" + v.Encode()
}

// Paginate keeps the rows for which keep is true (all rows when keep is nil) and returns the current
// page of them. A page beyond the last shows the last.
func Paginate[T any](l *List, rows []T, keep func(T) bool) []T {
	var hits []T
	for _, x := range rows {
		if keep == nil || keep(x) {
			hits = append(hits, x)
		}
	}
	size := l.PageSize
	if size <= 0 {
		size = DefaultPageSize
	}
	l.Total = len(hits)
	l.Pages = max(1, (l.Total+size-1)/size)
	l.Page = min(max(l.Page, 1), l.Pages)
	from := (l.Page - 1) * size
	to := min(from+size, l.Total)
	l.From, l.To = 0, 0
	if l.Total > 0 {
		l.From, l.To = from+1, to
	}
	return hits[from:to]
}

// Contains reports whether every word of the query occurs in one of the fields, ignoring case.
// Apps with their own matching rules (umlauts, spaces) use their own function instead.
func Contains(query string, fields ...string) bool {
	text := strings.ToLower(strings.Join(fields, " "))
	for _, w := range strings.Fields(strings.ToLower(query)) {
		if !strings.Contains(text, w) {
			return false
		}
	}
	return true
}

// Empty is the text for a table without rows: NoResults when search or filters are active,
// otherwise the list's own text (e.g. "No customers yet.").
func (l *List) Empty(none, noResults string) string {
	if l.Active() {
		return noResults
	}
	return none
}
