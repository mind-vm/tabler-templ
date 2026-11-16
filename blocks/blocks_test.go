package blocks_test

import (
	"context"
	"io"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/alert"
	"github.com/mind-vm/tabler-templ/blocks/create"
	"github.com/mind-vm/tabler-templ/blocks/detail"
	"github.com/mind-vm/tabler-templ/blocks/history"
	"github.com/mind-vm/tabler-templ/blocks/landing"
	"github.com/mind-vm/tabler-templ/blocks/listpage"
	"github.com/mind-vm/tabler-templ/blocks/metric"
	"github.com/mind-vm/tabler-templ/blocks/shell"
	"github.com/mind-vm/tabler-templ/blocks/signin"
	"github.com/mind-vm/tabler-templ/blocks/stats"
	"github.com/mind-vm/tabler-templ/chart/svg"
	"github.com/mind-vm/tabler-templ/kit"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(ctx, &b))
	return b.String()
}

// markup is render without the blocks' inline <style> elements, for assertions on the markup alone.
func markup(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	return styleRE.ReplaceAllString(render(t, ctx, c), "")
}

var styleRE = regexp.MustCompile(`(?s)<style[^>]*>.*?</style>`)

func mit(c templ.Component, kinder string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, templ.Raw(kinder)), w)
	})
}

var status = listpage.Filter{Name: "status", Label: "All states", Options: []listpage.Option{{Value: "open", Label: "Open"}, {Value: "done", Label: "Done"}}}

func TestListFromRequest(t *testing.T) {
	l := listpage.New(httptest.NewRequest("GET", "/customers?q=+berlin+&status=done&page=2", nil), listpage.Config{}, status)
	assert.Equal(t, "berlin", l.Query)
	assert.Equal(t, "done", l.Value("status"))
	assert.True(t, l.Active())
	assert.Equal(t, "/customers?page=3&q=berlin&status=done", l.Link(3))
	assert.Equal(t, "/customers?q=berlin&status=done", l.Link(1), "page 1 is the plain URL")

	// Unknown filter values count as "all"; parameter names are configurable.
	l = listpage.New(httptest.NewRequest("GET", "/x?status=hacked&seite=2", nil), listpage.Config{PageParam: "seite", PageSize: 2}, status)
	assert.Equal(t, "", l.Value("status"))
	assert.False(t, l.Active())
	assert.Equal(t, "/x?seite=2", l.Link(2))
	assert.Equal(t, "nothing here", l.Empty("nothing here", "no results"))
}

func TestPaginate(t *testing.T) {
	rows := []int{1, 2, 3, 4, 5}
	l := listpage.New(httptest.NewRequest("GET", "/x?page=2", nil), listpage.Config{PageSize: 2})
	assert.Equal(t, []int{3, 4}, listpage.Paginate(l, rows, nil))
	assert.Equal(t, 5, l.Total)
	assert.Equal(t, 3, l.Pages)
	assert.Equal(t, [2]int{3, 4}, [2]int{l.From, l.To})

	// A page beyond the last shows the last; keep filters.
	l = listpage.New(httptest.NewRequest("GET", "/x?page=99", nil), listpage.Config{PageSize: 2})
	assert.Equal(t, []int{5}, listpage.Paginate(l, rows, func(n int) bool { return n%2 == 1 }))
	assert.Equal(t, 2, l.Page)

	l = listpage.New(httptest.NewRequest("GET", "/x", nil), listpage.Config{})
	assert.Empty(t, listpage.Paginate(l, rows, func(int) bool { return false }))
	assert.Equal(t, [3]int{0, 0, 1}, [3]int{l.From, l.To, l.Pages})
}

