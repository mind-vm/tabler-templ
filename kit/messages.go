package kit

import "context"

// Messages are the few texts the kit renders itself (mostly aria-labels). An app sets them once per
// request with WithMessages; a component's own props still override them per call. The kit has no
// translation files: an app with several languages picks the Messages for the request's language.
type Messages struct {
	Close    string
	Pages    string
	PrevPage string
	NextPage string
	// Sidebar: the button that collapses it to icons (desktop) and the one that opens it (mobile).
	ToggleSidebar string
	OpenMenu      string
	// Blocks: list page (search, reset, "1–25 of 80", empty result) and page shell (breadcrumb label).
	Search      string
	SearchLabel string
	Reset       string
	Of          string
	NoResults   string
	Breadcrumbs string
}

var (
	English = Messages{Close: "Close", Pages: "Pages", PrevPage: "Previous page", NextPage: "Next page",
		ToggleSidebar: "Collapse or expand the sidebar", OpenMenu: "Open menu",
		Search: "Search", SearchLabel: "Search", Reset: "Reset", Of: "of", NoResults: "No results.", Breadcrumbs: "Breadcrumbs"}
	German = Messages{Close: "Schließen", Pages: "Seiten", PrevPage: "Vorige Seite", NextPage: "Nächste Seite",
		ToggleSidebar: "Seitenleiste ein- oder ausklappen", OpenMenu: "Menü öffnen",
		Search: "Suchen", SearchLabel: "Suchen", Reset: "Zurücksetzen", Of: "von", NoResults: "Keine Treffer.", Breadcrumbs: "Brotkrumen"}
)

type messagesKey struct{}

// WithMessages returns a context whose components render m.
func WithMessages(ctx context.Context, m Messages) context.Context {
	return context.WithValue(ctx, messagesKey{}, m)
}

// T returns the context's Messages, English by default.
func T(ctx context.Context) Messages {
	if m, ok := ctx.Value(messagesKey{}).(Messages); ok {
		return m
	}
	return English
}

// Or returns s, or fallback when s is empty.
func Or(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}
