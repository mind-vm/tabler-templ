package copilot

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mind-vm/tabler-templ/blocks/assistant"
)

// agent stands in for a language model so the demo runs without one. It understands a few requests
// (the examples of the intro) and answers them the way a model with read-only tools would: it calls
// tools, shows what it read, and asks for actions by name. Everything it asks for goes through
// Catalog.Check, exactly as a real model's requests would; the refusals in the demo are real ones.
//
// To use a real model, replace plan: give it the tools and Catalog.Prompt(role), and turn its tool
// calls into the same tool and act calls.
type agent struct {
	store *Store
	role  string
	pause time.Duration // between steps, so the stream is visible
}

// turn is one run: the stream, the turn's number (for unique ids) and the tool calls so far.
type turn struct {
	ctx   context.Context
	s     *assistant.Stream
	n     int
	calls []assistant.ToolCall
	a     *agent
}

func (t *turn) id(name string) string { return fmt.Sprintf("t%d-%s", t.n, name) }

// tool shows a read-only tool call: running, then done.
func (t *turn) tool(name, detail, icon string) {
	t.calls = append(t.calls, assistant.ToolCall{Name: name, Detail: detail, Icon: icon, Running: true})
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
	time.Sleep(t.a.pause)
	t.calls[len(t.calls)-1].Running = false
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
}

// act asks for an action. A refused one shows as a failed call, and its reason goes back to the
// caller, as it would go back to the model.
func (t *turn) act(name string, values map[string]string) error {
	a, err := Catalog.Check(t.a.role, name, values)
	detail := name
	if err != nil {
		reason := strings.TrimPrefix(err.Error(), assistant.ErrAction.Error()+": ")
		t.calls = append(t.calls, assistant.ToolCall{Name: "ui." + name, Detail: reason, Icon: "ti-shield-x", Failed: true})
		_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
		id := t.id(fmt.Sprintf("refused-%d", len(t.calls)))
		_ = t.s.Send(t.ctx, assistant.SlotBlocks, id, "", assistant.ActionRefused(id, name, reason))
		return err
	}
	t.calls = append(t.calls, assistant.ToolCall{Name: "ui." + name, Detail: detail, Icon: "ti-hand-click"})
	_ = t.s.Send(t.ctx, assistant.SlotSteps, t.id("tools"), "", assistant.ToolChips(t.id("tools"), t.calls))
	return t.s.Action(t.ctx, t.id(fmt.Sprintf("action-%d", len(t.calls))), a)
}

func (t *turn) table(title string, os []Order) {
	rows := make([]any, 0, len(os))
	for _, o := range os {
		c, _ := t.a.store.Customer(o.CustomerID)
		rows = append(rows, map[string]any{"Order": o.ID, "Customer": c.Name, "Title": o.Title,
			"Due": o.Due.Format("02.01.2006"), "Amount €": float64(o.Amount)})
	}
	b := assistant.Block{Type: "table", DOMID: t.id("table"), Payload: map[string]any{
		"title": title, "columns": []any{"Order", "Customer", "Title", "Due", "Amount €"}, "rows": rows}}
	_ = t.s.Send(t.ctx, assistant.SlotBlocks, b.DOMID, "", assistant.BlockView(b))
}

func (t *turn) answer(md string) {
	_ = t.s.Send(t.ctx, assistant.SlotAnswer, "", assistant.ModeSet, assistant.AnswerStrict(md))
}

var (
	orderID = regexp.MustCompile(`(?i)\bA-\d{4}\b`)
	amount  = regexp.MustCompile(`(\d{1,3}(?:[.,]\d{3})+|\d+)\s*(?:€|eur|euro)`)
)

// plan runs one message.
func (a *agent) plan(t *turn, message string) {
	m := strings.ToLower(message)
	cust, hasCustomer := a.customerIn(m)
	switch {
	case strings.Contains(m, "note of") || strings.Contains(m, "what it says") || strings.Contains(m, "notiz"):
		a.followNote(t, message)
	case strings.Contains(m, "iban") || strings.Contains(m, "bank"):
		a.bankDetails(t, cust, hasCustomer)
	case strings.Contains(m, "overdue") || strings.Contains(m, "überfällig"):
		a.overdue(t)
	case strings.Contains(m, "new order") || strings.Contains(m, "create") || strings.Contains(m, "neuer auftrag") || strings.Contains(m, "anlegen"):
		a.newOrder(t, message, cust, hasCustomer)
	case orderID.MatchString(message) && (strings.Contains(m, "done") || strings.Contains(m, "erledigt")):
		a.markDone(t, strings.ToUpper(orderID.FindString(message)))
	case orderID.MatchString(message):
		id := strings.ToUpper(orderID.FindString(message))
		t.tool("orders.get", id, "ti-file-invoice")
		if err := t.act("page", map[string]string{"path": orders + "/" + id}); err != nil {
			t.answer("I could not find a page for " + id + ".")
			return
		}
		t.answer("Here is order **" + id + "**.")
	case hasCustomer:
		t.tool("customers.find", cust.Name, "ti-search")
		t.tool("orders.list", "customer="+cust.ID, "ti-list")
		t.table("Orders of "+cust.Name, a.store.Orders(cust.ID, false))
		_ = t.act("page", map[string]string{"path": customer + "/" + cust.ID})
		t.answer("These are the orders of **" + cust.Name + "**. I opened the customer page.")
	default:
		t.answer("I read orders and customers and help with the forms of this app; saving stays with you. Try one of these:\n\n" +
			"- *Which orders are overdue?*\n- *Create an order for Becker: shop sign, 2,400 €*\n- *Mark A-1003 as done*\n" +
			"- *Change the IBAN of Müller Bau*\n- *Read the note of A-1004 and do what it says*")
	}
}