// GS-51: the database does search, filter, page and count; the list says what to fetch and takes the total.
func TestGS51_ListPagedByTheDatabase(t *testing.T) {
	l := listpage.New(httptest.NewRequest("GET", "/x?q=M%C3%BChl+Weg&page=3", nil), listpage.Config{PageSize: 10})
	assert.Equal(t, []string{"muehl", "weg"}, l.Words())
	off, lim := l.Window()
	assert.Equal(t, [2]int{20, 10}, [2]int{off, lim})
	assert.True(t, l.SetTotal(45))
	assert.Equal(t, [4]int{3, 5, 21, 30}, [4]int{l.Page, l.Pages, l.From, l.To})

	// Beyond the last page: the last one, and fetch again.
	l = listpage.New(httptest.NewRequest("GET", "/x?page=9", nil), listpage.Config{PageSize: 10})
	assert.False(t, l.SetTotal(12))
	off, _ = l.Window()
	assert.Equal(t, [4]int{2, 10, 11, 12}, [4]int{l.Page, off, l.From, l.To})

	l = listpage.New(httptest.NewRequest("GET", "/x", nil), listpage.Config{})
	assert.True(t, l.SetTotal(0))
	assert.Empty(t, l.Words())
	assert.Equal(t, [3]int{0, 0, 1}, [3]int{l.From, l.To, l.Pages})
}

// The flow of Window's example with count(*) OVER (): page 9 brings no row and a total of 0; the second
// fetch brings page 1 and the real total.
func TestGS51_PageBeyondTheLastWithCountOver(t *testing.T) {
	all := 12
	fetch := func(offset, limit int) (rows, total int) { // count(*) OVER (): 0 when no row comes back
		rows = max(0, min(limit, all-offset))
		if rows == 0 {
			return 0, 0
		}
		return rows, all
	}
	l := listpage.New(httptest.NewRequest("GET", "/x?page=9", nil), listpage.Config{PageSize: 10})
	off, lim := l.Window()
	rows, total := fetch(off, lim)
	require.False(t, l.SetTotal(total))
	off, lim = l.Window()
	rows, total = fetch(off, lim)
	assert.True(t, l.SetTotal(total))
	assert.Equal(t, 10, rows)
	assert.Equal(t, [5]int{1, 2, 12, 1, 10}, [5]int{l.Page, l.Pages, l.Total, l.From, l.To})
}

// FoldSQL folds like Fold: one replace per letter of the same table, inside lower().
func TestGS50_FoldSQLMatchesFold(t *testing.T) {
	sql := listpage.FoldSQL("c.name")
	assert.Equal(t, `replace(replace(replace(replace(lower(c.name), 'ä', 'ae'), 'ö', 'oe'), 'ü', 'ue'), 'ß', 'ss')`, sql)
	// Simulate the expression in Go: lower, then the replaces in their order, must equal Fold.
	re := regexp.MustCompile(`'(.)', '([a-z]+)'`)
	for _, in := range []string{"Mühlenweg", "GRÜẞE Straße", "Öl und Äpfel", "plain"} {
		got := strings.ToLower(in)
		for _, m := range re.FindAllStringSubmatch(sql, -1) {
			got = strings.ReplaceAll(got, m[1], m[2])
		}
		assert.Equal(t, listpage.Fold(in), got, in)
	}
}

func TestLikePattern(t *testing.T) {
	assert.Equal(t, `%muehl%`, listpage.LikePattern("muehl"))
	assert.Equal(t, `%50\%\_a\\b%`, listpage.LikePattern(`50%_a\b`))
	l := listpage.New(httptest.NewRequest("GET", "/x?q=M%C3%BChl+10%25", nil), listpage.Config{})
	assert.Equal(t, []string{"%muehl%", `%10\%%`}, l.Patterns())
}

func TestContains(t *testing.T) {
	assert.True(t, listpage.Contains("BERLIN mitte", "Mitte", "Berlin"))
	assert.False(t, listpage.Contains("berlin süd", "Berlin Mitte"))
	assert.True(t, listpage.Contains("", "x"))
}

// GS-50: the search ignores case and umlauts (ä = ae), both ways.
func TestGS50_ContainsFoldsUmlauts(t *testing.T) {
	assert.True(t, listpage.Contains("muehl", "Mühlhausen"))
	assert.True(t, listpage.Contains("Mühl", "Muehlenweg"))
	assert.True(t, listpage.Contains("STRASSE", "Hauptstraße"))
	assert.True(t, listpage.Contains("goethe", "Göthe"))
	assert.False(t, listpage.Contains("mühl nord", "Mühlhausen"))
	assert.Equal(t, "muellerstrasse", listpage.Fold("Müllerstraße"))
}

