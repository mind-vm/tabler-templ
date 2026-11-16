// Package mcpapp is the demo of package mcpapp: a small order list as an MCP App. Its tools are plain Go
// functions with the shapes of MCP (Tools, Resource), so an MCP server (e.g. the official Go SDK) can
// serve them as they are:
//
//   - list_orders (model and view): the open orders; shows the resource ui://orders/list, and its result
//     carries the table as HTML for the view and a summary as text for the model.
//   - complete_order (view only): marks an order as done; the result is the order's new row.
//   - find_orders (view only): the search form's tool; the result is the table's body.
//
// cmd/demo mounts it at /mcpapp with a host simulation: the page plays the host (Claude, ChatGPT …), loads
// the resource into a sandboxed iframe, answers the bridge's requests and calls the tools through
// POST /mcpapp/call.
package mcpapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/mcpapp"
)

// Home is the host simulation.
const Home = "/mcpapp"

// ResourceURI is the order list's UI resource.
const ResourceURI = "ui://orders/list"

// Order is one order.
type Order struct {
	ID       string
	Customer string
	Item     string
	Amount   int // euros
	Due      string
	Overdue  bool
	Done     bool
}

// Store holds the orders in memory.
type Store struct {
	mu     sync.Mutex
	orders []Order
}

// NewStore is a store with sample orders.
func NewStore() *Store {
	return &Store{orders: []Order{
		{ID: "A-1001", Customer: "Becker GmbH", Item: "Shop sign", Amount: 2400, Due: "2026-11-10", Overdue: true},
		{ID: "A-1002", Customer: "Müller Bau", Item: "Site fence banners", Amount: 860, Due: "2026-11-14", Overdue: true},
		{ID: "A-1003", Customer: "Café Lindenhof", Item: "Menu boards", Amount: 540, Due: "2026-11-20"},
		{ID: "A-1004", Customer: "Nordwind Reisen", Item: "Window film", Amount: 1290, Due: "2026-11-24"},
		{ID: "A-1005", Customer: "Grünwerk", Item: "Vehicle lettering", Amount: 3150, Due: "2026-12-02"},
	}}
}

func (s *Store) find(q string) []Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	q = strings.ToLower(strings.TrimSpace(q))
	var out []Order
	for _, o := range s.orders {
		if q == "" || strings.Contains(strings.ToLower(o.ID+" "+o.Customer+" "+o.Item), q) {
			out = append(out, o)
		}
	}
	return out
}

func (s *Store) complete(id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := slices.IndexFunc(s.orders, func(o Order) bool { return o.ID == id })
	if i < 0 {
		return Order{}, false
	}
	s.orders[i].Done, s.orders[i].Overdue = true, false
	return s.orders[i], true
}

// Tool is a tool's declaration as in tools/list.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
	Meta        map[string]any `json:"_meta,omitempty"`
}

// Content is a text content block.
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Result is a CallToolResult.
type Result struct {
	Content []Content      `json:"content"`
	Meta    map[string]any `json:"_meta,omitempty"`
	IsError bool           `json:"isError,omitempty"`
}

// App is the demo's MCP side: its tools and its resource.
type App struct{ store *Store }

// New is the app on store.
func New(store *Store) *App { return &App{store: store} }

func object(props map[string]any, required ...string) map[string]any {
	s := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		s["required"] = required
	}
	return s
}

// Tools are the declarations for tools/list.
func (a *App) Tools() []Tool {
	str := map[string]any{"type": "string"}
	return []Tool{
		{Name: "list_orders", Description: "Lists the orders, optionally filtered by a search text, and shows them as a table.",
			InputSchema: object(map[string]any{"query": str}), Meta: mcpapp.ToolMeta(ResourceURI)},
		{Name: "find_orders", Description: "Search form of the order table.",
			InputSchema: object(map[string]any{"query": str}), Meta: mcpapp.AppToolMeta()},
		{Name: "complete_order", Description: "Marks an order as done.",
			InputSchema: object(map[string]any{"id": str}, "id"), Meta: mcpapp.AppToolMeta()},
	}
}

// Resource is the content of resources/read for uri.
func (a *App) Resource(ctx context.Context, uri string) (map[string]any, error) {
	if uri != ResourceURI {
		return nil, fmt.Errorf("unknown resource %q", uri)
	}
	html, err := mcpapp.Page(ctx, mcpapp.Props{Title: "Orders", App: "orders", Version: "1.0.0", HostStyles: true}, shell())
	if err != nil {
		return nil, err
	}
	return map[string]any{"uri": uri, "mimeType": mcpapp.MIMEType, "text": html, "_meta": mcpapp.ResourceMeta{}.Meta()}, nil
}

// Call runs the tool name with args.
func (a *App) Call(ctx context.Context, name string, args map[string]any) Result {
	str := func(k string) string { s, _ := args[k].(string); return s }
	html := func(c templ.Component, text string) Result {
		h, err := mcpapp.Render(ctx, c)
		if err != nil {
			return failed(err.Error())
		}
		return Result{Content: []Content{{Type: "text", Text: text}}, Meta: mcpapp.ResultMeta(h)}
	}
	switch name {
	case "list_orders":
		q := str("query")
		os := a.store.find(q)
		return html(table(q, os), summary(os))
	case "find_orders":
		os := a.store.find(str("query"))
		return html(rows(os), summary(os))
	case "complete_order":
		o, ok := a.store.complete(str("id"))
		if !ok {
			return failed("no order " + str("id"))
		}
		return html(row(o), o.ID+" is done.")
	}
	return failed("unknown tool " + name)
}

func failed(msg string) Result {
	return Result{Content: []Content{{Type: "text", Text: msg}}, IsError: true}
}

// summary is what the model reads instead of the table.
func summary(os []Order) string {
	if len(os) == 0 {
		return "No orders."
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d orders:\n", len(os))
	for _, o := range os {
		state := "open"
		switch {
		case o.Done:
			state = "done"
		case o.Overdue:
			state = "overdue"
		}
		fmt.Fprintf(&b, "- %s %s, %s, %d EUR, due %s, %s\n", o.ID, o.Customer, o.Item, o.Amount, o.Due, state)
	}
	return b.String()
}

// Register mounts the host simulation and the endpoints it calls on mux.
func Register(mux *http.ServeMux) *App {
	a := New(NewStore())
	mux.Handle("GET "+Home, templ.Handler(hostPage()))
	mux.HandleFunc("GET "+Home+"/tools", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, a.Tools()) })
	mux.HandleFunc("GET "+Home+"/resource", func(w http.ResponseWriter, r *http.Request) {
		res, err := a.Resource(r.Context(), r.URL.Query().Get("uri"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJSON(w, res)
	})
	mux.HandleFunc("POST "+Home+"/call", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, a.Call(r.Context(), req.Name, req.Arguments))
	})
	return a
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
