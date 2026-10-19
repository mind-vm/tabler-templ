# tablerkit

Server-rendered UI components for [Tabler](https://tabler.io) in Go, written with [templ](https://templ.guide).
The API follows [shadcn-templ](https://github.com/axadrn/shadcn-templ): one package per component, a `Props`
struct, `Class` and `Attributes` on every component, children via `{ children... }`. Unlike shadcn-templ it
needs no Tailwind, no Node and (except for the sidebar and the detail dialogs) no JavaScript of its own: components emit Tabler's classes and
`data-bs-*` attributes, and Tabler's bundled script does the rest (modals, dropdowns, dismissible alerts).

> Status: early (v0.x). The API may still change between minor versions until v1.

```sh
go get github.com/mind-vm/tablerkit@latest
```

## Use

```go
import (
	"github.com/mind-vm/tablerkit/button"
	"github.com/mind-vm/tablerkit/card"
	"github.com/mind-vm/tablerkit/modal"
)

templ Customer(c Customer) {
	@card.Card(card.Props{Section: true}) {
		@card.Header() {
			@card.Title() { Customer }
			@card.Actions() {
				@button.Button(button.Props{Size: button.SizeSmall, Icon: "ti-pencil", Type: "button",
					Attributes: modal.Trigger("edit-customer")}) { Edit }
			}
		}
		@card.Body() { { c.Name } }
	}
	@modal.Modal(modal.Props{ID: "edit-customer", Title: "Customer", Centered: true, Fade: true}) {
		…form…
	}
}
```

The page must load Tabler's CSS and JS (`@tabler/core`) and, for icons, the Tabler icons webfont. The kit
does not ship or link them: serve them yourself (embedded in your binary or from a CDN, as you prefer).

## Components

| Package | Components |
|---|---|
| `button` | `Button` (`<button>` or `<a>` with `Href`), variants, outline, sizes, leading icon, icon-only; `Classes(p)` |
| `card` | `Card`, `Header`, `Title`, `Actions`, `Body`, `Footer` |
| `modal` | `Modal` (sizes, centered, scrollable, custom header, footer, `<form>` as content); `Trigger(id)` attributes for any button |
| `form` | `Field` (label, hint, error), `Input`, `Select`, `Textarea`, `Check` (checkbox, radio, switch); `AriaLabel` for controls without a visible label |
| `table` | `Table` (optionally responsive, wrapper classes), `Empty` row |
| `pagination` | `Pagination` with first/last/neighbour pages; `Numbers(page, pages)` |
| `alert` | `Alert` (variant, icon, dismissible) |
| `badge` | `Badge` (Tabler color, light) |
| `sidebar` | `Sidebar` (fixed left, collapsible to icons, a sheet on small screens), `Header`, `Brand`, `Content`, `Group`, `Item` (icon, active, count badge), `Separator`, `Footer`, `User` (dropdown), `Trigger` (menu and collapse buttons for the page header), `Styles`, `Script`, `Collapsed(r)` |
| `icon` | `Icon("ti-…")` – full class name, so icon-font subsetters can find it in the source |
| `kit` | `Classes`, `If`, `Attrs`, `Default`, `Messages` |

Not yet: dropdown, tabs, avatar, toast, datagrid, steps.

## Blocks

Components are small and know nothing about the app. Blocks are bigger pieces built from them, the
way [shadcn/ui](https://ui.shadcn.com/blocks) separates the two: a whole page frame or a whole list
page, so a new app starts with the parts every app needs. Like components they are packages with
props; what differs per app goes into slots (`templ.Component` props) or the children.

| Package | Block |
|---|---|
| `blocks/shell` | `Page`: the HTML document with sidebar, sticky header (sidebar buttons, breadcrumbs, actions), page header (heading, extra, subtitle, actions), flash message and content; `Header`, `Breadcrumbs`, `PageHeader`, `Styles` for own frames |
| `blocks/listpage` | `New(r, Config, filters…)` reads search, filters and page from the URL; `Paginate` cuts out a page (25 rows); `Card` = `Search` (search field, filter selects, Reset) + children + `Footer` ("1–25 of 80" and page numbers) |
| `blocks/detail` | `Section` (card with title, extra and at most one action), `Actions`, `EditButton`, `Fields`/`Field`; `Dialog` (a section's edit dialog: context line, error, `Section` hint, `Large` for lists) and `Script(ScriptConfig)`: hidden fields in its forms, footer with Cancel, focus, reopening with error and values after a failed save, reopening a large dialog after saving a row, values from the opening button; `window.kitDialog.open/matchForm` |
| `blocks/create` | `Form` (post form of a create page) and `Footer` (hint, Cancel back to the list, primary Create button): creating is a page of its own, not a dialog |
| `blocks/history` | `Timeline([]Entry)`: when, who, what and the changed fields, newest first; the app phrases the text |
| `blocks/signin` | `Page` (centered sign-in page without sidebar: brand, heading, intro, flash, footer), `Password` (user name and password with the browser's autocomplete names, error, "stay signed in", hidden fields) and `Picker` (accounts to choose, e.g. for demos); credentials and sessions stay in the app |

```go
func customers(w http.ResponseWriter, r *http.Request) {
	l := listpage.New(r, listpage.Config{Placeholder: "Name, city …"})
	rows := listpage.Paginate(l, allCustomers(), func(c Customer) bool { return listpage.Contains(l.Query, c.Name, c.City) })
	customersPage(r, l, rows).Render(r.Context(), w)
}

templ customersPage(r *http.Request, l *listpage.List, rows []Customer) {
	@shell.Page(shell.Props{Title: "Customers", AppName: "Acme", Head: head(), Scripts: scripts(),
		Sidebar: nav(r), Collapsed: sidebar.Collapsed(r), Breadcrumbs: []shell.Crumb{{Text: "Customers"}}}) {
		@listpage.Card(l) {
			@table.Table() { … rows … }
		}
	}
}
```

`shell.Page` renders the kit's CSS before `Head`, so the app's theme overrides it, and the sidebar
script after `Scripts`. Texts come from `kit.Messages` (Search, Reset, "of", Breadcrumbs, Edit, Cancel, Create, NoEntries, SignIn, Username, Password …).

The detail script is the kit's second script (after the sidebar's). The server and the script exchange
the dialog's state after a failed save: the page sets a global variable (default `kitDialogState`) to
`{"id": …, "error": …, "values": {field: [values]}}`, and the forms carry hidden fields with the dialog's
ID and the page to return to. `ScriptConfig` renames all of them for apps that already have their own.

### Sidebar

Tabler has no sidebar that collapses to icons, so `sidebar` is the one package with CSS and script of its
own: render `@sidebar.Styles()` in `<head>` and `@sidebar.Script()` after Tabler's script. On large screens
the sidebar is fixed on the left and the `Trigger` collapses it to icons (items show their label as a
tooltip then); the script remembers the choice in the cookie `kit_sidebar`, and `sidebar.Collapsed(r)`
reads it so the server renders the same state. On small screens the sidebar is Tabler's offcanvas
(`offcanvas-lg`), opened by the menu button of the `Trigger`. Theme it with the CSS variables listed in
`sidebar/sidebar.go` (`--kit-sidebar-width`, `--kit-sidebar-bg`, …). The page's `.page-wrapper` must be the
sidebar's next sibling:

```go
templ Layout(r *http.Request) {
	<div class="page">
		@sidebar.Sidebar(sidebar.Props{ID: "sidebar", Collapsed: sidebar.Collapsed(r)}) {
			@sidebar.Header() {
				@sidebar.Brand(sidebar.BrandProps{Href: "/", Title: "Acme"})
			}
			@sidebar.Content() {
				@sidebar.Group(sidebar.GroupProps{Label: "Work"}) {
					@sidebar.Item(sidebar.ItemProps{Href: "/customers", Icon: "ti-building", Label: "Customers", Active: true})
				}
			}
			@sidebar.Footer() {
				@sidebar.User(sidebar.UserProps{Initials: "AA", Name: "Ann Example"}) {
					<a class="dropdown-item" href="/logout">Log out</a>
				}
			}
		}
		<div class="page-wrapper">
			<header>@sidebar.Trigger(sidebar.TriggerProps{Target: "sidebar", Collapsed: sidebar.Collapsed(r)})</header>
			…
		</div>
	</div>
}
```

## Conventions

- Every component takes `Class` (appended) and `Attributes` (spread last; an attribute you pass wins over a
  default, e.g. `title`, `id`, `type`).
- Components with only optional props take `...Props`, so `@card.Body()` works; components with a required
  prop (`modal.Modal`, `pagination.Pagination`, form controls) take `Props`.
- templ keeps one space where the source has a line break between text and an element or `if` (HTML
  whitespace normalisation; there is no `{{- -}}`-style trim, see a-h/templ#1205). Where that matters
  (punctuation), put the punctuation right after the element on the same line, build the sentence in
  Go, or write an intended space as `{ " " }`.
- Generated `*_templ.go` files are committed, so the module works with a plain `go get`.

## Theming

The kit never hardcodes colors, radii or fonts. All appearance comes from CSS:

1. **Tabler's custom properties.** Override `--tblr-primary`, `--tblr-border-radius`, `--tblr-body-bg`,
   `--tblr-font-sans-serif` … in your own stylesheet loaded after Tabler. This restyles every component.
2. **Dark mode.** Tabler's `data-bs-theme="dark"` on `<html>` or any container.
3. **Your own classes.** `Class` on any component (e.g. `card.Card(card.Props{Class: "island"})`) for app
   specific looks, styled in your CSS.

## Texts and languages

The kit renders only a few texts itself (aria-labels such as "Close", "Previous page"). Set them once per
request; props still override them per call:

```go
ctx = kit.WithMessages(r.Context(), kit.German) // or your own kit.Messages{…}
```

There are no translation files in the kit. An app with several languages picks the `Messages` for the
request's language (and translates its own texts with whatever it uses, e.g. go-i18n).

## Develop

```sh
go tool templ generate   # after editing .templ files (templ is a Go tool dependency)
go test ./...
```

## License

MIT, see [LICENSE](LICENSE). Tabler and the Tabler icons (also MIT) are not part of this module.
