package assistant

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

var testCatalog = &Catalog{
	Pages: []string{"orders", "customers"},
	Actions: []ActionSpec{
		{Name: "order-new", Title: "New order", Kind: ActionFill, Path: "/orders", Form: "order-new",
			Fields: []string{"customer", "title", "amount"}},
		{Name: "order-edit", Title: "Edit order", Kind: ActionFill, Path: "/orders/{id}", Form: "order-edit",
			Params: []string{"id"}, Fields: []string{"status", "note"}, Roles: []string{"office"}},
		{Name: "page", Title: "Open page", Kind: ActionOpen, Path: "{path}", Params: []string{"path"}},
		{Name: "overdue-filter", Title: "Overdue filter", Kind: ActionShow, Path: "/orders", Anchor: "orders.filter.overdue"},
	},
}

func TestCatalogIsValid(t *testing.T) {
	if p := testCatalog.Problems(); len(p) > 0 {
		t.Fatal(p)
	}
}

func TestCheckFill(t *testing.T) {
	a, err := testCatalog.Check("office", "order-edit", map[string]string{"id": "A-1003", "status": "done", "note": " "})
	if err != nil {
		t.Fatal(err)
	}
	if a.Path != "/orders/A-1003" || a.Form != "order-edit" || a.Kind != ActionFill {
		t.Errorf("unexpected action %+v", a)
	}
	if len(a.Fields) != 1 || a.Fields["status"] != "done" {
		t.Errorf("empty values must be dropped, ids not filled as fields: %v", a.Fields)
	}
}

func TestCheckRefuses(t *testing.T) {
	for _, c := range []struct {
		role, name string
		values     map[string]string
		want       string
	}{
		{"office", "order-delete", nil, "no action"},
		{"field", "order-edit", map[string]string{"id": "A-1"}, "not available to the role"},
		{"office", "order-edit", map[string]string{"id": "../admin"}, "not a valid ID"},
		{"office", "order-new", map[string]string{"price": "1"}, "no field"},
		{"office", "order-new", map[string]string{"customerEmail": "x"}, "does not fill in"},
		{"office", "order-new", map[string]string{"iban": "DE02120300000000202051"}, "does not fill in"},
		{"office", "order-new", map[string]string{"title": "pay to DE02 1203 0000 0000 2020 51"}, "looks like an IBAN"},
		{"office", "order-new", map[string]string{"title": "write to boss@example.com"}, "looks like an IBAN or an e-mail"},
		{"office", "order-new", map[string]string{"title": strings.Repeat("x", 501)}, "longer than"},
		{"office", "page", map[string]string{"path": "//evil.example"}, "not a page"},
		{"office", "page", map[string]string{"path": "/admin"}, "not a page"},
		{"office", "page", map[string]string{"path": "/orders/A-1/delete"}, "not a page"},
		{"office", "page", map[string]string{"path": "https://evil.example/orders"}, "not a page"},
	} {
		_, err := testCatalog.Check(c.role, c.name, c.values)
		if err == nil || !errors.Is(err, ErrAction) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %v: got %v, want %q", c.name, c.values, err, c.want)
		}
	}
}

func TestCheckPage(t *testing.T) {
	for _, p := range []string{"/orders", "/customers/K-7"} {
		if a, err := testCatalog.Check("", "page", map[string]string{"path": p}); err != nil || a.Path != p {
			t.Errorf("%s: %+v %v", p, a, err)
		}
	}
}

func TestNamesAndPrompt(t *testing.T) {
	if got := strings.Join(testCatalog.Names("field"), ","); got != "order-new,overdue-filter,page" {
		t.Errorf("names for role field: %s", got)
	}
	p := testCatalog.Prompt("office")
	if !strings.Contains(p, "- order-edit:") || !strings.Contains(p, "(fields: id, status, note)") {
		t.Errorf("prompt: %s", p)
	}
}

func TestProblems(t *testing.T) {
	bad := &Catalog{Pages: []string{"Orders"}, Actions: []ActionSpec{
		{Name: "x", Title: "X", Kind: ActionFill, Path: "/x", Form: "f", Fields: []string{"iban"}},
		{Name: "x", Kind: "delete", Path: "x/{id}"},
		{Name: "y", Title: "Y", Kind: ActionShow, Path: "/y", Anchor: "Bad Anchor"},
	}}
	got := bad.Problems()
	for _, want := range []string{`page "Orders"`, "sensitive", "twice", "no title", "unknown kind", "does not start with /", "data-tour"} {
		if !strings.Contains(errors.Join(got...).Error(), want) {
			t.Errorf("missing problem %q in %v", want, got)
		}
	}
}

