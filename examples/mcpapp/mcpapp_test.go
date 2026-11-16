package mcpapp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mind-vm/tabler-templ/mcpapp"
)

func html(t *testing.T, r Result) string {
	t.Helper()
	if r.IsError {
		t.Fatalf("error result: %+v", r)
	}
	h, _ := r.Meta[mcpapp.MetaHTML].(string)
	if h == "" {
		t.Fatalf("no HTML in %+v", r)
	}
	return h
}

func TestTools(t *testing.T) {
	a := New(NewStore())
	byName := map[string]Tool{}
	for _, tool := range a.Tools() {
		byName[tool.Name] = tool
	}
	ui := func(name string) map[string]any { return byName[name].Meta["ui"].(map[string]any) }
	if ui("list_orders")["resourceUri"] != ResourceURI || ui("list_orders")["visibility"] != nil {
		t.Errorf("list_orders: %v", ui("list_orders"))
	}
	for _, name := range []string{"find_orders", "complete_order"} {
		if v, _ := ui(name)["visibility"].([]string); len(v) != 1 || v[0] != mcpapp.VisibilityApp {
			t.Errorf("%s is not only for the view: %v", name, ui(name))
		}
	}

	res, err := a.Resource(context.Background(), ResourceURI)
	if err != nil {
		t.Fatal(err)
	}
	if res["mimeType"] != mcpapp.MIMEType || !strings.Contains(res["text"].(string), "data-mcp-result") {
		t.Errorf("resource: %v", res["mimeType"])
	}
	if _, err := a.Resource(context.Background(), "ui://other"); err == nil {
		t.Error("unknown resource served")
	}
}

func TestCall(t *testing.T) {
	ctx := context.Background()
	a := New(NewStore())

	r := a.Call(ctx, "list_orders", map[string]any{})
	if h := html(t, r); strings.Count(h, `<tr id="order-`) != 5 || !strings.Contains(h, `data-mcp-tool="find_orders"`) {
		t.Errorf("list_orders HTML:\n%s", h)
	}
	if !strings.Contains(r.Content[0].Text, "A-1001 Becker GmbH") || strings.Contains(r.Content[0].Text, "<") {
		t.Errorf("the model's text: %q", r.Content[0].Text)
	}

	if h := html(t, a.Call(ctx, "find_orders", map[string]any{"query": "Lindenhof"})); strings.Count(h, "<tr") != 1 {
		t.Errorf("find_orders HTML:\n%s", h)
	}

	h := html(t, a.Call(ctx, "complete_order", map[string]any{"id": "A-1001"}))
	if !strings.HasPrefix(h, `<tr id="order-A-1001"`) || !strings.Contains(h, ">done<") || strings.Contains(h, "complete_order") {
		t.Errorf("complete_order HTML:\n%s", h)
	}
	if !strings.Contains(a.Call(ctx, "list_orders", nil).Content[0].Text, "2400 EUR, due 2026-11-10, done") {
		t.Error("completed order not kept")
	}

	for _, c := range []struct {
		name string
		args map[string]any
	}{{"complete_order", map[string]any{"id": "A-9"}}, {"drop_tables", nil}} {
		if r := a.Call(ctx, c.name, c.args); !r.IsError {
			t.Errorf("%s: want an error, got %+v", c.name, r)
		}
	}
}

func TestHost(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	res, err := http.Get(srv.URL + Home)
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("host page: %v %v", err, res.Status)
	}
	res.Body.Close()

	res, err = http.Post(srv.URL+Home+"/call", "application/json", strings.NewReader(`{"name":"complete_order","arguments":{"id":"A-1002"}}`))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var r Result
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html(t, r), "order-A-1002") {
		t.Errorf("call: %+v", r)
	}

	res, err = http.Get(srv.URL + Home + "/resource?uri=ui://nope")
	if err != nil || res.StatusCode != http.StatusNotFound {
		t.Errorf("unknown resource: %v %v", err, res.Status)
	}
}
