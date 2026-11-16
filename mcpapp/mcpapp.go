// Package mcpapp renders tabler-templ pages as MCP Apps (https://modelcontextprotocol.io/docs/extensions/apps,
// formerly MCP-UI): a tool of an MCP server returns a UI resource, and the host (Claude, ChatGPT, VS Code …)
// shows it in a sandboxed iframe beside the conversation. The page is rendered on the server, like every
// tabler-templ page; what is new is the way back: in the sandbox there is no request to the app's own
// server (no htmx), only messages to the host, which calls the server's tools.
//
// The package has no dependency on an MCP SDK. It gives the strings and the _meta maps the SDK needs, the
// document, and a small script (the bridge) in it that speaks the MCP Apps protocol (JSON-RPC over
// postMessage, protocol version ProtocolVersion):
//
//   - Document is the UI resource: a whole HTML document with Tabler, the kit's CSS and the bridge.
//     Serve it from resources/read with MIMEType, and point the tool at it with ToolMeta.
//   - Elements carry what they do as attributes: Tool on a button or a form calls a tool of the server
//     (tools/call) and swaps the HTML of its result into Target; Prompt sends a message to the
//     conversation (ui/message); links open through the host (ui/open-link).
//   - A tool answers with HTML in its result's _meta (ResultMeta, key MetaHTML). The bridge puts the HTML
//     of the tool that opened the view into the element marked Result, and that of a Tool call into its
//     Target. The model reads the result's text content, not the HTML.
//   - The host's theme (light or dark) is applied; with Props.HostStyles its colors, font and radius too.
//
// A resource is often loaded once and cached by the host, so it should be a shell (title, a spinner in the
// Result element) and the data should come with the tool's result:
//
//	// resources/read "ui://orders/list"
//	html, _ := mcpapp.Page(ctx, mcpapp.Props{Title: "Orders"}, shell())
//	// → contents: [{uri, mimeType: mcpapp.MIMEType, text: html, _meta: mcpapp.ResourceMeta{}.Meta()}]
//
//	// the tool "list_orders", declared with _meta: mcpapp.ToolMeta("ui://orders/list")
//	table, _ := mcpapp.Render(ctx, orderTable(orders))
//	// → CallToolResult{content: [{type: "text", text: summary}], _meta: mcpapp.ResultMeta(table)}
//
// Scripts inside swapped HTML do not run; Tabler's own components (dropdowns, modals) work with
// Props.Scripts. Everything the page loads from outside must be declared in the resource's CSP
// (ResourceMeta); the default declares the origin of cdn.Base.
package mcpapp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/cdn"
)

// The protocol, as of the stable MCP Apps specification.
const (
	// MIMEType is the UI resource's mimeType.
	MIMEType = "text/html;profile=mcp-app"
	// Extension is the id of the MCP Apps extension in the capabilities.
	Extension = "io.modelcontextprotocol/ui"
	// ProtocolVersion is the version the bridge asks for in ui/initialize.
	ProtocolVersion = "2026-01-26"
	// MetaHTML is the key of the HTML in a tool result's _meta (ResultMeta). The bridge also reads it from
	// structuredContent, for hosts that do not pass _meta on.
	MetaHTML = "io.github.mind-vm/tabler-templ/html"
)

// Visibility of a tool (ToolMeta): who may call it.
const (
	VisibilityModel = "model" // the model, as any tool
	VisibilityApp   = "app"   // the view (Tool attributes); a tool only for the view is hidden from the model
)

// ToolMeta is the _meta of a tool that shows the UI resource uri ("ui://…"). Without visibility, model and
// view may call it; ToolMeta(uri, VisibilityApp) makes a tool only the view calls (a button's action).
func ToolMeta(uri string, visibility ...string) map[string]any {
	ui := map[string]any{"resourceUri": uri}
	if len(visibility) > 0 {
		ui["visibility"] = visibility
	}
	return map[string]any{"ui": ui}
}

// AppToolMeta is the _meta of a tool only the view calls and that shows no resource of its own: a
// button's action whose result the bridge swaps into the button's Target.
func AppToolMeta() map[string]any {
	return map[string]any{"ui": map[string]any{"visibility": []string{VisibilityApp}}}
}

// ResultMeta is the _meta of a tool result that carries html for the view.
func ResultMeta(html string) map[string]any {
	return map[string]any{MetaHTML: html}
}