func TestRequireRoles(t *testing.T) {
	c := &Catalog{RequireRoles: true, Pages: testCatalog.Pages, Actions: testCatalog.Actions}
	if got := errors.Join(c.Problems()...); got == nil || !strings.Contains(got.Error(), `action "order-new": no role`) {
		t.Errorf("an action without roles must be a problem: %v", got)
	}
	if _, err := c.Check("office", "order-new", map[string]string{"title": "x"}); !errors.Is(err, ErrAction) {
		t.Errorf("an action without roles must be refused, got %v", err)
	}
	if _, err := c.Check("office", "order-edit", map[string]string{"id": "A-1"}); err != nil {
		t.Errorf("an action of the role: %v", err)
	}
	if got := strings.Join(c.Names("office"), ","); got != "order-edit" {
		t.Errorf("names: %s", got)
	}
	if p := c.Prompt("office"); strings.Contains(p, "order-new") {
		t.Errorf("prompt offers an action without roles: %s", p)
	}
}

func TestActionOptions(t *testing.T) {
	for _, c := range []struct{ attr, want string }{
		{"", `"formAttribute":"data-form"`},
		{"data-formular", `"formAttribute":"data-formular"`},
		{`x"]`, `"formAttribute":"data-form"`}, // not an attribute name: the default
	} {
		var b strings.Builder
		ctx := WithActionOptions(context.Background(), ActionOptions{FormAttribute: c.attr})
		if err := ActionScript().Render(ctx, &b); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(b.String(), c.want) {
			t.Errorf("%q: missing %s", c.attr, c.want)
		}
	}
}

func TestStreamAction(t *testing.T) {
	rec := httptest.NewRecorder()
	s := Open(rec)
	a, err := testCatalog.Check("", "order-new", map[string]string{"customer": "K-1", "title": "Relaunch <b>"})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Action(context.Background(), "action-1", a); err != nil {
		t.Fatal(err)
	}
	out := rec.Body.String()
	for _, want := range []string{
		`"type":"html"`, `"slot":"blocks"`, `"id":"action-1"`, // the card
		`data: {"args":{"id":"action-1","name":"order-new","kind":"fill","title":"New order","path":"/orders","form":"order-new","fields":{"customer":"K-1","title":"Relaunch \u003cb\u003e"}},"type":"action"}`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
	if strings.Contains(out, "<b>") {
		t.Errorf("a value reached the page unescaped:\n%s", out)
	}
}

func TestActionCard(t *testing.T) {
	a := Action{ID: "action-2", Kind: ActionFill, Title: "New order", Path: "/orders", Form: "order-new",
		Fields: map[string]string{"title": "Relaunch", "amount": "4800"}}
	var b strings.Builder
	if err := ActionCard(a).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{`id="action-2"`, "Fill in", "New order", "<dt>amount</dt><dd>4800</dd><dt>title</dt>", "data-assistant-action=", "data-action-note"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
}

func TestAnswerStrict(t *testing.T) {
	var b strings.Builder
	if err := AnswerStrict("see [here](https://evil.example) ![x](https://t.example/p.png)").Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if out := b.String(); strings.Contains(out, "<a") || strings.Contains(out, "<img") {
		t.Errorf("strict answer has a link or image: %s", out)
	}
}

func TestTranscript(t *testing.T) {
	tr := &Transcript{}
	s := Open(httptest.NewRecorder()).Record(tr)
	ctx := context.Background()
	_ = s.Turn(ctx, "first <question>")
	_ = s.SendHTML(SlotBlocks, "b1", "", "<p>one</p>")
	_ = s.SendHTML(SlotBlocks, "b1", "", "<p>one revised</p>") // same turn: in place
	_ = s.SendHTML(SlotStatus, "", ModeSet, "<i>working</i>")
	_ = s.SendHTML(SlotStatus, "", ModeSet, "<i>done</i>")
	_ = s.Turn(ctx, "second")
	_ = s.SendHTML(SlotBlocks, "b1", "", "<p>moved</p>") // earlier turn: moves here
	_ = s.SendHTML(SlotAnswer, "", ModeSet, "<p>answer</p>")
	var b strings.Builder
	if err := tr.View().Render(ctx, &b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if tr.Len() != 2 || strings.Count(out, `class="assistant-turn`) != 2 {
		t.Fatalf("want two turns: %s", out)
	}
	for _, want := range []string{"first &lt;question&gt;", "<i>done</i>", "<p>moved</p>", "<p>answer</p>"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
	for _, gone := range []string{"<p>one</p>", "one revised", "working"} {
		if strings.Contains(out, gone) {
			t.Errorf("%q should be replaced: %s", gone, out)
		}
	}
	tr.Max = 1
	_ = s.Turn(ctx, "third")
	if tr.Len() != 1 {
		t.Errorf("Max not applied: %d", tr.Len())
	}
}
