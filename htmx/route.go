package htmx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Endpoint is one route of a component: its method and path pattern in one value, from which both the
// mux pattern (Pattern) and the URLs (URL) are made. The path is written once, so the router and the
// markup cannot drift apart.
//
//	var Toggle = htmx.Route("POST", "/ui/todolist/{id}/toggle")
//
//	mux.HandleFunc(Toggle.Pattern(), h.Toggle)       // "POST /ui/todolist/{id}/toggle"
//	<input hx-post={ Toggle.URL(item.ID) } …/>        // "/ui/todolist/42/toggle"
type Endpoint struct {
	method string
	path   string
	params []string // names of the {…} segments, in order
}

// Route defines an endpoint. path is a net/http pattern path: segments like {id} take one value,
// {rest...} the rest of the path, and {$} marks the end. It panics on a malformed pattern, so a bad
// route fails when the package loads.
func Route(method, path string) Endpoint {
	if !strings.HasPrefix(path, "/") {
		panic("htmx: route path must start with /: " + path)
	}
	e := Endpoint{method: strings.ToUpper(method), path: path}
	for seg := range strings.SplitSeq(path, "/") {
		if !strings.HasPrefix(seg, "{") {
			continue
		}
		if !strings.HasSuffix(seg, "}") {
			panic("htmx: a {…} must be a whole segment: " + path)
		}
		name := strings.TrimSuffix(strings.Trim(seg, "{}"), "...")
		if name != "$" {
			e.params = append(e.params, name)
		}
	}
	return e
}

// Method is the endpoint's HTTP method, e.g. "POST".
func (e Endpoint) Method() string { return e.method }

// Path is the path pattern, e.g. "/ui/todolist/{id}/toggle".
func (e Endpoint) Path() string { return e.path }

// Params are the names of the path's parameters, in order.
func (e Endpoint) Params() []string { return append([]string(nil), e.params...) }

// Pattern is the pattern for http.ServeMux, method included: "POST /ui/todolist/{id}/toggle".
func (e Endpoint) Pattern() string {
	if e.method == "" {
		return e.path
	}
	return e.method + " " + e.path
}

// URL fills the path's parameters in order. Values are escaped (a {rest...} parameter keeps its
// slashes). It panics when the number of values does not match the parameters: a programming error
// that any render test of the template finds.
func (e Endpoint) URL(values ...string) string {
	if len(values) != len(e.params) {
		panic(fmt.Sprintf("htmx: %s takes %d values (%s), got %d", e.path, len(e.params), strings.Join(e.params, ", "), len(values)))
	}
	var b strings.Builder
	i := 0
	for n, seg := range strings.Split(e.path, "/") {
		if n > 0 {
			b.WriteByte('/')
		}
		switch {
		case seg == "{$}":
		case strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "...}"):
			parts := strings.Split(values[i], "/")
			for j, p := range parts {
				parts[j] = url.PathEscape(p)
			}
			b.WriteString(strings.Join(parts, "/"))
			i++
		case strings.HasPrefix(seg, "{"):
			b.WriteString(url.PathEscape(values[i]))
			i++
		default:
			b.WriteString(seg)
		}
	}
	return b.String()
}

// Query appends query parameters to a URL; empty values are left out.
//
//	List.URL() + htmx.Query("q", query, "page", "2")
func Query(kv ...string) string {
	v := url.Values{}
	for i := 0; i+1 < len(kv); i += 2 {
		if kv[i+1] != "" {
			v.Set(kv[i], kv[i+1])
		}
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

// ID is an element id for a component instance: base for the default instance (""), base-inst
// otherwise, with inst reduced to letters, digits, "-" and "_". Components define their IDs as
// functions on top of it, so they can be embedded more than once on a page.
//
//	func RootID(inst string) string { return htmx.ID("todolist", inst) }
func ID(base, inst string) string {
	if inst == "" {
		return base
	}
	b := []byte(inst)
	for i, c := range b {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
			b[i] = '-'
		}
	}
	return base + "-" + string(b)
}

// Target is the CSS selector of an id, for hx-target: "#" + id.
func Target(id string) string { return "#" + id }

// Vals is the JSON for hx-vals from key-value pairs: Vals("instance", "b") is {"instance":"b"}.
func Vals(kv ...string) string {
	m := map[string]string{}
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	b, _ := json.Marshal(m)
	return string(b)
}

// TB is the part of testing.TB that CheckRoutes needs.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
}

// CheckRoutes is a contract test: every endpoint must be registered in mux under its own pattern, and
// its URL (with sample values) must be routed to that pattern. It catches a route that Register
// forgot, registered under another pattern, or that another pattern shadows.
//
//	func TestContract(t *testing.T) {
//		mux := http.NewServeMux()
//		todolist.Register(mux, fakeService{})
//		htmx.CheckRoutes(t, mux, todolist.List, todolist.Create, todolist.Toggle)
//	}
func CheckRoutes(t TB, mux *http.ServeMux, endpoints ...Endpoint) {
	t.Helper()
	for _, e := range endpoints {
		vals := make([]string, len(e.params))
		for i := range vals {
			vals[i] = "x" + fmt.Sprint(i+1)
		}
		method := e.method
		if method == "" {
			method = http.MethodGet
		}
		r, err := http.NewRequest(method, e.URL(vals...), nil)
		if err != nil {
			t.Errorf("htmx: %s: %v", e.Pattern(), err)
			continue
		}
		if _, pattern := mux.Handler(r); pattern != e.Pattern() {
			if pattern == "" {
				t.Errorf("htmx: %s is not registered", e.Pattern())
			} else {
				t.Errorf("htmx: %s %s goes to %q, want %q", method, r.URL.Path, pattern, e.Pattern())
			}
		}
	}
}