func TestListCard(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	l := listpage.New(httptest.NewRequest("GET", "/kunden?q=a", nil), listpage.Config{Placeholder: "Name …"}, status)
	listpage.Paginate(l, make([]int, 30), nil)
	h := render(t, ctx, mit(listpage.Card(l), "<table></table>"))
	assert.Contains(t, h, `<form method="get" action="/kunden" class="kit-list-search" role="search">`)
	assert.Contains(t, h, `value="a" aria-label="Suchen" name="q" placeholder="Name …"`)
	assert.Contains(t, h, `aria-label="Suchen"`)
	assert.Contains(t, h, `>Zurücksetzen</a>`, "reset while the search is active")
	assert.Contains(t, h, "<table></table>")
	assert.Contains(t, h, "1–25 von 30")
	assert.Contains(t, h, `href="/kunden?page=2&amp;q=a"`)

	// One page: no page numbers; empty list: no footer.
	l = listpage.New(httptest.NewRequest("GET", "/kunden", nil), listpage.Config{})
	listpage.Paginate(l, make([]int, 3), nil)
	h = render(t, ctx, listpage.Footer(l))
	assert.Contains(t, h, "1–3 von 3")
	assert.NotContains(t, h, "pagination")
	listpage.Paginate(l, []int{}, nil)
	assert.Empty(t, markup(t, ctx, listpage.Footer(l)))
}

func TestShellPage(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	p := shell.Props{Lang: "de", Title: "Kunden", AppName: "Acme",
		Head:        templ.Raw(`<link rel="stylesheet" href="/tabler.css">`),
		Scripts:     templ.Raw(`<script src="/tabler.js"></script>`),
		Sidebar:     templ.Raw(`<aside id="sidebar" class="kit-sidebar"></aside>`),
		Breadcrumbs: []shell.Crumb{{Text: "Ablauf"}, {Text: "Kunden", Href: "/kunden"}, {Text: "Stadt"}},
		Subtitle:    templ.Raw("6 Kunden"),
		Actions:     templ.Raw(`<a class="btn">Neu</a>`),
		Flash:       &shell.Flash{Variant: alert.VariantSuccess, Text: "Gespeichert"},
	}
	h := render(t, ctx, mit(shell.Page(p), "<p>Inhalt</p>"))
	assert.True(t, strings.HasPrefix(h, "<!doctype html><html lang=\"de\">"), h[:40])
	assert.Contains(t, h, "<title>Kunden · Acme</title>")
	assert.Contains(t, h, `</style><link rel="stylesheet" href="/tabler.css">`, "the kit's CSS first, so the app's theme overrides it")
	assert.Contains(t, h, `data-bs-target="#sidebar"`, "sidebar buttons in the header")
	assert.Contains(t, h, `aria-label="Brotkrumen"`)
	assert.Contains(t, h, `<li class="breadcrumb-item">Ablauf</li><li class="breadcrumb-item"><a href="/kunden">Kunden</a></li><li class="breadcrumb-item active" aria-current="page">Stadt</li>`)
	assert.Contains(t, h, `<h2 class="page-title">Kunden</h2>`)
	assert.Contains(t, h, `<div class="kit-shell-subtitle">6 Kunden</div>`)
	assert.Contains(t, h, `<div class="btn-list"><a class="btn">Neu</a></div>`)
	assert.Contains(t, h, "alert-success")
	assert.Contains(t, h, "<p>Inhalt</p>")
	assert.Contains(t, h, `<script src="/tabler.js"></script><script>`, "the sidebar script after the app's")

	// Without sidebar (sign-in): the brand instead of the buttons, no sidebar script.
	h = render(t, ctx, shell.Page(shell.Props{Title: "Anmelden", Brand: templ.Raw(`<a class="brand">Acme</a>`)}))
	assert.Contains(t, h, `<a class="brand">Acme</a>`)
	assert.NotContains(t, h, "data-kit-sidebar-toggle")
	assert.NotContains(t, h, "kit_sidebar=")
}

