package reference

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

// The data of the reference app lives in memory (GS-75: sample data without a recognisable trade). Today
// is fixed, so pages and screenshots stay the same from day to day.

// Today is the date the app pretends it is.
var Today = time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC)

// Signal is the traffic light of a customer's open items (GS-52: a dot with a word).
type Signal string

const (
	SignalOverdue Signal = "rot"
	SignalDue     Signal = "gelb"
	SignalNone    Signal = "gruen"
)

// Signals in the order of the filter and the sorting (GS-53: what needs action first).
var Signals = []Signal{SignalOverdue, SignalDue, SignalNone}

func (s Signal) Label() string {
	switch s {
	case SignalOverdue:
		return "überfällig"
	case SignalDue:
		return "fällig"
	}
	return "keine"
}

func (s Signal) color() string {
	switch s {
	case SignalOverdue:
		return "red"
	case SignalDue:
		return "yellow"
	}
	return "green"
}

func (s Signal) rank() int { return slices.Index(Signals, s) }

// CustomerStatus is shown as a pill.
type CustomerStatus string

const (
	CustomerActive  CustomerStatus = "aktiv"
	CustomerNew     CustomerStatus = "neu"
	CustomerDormant CustomerStatus = "ruhend"
)

func (s CustomerStatus) Label() string {
	switch s {
	case CustomerNew:
		return "Neu"
	case CustomerDormant:
		return "Ruhend"
	}
	return "Aktiv"
}

type Note struct {
	Date time.Time
	Text string
}

type Customer struct {
	No      string // K-1001
	Name    string
	Street  string
	Zip     string
	City    string
	Status  CustomerStatus
	Signal  Signal
	Revenue int64 // cents, this year
	Contact Contact
	Notes   []Note
	Payment string // payment terms, read-only (no use case to change them here)
	Overdue int64  // cents of overdue items, from the accounting app
	Created time.Time
}

type Contact struct {
	Name, Phone, Email string
}

// OfferStatus is shown as a pill.
type OfferStatus string

const (
	OfferDraft     OfferStatus = "entwurf"
	OfferSubmitted OfferStatus = "abgegeben"
	OfferAccepted  OfferStatus = "angenommen"
	OfferDeclined  OfferStatus = "abgelehnt"
)

var OfferStatuses = []OfferStatus{OfferDraft, OfferSubmitted, OfferAccepted, OfferDeclined}

func (s OfferStatus) Label() string {
	switch s {
	case OfferSubmitted:
		return "Abgegeben"
	case OfferAccepted:
		return "Angenommen"
	case OfferDeclined:
		return "Abgelehnt"
	}
	return "Entwurf"
}

type Position struct {
	ID    string
	No    string // 01.01.0010
	Text  string
	Qty   Decimal
	Unit  string
	Price int64 // cents per unit
}

// Amount is quantity × price, in cents.
func (p Position) Amount() int64 { return p.Qty.Times(p.Price) }

type Offer struct {
	No        string // A-2026-001
	Customer  string // customer number
	Title     string
	Due       time.Time // to be submitted by
	Valid     time.Time // valid until
	Status    OfferStatus
	Submitted time.Time
	Positions []Position
}

// Total is the sum of the positions, in cents.
func (o Offer) Total() int64 {
	var t int64
	for _, p := range o.Positions {
		t += p.Amount()
	}
	return t
}

// Locked: the positions can no longer change (GS-59).
func (o Offer) Locked() bool { return o.Status != OfferDraft }

// NextNo proposes the number of a new position: the last one plus 10 (GS-57).
func (o Offer) NextNo() string {
	if len(o.Positions) == 0 {
		return "01.01.0010"
	}
	last := o.Positions[len(o.Positions)-1].No
	var a, b, c int
	if _, err := fmt.Sscanf(last, "%d.%d.%d", &a, &b, &c); err != nil {
		return ""
	}
	return fmt.Sprintf("%02d.%02d.%04d", a, b, c+10)
}

// Errors of the store; FieldErrors are per field and shown at the field (GS-43).
var ErrNotFound = errors.New("nicht gefunden")

type FieldErrors map[string]string

func (e FieldErrors) Error() string { return "Bitte prüfen Sie die markierten Felder." }

// Detail is the messages of all fields in a fixed order, for a row of a list dialog, where there is no
// room for a message at the field.
func (e FieldErrors) Detail() string {
	var out []string
	for _, k := range []string{"no", "text", "qty", "unit", "price"} {
		if m := e[k]; m != "" {
			out = append(out, m)
		}
	}
	return strings.Join(out, " ")
}

