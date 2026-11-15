package reference

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/mind-vm/tabler-templ/blocks/assistant"
)

// The AI assistant of B-09 (concept Assistent §6): a tab of the panel that reads, shows what it read and
// fills in forms; it never saves (KA-01, KA-02). An agent with a fixed script stands in for the language
// model, so the reference runs without one; what it asks to do still goes through the Catalog, exactly as
// a model's requests would.

// examples are the questions the panel offers.
var examples = []string{
	"Welche Kunden haben überfällige Posten?",
	"Trage für die Wohnungsgesellschaft Nordstadt die Anschrift Hafenstraße 12 ein",
	"Was steht im Angebot A-2026-001?",
	"Ändere die IBAN der Wohnungsgesellschaft Nordstadt",
}

// catalog is the allowlist: the pages and forms the assistant may use. Bank details and e-mail addresses
// it never fills in (KA-20, assistant.DefaultSensitiveFields).
var catalog = &assistant.Catalog{
	Pages: []string{"reference/kunden", "reference/angebote"},
	Actions: []assistant.ActionSpec{
		{Name: "kunde-stammdaten", Title: "Stammdaten", Kind: assistant.ActionFill, Path: pathCustomers + "/{id}", Form: "kunde-stammdaten",
			Description: "den Dialog Stammdaten eines Kunden öffnen und vorausfüllen", Params: []string{"id"}, Fields: []string{"name", "street", "zip", "city"}},
		{Name: "kunden-ueberfaellig", Title: "Kunden mit überfälligen Posten", Kind: assistant.ActionOpen, Path: pathCustomers + "?posten=" + string(SignalOverdue),
			Description: "die Liste der Kunden mit überfälligen Posten öffnen"},
		{Name: "seite", Title: "Seite öffnen", Kind: assistant.ActionOpen, Path: "{path}", Params: []string{assistant.ParamPath},
			Description: "eine Seite öffnen: " + pathCustomers + "/{nr}, " + pathOffers + "/{nr}"},
	},
}

type turn struct {
	ctx   context.Context
	s     *assistant.Stream
	n     int
	calls []assistant.ToolCall
	pause time.Duration
}

func (t *turn) id(name string) string { return fmt.Sprintf("t%d-%s", t.n, name) }

// tool shows a reading tool call: running, then done.
func (t *turn) tool(name, detail string) {
	t.calls = append(t.calls, assistant.ToolCall{Name: name, Detail: detail, Icon: "ti-database-search", Running: true})
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
	time.Sleep(t.pause)
	t.calls[len(t.calls)-1].Running = false
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
}

// act asks for an action; a refused one shows as a refused card with a German reason.
func (t *turn) act(name, refusedText string, values map[string]string) error {
	a, err := catalog.Check("buero", name, values)
	if err != nil {
		t.calls = append(t.calls, assistant.ToolCall{Name: "ui." + name, Detail: "abgelehnt", Icon: "ti-shield-x", Failed: true})
		_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
		id := t.id(fmt.Sprintf("refused-%d", len(t.calls)))
		_ = t.s.Send(t.ctx, assistant.SlotBlocks, id, "", assistant.ActionRefused(id, name, refusedText))
		return err
	}
	t.calls = append(t.calls, assistant.ToolCall{Name: "ui." + name, Detail: a.Title, Icon: "ti-hand-click"})
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
	return t.s.Action(t.ctx, t.id(fmt.Sprintf("action-%d", len(t.calls))), a)
}

func (t *turn) table(title string, columns []any, rows []any) {
	b := assistant.Block{Type: "table", DOMID: t.id("table"), Payload: map[string]any{"title": title, "columns": columns, "rows": rows}}
	_ = t.s.Send(t.ctx, assistant.SlotBlocks, b.DOMID, "", assistant.BlockView(b))
}

func (t *turn) answer(md string) {
	_ = t.s.Send(t.ctx, assistant.SlotAnswer, "", assistant.ModeSet, assistant.AnswerStrict(md))
}

var (
	offerNo  = regexp.MustCompile(`(?i)\bA-\d{4}-\d{3}\b`)
	streetRE = regexp.MustCompile(`([A-ZÄÖÜ][\p{L}-]*(?:straße|strasse|weg|platz|allee|ring|gasse)\s*\d+[a-z]?)\b`)
)

