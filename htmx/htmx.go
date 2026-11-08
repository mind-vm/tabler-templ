// Package htmx is the kit's support for htmx (https://htmx.org): it loads htmx (Library), sets up the
// page's client (Setup: toasts, a dialog, 422 handling, Tabler after swaps), and gives components a
// contract and their handlers a way to answer:
//
//   - Route defines an endpoint once; Pattern is its mux pattern, URL builds its links.
//   - ID, Target and Vals name a component instance's elements.
//   - Respond writes a fragment with its headers first (Vary, HX-Trigger, status …), and the whole page
//     for requests without htmx (FullPage); Error and Redirect answer failures and redirects.
//   - CheckRoutes is the contract test.
//
// Library goes in the head. It turns off htmx's inline indicator styles, which a Content-Security-Policy
// with nonces would block, and puts the same rules into a <style> with the nonce instead.
//
//	<head> … @htmx.Library() </head>
//	<body> … @assets.Scripts() @htmx.Setup() </body>
package htmx

import (
	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/cdn"
	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/internal/plugin"
	"github.com/mind-vm/tabler-templ/versions"
)

// Version is the htmx version Library loads.
const Version = versions.HTMX

// LibraryURL is htmx on cdn.Base.
func LibraryURL() string { return cdn.URL("htmx.org", Version, "dist/htmx.min.js") }

const config = `<meta name="htmx-config" content='{"includeIndicatorStyles":false}'/>`

var indicator = csp.Sheet("htmx", `.htmx-indicator{opacity:0}.htmx-request .htmx-indicator,.htmx-request.htmx-indicator{opacity:1;transition:opacity .2s ease-in}`)

// Library loads htmx with the kit's configuration; put it in the head.
func Library() templ.Component {
	return templ.Join(templ.Raw(config), indicator, plugin.Script(LibraryURL()))
}

func init() {
	cdn.Register("htmx", "htmx.org", Version, "dist/htmx.min.js")
}
