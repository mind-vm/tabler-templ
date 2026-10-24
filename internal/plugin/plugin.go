// Package plugin starts the JavaScript plugins of the plugin packages (datepicker, selectbox, dropzone …):
// each package hands Start a function that sets up the plugin's elements below a root, and Start renders it
// once per page, runs it when the page is ready and again on content swapped in by htmx.
package plugin

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
)

var (
	mu      sync.Mutex
	handles = map[string]*templ.OnceHandle{}
)

func handle(name string) *templ.OnceHandle {
	mu.Lock()
	defer mu.Unlock()
	h, ok := handles[name]
	if !ok {
		h = templ.NewOnceHandle()
		handles[name] = h
	}
	return h
}

// Start renders, once per page and name, a script that calls setup (a JavaScript function expression
// taking a root element) with the document when the page is ready, and with the new content after each
// htmx swap. setup must skip elements it has set up already.
func Start(name, setup string) templ.Component {
	script := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, csp.Open(ctx, "script")+"(function () {\n\tvar setup = "+setup+";\n"+
			"\tdocument.readyState !== 'loading' ? setup(document) : document.addEventListener('DOMContentLoaded', function () { setup(document); }, { once: true });\n"+
			"\tdocument.addEventListener('htmx:afterSettle', function (e) { setup(e.target); });\n})();</script>")
		return err
	})
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return handle(name).Once().Render(templ.WithChildren(ctx, script), w)
	})
}

// Tabler starts the elements matching selector with one of Tabler's own components (window.tabler[class]),
// for content swapped in after Tabler's script ran; Tabler starts those on the page itself.
func Tabler(class, selector string) templ.Component {
	return Start("tabler-"+class, `function (root) {
		if (!window.tabler || !window.tabler.`+class+`) return;
		root.querySelectorAll(`+JSON(selector)+`).forEach(function (el) { window.tabler.`+class+`.getOrCreateInstance(el); });
	}`)
}

// JSON is v as JSON, for attributes and scripts.
func JSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Script loads a script from src, without defer: Tabler's script, which comes later on the page, starts its
// components when it runs and needs the plugins loaded by then.
func Script(src string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<script`+csp.Attr(ctx)+` src="`+templ.EscapeString(src)+`"></script>`)
		return err
	})
}

// Stylesheet links a stylesheet.
func Stylesheet(href string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<link rel="stylesheet" href="`+templ.EscapeString(href)+`"/>`)
		return err
	})
}