// CSP declares the origins the view may reach (the host builds the iframe's Content-Security-Policy from
// it). Without a declaration only inline scripts and styles and data: images are allowed.
type CSP struct {
	ConnectDomains  []string // fetch, XHR, WebSocket
	ResourceDomains []string // scripts, stylesheets, fonts, images, media
	FrameDomains    []string // nested iframes
	BaseURIDomains  []string // <base href>
}

func (c CSP) empty() bool {
	return len(c.ConnectDomains)+len(c.ResourceDomains)+len(c.FrameDomains)+len(c.BaseURIDomains) == 0
}

// DefaultCSP allows the origin of cdn.Base (jsDelivr by default), where Document loads Tabler from. With
// a relative cdn.Base (the app serves Tabler itself) it is empty: the sandbox cannot reach the app by a
// relative path, so set cdn.Base to the app's absolute URL or inline the CSS (Props.Head).
func DefaultCSP() CSP {
	u, err := url.Parse(cdn.Base)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return CSP{}
	}
	return CSP{ResourceDomains: []string{u.Scheme + "://" + u.Host}}
}

// ResourceMeta is the _meta of a UI resource (in resources/list and in the content of resources/read).
type ResourceMeta struct {
	// CSP is what the view may load; the zero value means DefaultCSP. NoCSP declares nothing.
	CSP   CSP
	NoCSP bool
	// PrefersBorder asks the host for a border and a background around the view (true) or none (false);
	// nil leaves it to the host.
	PrefersBorder *bool
	// Domain is a dedicated origin for the sandbox; its format depends on the host.
	Domain string
}

// Meta is the map for the resource's _meta.
func (m ResourceMeta) Meta() map[string]any {
	ui := map[string]any{}
	c := m.CSP
	if c.empty() && !m.NoCSP {
		c = DefaultCSP()
	}
	if !c.empty() {
		csp := map[string]any{}
		add := func(k string, v []string) {
			if len(v) > 0 {
				csp[k] = v
			}
		}
		add("connectDomains", c.ConnectDomains)
		add("resourceDomains", c.ResourceDomains)
		add("frameDomains", c.FrameDomains)
		add("baseUriDomains", c.BaseURIDomains)
		ui["csp"] = csp
	}
	if m.PrefersBorder != nil {
		ui["prefersBorder"] = *m.PrefersBorder
	}
	if m.Domain != "" {
		ui["domain"] = m.Domain
	}
	return map[string]any{"ui": ui}
}

// Render renders c to a string: a Document for resources/read, or a fragment for ResultMeta.
func Render(ctx context.Context, c templ.Component) (string, error) {
	var b bytes.Buffer
	if err := c.Render(ctx, &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// Page renders Document(p) with body as its page, for resources/read.
func Page(ctx context.Context, p Props, body templ.Component) (string, error) {
	return Render(templ.WithChildren(ctx, body), Document(p))
}

// Tool makes a button (or a link) call the tool name with args when clicked, or a form call it when
// submitted, with args and the form's fields (later ones win; numbers and range inputs are sent as
// numbers, a checkbox without value as true). While the call runs the button shows Tabler's spinner.
// Combine with Target to show the HTML of the result:
//
//	<button class="btn" { mcpapp.Tool("complete_order", map[string]any{"id": o.ID}, mcpapp.Target("#order-"+o.ID, true))... }>
func Tool(name string, args map[string]any, more ...templ.Attributes) templ.Attributes {
	a := templ.Attributes{"data-mcp-tool": name}
	if len(args) > 0 {
		b, _ := json.Marshal(args)
		a["data-mcp-args"] = string(b)
	}
	return merge(a, more)
}

// Target is where the HTML of a Tool's result goes: into the element matched by selector, or in its place
// (outer). Without a Target the result is only announced (event mcp:tool-result on the element).
func Target(selector string, outer bool) templ.Attributes {
	a := templ.Attributes{"data-mcp-target": selector}
	if outer {
		a["data-mcp-swap"] = "outer"
	}
	return a
}

// Prompt makes a button send text to the conversation as the user's message (ui/message), e.g.
// "Summarize the overdue orders".
func Prompt(text string) templ.Attributes {
	return templ.Attributes{"data-mcp-prompt": text}
}

// Result marks the element that receives the HTML of the result of the tool that opened the view.
func Result() templ.Attributes {
	return templ.Attributes{"data-mcp-result": true}
}

func merge(a templ.Attributes, more []templ.Attributes) templ.Attributes {
	for _, m := range more {
		for k, v := range m {
			a[k] = v
		}
	}
	return a
}
