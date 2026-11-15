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
// page of them. A page beyond the last shows the last. For lists the database pages, see Window.
func Paginate[T any](l *List, rows []T, keep func(T) bool) []T {
	var hits []T
	for _, x := range rows {
		if keep == nil || keep(x) {
			hits = append(hits, x)
		}
	}
	l.SetTotal(len(hits))
	from, size := l.Window()
	return hits[from:min(from+size, l.Total)]
}

func (l *List) size() int {
	if l.PageSize <= 0 {
		return DefaultPageSize
	}
	return l.PageSize
}

// Window is the part of the list the database fetches for the page asked for: OFFSET and LIMIT. With
// Patterns (or Words), Value and SetTotal, the database does the search, the filter, the page and the
// count in fixed queries instead of the app holding every row (GS-51 of the design concept):
//
//	l := listpage.New(r, cfg, statusFilter)
//	offset, limit := l.Window()
//	rows, total := store.Customers(ctx, l.Patterns(), l.Value("status"), offset, limit)
//	if !l.SetTotal(total) { // the page asked for lies beyond the last one
//		offset, limit = l.Window()
//		rows, total = store.Customers(ctx, l.Patterns(), l.Value("status"), offset, limit)
//		l.SetTotal(total)
//	}
//
// SetTotal is called again after the second fetch: with count(*) OVER (), a page beyond the last brings no
// row and so no count (0), and the list goes to page 1. A separate count query avoids the second fetch.
// The query, e.g. for sqlc with PostgreSQL, where every pattern must match the folded text:
//
//	SELECT c.*, count(*) OVER () AS total FROM customers c
//	WHERE NOT EXISTS (SELECT 1 FROM unnest(@patterns::text[]) AS p WHERE <FoldSQL("c.name || ' ' || c.city")> NOT LIKE p)
//	  AND (@status::text = '' OR c.status = @status)
//	ORDER BY c.name LIMIT @lim OFFSET @off
func (l *List) Window() (offset, limit int) {
	return (max(l.Page, 1) - 1) * l.size(), l.size()
}

// Words are the words of the query as Contains compares them (Fold); a row matches when every word occurs
// in its searched text, folded the same way. Empty without a query.
func (l *List) Words() []string {
	return strings.Fields(Fold(l.Query))
}

// Patterns are Words as LIKE patterns (LikePattern): a row matches when its folded text is LIKE every one.
func (l *List) Patterns() []string {
	ws := l.Words()
	ps := make([]string, len(ws))
	for i, w := range ws {
		ps[i] = LikePattern(w)
	}
	return ps
}

// SetTotal records the number of all rows that match (from the database), and sets Pages, Page, From and
// To for the footer. A page beyond the last becomes the last; then it reports false, and Window gives the
// rows to fetch instead.
func (l *List) SetTotal(total int) bool {
	asked := max(l.Page, 1)
	l.Total = max(total, 0)
	l.Pages = max(1, (l.Total+l.size()-1)/l.size())
	l.Page = min(asked, l.Pages)
	from := (l.Page - 1) * l.size()
	l.From, l.To = 0, 0
	if l.Total > 0 {
		l.From, l.To = from+1, min(from+l.size(), l.Total)
	}
	return l.Page == asked
}

// Contains reports whether every word of the query occurs in one of the fields, ignoring case and German
// umlauts: ä, ö, ü and ß match ae, oe, ue and ss both ways ("muehl" finds "Mühlhausen", "Müller" finds
// "Mueller"). Apps with other matching rules use their own function instead.
func Contains(query string, fields ...string) bool {
	text := Fold(strings.Join(fields, " "))
	for _, w := range strings.Fields(Fold(query)) {
		if !strings.Contains(text, w) {
			return false
		}
	}
	return true
}

// Fold is text as Contains compares it: lower case, with ä, ö, ü, ß written as ae, oe, ue, ss. Apps that
// search in SQL fold the column the same way with FoldSQL.
func Fold(s string) string { return umlauts.Replace(strings.ToLower(s)) }

// folds are the letters Fold writes out; Fold and FoldSQL both come from this table.
var folds = [][2]string{{"ä", "ae"}, {"ö", "oe"}, {"ü", "ue"}, {"ß", "ss"}}

var umlauts = func() *strings.Replacer {
	var args []string
	for _, f := range folds {
		args = append(args, f[0], f[1])
	}
	return strings.NewReplacer(args...)
}()

// FoldSQL is the SQL expression that folds expr like Fold, e.g. FoldSQL("c.name || ' ' || c.city") for
// PostgreSQL or SQLite: lower() and one replace() per letter. lower() must know the umlauts (PostgreSQL
// with a UTF-8 database does; SQLite's built-in lower() does not, so fold Ä, Ö, Ü there as well or keep a
// folded column). expr is SQL written by the app, never input.
func FoldSQL(expr string) string {
	out := "lower(" + expr + ")"
	for _, f := range folds {
		out = "replace(" + out + ", '" + f[0] + "', '" + f[1] + "')"
	}
	return out
}

// LikePattern is the LIKE pattern that finds the word anywhere: % and _ in the word stand for themselves
// (escaped with a backslash, the default escape of PostgreSQL's LIKE; SQLite needs ESCAPE '\').
func LikePattern(word string) string {
	return "%" + likeEscaper.Replace(word) + "%"
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

// Empty is the text for a table without rows: NoResults when search or filters are active,
// otherwise the list's own text (e.g. "No customers yet.").
func (l *List) Empty(none, noResults string) string {
	if l.Active() {
		return noResults
	}
	return none
}
