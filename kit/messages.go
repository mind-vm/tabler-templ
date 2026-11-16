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
	// Page shell: the link that skips to the content, first on the page.
	SkipToContent string
	// Detail block: the edit button, the dialog's cancel button and its hint for a section dialog.
	Edit        string
	Cancel      string
	SectionHint string
	// Row edit block: the hint above the rows and the buttons of a row.
	RowsHint  string
	SaveRow   string
	DeleteRow string
	AddRow    string
	// Package deadline: the time left or overdue, with the units (fmt with a number).
	DeadlineLeft    string // "%s left"
	DeadlineOverdue string // "overdue by %s"
	Minutes         string // "%d min"
	Hours           string // "%d h"
	Days            string // "%d days"
	DaysSince       string // days after DeadlineOverdue where the language needs another form; empty: Days
	DateTime        string // Go layout of a date with time, e.g. "02.01.2006, 15:04 Uhr" (GS-73)
	// Package mark: text from outside and internal notes.
	MarkExternal string
	MarkInternal string
	// Create page and history.
	Create    string
	NoEntries string
	// Sign-in page.
	SignIn   string
	Username string
	Password string
	Remember string
	// Password field: the button that shows or hides the password.
	ShowPassword string
	HidePassword string
	// Password strength meter.
	PasswordStrength string
	// Read out after the label of a required field, e.g. " (required)".
	Required string
	// Split view and action bar: the back link on small screens, the inspector toggle, unsaved changes.
	Back           string
	ToggleDetails  string
	UnsavedChanges string
	// htmx: the message for a failed request whose details stay on the server, and for a lost connection.
	ErrorGeneric string
	Offline      string
}

var (
	English = Messages{Close: "Close", Pages: "Pages", PrevPage: "Previous page", NextPage: "Next page",
		Previous: "Previous", Next: "Next", First: "First", Last: "Last", FirstPage: "First page", LastPage: "Last page", Page: "Page",
		ToggleSidebar: "Collapse or expand the sidebar", OpenMenu: "Open menu",
		Search: "Search", SearchLabel: "Search", Reset: "Reset", Of: "of", NoResults: "No results.", Breadcrumbs: "Breadcrumbs", SkipToContent: "Skip to content",
		Edit: "Edit", Cancel: "Cancel", SectionHint: "Only this section is saved.",
		RowsHint: "Each row is saved on its own.", SaveRow: "Save row", DeleteRow: "Delete row", AddRow: "Add",
		MarkExternal: "from outside", MarkInternal: "internal, never sent",
		DeadlineLeft: "%s left", DeadlineOverdue: "overdue by %s", Minutes: "%d min", Hours: "%d h", Days: "%d days", DateTime: "2006-01-02 15:04",
		Create: "Create", NoEntries: "No entries.",
		SignIn: "Sign in", Username: "User name", Password: "Password", Remember: "Stay signed in",
		ShowPassword: "Show password", HidePassword: "Hide password", PasswordStrength: "Password strength",
		Required: " (required)",
		Back:     "Back", ToggleDetails: "Show or hide details", UnsavedChanges: "Unsaved changes",
		ErrorGeneric: "Something went wrong. Please try again.", Offline: "No connection to the server."}
	German = Messages{Close: "Schließen", Pages: "Seiten", PrevPage: "Vorige Seite", NextPage: "Nächste Seite",
		Previous: "Zurück", Next: "Weiter", First: "Erste", Last: "Letzte", FirstPage: "Erste Seite", LastPage: "Letzte Seite", Page: "Seite",
		ToggleSidebar: "Seitenleiste ein- oder ausklappen", OpenMenu: "Menü öffnen",
		Search: "Suchen", SearchLabel: "Suchen", Reset: "Zurücksetzen", Of: "von", NoResults: "Keine Treffer.", Breadcrumbs: "Brotkrumen", SkipToContent: "Zum Inhalt",
		Edit: "Bearbeiten", Cancel: "Abbrechen", SectionHint: "Gespeichert wird nur dieser Abschnitt.",
		RowsHint: "Jede Zeile wird für sich gespeichert.", SaveRow: "Zeile speichern", DeleteRow: "Zeile löschen", AddRow: "Hinzufügen",
		MarkExternal: "von außen", MarkInternal: "intern, nie gesendet",
		DeadlineLeft: "noch %s", DeadlineOverdue: "überfällig seit %s", Minutes: "%d Min.", Hours: "%d Std.", Days: "%d Tage", DaysSince: "%d Tagen", DateTime: "02.01.2006, 15:04 Uhr",
		Create: "Anlegen", NoEntries: "Keine Einträge.",
		SignIn: "Anmelden", Username: "Benutzername", Password: "Passwort", Remember: "Angemeldet bleiben",
		ShowPassword: "Passwort zeigen", HidePassword: "Passwort verbergen", PasswordStrength: "Passwortstärke",
		Required: " (Pflichtfeld)",
		Back:     "Zurück", ToggleDetails: "Details ein- oder ausblenden", UnsavedChanges: "Ungespeicherte Änderungen",
		ErrorGeneric: "Etwas ist schiefgegangen. Bitte versuche es noch einmal.", Offline: "Keine Verbindung zum Server."}
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
