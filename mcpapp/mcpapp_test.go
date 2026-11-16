package mcpapp

import (
	"context"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/cdn"
	"github.com/mind-vm/tabler-templ/internal/csp"
)

func jsonOf(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestToolMeta(t *testing.T) {
	if got := jsonOf(t, ToolMeta("ui://orders/list")); got != `{"ui":{"resourceUri":"ui://orders/list"}}` {
		t.Errorf("ToolMeta = %s", got)
	}
	if got := jsonOf(t, ToolMeta("ui://x", VisibilityApp)); got != `{"ui":{"resourceUri":"ui://x","visibility":["app"]}}` {
		t.Errorf("ToolMeta app = %s", got)
	}
	if got := jsonOf(t, AppToolMeta()); got != `{"ui":{"visibility":["app"]}}` {
		t.Errorf("AppToolMeta = %s", got)
	}
	if m := ResultMeta("<p>x</p>"); m[MetaHTML] != "<p>x</p>" || len(m) != 1 {
		t.Errorf("ResultMeta = %v", m)
	}
}

func TestResourceMeta(t *testing.T) {
	defer func(b string) { cdn.Base = b }(cdn.Base)

	cdn.Base = cdn.JSDelivr
	if got := jsonOf(t, ResourceMeta{}.Meta()); got != `{"ui":{"csp":{"resourceDomains":["https://cdn.jsdelivr.net"]}}}` {
		t.Errorf("default = %s", got)
	}
	no := false
	got := jsonOf(t, ResourceMeta{NoCSP: true, PrefersBorder: &no, Domain: "orders.example"}.Meta())
	if got != `{"ui":{"domain":"orders.example","prefersBorder":false}}` {
		t.Errorf("NoCSP = %s", got)
	}
	got = jsonOf(t, ResourceMeta{CSP: CSP{ConnectDomains: []string{"https://api.example.com"}}}.Meta())
	if got != `{"ui":{"csp":{"connectDomains":["https://api.example.com"]}}}` {
		t.Errorf("own CSP = %s", got)
	}

	cdn.Base = "/static/npm/"
	if c := DefaultCSP(); !reflect.DeepEqual(c, CSP{}) {
		t.Errorf("relative cdn.Base: DefaultCSP = %+v, want none", c)
	}
	cdn.Base = "https://app.example.com:8443/npm/"
	if c := DefaultCSP(); !reflect.DeepEqual(c.ResourceDomains, []string{"https://app.example.com:8443"}) {
		t.Errorf("absolute cdn.Base: DefaultCSP = %+v", c)
	}
}

func TestAttributes(t *testing.T) {
	a := Tool("complete_order", map[string]any{"id": "A-1"}, Target("#order-A-1", true))
	want := templ.Attributes{"data-mcp-tool": "complete_order", "data-mcp-args": `{"id":"A-1"}`,
		"data-mcp-target": "#order-A-1", "data-mcp-swap": "outer"}
	if !reflect.DeepEqual(a, want) {
		t.Errorf("Tool = %v", a)
	}
	if a := Tool("find", nil); len(a) != 1 {
		t.Errorf("Tool without args = %v", a)
	}
	if a := Target("#x", false); len(a) != 1 {
		t.Errorf("inner Target = %v", a)
	}
	if a := Prompt("Hi"); a["data-mcp-prompt"] != "Hi" {
		t.Errorf("Prompt = %v", a)
	}
}

func page(t *testing.T, ctx context.Context, p Props) string {
	t.Helper()
	html, err := Page(ctx, p, templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<div id="slot" data-mcp-result></div>`)
		return err
	}))
	if err != nil {
		t.Fatal(err)
	}
	return html
}

func TestDocument(t *testing.T) {
	html := page(t, context.Background(), Props{Title: "Orders <1>", App: "orders", Version: "1.2.0"})
	for _, want := range []string{
		"<!doctype html>",
		"<title>Orders &lt;1&gt;</title>",
		"tabler.min.css",
		`<div class="p-3" data-mcp-root><div id="slot" data-mcp-result></div></div>`,
		`data-key="` + MetaHTML + `"`,
		`data-protocol="` + ProtocolVersion + `"`,
		`data-name="orders"`,
		`data-version="1.2.0"`,
		"ui/initialize",
	} {
		if !strings.Contains(strings.ToLower(html), strings.ToLower(want)) {
			t.Errorf("document lacks %q", want)
		}
	}
	for _, not := range []string{"data-host-styles>", "tabler.min.js"} {
		if strings.Contains(html, not) {
			t.Errorf("document has %q", not)
		}
	}
	if strings.Count(html, "<script") != 1 {
		t.Errorf("want only the bridge's script, got %d", strings.Count(html, "<script"))
	}

	html = page(t, templ.WithNonce(context.Background(), "n0nce"), Props{Title: "x", HostStyles: true, Scripts: true})
	if !strings.Contains(html, `<script nonce="n0nce" data-key=`) || !strings.Contains(html, "data-host-styles>") {
		t.Errorf("bridge without nonce or host styles:\n%s", html)
	}
	if !strings.Contains(html, "tabler.min.js") {
		t.Error("Scripts does not load Tabler's script")
	}
}

// With the kit's sheets served as a file, the sandbox could not load them: the document inlines them.
func TestDocumentInlinesExternalSheets(t *testing.T) {
	csp.SetExternal(true)
	defer csp.SetExternal(false)
	html := page(t, context.Background(), Props{Title: "x"})
	if !strings.Contains(html, "/* icon */") {
		t.Error("kit CSS not inlined with csp.SetExternal")
	}
}

func TestBridgeIsTheEmbeddedScript(t *testing.T) {
	for _, want := range []string{"ui/notifications/initialized", "ui/notifications/tool-result", "ui/notifications/size-changed",
		"ui/notifications/host-context-changed", "ui/resource-teardown", "tools/call", "ui/message", "ui/open-link"} {
		if !strings.Contains(bridgeJS, want) {
			t.Errorf("bridge lacks %s", want)
		}
	}
}