// A link to the content comes first, and a lasting notice stands above the message (GS-63).
func TestShellSkipLinkAndNotice(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, shell.Page(shell.Props{Title: "x", Notice: shell.Notice(shell.NoticeProps{Text: "Demo-Instanz", Action: "Mehr", ActionHref: "/demo"}),
		Flash: &shell.Flash{Text: "Gespeichert."}}))
	body := h[strings.Index(h, "<body"):]
	assert.Regexp(t, `^<body[^>]*><a class="kit-skip-link" href="#kit-content">Zum Inhalt</a>`, body)
	assert.Contains(t, h, `id="kit-content"`)
	assert.Less(t, strings.Index(h, "Demo-Instanz"), strings.Index(h, "Gespeichert."))
	assert.Contains(t, h, "alert-info")
	assert.Contains(t, h, `class="alert-action">Mehr</a>`)
	assert.NotContains(t, h[strings.Index(h, "Demo-Instanz")-400:strings.Index(h, "Demo-Instanz")], "alert-dismissible")
}

// GS-10: the brand stands in the header on small screens, also with a sidebar (hidden there).
func TestShellBrandWithSidebarOnPhones(t *testing.T) {
	h := markup(t, context.Background(), shell.Page(shell.Props{Title: "x", Sidebar: templ.Raw(`<aside id="sidebar"></aside>`), Brand: templ.Raw(`<a class="brand">Acme</a>`)}))
	assert.Contains(t, h, `<div class="kit-shell-brand d-lg-none"><a class="brand">Acme</a></div>`)
	h = markup(t, context.Background(), shell.Page(shell.Props{Title: "x", Sidebar: templ.Raw(`<aside id="sidebar"></aside>`)}))
	assert.NotContains(t, h, "kit-shell-brand")
}

// GS-60: the message after an action offers the next step as a link.
func TestGS60_ShellFlashWithNextStep(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, shell.Page(shell.Props{Title: "x", Flash: &shell.Flash{Variant: alert.VariantSuccess, Text: "Einsatz disponiert.", Next: "Zur Plantafel", NextHref: "/plantafel"}}))
	assert.Regexp(t, `Einsatz disponiert\.\s*<a href="/plantafel" class="alert-action">Zur Plantafel</a>`, h)
	h = render(t, ctx, shell.Page(shell.Props{Title: "x", Flash: &shell.Flash{Text: "Gespeichert.", Next: "Weiter"}}))
	assert.NotContains(t, h, "alert-action", "no link without an address")
}

func TestDetailSection(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	edit := detail.EditButton(detail.EditButtonProps{Dialog: "dialog-kunde"})
	h := render(t, ctx, mit(detail.Section(detail.SectionProps{Title: "Kunde", Extra: "3", Action: edit, Class: "h-100"}), "<p>x</p>"))
	assert.Contains(t, h, `<section class="card kit-section h-100">`)
	assert.Contains(t, h, `<h3 class="kit-section-title">Kunde <small>3</small></h3>`)
	assert.Contains(t, h, `data-bs-target="#dialog-kunde"`)
	assert.Contains(t, h, `title="Bearbeiten"`)
	assert.Contains(t, h, `<span class="kit-section-action-text">Bearbeiten</span>`)
	assert.Contains(t, h, "<p>x</p>")
	assert.NotContains(t, markup(t, ctx, detail.Section(detail.SectionProps{})), "kit-section-head", "no head without title and action")

	h = render(t, ctx, mit(detail.Fields(), render(t, ctx, detail.Field("Telefon", "", "–"))+render(t, ctx, detail.Field("Ort", "Bonn", "–"))))
	assert.Equal(t, `<dl class="kit-fields"><dt>Telefon</dt><dd class="kit-field-empty">–</dd><dt>Ort</dt><dd>Bonn</dd></dl>`, h)
}

