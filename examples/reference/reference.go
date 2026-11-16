// Package reference is the living reference of the design concept „Gestaltung“ (rules GS-01 to GS-75):
// every pattern of the concept as a page built from the real blocks of this module, with a short German
// text on top that names the rules and the block. The pages and their sample data are German like the
// products that follow the concept; the code is English like the rest of the module.
//
// The pages (§11 of the concept):
//
//	/reference                  overview of all patterns
//	/reference/rahmen           B-01 frame: header, grouped sidebar, content, panel on the right
//	/reference/start            B-02 start page: figures that lead to filtered lists, the list „Brennt“
//	/reference/kunden           B-03 list: search, one filter, 25 rows, pages, state in the URL
//	/reference/kunden/K-1001    B-04 detail page: islands with one Edit each, dialogs, errors, sub-list
//	/reference/angebote/A-2026-001  B-05 edit a list: the full-page dialog with one form per row
//	/reference/kunden/neu       B-06 create: a page of its own, then the detail page with a message
//	/reference/zustaende        B-07 states: message, empty, empty search, error pages, read-only, loading
//	/reference/marke            B-08 brand and form: tokens, signal colors, buttons, pills, formats
//	/reference/assistent        B-09 the AI assistant in the panel (proposal, from the concept Assistent)
//
// Saving writes to a store in the process, so messages, errors in dialogs and editing the list really
// run. Nothing loads from a CDN: Tabler, the icons and the Geist fonts are embedded (static/).
//
// cmd/demo mounts it at /reference. To run it alone:
//
//	mux := http.NewServeMux()
//	reference.Register(mux)
//	http.ListenAndServe("localhost:8080", http.NewCrossOriginProtection().Handler(mux))
package reference

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/alert"
	"github.com/mind-vm/tabler-templ/blocks/assistant"
	"github.com/mind-vm/tabler-templ/blocks/listpage"
	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/blocks/sidepanel"
	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/static"
)

// Home is the overview; every page lives below it.
const Home = "/reference"

// Paths of the pages.
const (
	pathStart     = Home + "/start"
	pathFrame     = Home + "/rahmen"
	pathCustomers = Home + "/kunden"
	pathNew       = Home + "/kunden/neu"
	pathOffers    = Home + "/angebote"
	pathStates    = Home + "/zustaende"
	pathBrand     = Home + "/marke"
	pathAssistant = Home + "/assistent"
	pathRequests  = Home + "/anfragen"
	pathStatic    = Home + "/static/"
)

func customerURL(no string) string { return pathCustomers + "/" + no }
func offerURL(no string) string    { return pathOffers + "/" + no }
func requestURL(no string) string  { return pathRequests + "/" + no }

//go:embed all:static
var embedded embed.FS

// App is the reference app with its data and the assistant's conversations.
type App struct {
	store *Store
	files *static.Files

	mu       sync.Mutex
	sessions map[string]*session
}

type session struct {
	transcript *assistant.Transcript
	busy       bool
	turns      int
}

// Register mounts the app with fresh sample data.
func Register(mux *http.ServeMux) *App { return RegisterWith(mux, NewStore()) }