// Store holds customers and offers.
type Store struct {
	mu        sync.Mutex
	customers []*Customer
	offers    []*Offer
	nextID    int
	readOnly  bool
}

// NewStore returns the sample data.
func NewStore() *Store {
	s := &Store{}
	s.seed()
	return s
}

// Reset puts the sample data back.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.customers, s.offers, s.nextID, s.readOnly = nil, nil, 0, false
	s.seed()
}

// ReadOnly: the company is locked, everything can be read but nothing changed (GS-63).
func (s *Store) ReadOnly() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.readOnly
}

func (s *Store) SetReadOnly(on bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.readOnly = on
}

func (s *Store) Customers() []Customer {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Customer, len(s.customers))
	for i, c := range s.customers {
		out[i] = *c
	}
	// GS-53: what needs action first, then by name.
	slices.SortStableFunc(out, func(a, b Customer) int {
		if d := a.Signal.rank() - b.Signal.rank(); d != 0 {
			return d
		}
		return strings.Compare(fold(a.Name), fold(b.Name))
	})
	return out
}

func (s *Store) Customer(no string) (Customer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c := s.customer(no); c != nil {
		return *c, true
	}
	return Customer{}, false
}

func (s *Store) customer(no string) *Customer {
	for _, c := range s.customers {
		if c.No == no {
			return c
		}
	}
	return nil
}

// CustomerInput is the form of the island „Stammdaten“ and of the create page.
type CustomerInput struct {
	Name, Street, Zip, City string
}

