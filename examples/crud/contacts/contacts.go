// Package contacts is the CRUD demo's domain: contacts and an in-memory store. It knows nothing about
// HTTP or HTML; the UI components reach it through the small interfaces they define themselves.
package contacts

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	Lead     Status = "lead"
	Active   Status = "active"
	Inactive Status = "inactive"
)

// Statuses in the order the UI offers them.
var Statuses = []Status{Lead, Active, Inactive}

func (s Status) Label() string {
	switch s {
	case Lead:
		return "Lead"
	case Active:
		return "Active"
	case Inactive:
		return "Inactive"
	}
	return string(s)
}

func (s Status) Valid() bool { return slices.Contains(Statuses, s) }

type Contact struct {
	ID      string
	Name    string
	Company string
	Email   string
	Status  Status
	Tags    []string
	Notes   string
	Updated time.Time
	// Version grows with every change; an update with an older version is a conflict.
	Version int
}

// Event is an entry of a contact's activity.
type Event struct {
	When time.Time
	Icon string
	Text string
}

// Query narrows List: a search over name, company and e-mail, a status, and the order
// ("name", "company" or "updated").
type Query struct {
	Search string
	Status Status
	Sort   string
}

var (
	ErrNotFound = errors.New("contact not found")
	ErrConflict = errors.New("contact was changed in the meantime")
)

// FieldErrors are validation errors by field name.
type FieldErrors map[string]string

func (e FieldErrors) Error() string { return fmt.Sprintf("%d invalid fields", len(e)) }

// Validate checks a contact before it is stored.
func Validate(c Contact) FieldErrors {
	e := FieldErrors{}
	if strings.TrimSpace(c.Name) == "" {
		e["name"] = "Enter a name."
	}
	if c.Email != "" {
		if a, err := mail.ParseAddress(c.Email); err != nil || a.Address != c.Email {
			e["email"] = "Enter an e-mail address like ada@example.com."
		}
	}
	if !c.Status.Valid() {
		e["status"] = "Pick a status."
	}
	if len(e) == 0 {
		return nil
	}
	return e
}

// Store keeps the contacts in memory; it is safe for concurrent use.
type Store struct {
	mu       sync.Mutex
	byID     map[string]Contact
	deleted  map[string]Contact
	activity map[string][]Event
	seq      int
	now      func() time.Time
}

// NewStore returns a store with the demo's contacts.
func NewStore() *Store {
	s := &Store{byID: map[string]Contact{}, deleted: map[string]Contact{}, activity: map[string][]Event{}, now: time.Now}
	base := time.Date(2026, 11, 2, 9, 0, 0, 0, time.UTC)
	for i, c := range seed {
		s.seq++
		c.ID = fmt.Sprint(s.seq)
		c.Version = 1
		c.Updated = base.Add(-time.Duration(i) * 7 * time.Hour)
		s.byID[c.ID] = c
		s.activity[c.ID] = []Event{{When: c.Updated, Icon: "ti-user-plus", Text: "Contact created"}}
	}
	return s
}