func TestDetailDialog(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, mit(detail.Dialog(detail.DialogProps{ID: "dialog-kunde", Title: "Kunde", Tooltip: "Kunde (E-32)", Context: "10001 · Stadt",
		Section: true, Error: "Kundennummer ist fest"}), "<form></form>"))
	assert.Contains(t, h, `class="modal fade kit-dialog" id="dialog-kunde"`)
	assert.Contains(t, h, `data-kit-dialog`)
	assert.Contains(t, h, `data-kit-section`)
	assert.Contains(t, h, `<div class="kit-dialog-context">10001 · Stadt</div>`)
	assert.Contains(t, h, `<h2 class="modal-title" id="dialog-kunde-title" title="Kunde (E-32)">Kunde</h2>`)
	assert.Contains(t, h, `<div class="alert kit-dialog-error alert-danger" role="alert">`)
	assert.Contains(t, h, "<form></form>")

	h = render(t, ctx, detail.Dialog(detail.DialogProps{ID: "dialog-liste", Title: "Liste", Large: true}))
	assert.Contains(t, h, "kit-dialog-large")
	assert.Contains(t, h, "modal-fullscreen")
	assert.NotContains(t, h, "data-kit-section")
	assert.NotContains(t, h, "kit-dialog-error")
}

func TestDetailScript(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, detail.Script(detail.ScriptConfig{DialogField: "_drawer", StateVar: "drawerZustand"}))
	assert.True(t, strings.HasPrefix(h, "<script>(function(cfg){"))
	assert.Contains(t, h, `"dialogField":"_drawer"`)
	assert.Contains(t, h, `"returnField":"_return"`, "defaults for what is not set")
	assert.Contains(t, h, `"stateVar":"drawerZustand"`)
	assert.Contains(t, h, `"cancel":"Abbrechen"`)
	assert.Contains(t, h, `"sectionHint":"Gespeichert wird nur dieser Abschnitt."`)
	assert.Equal(t, 1, strings.Count(h, "</script>"))

	// A configured name cannot end the script.
	h = render(t, ctx, detail.Script(detail.ScriptConfig{StateVar: "</script><script>alert(1)"}))
	assert.Equal(t, 1, strings.Count(h, "</script>"))
}

func TestCreate(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := markup(t, ctx, mit(create.Form(create.FormProps{Action: "/kunden", Attributes: templ.Attributes{"data-formular": "kunde-anlegen"}}), "<input name=x>"))
	assert.Equal(t, `<form method="post" action="/kunden" class="kit-create" data-formular="kunde-anlegen"><input name=x></form>`, h)

	h = render(t, ctx, create.Footer(create.FooterProps{Cancel: "/kunden", Submit: "Kunde anlegen", Hint: "Danach öffnet sich die Akte."}))
	assert.Contains(t, h, `<span class="kit-create-hint">Danach öffnet sich die Akte.</span>`)
	assert.Contains(t, h, `<a href="/kunden" class="btn">Abbrechen</a>`)
	assert.Contains(t, h, `<button class="btn btn-primary"><i class="ti ti-plus icon" aria-hidden="true"></i>Kunde anlegen</button>`)

	h = render(t, ctx, create.Footer(create.FooterProps{}))
	assert.Contains(t, h, ">Anlegen</button>", "default label")
	assert.NotContains(t, h, "Abbrechen", "no cancel without a target")
}

func TestHistory(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, history.Timeline([]history.Entry{
		{Time: "19.10.2026 14:03", Who: "Birgit Büro", Text: "hat Felder geändert", Kind: history.Changed, Details: []string{"ort: Bonn → Köln", "plz: 1 → 2"}},
		{Time: "18.10.2026 09:00", Text: "angelegt", Kind: history.Created},
	}))
	assert.Contains(t, h, `<li class="kit-history-changed"><span class="kit-history-time">19.10.2026 14:03</span>`)
	assert.Contains(t, h, `<strong>Birgit Büro</strong> hat Felder geändert`)
	assert.Contains(t, h, `<div class="kit-history-details">ort: Bonn → Köln · plz: 1 → 2</div>`)
	assert.Contains(t, h, `<div>angelegt</div>`)
	assert.Contains(t, render(t, ctx, history.Timeline(nil)), `<li class="kit-history-empty">Keine Einträge.</li>`)
	assert.Contains(t, render(t, ctx, history.Timeline(nil, history.Props{Empty: "Noch nichts."})), "Noch nichts.")
}

