package split

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func text(s string) templ.Component { return templ.Raw(s) }

func TestSplitSelected(t *testing.T) {
	out := render(t, Split(Props{List: text("LIST"), Detail: text("DETAIL"), Selected: true, BackHref: "/tickets",
		Aside: text("FACTS"), AsideTitle: "Details", ListWidth: "20rem"}))
	for _, want := range []string{`class="kit-split kit-split-has-list kit-split-has-aside kit-split-selected"`, "LIST", "DETAIL",
		`id="split-detail"`, `id="split-aside"`, `href="/tickets"`, "Back", "Details", "--kit-split-list:20rem"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
	if strings.Count(out, "<style") != 1 || strings.Count(out, "<script") != 1 {
		t.Errorf("styles or script not exactly once:\n%s", out)
	}
}

func TestSplitEmptyAndClosedAside(t *testing.T) {
	out := render(t, Split(Props{ID: "s", List: text("LIST"), Empty: text("NOTHING"), Aside: text("A"), AsideClosed: true}))
	if !strings.Contains(out, "kit-split-empty") || !strings.Contains(out, "NOTHING") || !strings.Contains(out, "data-aside-closed") {
		t.Fatalf("empty state:\n%s", out)
	}
	if strings.Contains(out, `kit-split-has-aside kit-split-selected"`) || strings.Contains(out, `class="kit-split-back"`) {
		t.Fatalf("no selection, no back link:\n%s", out)
	}
	tg := render(t, AsideToggle(ToggleProps{Target: "s", Closed: true}))
	if !strings.Contains(tg, `data-kit-split-toggle="s"`) || !strings.Contains(tg, `aria-expanded="false"`) || !strings.Contains(tg, `aria-controls="s-aside"`) {
		t.Fatalf("toggle:\n%s", tg)
	}
}

func TestListItem(t *testing.T) {
	out := render(t, ListItem(ListItemProps{Href: "/t/1", Title: "Export stops", Meta: "09:12", Text: "Since Monday", Active: true, Unread: true}))
	for _, want := range []string{`href="/t/1"`, "kit-split-item active unread", `aria-current="page"`, "Export stops", "09:12", "Since Monday"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
}