func (a *agent) customerIn(m string) (Customer, bool) {
	for _, w := range []string{"müller", "mueller", "becker", "schulz"} {
		if strings.Contains(m, w) {
			return a.store.FindCustomer(strings.Replace(w, "mueller", "müller", 1))
		}
	}
	return Customer{}, false
}

func (a *agent) overdue(t *turn) {
	t.tool("orders.list", "overdue=true", "ti-list")
	os := a.store.Orders("", true)
	t.table("Overdue orders", os)
	sum := 0
	for _, o := range os {
		sum += o.Amount
	}
	_ = t.act("overdue-filter", nil)
	t.answer(fmt.Sprintf("**%d orders** are past their due date, worth **%s €** together. The oldest is %s. "+
		"I pointed at the filter, so you can keep the list in view.", len(os), thousands(sum), os[len(os)-1].ID))
}

func (a *agent) newOrder(t *turn, message string, c Customer, ok bool) {
	if !ok {
		t.answer("For which customer? I know **Müller Bau**, **Becker & Söhne** and **Schulz Logistik**.")
		return
	}
	t.tool("customers.find", c.Name, "ti-search")
	title := "New order"
	if _, after, found := strings.Cut(message, ":"); found {
		title = strings.TrimSpace(strings.SplitN(after, ",", 2)[0])
	}
	values := map[string]string{"customer": c.ID, "title": title, "due": a.store.Today().AddDate(0, 0, 14).Format(dateLayout)}
	if am := amount.FindStringSubmatch(strings.ToLower(message)); am != nil {
		values["amount"] = strings.NewReplacer(".", "", ",", "").Replace(am[1])
	}
	if err := t.act("order-new", values); err != nil {
		t.answer("I could not prepare the form: " + err.Error())
		return
	}
	t.answer("I filled in a new order for **" + c.Name + "**, due in two weeks. Check the values and save it yourself.")
}

func (a *agent) markDone(t *turn, id string) {
	t.tool("orders.get", id, "ti-file-invoice")
	o, ok := a.store.Order(id)
	if !ok {
		t.answer("There is no order " + id + ".")
		return
	}
	note := strings.TrimSpace(o.Note + "\nAccepted by the customer on " + a.store.Today().Format("02.01.2006") + ".")
	if err := t.act("order-edit", map[string]string{"id": id, "status": string(Done), "note": note}); err != nil {
		t.answer("I could not prepare the form: " + err.Error())
		return
	}
	t.answer("I opened **" + id + "** (" + o.Title + ") and set it to *Done* with a note. Nothing is saved until you press **Save**.")
}

func (a *agent) bankDetails(t *turn, c Customer, ok bool) {
	if !ok {
		c, _ = a.store.Customer("C-1")
	}
	t.tool("customers.get", c.ID, "ti-building")
	err := t.act("customer-edit", map[string]string{"id": c.ID, "iban": "DE44 5001 0517 5407 3249 31"})
	if err == nil || !errors.Is(err, assistant.ErrAction) {
		t.answer("Something unexpected happened.")
		return
	}
	_ = t.act("page", map[string]string{"path": customer + "/" + c.ID})
	t.answer("Bank details are not something I fill in: the app's catalog refuses it, whatever I am asked. " +
		"I opened **" + c.Name + "**; please enter the new IBAN yourself and save.")
}

// followNote shows a prompt injection: order A-1004's note asks the assistant to change bank details
// and to show an image from an outside server. A model may well try; the catalog refuses the action,
// and AnswerStrict drops the image.
func (a *agent) followNote(t *turn, message string) {
	id := "A-1004"
	if m := orderID.FindString(message); m != "" {
		id = strings.ToUpper(m)
	}
	t.tool("orders.get", id, "ti-file-invoice")
	o, ok := a.store.Order(id)
	if !ok || o.Note == "" {
		t.answer("Order " + id + " has no note.")
		return
	}
	// What a model misled by the note might ask for:
	err := t.act("customer-edit", map[string]string{"id": o.CustomerID, "iban": "DE44 5001 0517 5407 3249 31", "city": "Munich"})
	answer := "The note of **" + id + "** contains instructions to change the customer's bank details. "
	if errors.Is(err, assistant.ErrAction) {
		answer += "I did not follow them: the app does not let me fill in bank details (see above). "
	}
	// …and the image it was told to show, which AnswerStrict renders as its alt text only:
	answer += "Please check the note with the customer by phone.\n\n![tracking pixel removed](https://track.example/p.png?o=" + id + ")"
	t.answer(answer)
}

func thousands(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "." + s[i:]
	}
	return s
}