func TestSignin(t *testing.T) {
	ctx := kit.WithMessages(context.Background(), kit.German)
	h := render(t, ctx, mit(signin.Page(signin.Props{Lang: "de", AppName: "Acme", Head: templ.Raw(`<link rel="stylesheet" href="/t.css">`),
		Intro: "Wähle ein Konto.", Flash: &signin.Flash{Variant: alert.VariantInfo, Text: "Abgemeldet"}}), "<p>Formular</p>"))
	assert.True(t, strings.HasPrefix(h, `<!doctype html><html lang="de">`))
	assert.Contains(t, h, "<title>Anmelden · Acme</title>")
	assert.Contains(t, h, `</style><link rel="stylesheet" href="/t.css">`, "kit CSS first, then the app's")
	assert.Contains(t, h, `<span class="kit-signin-name">Acme</span>`, "app name as brand by default")
	assert.Contains(t, h, `<h2 class="h2 text-center mb-2">Anmelden</h2>`)
	assert.Contains(t, h, "Wähle ein Konto.")
	assert.Contains(t, h, "Abgemeldet")
	assert.Contains(t, h, "<p>Formular</p>")
	assert.NotContains(t, h, "kit-sidebar", "no sidebar on the sign-in page")

	h = render(t, ctx, signin.Password(signin.PasswordProps{Action: "/login", Username: "ann", Error: "Benutzername oder Passwort falsch",
		RememberName: "remember", Hidden: []signin.HiddenField{{Name: "next", Value: "/orders"}}}))
	assert.Contains(t, h, `<input type="hidden" name="next" value="/orders">`)
	assert.Contains(t, h, "alert-danger")
	assert.Contains(t, h, `<label class="form-label" for="signin-username">Benutzername</label>`)
	assert.Contains(t, h, `value="ann"`)
	assert.Contains(t, h, `autocomplete="current-password"`)
	assert.Contains(t, h, "Angemeldet bleiben")
	assert.Regexp(t, `id="signin-password"[^>]*autofocus|autofocus[^>]*id="signin-password"`, h, "with a user name, focus the password")

	h = render(t, ctx, signin.Picker(signin.PickerProps{Action: "/anmelden"}, []signin.Choice{
		{Value: "1", Name: "Birgit Büro", Detail: "Büro", Icon: "ti-building", Color: "blue"},
		{Value: "2", Name: "Alt", Initials: "AL", Disabled: true}}))
	assert.Contains(t, h, `<input type="hidden" name="id" value="1">`)
	assert.Contains(t, h, `avatar bg-blue-lt`)
	assert.Contains(t, h, `<span class="fw-bold text-truncate">Birgit Büro</span>`)
	assert.Contains(t, h, "opacity-50")
	assert.Contains(t, h, " disabled>")
}

func TestStats(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, stats.Row([]stats.Stat{
		{Label: "Ist gesamt", Value: "4.595,00 €", Extra: "26 %"},
		{Label: "Rote Positionen", Marker: templ.Raw(`<span class="dot"></span>`), Value: "2", Link: "/freigabe", LinkText: "Prüfen"},
	}, stats.Props{Class: "mb-4"}))
	assert.Contains(t, h, `<div class="kit-stats mb-4">`)
	assert.Regexp(t, `<div class="kit-stat-value"><span>4.595,00 €</span>\s*<small>26 %</small>`, h)
	assert.Contains(t, h, `<div class="kit-stat-label"><span class="dot"></span><span>Rote Positionen</span></div>`)
	assert.Contains(t, h, `<a href="/freigabe"><span>Prüfen</span><i class="ti ti-arrow-right icon" aria-hidden="true"></i></a>`)
}

