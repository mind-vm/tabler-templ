package reference

import (
	"errors"
	"strings"
	"time"
)

// Requests from customers, each with a proposal that an AI prepared (B-10, GS-34, GS-44; concept Assistent
// KA-01 to KA-04 in spirit): the AI drafts, the person decides and sends. Nothing leaves without a person.

type RequestStatus string

const (
	RequestOpen      RequestStatus = "offen"
	RequestAnswered  RequestStatus = "beantwortet"
	RequestDiscarded RequestStatus = "verworfen"
)

var RequestStatuses = []RequestStatus{RequestOpen, RequestAnswered, RequestDiscarded}

func (s RequestStatus) Label() string {
	switch s {
	case RequestAnswered:
		return "Beantwortet"
	case RequestDiscarded:
		return "Verworfen"
	}
	return "Offen"
}

// Event is an entry of a request's history: a message from outside, an internal note or a decision.
type Event struct {
	At       time.Time
	Who      string
	Text     string
	External bool // from the customer (GS-76); else internal, never sent
}

// Proposal is what the AI prepared: the reply, how sure it is, what it read and why.
type Proposal struct {
	Reply      string
	Confidence float64
	Sources    []Source
	Reasons    []string
}

type Source struct {
	Label, Href string
}

type Request struct {
	No       string // R-2026-001
	Customer string
	Subject  string
	Received time.Time
	Due      time.Time // to be answered by
	Status   RequestStatus
	History  []Event
	Proposal Proposal
}

// DiscardReasons are the reasons to choose from when a proposal is discarded (GS-44: a required choice).
var DiscardReasons = []string{"Inhalt falsch", "Ton unpassend", "Erst intern klären", "Kunde bekommt einen Anruf"}

func (s *Store) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Request, len(s.requests))
	for i, r := range s.requests {
		out[i] = *r
	}
	return out
}

func (s *Store) Request(no string) (Request, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.requests {
		if r.No == no {
			return *r, true
		}
	}
	return Request{}, false
}

func (s *Store) request(no string) *Request {
	for _, r := range s.requests {
		if r.No == no {
			return r
		}
	}
	return nil
}

var errDecided = errors.New("Über diese Anfrage ist schon entschieden.")

// Send sends the reply the person checked and maybe changed.
func (s *Store) Send(no, reply string) (Request, error) {
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return Request{}, FieldErrors{"reply": "Die Antwort ist leer."}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.request(no)
	if r == nil {
		return Request{}, ErrNotFound
	}
	if r.Status != RequestOpen {
		return *r, errDecided
	}
	r.Status = RequestAnswered
	changed := ""
	if reply != r.Proposal.Reply {
		changed = " (geändert)"
	}
	r.History = append(r.History, Event{At: Now, Who: "Sie", Text: "Antwort gesendet" + changed + ":\n" + reply})
	return *r, nil
}

// Discard drops the proposal; it needs a reason from the list (GS-44).
func (s *Store) Discard(no, reason, note string) (Request, error) {
	ok := false
	for _, x := range DiscardReasons {
		ok = ok || x == reason
	}
	if !ok {
		return Request{}, FieldErrors{"reason": "Bitte wählen Sie einen Grund."}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.request(no)
	if r == nil {
		return Request{}, ErrNotFound
	}
	if r.Status != RequestOpen {
		return *r, errDecided
	}
	r.Status = RequestDiscarded
	text := "Vorschlag verworfen: " + reason
	if n := strings.TrimSpace(note); n != "" {
		text += ". " + n
	}
	r.History = append(r.History, Event{At: Now, Who: "Sie", Text: text})
	return *r, nil
}

func at(d, h, m int) time.Time { return time.Date(2026, 11, d, h, m, 0, 0, time.UTC) }

func (s *Store) seedRequests() {
	s.requests = []*Request{
		{No: "R-2026-001", Customer: "K-1001", Subject: "Drei zusätzliche Termine im Dezember", Received: at(14, 16, 2), Due: at(15, 10, 0), Status: RequestOpen,
			History: []Event{
				{At: at(14, 16, 2), Who: "Erika Beispiel", External: true, Text: "Guten Tag,\nwir brauchen im Dezember drei zusätzliche Einsätze wie im Winterangebot, am 4., 11. und 18.12. Können Sie das bis Ende der Woche bestätigen?\nViele Grüße\nErika Beispiel"},
				{At: at(14, 16, 40), Who: "Max Mustermann", Text: "Die Termine am 4. und 11. sind frei, der 18. noch offen."},
			},
			Proposal: Proposal{
				Reply:      "Guten Tag Frau Beispiel,\n\nvielen Dank für Ihre Anfrage. Die Einsätze am 4. und 11.12. bestätigen wir gern zu den Preisen aus unserem Angebot A-2026-002. Für den 18.12. melden wir uns bis Mittwoch.\n\nMit freundlichen Grüßen\nMusterdienst GmbH",
				Confidence: 0.82,
				Sources:    []Source{{"Wohnungsgesellschaft Nordstadt", customerURL("K-1001")}, {"A-2026-002 · Zusatzleistungen Winter", offerURL("A-2026-002")}},
				Reasons:    []string{"Das Winterangebot A-2026-002 enthält dieselbe Leistung je Einsatz.", "Die interne Notiz nennt den 4. und 11.12. als frei; den 18.12. sagt der Vorschlag nicht zu."},
			}},
		{No: "R-2026-002", Customer: "K-1008", Subject: "Frage zur Rechnung November", Received: at(15, 7, 45), Due: at(16, 12, 0), Status: RequestOpen,
			History: []Event{{At: at(15, 7, 45), Who: "Kundin", External: true, Text: "Hallo, auf der Rechnung steht ein Zuschlag Wochenende. Wir hatten keinen Einsatz am Wochenende. Bitte prüfen."}},
			Proposal: Proposal{
				Reply:      "Guten Tag,\n\nwir prüfen den Zuschlag und melden uns bis morgen Mittag.\n\nMit freundlichen Grüßen\nMusterdienst GmbH",
				Confidence: 0.46,
				Sources:    []Source{{"Kunde K-1008", customerURL("K-1008")}},
				Reasons:    []string{"Zu Einsätzen am Wochenende liegen keine Daten vor; der Vorschlag sagt deshalb nur eine Prüfung zu."},
			}},
		{No: "R-2026-003", Customer: "K-1015", Subject: "Bitte um neue Ansprechpartnerin", Received: at(12, 9, 10), Due: at(13, 12, 0), Status: RequestAnswered,
			History: []Event{
				{At: at(12, 9, 10), Who: "Kunde", External: true, Text: "Ab Dezember ist Frau Probe unsere Ansprechpartnerin."},
				{At: at(12, 11, 0), Who: "Sie", Text: "Antwort gesendet:\nVielen Dank, wir haben Frau Probe eingetragen."},
			},
			Proposal: Proposal{Reply: "Vielen Dank, wir haben Frau Probe eingetragen.", Confidence: 0.9}},
	}
}