// plan answers one message.
func (a *App) plan(t *turn, message string) {
	m := fold(message)
	c, hasCustomer := a.customerIn(m)
	switch {
	case strings.Contains(m, "iban") || strings.Contains(m, "bank"):
		if !hasCustomer {
			c, _ = a.store.Customer("K-1001")
		}
		t.tool("kunden.lesen", c.No)
		_ = t.act("kunde-stammdaten", "Bankverbindungen fülle ich nie aus. Sie tragen sie selbst im Formular ein.", map[string]string{"id": c.No, "iban": "DE00 0000 0000 0000 0000 00"})
		t.answer("Eine Bankverbindung ändere oder fülle ich nicht aus. Die IBAN von **" + c.Name + "** tragen Sie selbst ein; ich öffne dafür keine Formulare.")
	case strings.Contains(m, "ueberfaellig") || strings.Contains(m, "offene posten"):
		t.tool("kunden.liste", "offene Posten: überfällig")
		var rows []any
		for _, k := range a.store.Customers() {
			if k.Signal == SignalOverdue {
				rows = append(rows, map[string]any{"Kunde": k.Name, "Ort": k.City, "Überfällig": Euro(k.Overdue)})
			}
		}
		t.table("Kunden mit überfälligen Posten", []any{"Kunde", "Ort", "Überfällig"}, rows)
		_ = t.act("kunden-ueberfaellig", "", nil)
		t.answer(fmt.Sprintf("%d Kunden haben überfällige Posten. Ich habe die gefilterte Liste geöffnet.", len(rows)))
	case offerNo.MatchString(message):
		no := strings.ToUpper(offerNo.FindString(message))
		t.tool("angebote.lesen", no)
		o, ok := a.store.Offer(no)
		if !ok {
			t.answer("Ein Angebot " + no + " finde ich nicht.")
			return
		}
		var rows []any
		for _, p := range o.Positions {
			rows = append(rows, map[string]any{"Pos.-Nr.": p.No, "Bezeichnung": p.Text, "Betrag": Euro(p.Amount())})
		}
		t.table("Positionen von "+o.No, []any{"Pos.-Nr.", "Bezeichnung", "Betrag"}, rows)
		_ = t.act("seite", "", map[string]string{"path": offerURL(o.No)})
		t.answer(fmt.Sprintf("**%s** hat %d Positionen über **%s**, Stand: %s. Ich habe das Angebot geöffnet.", o.Title, len(o.Positions), Euro(o.Total()), o.Status.Label()))
	case (strings.Contains(m, "anschrift") || strings.Contains(m, "adresse") || strings.Contains(m, "strasse")) && hasCustomer:
		street := streetRE.FindString(message)
		t.tool("kunden.lesen", c.No)
		if street == "" {
			t.answer("Welche Anschrift soll ich eintragen? Nennen Sie Straße und Hausnummer.")
			return
		}
		if err := t.act("kunde-stammdaten", "", map[string]string{"id": c.No, "street": street}); err != nil {
			t.answer("Das Formular kann ich nicht öffnen.")
			return
		}
		t.answer("Ich habe den Dialog **Stammdaten** von " + c.Name + " geöffnet und „" + street + "“ eingetragen. Bitte prüfen Sie die Werte und speichern Sie selbst.")
	default:
		t.answer("Ich lese Kunden und Angebote und fülle Formulare vor; speichern tun Sie. Probieren Sie:\n\n- *" + strings.Join(examples, "*\n- *") + "*")
	}
}

func (a *App) customerIn(folded string) (Customer, bool) {
	for _, c := range a.store.Customers() {
		if strings.Contains(folded, fold(c.Name)) || strings.Contains(folded, strings.ToLower(c.No)) {
			return c, true
		}
	}
	if strings.Contains(folded, "nordstadt") {
		return a.store.Customer("K-1001")
	}
	return Customer{}, false
}

// pause is the agent's pause between steps, so the stream is visible; tests set it to 0.
var pause = 300 * time.Millisecond

// run answers one message as a stream of fragments and actions, recorded in the session's transcript.
func (a *App) run(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&in); err != nil || strings.TrimSpace(in.Message) == "" || len(in.Message) > 4000 {
		http.Error(w, "Nachricht fehlt oder ist zu lang.", http.StatusBadRequest)
		return
	}
	sess := a.session(w, r)
	a.mu.Lock()
	if sess.busy {
		a.mu.Unlock()
		http.Error(w, "Der Assistent antwortet noch.", http.StatusConflict)
		return
	}
	sess.busy = true
	sess.turns++
	n := sess.turns
	a.mu.Unlock()
	defer func() { a.mu.Lock(); sess.busy = false; a.mu.Unlock() }()

	s := assistant.Open(w)
	if s == nil {
		http.Error(w, "Streaming geht hier nicht.", http.StatusInternalServerError)
		return
	}
	s.Record(sess.transcript)
	ctx := germanContext(r.Context())
	_ = s.Turn(ctx, in.Message)
	t := &turn{ctx: ctx, s: s, n: n, pause: pause}
	start := time.Now()
	a.plan(t, in.Message)
	_ = s.Send(ctx, assistant.SlotStatus, "", assistant.ModeSet, assistant.Usage(assistant.UsageProps{
		Text: time.Since(start).Round(100*time.Millisecond).String() + " · festes Skript, kein Sprachmodell"}))
}
