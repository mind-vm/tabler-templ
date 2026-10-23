package kit

import (
	"context"
	"reflect"
)

// Messages are the few texts the kit renders itself (mostly aria-labels). An app sets them once per
// request with WithMessages; a component's own props still override them per call. The kit has no
// translation files: an app with several languages picks the Messages for the request's language.
type Messages struct {
	Close    string
	Pages    string
	PrevPage string
	NextPage string
	// Pagination with text labels, first/last links and page numbers ("Page 3").
	Previous  string
	Next      string
	First     string
	Last      string
	FirstPage string
	LastPage  string
	Page      string
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
	// Detail block: the edit button, the dialog's cancel button and its hint for a section dialog.
	Edit        string
	Cancel      string
	SectionHint string
	// Create page and history.
	Create    string
	NoEntries string
	// Sign-in page.
	SignIn   string
	Username string
	Password string
	Remember string
}

var (
	English = Messages{Close: "Close", Pages: "Pages", PrevPage: "Previous page", NextPage: "Next page",
		Previous: "Previous", Next: "Next", First: "First", Last: "Last", FirstPage: "First page", LastPage: "Last page", Page: "Page",
		ToggleSidebar: "Collapse or expand the sidebar", OpenMenu: "Open menu",
		Search: "Search", SearchLabel: "Search", Reset: "Reset", Of: "of", NoResults: "No results.", Breadcrumbs: "Breadcrumbs",
		Edit: "Edit", Cancel: "Cancel", SectionHint: "Only this section is saved.",
		Create: "Create", NoEntries: "No entries.",
		SignIn: "Sign in", Username: "User name", Password: "Password", Remember: "Stay signed in"}
	German = Messages{Close: "Schließen", Pages: "Seiten", PrevPage: "Vorige Seite", NextPage: "Nächste Seite",
		Previous: "Zurück", Next: "Weiter", First: "Erste", Last: "Letzte", FirstPage: "Erste Seite", LastPage: "Letzte Seite", Page: "Seite",
		ToggleSidebar: "Seitenleiste ein- oder ausklappen", OpenMenu: "Menü öffnen",
		Search: "Suchen", SearchLabel: "Suchen", Reset: "Zurücksetzen", Of: "von", NoResults: "Keine Treffer.", Breadcrumbs: "Brotkrumen",
		Edit: "Bearbeiten", Cancel: "Abbrechen", SectionHint: "Gespeichert wird nur dieser Abschnitt.",
		Create: "Anlegen", NoEntries: "Keine Einträge.",
		SignIn: "Anmelden", Username: "Benutzername", Password: "Passwort", Remember: "Angemeldet bleiben"}
)

type messagesKey struct{}

// WithMessages returns a context whose components render m.
func WithMessages(ctx context.Context, m Messages) context.Context {
	return context.WithValue(ctx, messagesKey{}, m)
}

// T returns the context's Messages, English by default. Empty fields fall back to English, so an
// app's own Messages keep working when a new version of the kit adds texts.
func T(ctx context.Context) Messages {
	m, ok := ctx.Value(messagesKey{}).(Messages)
	if !ok {
		return English
	}
	v, en := reflect.ValueOf(&m).Elem(), reflect.ValueOf(English)
	for i := range v.NumField() {
		if f := v.Field(i); f.Kind() == reflect.String && f.String() == "" {
			f.SetString(en.Field(i).String())
		}
	}
	return m
}

// Or returns s, or fallback when s is empty.
func Or(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}