func (s *Store) List(_ context.Context, q Query) ([]Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Contact
	words := strings.Fields(strings.ToLower(q.Search))
	for _, c := range s.byID {
		if q.Status != "" && c.Status != q.Status {
			continue
		}
		text := strings.ToLower(c.Name + " " + c.Company + " " + c.Email + " " + strings.Join(c.Tags, " "))
		match := true
		for _, w := range words {
			if !strings.Contains(text, w) {
				match = false
				break
			}
		}
		if match {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		switch q.Sort {
		case "company":
			if a.Company != b.Company {
				return strings.ToLower(a.Company) < strings.ToLower(b.Company)
			}
		case "updated":
			if !a.Updated.Equal(b.Updated) {
				return a.Updated.After(b.Updated)
			}
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	return out, nil
}

func (s *Store) Get(_ context.Context, id string) (Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return Contact{}, ErrNotFound
	}
	return c, nil
}

// Count is the number of contacts with the status, all for "".
func (s *Store) Count(_ context.Context, status Status) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, c := range s.byID {
		if status == "" || c.Status == status {
			n++
		}
	}
	return n, nil
}

func (s *Store) Create(_ context.Context, c Contact) (Contact, error) {
	c = clean(c)
	if e := Validate(c); e != nil {
		return c, e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	c.ID, c.Version, c.Updated = fmt.Sprint(s.seq), 1, s.now()
	s.byID[c.ID] = c
	s.log(c.ID, "ti-user-plus", "Contact created")
	return c, nil
}

// Update stores c if c.Version is the stored version; otherwise ErrConflict.
func (s *Store) Update(_ context.Context, c Contact) (Contact, error) {
	c = clean(c)
	if e := Validate(c); e != nil {
		return c, e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.byID[c.ID]
	if !ok {
		return c, ErrNotFound
	}
	if old.Version != c.Version {
		return c, ErrConflict
	}
	c.Version, c.Updated = old.Version+1, s.now()
	s.byID[c.ID] = c
	s.log(c.ID, "ti-pencil", "Details edited")
	return c, nil
}

func (s *Store) SetStatus(_ context.Context, id string, st Status) (Contact, error) {
	if !st.Valid() {
		return Contact{}, FieldErrors{"status": "Pick a status."}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return c, ErrNotFound
	}
	if c.Status != st {
		s.log(id, "ti-flag", "Status changed from "+c.Status.Label()+" to "+st.Label())
	}
	c.Status, c.Version, c.Updated = st, c.Version+1, s.now()
	s.byID[id] = c
	return c, nil
}

// Delete removes a contact; Restore brings it back (undo).
func (s *Store) Delete(_ context.Context, id string) (Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.byID[id]
	if !ok {
		return c, ErrNotFound
	}
	delete(s.byID, id)
	s.deleted[id] = c
	return c, nil
}

func (s *Store) Restore(_ context.Context, id string) (Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.deleted[id]
	if !ok {
		return c, ErrNotFound
	}
	delete(s.deleted, id)
	s.byID[id] = c
	s.log(id, "ti-arrow-back-up", "Restored after deletion")
	return c, nil
}

// Activity is the contact's history, newest first.
func (s *Store) Activity(_ context.Context, id string) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.byID[id]; !ok {
		return nil, ErrNotFound
	}
	ev := slices.Clone(s.activity[id])
	slices.Reverse(ev)
	return ev, nil
}

func (s *Store) log(id, icon, text string) {
	s.activity[id] = append(s.activity[id], Event{When: s.now(), Icon: icon, Text: text})
}

func clean(c Contact) Contact {
	c.Name, c.Company, c.Email = strings.TrimSpace(c.Name), strings.TrimSpace(c.Company), strings.TrimSpace(c.Email)
	tags := c.Tags[:0:0]
	for _, t := range c.Tags {
		if t = strings.ToLower(strings.TrimSpace(t)); t != "" && !slices.Contains(tags, t) {
			tags = append(tags, t)
		}
	}
	c.Tags = tags
	return c
}

// ParseTags splits a comma-separated input into tags.
func ParseTags(s string) []string { return strings.Split(s, ",") }

var seed = []Contact{
	{Name: "Ada Lovelace", Company: "Analytical Engines", Email: "ada@engines.example", Status: Active, Tags: []string{"vip", "math"}},
	{Name: "Grace Hopper", Company: "Navy Labs", Email: "grace@navy.example", Status: Active, Tags: []string{"compiler"}},
	{Name: "Alan Turing", Company: "Bletchley Park", Email: "alan@bletchley.example", Status: Inactive, Tags: []string{"math"}},
	{Name: "Katherine Johnson", Company: "NASA Langley", Email: "katherine@nasa.example", Status: Active, Tags: []string{"orbits"}},
	{Name: "Linus Torvalds", Company: "Kernel Org", Email: "linus@kernel.example", Status: Lead, Tags: []string{"os"}},
	{Name: "Margaret Hamilton", Company: "Apollo Software", Email: "margaret@apollo.example", Status: Active, Tags: []string{"vip", "space"}},
	{Name: "Dennis Ritchie", Company: "Bell Labs", Email: "dmr@bell.example", Status: Inactive, Tags: []string{"c", "os"}},
	{Name: "Barbara Liskov", Company: "MIT CSAIL", Email: "barbara@mit.example", Status: Lead, Tags: []string{"types"}},
	{Name: "Ken Thompson", Company: "Bell Labs", Email: "ken@bell.example", Status: Active, Tags: []string{"go", "os"}},
	{Name: "Radia Perlman", Company: "Network Works", Email: "radia@net.example", Status: Lead, Tags: []string{"network"}},
	{Name: "Rob Pike", Company: "Plan 9", Email: "rob@plan9.example", Status: Active, Tags: []string{"go"}},
	{Name: "Frances Allen", Company: "IBM Research", Email: "fran@ibm.example", Status: Inactive, Tags: []string{"compiler"}},
	{Name: "Edsger Dijkstra", Company: "Eindhoven Tech", Email: "edsger@tue.example", Status: Lead, Tags: []string{"algorithms"}},
	{Name: "Hedy Lamarr", Company: "Frequency Hopping", Email: "hedy@hop.example", Status: Active, Tags: []string{"radio", "vip"}},
	{Name: "Donald Knuth", Company: "Stanford", Email: "don@stanford.example", Status: Inactive, Tags: []string{"tex", "algorithms"}},
	{Name: "Annie Easley", Company: "NASA Glenn", Email: "annie@glenn.example", Status: Lead, Tags: []string{"space"}},
	{Name: "John McCarthy", Company: "Lisp Lab", Email: "john@lisp.example", Status: Inactive, Tags: []string{"ai"}},
	{Name: "Sophie Wilson", Company: "Acorn", Email: "sophie@acorn.example", Status: Active, Tags: []string{"arm"}},
	{Name: "Robert Griesemer", Company: "Gopher Inc", Email: "gri@gopher.example", Status: Lead, Tags: []string{"go"}},
	{Name: "Joan Clarke", Company: "Hut 8", Email: "joan@hut8.example", Status: Active, Tags: []string{"crypto"}},
	{Name: "Tim Berners-Lee", Company: "CERN", Email: "tim@cern.example", Status: Active, Tags: []string{"web", "vip"}},
	{Name: "Shafi Goldwasser", Company: "Weizmann", Email: "shafi@weizmann.example", Status: Lead, Tags: []string{"crypto"}},
	{Name: "Niklaus Wirth", Company: "ETH Zürich", Email: "niklaus@eth.example", Status: Inactive, Tags: []string{"pascal"}},
	{Name: "Mary Kenneth Keller", Company: "Clarke College", Email: "mary@clarke.example", Status: Lead, Tags: []string{"basic"}},
	{Name: "Guido van Rossum", Company: "Python Works", Email: "guido@python.example", Status: Active, Tags: []string{"python"}},
	{Name: "Adele Goldberg", Company: "Xerox PARC", Email: "adele@parc.example", Status: Active, Tags: []string{"smalltalk"}},
	{Name: "Bjarne Stroustrup", Company: "Bell Labs", Email: "bjarne@bell.example", Status: Lead, Tags: []string{"c++"}},
	{Name: "Lynn Conway", Company: "VLSI Group", Email: "lynn@vlsi.example", Status: Inactive, Tags: []string{"chips"}},
}
