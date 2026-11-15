// Package sidepanel is a work area on the right of the page with tabs, e.g. notes and an assistant.
// Buttons in the page's header open it on their tab; a second click closes it. It stays open across
// pages: the state is a cookie, so the server renders the next page with the panel open already (like
// the sidebar). Below xl it is a sheet over the page and does not open by itself after a page change.
// It is not a place for edit forms; those belong in the dialogs of the page.
//
//	open, tab := sidepanel.Read(r, "notes", "assistant")
//	p := &sidepanel.Panel{Open: open, Active: tab, Tabs: []sidepanel.Tab{
//		{Name: "notes", Icon: "ti-notes", Title: "Notes", Content: notes()},
//		{Name: "assistant", Icon: "ti-sparkles", Title: "Assistant", Content: chat()},
//	}}
//	@shell.Page(shell.Props{…, HeaderActions: sidepanel.Buttons(p), BodyAttributes: sidepanel.BodyAttributes(p, nil)}) { … }
//	… at the end of <body>, after Tabler's script:
//	@sidepanel.Aside(p)
//	@sidepanel.Script()
package sidepanel

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// Cookie holds the state: the name of the open tab, or "closed" (or nothing) when the panel is closed.
const Cookie = "kit-sidepanel"

// ID is the panel's id; the header buttons point to it with aria-controls.
const ID = "kit-sidepanel"

// Tab is one tab of the panel.
type Tab struct {
	Name    string // id, e.g. "notes": in the cookie, in data attributes and in ids
	Icon    string // e.g. "ti-notes"
	Title   string // of the tab and its header button, e.g. "Notes"
	Content templ.Component
	// Head holds buttons in the panel's head that belong to this tab only (e.g. "New chat").
	Head templ.Component
	// Button replaces the default header button (e.g. one with a count); it carries ButtonAttributes.
	Button templ.Component
	// Forms: the tab has forms that are sent from the panel. On small screens the panel then stays open
	// after the page reloads, and opening it from the header focuses its first text area.
	Forms bool
}

// Panel is a page's panel. Without tabs there is none.
type Panel struct {
	Tabs   []Tab
	Open   bool
	Active string // the visible tab (when closed: the last one chosen); default the first
	Label  string // aria-label; default the tabs' titles
}

// Has reports whether the panel has at least one tab.
func (p *Panel) Has() bool { return p != nil && len(p.Tabs) > 0 }

func (p *Panel) hasTab(name string) bool {
	return p.Has() && slices.ContainsFunc(p.Tabs, func(t Tab) bool { return t.Name == name })
}

func (p *Panel) active() string {
	if p.hasTab(p.Active) {
		return p.Active
	}
	if p.Has() {
		return p.Tabs[0].Name
	}
	return ""
}

func (p *Panel) shows(name string) bool { return p.active() == name }

func (p *Panel) label() string {
	if p.Label != "" {
		return p.Label
	}
	var ts []string
	for _, t := range p.Tabs {
		ts = append(ts, t.Title)
	}
	return strings.Join(ts, ", ")
}

// Read reads the state from the cookie: open on one of the tabs names, else closed (tab "").
func Read(r *http.Request, names ...string) (open bool, tab string) {
	c, err := r.Cookie(Cookie)
	if err != nil || !slices.Contains(names, c.Value) {
		return false, ""
	}
	return true, c.Value
}

// Remember sets the cookie the way the script does: tab "" closes the panel.
func Remember(w http.ResponseWriter, tab string) {
	if tab == "" {
		tab = "closed"
	}
	http.SetCookie(w, &http.Cookie{Name: Cookie, Value: tab, Path: "/", MaxAge: 31536000, SameSite: http.SameSiteLaxMode})
}

// ButtonAttributes make a header button open the panel on the tab name (a second click closes it).
func ButtonAttributes(p *Panel, name string) templ.Attributes {
	return templ.Attributes{"data-kit-sidepanel-tab": name, "aria-controls": ID,
		"aria-expanded": strconv.FormatBool(p != nil && p.Open && p.active() == name)}
}

// BodyAttributes adds the class that makes room for an open panel, so the page does not jump when it
// comes open from the server.
func BodyAttributes(p *Panel, attrs templ.Attributes) templ.Attributes {
	if !p.Has() || !p.Open {
		return attrs
	}
	out := templ.Attributes{}
	for k, v := range attrs {
		out[k] = v
	}
	if c, ok := out["class"].(string); ok && c != "" {
		out["class"] = c + " kit-sidepanel-open"
	} else {
		out["class"] = "kit-sidepanel-open"
	}
	return out
}