func TestLanding(t *testing.T) {
	ctx := context.Background()
	h := render(t, ctx, mit(landing.Page(landing.Props{Lang: "de", AppName: "Acme", Description: "Aufträge ohne Papier",
		Head: templ.Raw(`<link rel="stylesheet" href="/t.css">`), Links: []landing.Link{{Href: "#ablauf", Text: "Ablauf"}},
		Actions: templ.Raw(`<a class="btn" href="/anmelden">Anmelden</a>`),
		Footer:  landing.Footer(landing.FooterProps{Text: "© Acme", Links: []landing.Link{{Href: "/impressum", Text: "Impressum"}}})}), "<p>Inhalt</p>"))
	assert.Contains(t, h, `<html lang="de">`)
	assert.Contains(t, h, `<title>Acme</title>`, "app name as title by default")
	assert.Contains(t, h, `<meta name="description" content="Aufträge ohne Papier">`)
	assert.Less(t, strings.Index(h, ".kit-landing{"), strings.Index(h, `href="/t.css"`), "the app's stylesheets override the block's CSS")
	assert.Contains(t, h, `<a class="kit-landing-brand" href="/">Acme</a>`)
	assert.Contains(t, h, `<a href="#ablauf">Ablauf</a>`)
	assert.Contains(t, h, `<main><p>Inhalt</p></main>`)
	assert.Contains(t, h, `<a href="/impressum">Impressum</a>`)
	assert.NotContains(t, h, "<script", "no script of its own")

	h = render(t, ctx, landing.Hero(landing.HeroProps{Title: "Titel", Lead: "Text", Visual: templ.Raw(`<img src="/bild.png">`)}))
	assert.Contains(t, h, `<h1 class="kit-landing-title">Titel</h1>`)
	assert.Contains(t, h, `<div class="kit-landing-visual"><img src="/bild.png"></div>`)
	assert.NotContains(t, h, "kit-landing-hero-text")
	h = render(t, ctx, landing.Hero(landing.HeroProps{Title: "Titel"}))
	assert.Contains(t, h, "kit-landing-hero-text", "without a visual the text is centered")
	assert.NotContains(t, h, "kit-landing-visual")

	h = render(t, ctx, mit(landing.Section(landing.SectionProps{ID: "funktionen", Title: "Funktionen", Muted: true}),
		"<p>x</p>"))
	assert.Contains(t, h, `<section id="funktionen" class="kit-landing-section kit-landing-muted">`)
	h = render(t, ctx, landing.Features([]landing.Feature{{Icon: "ti-bolt", Title: "Schnell", Text: "In Sekunden"}}))
	assert.Contains(t, h, `<i class="ti ti-bolt icon" aria-hidden="true"></i>`)
	assert.Contains(t, h, `<h3>Schnell</h3>`)
	h = render(t, ctx, landing.Steps([]landing.Step{{Title: "Angebot"}, {Title: "Auftrag"}}))
	assert.Regexp(t, `>2</span><h3>Auftrag</h3>`, h)
	h = render(t, ctx, mit(landing.CallToAction(landing.CallToActionProps{Title: "Ausprobieren"}), `<a class="btn">Los</a>`))
	assert.Contains(t, h, `<div class="kit-landing-buttons"><a class="btn">Los</a></div>`)
}

// metric.Gauge with SVG draws its ring on the server: no ECharts option, the percentage in the
// given format.
func TestGaugeSVG(t *testing.T) {
	ctx := templ.InitializeContext(context.Background())
	p := metric.GaugeProps{ID: "g", Title: "Delivery", Value: 96.5, Label: "delivered", Count: 965, Total: 1000, Color: "green"}
	assert.Contains(t, render(t, ctx, metric.Gauge(p)), "data-tabler-chart")
	p.SVG, p.Format = true, svg.German
	html := render(t, ctx, metric.Gauge(p))
	assert.NotContains(t, html, "data-tabler-chart")
	for _, want := range []string{`id="chart-gauge-g"`, `aria-label="Delivery: 96.5%"`, `--s0: var(--tblr-green);`, `>96,5%</text>`, `>delivered</text>`} {
		assert.Contains(t, html, want)
	}
}
