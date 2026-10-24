// Package csp puts templ's nonce (templ.WithNonce) on the kit's own <script> and <style> elements, so
// pages can run under a Content-Security-Policy with 'nonce-…' sources. Without a nonce in the context the
// elements render as before.
package csp

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// Nonce is the nonce attribute for a templ element: <script { csp.Nonce(ctx)... }>. It is empty when the
// context carries no nonce.
func Nonce(ctx context.Context) templ.Attributes {
	if n := templ.GetNonce(ctx); n != "" {
		return templ.Attributes{"nonce": n}
	}
	return templ.Attributes{}
}

// Open is the opening tag of a script or style element (tag "script" or "style") with the nonce, for code
// that writes HTML directly.
func Open(ctx context.Context, tag string) string {
	return "<" + tag + Attr(ctx) + ">"
}

// Attr is the nonce attribute with a leading space (` nonce="…"`), or "" without a nonce.
func Attr(ctx context.Context) string {
	if n := templ.GetNonce(ctx); n != "" {
		return ` nonce="` + templ.EscapeString(n) + `"`
	}
	return ""
}

// Style renders css in a <style> element with the nonce. css is trusted (the kit's own CSS).
func Style(css string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, Open(ctx, "style")+css+"</style>")
		return err
	})
}