// RegisterWith mounts the app on store.
func RegisterWith(mux *http.ServeMux, store *Store) *App {
	sub, err := fs.Sub(embedded, "static")
	if err != nil {
		panic(err)
	}
	// The filled icons' webfont comes from the bundle too (GS-09).
	a := &App{store: store, files: static.New(pathStatic, sub, static.WithFilledIcons()), sessions: map[string]*session{}}
	mux.Handle("GET "+pathStatic, a.files)

	mux.HandleFunc("GET "+Home, a.overview)
	mux.HandleFunc("GET "+Home+"/{$}", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, Home, http.StatusMovedPermanently) })
	mux.HandleFunc("GET "+pathFrame, a.frame)
	mux.HandleFunc("GET "+pathStart, a.start)
	mux.HandleFunc("GET "+pathCustomers, a.customers)
	mux.HandleFunc("GET "+pathNew, a.newCustomer)
	mux.HandleFunc("POST "+pathCustomers, a.writes(a.createCustomer))
	mux.HandleFunc("GET "+pathCustomers+"/{no}", a.customer)
	mux.HandleFunc("POST "+pathCustomers+"/{no}/stammdaten", a.writes(a.saveMaster))
	mux.HandleFunc("POST "+pathCustomers+"/{no}/ansprechpartner", a.writes(a.saveContact))
	mux.HandleFunc("POST "+pathCustomers+"/{no}/vermerke", a.writes(a.addNote))
	mux.HandleFunc("POST "+pathCustomers+"/{no}/unterlagen", a.writes(a.addFile))
	mux.HandleFunc("POST "+pathCustomers+"/{no}/unterlagen/{id}/entfernen", a.writes(a.removeFile))
	mux.HandleFunc("GET "+pathOffers, a.offers)
	mux.HandleFunc("GET "+pathOffers+"/{no}", a.offer)
	mux.HandleFunc("POST "+pathOffers+"/{no}/angebot", a.writes(a.saveOffer))
	mux.HandleFunc("POST "+pathOffers+"/{no}/abgeben", a.writes(a.submitOffer))
	mux.HandleFunc("POST "+pathOffers+"/{no}/positionen", a.writes(a.addPosition))
	mux.HandleFunc("POST "+pathOffers+"/{no}/positionen/{id}", a.writes(a.savePosition))
	mux.HandleFunc("POST "+pathOffers+"/{no}/positionen/{id}/loeschen", a.writes(a.deletePosition))
	mux.HandleFunc("GET "+pathRequests, a.requests)
	mux.HandleFunc("GET "+pathRequests+"/{no}", a.request)
	mux.HandleFunc("POST "+pathRequests+"/{no}/senden", a.writes(a.sendReply))
	mux.HandleFunc("POST "+pathRequests+"/{no}/verwerfen", a.writes(a.discard))
	mux.HandleFunc("GET "+pathExternal, a.external)
	mux.HandleFunc("GET "+pathExternalOffers, a.externalOffers)
	mux.HandleFunc("GET "+pathExternalSignIn, a.externalSignIn)
	mux.HandleFunc("POST "+pathExternalSignIn, a.externalLogin)
	mux.HandleFunc("GET "+pathExternalSignUp, a.externalSignUp)
	mux.HandleFunc("POST "+pathExternalOut, a.externalLogout)
	mux.HandleFunc("GET "+pathStates, a.states)
	mux.HandleFunc("GET "+pathStates+"/zahlen", a.lateFigures)
	mux.HandleFunc("GET "+pathStates+"/keine-rechte", a.forbidden)
	mux.HandleFunc("GET "+pathStates+"/fehler", a.failure)
	mux.HandleFunc("POST "+pathStates+"/nur-lesen", a.toggleReadOnly)
	mux.HandleFunc("POST "+pathStates+"/demo", a.toggleDemo)
	mux.HandleFunc("POST "+pathStates+"/zuruecknehmen", a.writes(a.withdraw))
	mux.HandleFunc("POST "+pathStates+"/kunde-loeschen", a.writes(a.deleteCustomer))
	mux.HandleFunc("POST "+Home+"/zuruecksetzen", a.reset)
	mux.HandleFunc("GET "+pathBrand, a.brand)
	mux.HandleFunc("GET "+pathAssistant, a.assistantPage)
	mux.HandleFunc("GET "+pathConcept, a.concept(concepts[0]))
	mux.HandleFunc("GET "+pathConceptAssistant, a.concept(concepts[1]))
	mux.HandleFunc("POST "+pathAssistant, a.run)
	mux.HandleFunc("GET "+Home+"/{rest...}", a.notFound)
	return a
}

// --- rendering -----------------------------------------------------------------------------------

// context sets the German texts of the kit and the assistant (GS-70).
func germanContext(ctx context.Context) context.Context {
	return assistant.WithTexts(kit.WithMessages(ctx, kit.German), assistant.German)
}

func (a *App) render(w http.ResponseWriter, r *http.Request, status int, p page, body templ.Component) {
	ctx := germanContext(r.Context())
	p.readOnly = a.store.ReadOnly()
	p.demo = a.store.Demo()
	p.theme = themeOf(w, r)
	p.session = a.session(w, r)
	if p.Flash == nil {
		p.Flash = takeFlash(w, r)
	}
	open, tab := sidepanel.Read(r, tabNotes, tabAssistant)
	p.panel = a.panel(p, open, tab)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if status != 0 {
		w.WriteHeader(status)
	}
	_ = a.page(p).Render(templ.WithChildren(ctx, body), w)
}