func (in CustomerInput) check() FieldErrors {
	errs := FieldErrors{}
	if strings.TrimSpace(in.Name) == "" {
		errs["name"] = "Bitte geben Sie einen Namen ein."
	}
	if strings.TrimSpace(in.City) == "" {
		errs["city"] = "Bitte geben Sie einen Ort ein."
	}
	if z := strings.TrimSpace(in.Zip); z != "" && (len(z) != 5 || strings.Trim(z, "0123456789") != "") {
		errs["zip"] = "Die Postleitzahl hat fünf Ziffern."
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (s *Store) CreateCustomer(in CustomerInput) (Customer, error) {
	if errs := in.check(); errs != nil {
		return Customer{}, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c := &Customer{No: fmt.Sprintf("K-%d", 1001+len(s.customers)), Name: strings.TrimSpace(in.Name), Street: strings.TrimSpace(in.Street),
		Zip: strings.TrimSpace(in.Zip), City: strings.TrimSpace(in.City), Status: CustomerNew, Signal: SignalNone, Payment: "14 Tage netto", Created: Today}
	s.customers = append(s.customers, c)
	return *c, nil
}

func (s *Store) UpdateCustomer(no string, in CustomerInput) (Customer, error) {
	if errs := in.check(); errs != nil {
		return Customer{}, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.customer(no)
	if c == nil {
		return Customer{}, ErrNotFound
	}
	c.Name, c.Street, c.Zip, c.City = strings.TrimSpace(in.Name), strings.TrimSpace(in.Street), strings.TrimSpace(in.Zip), strings.TrimSpace(in.City)
	return *c, nil
}

func (s *Store) UpdateContact(no string, k Contact) (Customer, error) {
	errs := FieldErrors{}
	if strings.TrimSpace(k.Name) == "" {
		errs["contact"] = "Bitte geben Sie einen Namen ein."
	}
	if e := strings.TrimSpace(k.Email); e != "" && !strings.Contains(e, "@") {
		errs["email"] = "Die E-Mail-Adresse braucht ein @."
	}
	if len(errs) > 0 {
		return Customer{}, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.customer(no)
	if c == nil {
		return Customer{}, ErrNotFound
	}
	c.Contact = Contact{Name: strings.TrimSpace(k.Name), Phone: strings.TrimSpace(k.Phone), Email: strings.TrimSpace(k.Email)}
	return *c, nil
}

// AddNote is a small addition to the customer (GS-44).
func (s *Store) AddNote(no, text string) (Customer, error) {
	if strings.TrimSpace(text) == "" {
		return Customer{}, FieldErrors{"text": "Bitte schreiben Sie einen Vermerk."}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.customer(no)
	if c == nil {
		return Customer{}, ErrNotFound
	}
	c.Notes = append([]Note{{Date: Today, Text: strings.TrimSpace(text)}}, c.Notes...)
	return *c, nil
}

func (s *Store) Offers() []Offer {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Offer, len(s.offers))
	for i, o := range s.offers {
		out[i] = *o
		out[i].Positions = slices.Clone(o.Positions)
	}
	// GS-53: drafts by due date first, then the others, newest first.
	slices.SortStableFunc(out, func(a, b Offer) int {
		if (a.Status == OfferDraft) != (b.Status == OfferDraft) {
			if a.Status == OfferDraft {
				return -1
			}
			return 1
		}
		if a.Status == OfferDraft {
			return a.Due.Compare(b.Due)
		}
		return strings.Compare(b.No, a.No)
	})
	return out
}

// OffersOf are a customer's offers in the order of Offers.
func (s *Store) OffersOf(customer string) []Offer {
	var out []Offer
	for _, o := range s.Offers() {
		if o.Customer == customer {
			out = append(out, o)
		}
	}
	return out
}

func (s *Store) Offer(no string) (Offer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if o := s.offer(no); o != nil {
		c := *o
		c.Positions = slices.Clone(o.Positions)
		return c, true
	}
	return Offer{}, false
}

func (s *Store) offer(no string) *Offer {
	for _, o := range s.offers {
		if o.No == no {
			return o
		}
	}
	return nil
}

// OfferInput is the form of the island „Angebot“.
type OfferInput struct {
	Title, Valid string
}

func (s *Store) UpdateOffer(no string, in OfferInput) (Offer, error) {
	errs := FieldErrors{}
	if strings.TrimSpace(in.Title) == "" {
		errs["title"] = "Bitte geben Sie einen Titel ein."
	}
	valid, err := parseDate(in.Valid)
	if err != nil {
		errs["valid"] = "Bitte geben Sie ein Datum als TT.MM.JJJJ ein."
	}
	if len(errs) > 0 {
		return Offer{}, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o := s.offer(no)
	if o == nil {
		return Offer{}, ErrNotFound
	}
	o.Title, o.Valid = strings.TrimSpace(in.Title), valid
	return *o, nil
}

// Submit makes the offer read-only.
func (s *Store) Submit(no string) (Offer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o := s.offer(no)
	if o == nil {
		return Offer{}, ErrNotFound
	}
	if o.Status != OfferDraft {
		return *o, errors.New("Das Angebot ist schon abgegeben.")
	}
	o.Status, o.Submitted = OfferSubmitted, Today
	return *o, nil
}

// PositionInput is one row of the dialog „Positionen bearbeiten“.
type PositionInput struct {
	No, Text, Qty, Unit, Price string
}

func (in PositionInput) parse() (Position, FieldErrors) {
	errs := FieldErrors{}
	p := Position{No: strings.TrimSpace(in.No), Text: strings.TrimSpace(in.Text), Unit: strings.TrimSpace(in.Unit)}
	var a, b, c int
	if n, _ := fmt.Sscanf(p.No, "%d.%d.%d", &a, &b, &c); n != 3 || len(p.No) != 10 {
		errs["no"] = "Die Positionsnummer hat die Form 01.01.0010."
	}
	if p.Text == "" {
		errs["text"] = "Bitte geben Sie eine Bezeichnung ein."
	}
	var err error
	if p.Qty, err = ParseDecimal(in.Qty); err != nil || p.Qty <= 0 {
		errs["qty"] = "Die Menge ist eine Zahl größer als 0, z. B. 12,5."
	}
	if p.Price, err = ParseCents(in.Price); err != nil || p.Price < 0 {
		errs["price"] = "Der Preis ist ein Betrag in Euro, z. B. 48,50."
	}
	if p.Unit == "" {
		p.Unit = "St"
	}
	if len(errs) > 0 {
		return p, errs
	}
	return p, nil
}

var errLocked = errors.New("Das Angebot ist abgegeben; die Positionen lassen sich nicht mehr ändern.")

func (s *Store) AddPosition(offer string, in PositionInput) (Position, error) {
	p, errs := in.parse()
	if errs != nil {
		return p, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o := s.offer(offer)
	if o == nil {
		return p, ErrNotFound
	}
	if o.Locked() {
		return p, errLocked
	}
	if slices.ContainsFunc(o.Positions, func(x Position) bool { return x.No == p.No }) {
		return p, FieldErrors{"no": "Die Position " + p.No + " gibt es schon."}
	}
	s.nextID++
	p.ID = fmt.Sprintf("p%d", 1000+s.nextID)
	o.Positions = append(o.Positions, p)
	sortPositions(o)
	return p, nil
}

func (s *Store) UpdatePosition(offer, id string, in PositionInput) (Position, error) {
	p, errs := in.parse()
	p.ID = id
	if errs != nil {
		return p, errs
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	o := s.offer(offer)
	if o == nil {
		return p, ErrNotFound
	}
	if o.Locked() {
		return p, errLocked
	}
	i := slices.IndexFunc(o.Positions, func(x Position) bool { return x.ID == id })
	if i < 0 {
		return p, ErrNotFound
	}
	if slices.ContainsFunc(o.Positions, func(x Position) bool { return x.No == p.No && x.ID != id }) {
		return p, FieldErrors{"no": "Die Position " + p.No + " gibt es schon."}
	}
	o.Positions[i] = p
	sortPositions(o)
	return p, nil
}

func (s *Store) DeletePosition(offer, id string) (Position, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o := s.offer(offer)
	if o == nil {
		return Position{}, ErrNotFound
	}
	if o.Locked() {
		return Position{}, errLocked
	}
	i := slices.IndexFunc(o.Positions, func(x Position) bool { return x.ID == id })
	if i < 0 {
		return Position{}, ErrNotFound
	}
	p := o.Positions[i]
	o.Positions = slices.Delete(o.Positions, i, i+1)
	return p, nil
}

func sortPositions(o *Offer) {
	slices.SortStableFunc(o.Positions, func(a, b Position) int { return strings.Compare(a.No, b.No) })
}

// Counts are the figures of the start page (GS-30).
type Counts struct {
	Drafts, DueThisWeek, Submitted int
	Overdue, Due, New              int
	OverdueSum                     int64
}

func (s *Store) Counts() Counts {
	var c Counts
	week := Today.AddDate(0, 0, 7)
	for _, o := range s.Offers() {
		if o.Status == OfferSubmitted {
			c.Submitted++
		}
		if o.Status == OfferDraft {
			c.Drafts++
			if !o.Due.After(week) {
				c.DueThisWeek++
			}
		}
	}
	for _, k := range s.Customers() {
		if k.Signal == SignalDue {
			c.Due++
		}
		if k.Signal == SignalOverdue {
			c.Overdue++
			c.OverdueSum += k.Overdue
		}
		if k.Status == CustomerNew {
			c.New++
		}
	}
	return c
}

// Burning are the drafts due within three days (GS-30: „Brennt“).
func (s *Store) Burning() []Offer {
	var out []Offer
	for _, o := range s.Offers() {
		if o.Status == OfferDraft && !o.Due.After(Today.AddDate(0, 0, 3)) {
			out = append(out, o)
		}
	}
	return out
}

func day(m time.Month, d int) time.Time { return time.Date(2026, m, d, 0, 0, 0, 0, time.UTC) }

func (s *Store) seed() {
	kinds := []string{"Hausverwaltung", "Genossenschaft", "Stiftung", "Bürgerhaus", "Verwaltung", "Zweckverband", "Wohnungsbau", "Gemeindehaus"}
	places := []string{"Südstadt", "Lindenhof", "Am Markt", "Mühlenweg", "Westend", "Höhenberg", "Seeufer", "Altstadt"}
	cities := []string{"Musterstadt", "Neustadt", "Altdorf", "Mühlhausen", "Bergheim", "Seefeld", "Talheim", "Großenbach"}
	streets := []string{"Hauptstraße", "Lindenallee", "Bahnhofstraße", "Am Anger", "Schulweg", "Gartenstraße", "Ringstraße", "Kirchplatz"}
	first := &Customer{No: "K-1001", Name: "Wohnungsgesellschaft Nordstadt", Street: "Nordring 4", Zip: "12345", City: "Musterstadt",
		Status: CustomerActive, Signal: SignalOverdue, Revenue: 18_640_000, Overdue: 412_350, Payment: "30 Tage netto",
		Contact: Contact{Name: "Erika Beispiel", Phone: "0123 456789-12", Email: "erika.beispiel@example.org"},
		Notes: []Note{{day(time.November, 12), "Rückfrage zum Jahresvertrag: Frau Beispiel ruft nach dem 20.11. zurück."},
			{day(time.October, 2), "Neue Ansprechpartnerin seit Oktober."}}, Created: day(time.February, 3)}
	s.customers = append(s.customers, first)
	for i := 1; i < 58; i++ {
		k := kinds[i%len(kinds)]
		p := places[(i/len(kinds)+i)%len(places)]
		sig := SignalNone
		switch {
		case i%7 == 0:
			sig = SignalOverdue
		case i%4 == 0:
			sig = SignalDue
		}
		st := CustomerActive
		switch {
		case i%11 == 0:
			st = CustomerDormant
		case i > 52:
			st = CustomerNew
		}
		c := &Customer{No: fmt.Sprintf("K-%d", 1001+i), Name: k + " " + p, Street: fmt.Sprintf("%s %d", streets[i%len(streets)], 1+i%37),
			Zip: fmt.Sprintf("%05d", 10000+i*731%89999), City: cities[(i*3)%len(cities)], Status: st, Signal: sig,
			Revenue: int64(120_000 + (i*7919)%4_800_000), Payment: "14 Tage netto", Created: day(time.Month(1+i%10), 1+i%27),
			Contact: Contact{Name: []string{"Max Mustermann", "Erika Muster", "Jonas Beispiel", "Lea Probe"}[i%4], Phone: fmt.Sprintf("0123 %06d", 100000+i*313)}}
		if sig == SignalOverdue {
			c.Overdue = int64(25_000 + (i*1301)%300_000)
		}
		if i == 13 {
			c.Contact = Contact{} // shows „–“ (GS-40)
			c.Street, c.Zip = "", ""
		}
		s.customers = append(s.customers, c)
	}

	texts := []struct {
		text, unit string
		qty        float64
		price      int64
	}{
		{"Grundleistung je Einsatz", "Einsatz", 52, 8_950},
		{"Zusatzleistung nach Aufwand", "Std", 40, 4_850},
		{"Stundensatz Fachkraft", "Std", 120, 5_600},
		{"Stundensatz Hilfskraft", "Std", 80, 3_900},
		{"Bereitschaft je Monat", "Monat", 12, 24_000},
		{"Anfahrt pauschal", "Einsatz", 52, 1_800},
		{"Material nach Verbrauch", "psch", 1, 150_000},
		{"Dokumentation je Einsatz", "Einsatz", 52, 950},
		{"Zuschlag Wochenende", "Std", 16, 1_450},
		{"Zuschlag Nacht", "Std", 8, 1_900},
		{"Abstimmung vor Ort", "Termin", 4, 12_000},
		{"Jahresbericht", "St", 1, 45_000},
	}
	nos := []string{"01.01.0010", "01.01.0020", "01.01.0030", "01.01.0040", "01.02.0010", "01.02.0020", "01.02.0030", "01.03.0010", "01.03.0020", "01.03.0030", "02.01.0010", "02.01.0020"}
	positions := func(n int) []Position {
		var ps []Position
		for i := range n {
			t := texts[i%len(texts)]
			s.nextID++
			ps = append(ps, Position{ID: fmt.Sprintf("p%d", 1000+s.nextID), No: nos[i%len(nos)], Text: t.text, Qty: DecimalOf(t.qty), Unit: t.unit, Price: t.price})
		}
		return ps
	}
	add := func(customer, title string, due, valid time.Time, st OfferStatus, n int) {
		o := &Offer{No: fmt.Sprintf("A-2026-%03d", len(s.offers)+1), Customer: customer, Title: title, Due: due, Valid: valid, Status: st, Positions: positions(n)}
		if st != OfferDraft {
			o.Submitted = due.AddDate(0, 0, -2)
		}
		s.offers = append(s.offers, o)
	}
	add("K-1001", "Jahresvertrag 2027", day(time.November, 18), day(time.December, 31), OfferDraft, 12)
	add("K-1001", "Zusatzleistungen Winter", day(time.November, 3), day(time.December, 15), OfferSubmitted, 5)
	titles := []string{"Rahmenvertrag", "Einzelauftrag", "Verlängerung", "Nachtrag", "Sonderleistung", "Jahresvertrag"}
	for i := range 12 {
		st := []OfferStatus{OfferAccepted, OfferDeclined, OfferSubmitted, OfferAccepted}[i%4]
		add("K-1001", titles[i%len(titles)]+" "+fmt.Sprint(2025-i/6), day(time.Month(1+i%10), 5+i), day(time.Month(1+i%10), 5+i).AddDate(0, 2, 0), st, 3+i%5)
	}
	for i := range 26 {
		c := fmt.Sprintf("K-%d", 1002+(i*5)%57)
		st := OfferStatuses[i%4]
		due := day(time.November, 16+i%14)
		if st != OfferDraft {
			due = day(time.Month(1+i%10), 2+i)
		}
		add(c, titles[i%len(titles)], due, due.AddDate(0, 1, 0), st, 2+i%6)
	}
}