// --- messages after an action (GS-60): a cookie read once by the next page -----------------------

const flashCookie = "reference-meldung"

type flashData struct {
	Variant  string `json:"v"`
	Text     string `json:"t"`
	Next     string `json:"n,omitempty"`
	NextHref string `json:"h,omitempty"`
}

func setFlash(w http.ResponseWriter, f flashData) {
	b, _ := json.Marshal(f)
	http.SetCookie(w, &http.Cookie{Name: flashCookie, Value: base64.RawURLEncoding.EncodeToString(b), Path: Home, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}

func takeFlash(w http.ResponseWriter, r *http.Request) *shell.Flash {
	c, err := r.Cookie(flashCookie)
	if err != nil {
		return nil
	}
	http.SetCookie(w, &http.Cookie{Name: flashCookie, Value: "", Path: Home, MaxAge: -1})
	b, err := base64.RawURLEncoding.DecodeString(c.Value)
	var f flashData
	if err != nil || json.Unmarshal(b, &f) != nil || f.Text == "" {
		return nil
	}
	if !strings.HasPrefix(f.NextHref, Home) {
		f.Next, f.NextHref = "", ""
	}
	v := alert.VariantSuccess
	if f.Variant == "danger" {
		v = alert.VariantDanger
	}
	return &shell.Flash{Variant: v, Text: f.Text, Next: f.Next, NextHref: f.NextHref}
}

// back is where a dialog returns to: the page it was opened on (the dialog script's _return) when that is
// the record's page (fallback) with its query, else fallback. After a failed save the address bar shows
// the form's action, which is no page to return to.
func back(r *http.Request, fallback string) string {
	ret := r.PostFormValue("_return")
	if u, err := url.Parse(ret); err == nil && u.Host == "" && u.Scheme == "" && u.Path == fallback {
		return ret
	}
	return fallback
}

func done(w http.ResponseWriter, r *http.Request, to string, f flashData) {
	if f.Variant == "" {
		f.Variant = "success"
	}
	setFlash(w, f)
	http.Redirect(w, r, to, http.StatusSeeOther)
}

// writes refuses every change while the company is read-only (GS-63): the buttons are gone, and the
// server checks anyway (GS-08).
func (a *App) writes(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.store.ReadOnly() {
			a.errorPage(w, r, http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

// --- the assistant's sessions --------------------------------------------------------------------

const sessionCookie = "reference-assistent"

func (a *App) session(w http.ResponseWriter, r *http.Request) *session {
	a.mu.Lock()
	defer a.mu.Unlock()
	if c, err := r.Cookie(sessionCookie); err == nil {
		if s, ok := a.sessions[c.Value]; ok {
			return s
		}
	}
	id := rand.Text()
	s := &session{transcript: &assistant.Transcript{Max: 20}}
	a.sessions[id] = s
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: id, Path: Home, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return s
}

// --- pages ---------------------------------------------------------------------------------------

func (a *App) overview(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 0, page{Key: "overview", Title: "Referenz der Gestaltung", Crumbs: []shell.Crumb{{Text: "Referenz"}},
		Explain: &overviewText}, overviewView())
}

func (a *App) frame(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 0, page{Key: "B-01", Title: "Rahmen", Crumbs: []shell.Crumb{{Text: "Referenz", Href: Home}, {Text: "Rahmen"}}}, frameView())
}

func (a *App) start(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 0, page{Key: "B-02", Title: "Start", Crumbs: []shell.Crumb{{Text: "Start"}}}, startView(a.store.Counts(), a.store.Burning(), a.customerNames()))
}

func (a *App) customerNames() map[string]string {
	m := map[string]string{}
	for _, c := range a.store.Customers() {
		m[c.No] = c.Name
	}
	return m
}

var signalFilter = listpage.Filter{Name: "posten", Label: "Alle offenen Posten", Options: func() []listpage.Option {
	var os []listpage.Option
	for _, s := range Signals {
		os = append(os, listpage.Option{Value: string(s), Label: "Offene Posten: " + s.Label()})
	}
	return os
}()}

func (a *App) customers(w http.ResponseWriter, r *http.Request) {
	l := listpage.New(r, listpage.Config{Placeholder: "Name, Ort oder Nummer"}, signalFilter)
	rows := listpage.Paginate(l, a.store.Customers(), func(c Customer) bool {
		return (l.Value("posten") == "" || string(c.Signal) == l.Value("posten")) && matches(l.Query, c.Name, c.City, c.No)
	})
	a.render(w, r, 0, page{Key: "B-03", Title: "Kunden", Crumbs: []shell.Crumb{{Text: "Kunden"}}, Actions: a.mainAction(newCustomerButton())},
		customersView(l, rows))
}

func (a *App) mainAction(c templ.Component) templ.Component {
	if a.store.ReadOnly() {
		return nil
	}
	return c
}

var statusFilter = listpage.Filter{Name: "status", Label: "Alle Stände", Options: func() []listpage.Option {
	var os []listpage.Option
	for _, s := range OfferStatuses {
		os = append(os, listpage.Option{Value: string(s), Label: s.Label()})
	}
	return os
}()}

func (a *App) offers(w http.ResponseWriter, r *http.Request) {
	names := a.customerNames()
	l := listpage.New(r, listpage.Config{Placeholder: "Nummer, Titel, Kunde oder Kundennummer"}, statusFilter)
	rows := listpage.Paginate(l, a.store.Offers(), func(o Offer) bool {
		return (l.Value("status") == "" || string(o.Status) == l.Value("status")) && matches(l.Query, o.No, o.Title, o.Customer, names[o.Customer])
	})
	a.render(w, r, 0, page{Key: "offers", Title: "Angebote", Crumbs: []shell.Crumb{{Text: "Angebote"}}}, offersView(l, rows, names))
}

func (a *App) customer(w http.ResponseWriter, r *http.Request) {
	a.customerWith(w, r, 0, dialogState{})
}

func (a *App) customerWith(w http.ResponseWriter, r *http.Request, status int, st dialogState) {
	c, ok := a.store.Customer(r.PathValue("no"))
	if !ok {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	offers := a.store.OffersOf(c.No)
	a.render(w, r, status, page{Key: "B-04", Title: c.Name, Crumbs: []shell.Crumb{{Text: "Tagesgeschäft"}, {Text: "Kunden", Href: pathCustomers}, {Text: c.Name}},
		HeadingExtra: customerHeading(c)}, customerView(c, offers, !a.store.ReadOnly(), st))
}

func (a *App) newCustomer(w http.ResponseWriter, r *http.Request) {
	a.newCustomerWith(w, r, 0, CustomerInput{}, nil)
}

func (a *App) newCustomerWith(w http.ResponseWriter, r *http.Request, status int, in CustomerInput, errs FieldErrors) {
	if a.store.ReadOnly() {
		a.errorPage(w, r, http.StatusForbidden)
		return
	}
	var f *shell.Flash
	if errs != nil {
		f = &shell.Flash{Variant: alert.VariantDanger, Text: errs.Error()}
	}
	a.render(w, r, status, page{Key: "B-06", Title: "Neuer Kunde", Flash: f, Crumbs: []shell.Crumb{{Text: "Tagesgeschäft"}, {Text: "Kunden", Href: pathCustomers}, {Text: "Neuer Kunde"}}},
		newCustomerView(in, errs))
}

func customerInput(r *http.Request) CustomerInput {
	return CustomerInput{Name: r.PostFormValue("name"), Street: r.PostFormValue("street"), Zip: r.PostFormValue("zip"), City: r.PostFormValue("city")}
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	in := customerInput(r)
	c, err := a.store.CreateCustomer(in)
	var errs FieldErrors
	if errors.As(err, &errs) {
		a.newCustomerWith(w, r, http.StatusUnprocessableEntity, in, errs)
		return
	}
	done(w, r, customerURL(c.No), flashData{Text: "Kunde " + c.Name + " angelegt.", Next: "Ansprechpartner eintragen", NextHref: customerURL(c.No) + "#ansprechpartner"})
}

// dialogState is a failed save: which dialog, the message, the values entered and the errors per field
// (GS-43). Row is the row of a list dialog that failed ("" for the new row).
type dialogState struct {
	ID     string
	Error  string
	Values url.Values
	Fields FieldErrors
	Row    string
	Kind   string
}

func (s dialogState) is(id string) bool { return s.ID == id }

// value is the value entered, or fallback when the dialog did not fail.
func (s dialogState) value(id, name, fallback string) string {
	if s.ID != id {
		return fallback
	}
	return s.Values.Get(name)
}

func (s dialogState) field(id, name string) string {
	if s.ID != id {
		return ""
	}
	return s.Fields[name]
}

func failed(id string, r *http.Request, err error) dialogState {
	st := dialogState{ID: id, Values: r.PostForm, Error: err.Error()}
	var errs FieldErrors
	if errors.As(err, &errs) {
		st.Fields = errs
	}
	return st
}

// rowFailed: a row of the list dialog failed; the message names the row and what is wrong (GS-43).
func rowFailed(r *http.Request, err error, kind, row, no string) dialogState {
	st := failed(dialogPositions, r, err)
	st.Kind, st.Row = kind, row
	if st.Fields != nil {
		label := "Neue Position"
		if row != "" {
			label = "Position " + no
		}
		st.Error = label + ": " + st.Fields.Detail()
	}
	return st
}

func (a *App) saveMaster(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	if _, err := a.store.UpdateCustomer(no, customerInput(r)); err != nil {
		a.customerFailed(w, r, err, dialogMaster)
		return
	}
	done(w, r, back(r, customerURL(no)), flashData{Text: "Stammdaten gespeichert."})
}

func (a *App) customerFailed(w http.ResponseWriter, r *http.Request, err error, dialog string) {
	if errors.Is(err, ErrNotFound) {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	a.customerWith(w, r, http.StatusUnprocessableEntity, failed(dialog, r, err))
}

func (a *App) saveContact(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	if _, err := a.store.UpdateContact(no, Contact{Name: r.PostFormValue("contact"), Phone: r.PostFormValue("phone"), Email: r.PostFormValue("email")}); err != nil {
		a.customerFailed(w, r, err, dialogContact)
		return
	}
	done(w, r, back(r, customerURL(no)), flashData{Text: "Ansprechpartner gespeichert."})
}

func (a *App) addNote(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	if _, err := a.store.AddNote(no, r.PostFormValue("text")); err != nil {
		a.customerFailed(w, r, err, dialogNote)
		return
	}
	done(w, r, back(r, customerURL(no)), flashData{Text: "Vermerk hinzugefügt."})
}

func (a *App) addFile(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize+1<<20)
	var name string
	var size int64
	if err := r.ParseMultipartForm(1 << 20); err == nil {
		for _, fh := range r.MultipartForm.File["file"] { // the button and, on phones, the camera
			if fh.Filename != "" {
				name, size = fh.Filename, fh.Size
				break
			}
		}
	} else {
		name, size = "zu groß", MaxFileSize+1
	}
	f, err := a.store.AddFile(no, name, size, r.PostFormValue("kind"), r.PostFormValue("ref"))
	if err != nil {
		st := failed(dialogFile, r, err)
		if st.Fields["file"] == "" {
			st.Error += " Bitte wählen Sie die Datei noch einmal."
		}
		a.customerFailedWith(w, r, err, st)
		return
	}
	done(w, r, back(r, customerURL(no)), flashData{Text: "Datei „" + f.Name + "“ hinzugefügt."})
}

func (a *App) removeFile(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	f, err := a.store.RemoveFile(no, r.PathValue("id"))
	if err != nil {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	done(w, r, customerURL(no), flashData{Text: "Datei „" + f.Name + "“ entfernt."})
}

func (a *App) customerFailedWith(w http.ResponseWriter, r *http.Request, err error, st dialogState) {
	if errors.Is(err, ErrNotFound) {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	a.customerWith(w, r, http.StatusUnprocessableEntity, st)
}

func (a *App) offer(w http.ResponseWriter, r *http.Request) { a.offerWith(w, r, 0, dialogState{}) }

func (a *App) offerWith(w http.ResponseWriter, r *http.Request, status int, st dialogState) {
	o, ok := a.store.Offer(r.PathValue("no"))
	if !ok {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	c, _ := a.store.Customer(o.Customer)
	edit := !a.store.ReadOnly()
	var act templ.Component
	if edit && !o.Locked() {
		act = submitButton(o)
	}
	a.render(w, r, status, page{Key: "B-05", Title: o.No + " · " + o.Title, Heading: o.Title, Actions: act, HeadingExtra: offerHeading(o),
		Crumbs: []shell.Crumb{{Text: "Tagesgeschäft"}, {Text: "Angebote", Href: pathOffers}, {Text: o.No}}}, offerView(o, c, edit, st))
}

func (a *App) offerFailed(w http.ResponseWriter, r *http.Request, err error, st dialogState) {
	if errors.Is(err, ErrNotFound) {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	a.offerWith(w, r, http.StatusUnprocessableEntity, st)
}

func (a *App) saveOffer(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	if _, err := a.store.UpdateOffer(no, OfferInput{Title: r.PostFormValue("title"), Valid: r.PostFormValue("valid")}); err != nil {
		a.offerFailed(w, r, err, failed(dialogOffer, r, err))
		return
	}
	done(w, r, back(r, offerURL(no)), flashData{Text: "Angebot gespeichert."})
}

func (a *App) submitOffer(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	o, err := a.store.Submit(no)
	if errors.Is(err, ErrNotFound) {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	if err != nil {
		done(w, r, offerURL(no), flashData{Variant: "danger", Text: err.Error()})
		return
	}
	done(w, r, offerURL(no), flashData{Text: "Angebot " + o.No + " abgegeben.", Next: "Zu den Angeboten", NextHref: pathOffers})
}

func positionInput(r *http.Request) PositionInput {
	return PositionInput{No: r.PostFormValue("no"), Text: r.PostFormValue("text"), Qty: r.PostFormValue("qty"), Unit: r.PostFormValue("unit"), Price: r.PostFormValue("price"), Long: r.PostFormValue("long")}
}

func (a *App) addPosition(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	p, err := a.store.AddPosition(no, positionInput(r))
	if err != nil {
		a.offerFailed(w, r, err, rowFailed(r, err, kindNewPosition, "", ""))
		return
	}
	done(w, r, back(r, offerURL(no)), flashData{Text: "Position " + p.No + " hinzugefügt."})
}

func (a *App) savePosition(w http.ResponseWriter, r *http.Request) {
	no, id := r.PathValue("no"), r.PathValue("id")
	p, err := a.store.UpdatePosition(no, id, positionInput(r))
	if err != nil {
		a.offerFailed(w, r, err, rowFailed(r, err, kindPosition, id, a.positionNo(no, id)))
		return
	}
	done(w, r, back(r, offerURL(no)), flashData{Text: "Position " + p.No + " gespeichert."})
}

// positionNo is the stored number of a position, for messages.
func (a *App) positionNo(offer, id string) string {
	o, _ := a.store.Offer(offer)
	for _, p := range o.Positions {
		if p.ID == id {
			return p.No
		}
	}
	return ""
}

func (a *App) deletePosition(w http.ResponseWriter, r *http.Request) {
	no, id := r.PathValue("no"), r.PathValue("id")
	p, err := a.store.DeletePosition(no, id)
	if err != nil {
		a.offerFailed(w, r, err, failed(dialogPositions, r, err))
		return
	}
	done(w, r, back(r, offerURL(no)), flashData{Text: "Position " + p.No + " gelöscht."})
}

func (a *App) states(w http.ResponseWriter, r *http.Request) { a.statesWith(w, r, 0, dialogState{}) }

func (a *App) statesWith(w http.ResponseWriter, r *http.Request, status int, st dialogState) {
	p := page{Key: "B-07", Title: "Zustände", Crumbs: []shell.Crumb{{Text: "Referenz", Href: Home}, {Text: "Zustände"}}}
	if r.URL.Query().Get("meldung") == "1" {
		p.Flash = &shell.Flash{Variant: alert.VariantSuccess, Text: "Angebot A-2026-001 abgegeben.", Next: "Zu den Angeboten", NextHref: pathOffers}
	}
	withdrawable, _ := a.store.Offer(offerToWithdraw)
	deletable, _ := a.store.Customer(customerToDelete)
	a.render(w, r, status, p, statesView(a.store.ReadOnly(), a.store.Demo(), withdrawable, deletable, st))
}

// The examples of B-07's actions with a reason and with a typed confirmation.
const (
	offerToWithdraw  = "A-2026-002"
	customerToDelete = "K-1058"
)

func (a *App) withdraw(w http.ResponseWriter, r *http.Request) {
	o, err := a.store.Withdraw(offerToWithdraw, r.PostFormValue("reason"))
	var errs FieldErrors
	switch {
	case errors.As(err, &errs):
		a.statesWith(w, r, http.StatusUnprocessableEntity, failed(dialogWithdraw, r, err))
	case err != nil:
		done(w, r, pathStates, flashData{Variant: "danger", Text: err.Error()})
	default:
		done(w, r, pathStates, flashData{Text: "Abgabe von " + o.No + " zurückgenommen; der Grund steht in den Vermerken des Kunden.", Next: "Zum Angebot", NextHref: offerURL(o.No)})
	}
}

func (a *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	c, err := a.store.DeleteCustomer(customerToDelete, r.PostFormValue("confirm"))
	var errs FieldErrors
	switch {
	case errors.As(err, &errs):
		a.statesWith(w, r, http.StatusUnprocessableEntity, failed(dialogDelete, r, err))
	case errors.Is(err, ErrNotFound):
		done(w, r, pathStates, flashData{Variant: "danger", Text: "Diesen Kunden gibt es nicht mehr. „Beispieldaten zurücksetzen“ legt ihn wieder an."})
	default:
		done(w, r, pathStates, flashData{Text: "Kunde " + c.Name + " endgültig gelöscht.", Next: "Zu den Kunden", NextHref: pathCustomers})
	}
}

func (a *App) toggleDemo(w http.ResponseWriter, r *http.Request) {
	on := !a.store.Demo()
	a.store.SetDemo(on)
	text := "Der Hinweis „Demo-Instanz“ ist aus."
	if on {
		text = "Der Hinweis „Demo-Instanz“ ist an und steht jetzt auf jeder Seite."
	}
	done(w, r, pathStates, flashData{Text: text})
}

// lateFigures is the part of B-07 that comes after the page (GS-64).
func (a *App) lateFigures(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = lateFiguresView(a.store.Counts()).Render(germanContext(r.Context()), w)
}

func (a *App) toggleReadOnly(w http.ResponseWriter, r *http.Request) {
	on := !a.store.ReadOnly()
	a.store.SetReadOnly(on)
	text := "Nur lesen ist aus: Bearbeiten geht wieder."
	if on {
		text = "Nur lesen ist an: Alle Knöpfe zum Bearbeiten sind ausgeblendet."
	}
	done(w, r, pathStates, flashData{Text: text})
}

func (a *App) reset(w http.ResponseWriter, r *http.Request) {
	a.store.Reset()
	done(w, r, back(r, Home), flashData{Text: "Beispieldaten zurückgesetzt."})
}

func (a *App) brand(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 0, page{Key: "B-08", Title: "Marke und Form", Crumbs: []shell.Crumb{{Text: "Referenz", Href: Home}, {Text: "Marke und Form"}}}, brandView())
}

func (a *App) assistantPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 0, page{Key: "B-09", Title: "KI-Assistent", Crumbs: []shell.Crumb{{Text: "Referenz", Href: Home}, {Text: "KI-Assistent"}}}, assistantView())
}

func (a *App) concept(c concept) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.render(w, r, 0, page{Key: "concept", Title: "Konzept " + c.Title, Crumbs: []shell.Crumb{{Text: "Referenz", Href: Home}, {Text: "Konzept " + c.Title}}},
			conceptView(c))
	}
}

var requestFilter = listpage.Filter{Name: "stand", Label: "Alle Stände", Options: func() []listpage.Option {
	var os []listpage.Option
	for _, s := range RequestStatuses {
		os = append(os, listpage.Option{Value: string(s), Label: s.Label()})
	}
	return os
}()}

func (a *App) requests(w http.ResponseWriter, r *http.Request) {
	names := a.customerNames()
	l := listpage.New(r, listpage.Config{Placeholder: "Betreff, Kunde oder Nummer"}, requestFilter)
	all := a.store.Requests()
	// GS-53: open ones first, by deadline.
	slices.SortStableFunc(all, func(x, y Request) int {
		if (x.Status == RequestOpen) != (y.Status == RequestOpen) {
			if x.Status == RequestOpen {
				return -1
			}
			return 1
		}
		return x.Due.Compare(y.Due)
	})
	rows := listpage.Paginate(l, all, func(q Request) bool {
		return (l.Value("stand") == "" || string(q.Status) == l.Value("stand")) && matches(l.Query, q.Subject, q.No, names[q.Customer])
	})
	a.render(w, r, 0, page{Key: "requests", Title: "Anfragen", Crumbs: []shell.Crumb{{Text: "Erledigen"}, {Text: "Anfragen"}}}, requestsView(l, rows, names))
}

func (a *App) request(w http.ResponseWriter, r *http.Request) { a.requestWith(w, r, 0, dialogState{}) }

func (a *App) requestWith(w http.ResponseWriter, r *http.Request, status int, st dialogState) {
	q, ok := a.store.Request(r.PathValue("no"))
	if !ok {
		a.errorPage(w, r, http.StatusNotFound)
		return
	}
	c, _ := a.store.Customer(q.Customer)
	var open []Offer
	for _, o := range a.store.OffersOf(c.No) {
		if o.Status == OfferDraft {
			open = append(open, o)
		}
	}
	a.render(w, r, status, page{Key: "B-10", Title: q.Subject, HeadingExtra: requestHeading(q),
		Crumbs: []shell.Crumb{{Text: "Erledigen"}, {Text: "Anfragen", Href: pathRequests}, {Text: q.No}}}, requestView(q, c, open, !a.store.ReadOnly(), st))
}

func (a *App) sendReply(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	q, err := a.store.Send(no, r.PostFormValue("reply"))
	var errs FieldErrors
	switch {
	case errors.Is(err, ErrNotFound):
		a.errorPage(w, r, http.StatusNotFound)
	case errors.As(err, &errs):
		a.requestWith(w, r, http.StatusUnprocessableEntity, failed("antwort", r, err))
	case err != nil:
		done(w, r, requestURL(no), flashData{Variant: "danger", Text: err.Error()})
	default:
		done(w, r, requestURL(no), flashData{Text: "Antwort an " + a.customerNames()[q.Customer] + " gesendet.", Next: "Zur nächsten Anfrage", NextHref: pathRequests + "?stand=" + string(RequestOpen)})
	}
}

func (a *App) discard(w http.ResponseWriter, r *http.Request) {
	no := r.PathValue("no")
	_, err := a.store.Discard(no, r.PostFormValue("reason"), r.PostFormValue("note"))
	var errs FieldErrors
	switch {
	case errors.Is(err, ErrNotFound):
		a.errorPage(w, r, http.StatusNotFound)
	case errors.As(err, &errs):
		a.requestWith(w, r, http.StatusUnprocessableEntity, failed(dialogDiscard, r, err))
	case err != nil:
		done(w, r, requestURL(no), flashData{Variant: "danger", Text: err.Error()})
	default:
		done(w, r, requestURL(no), flashData{Text: "Vorschlag verworfen; nichts wurde gesendet.", Next: "Zur nächsten Anfrage", NextHref: pathRequests + "?stand=" + string(RequestOpen)})
	}
}

// --- error pages in the frame (GS-62) -------------------------------------------------------------

func (a *App) notFound(w http.ResponseWriter, r *http.Request) {
	a.errorPage(w, r, http.StatusNotFound)
}

func (a *App) forbidden(w http.ResponseWriter, r *http.Request) {
	a.errorPage(w, r, http.StatusForbidden)
}

func (a *App) failure(w http.ResponseWriter, r *http.Request) {
	a.errorPage(w, r, http.StatusInternalServerError)
}

func (a *App) errorPage(w http.ResponseWriter, r *http.Request, status int) {
	e := errorTexts[status]
	a.render(w, r, status, page{Key: "error", Title: e.title, Crumbs: []shell.Crumb{{Text: e.title}}}, errorView(e))
}

type errorText struct{ title, text string }

var errorTexts = map[int]errorText{
	http.StatusNotFound:            {"Seite nicht gefunden", "Diese Seite gibt es nicht, oder sie wurde verschoben."},
	http.StatusForbidden:           {"Keine Berechtigung", "Für diese Seite fehlt Ihnen die Berechtigung. Fragen Sie bei Bedarf die Person, die Ihr Konto verwaltet."},
	http.StatusInternalServerError: {"Etwas ist schiefgegangen", "Die Seite konnte nicht geladen werden. Bitte versuchen Sie es gleich noch einmal."},
}

// pageKeys are the keys the sidebar marks; tests use them.
var pageKeys = []string{"B-11", "requests", "B-10", "concept", "overview", "B-01", "B-02", "B-03", "offers", "B-04", "B-05", "B-06", "B-07", "B-08", "B-09"}

func knownKey(k string) bool { return slices.Contains(pageKeys, k) }
