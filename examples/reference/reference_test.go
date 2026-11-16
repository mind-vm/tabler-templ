package reference

// Tests of the rules of the concept Gestaltung that can be checked on the server; each test carries the
// rule's number in its name (GSO-05). GS-14 (no page scrolls sideways) needs a browser:
// scripts/reference.sh checks it at the desktop and on a phone.

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/themes"
)

type client struct {
	t   *testing.T
	srv *httptest.Server
	c   *http.Client
	app *App
}

func newClient(t *testing.T) *client {
	t.Helper()
	pause = 0
	mux := http.NewServeMux()
	app := Register(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	jar, _ := cookiejar.New(nil)
	return &client{t: t, srv: srv, c: &http.Client{Jar: jar}, app: app}
}

// get returns status and body, following redirects.
func (c *client) get(path string) (int, string) {
	c.t.Helper()
	res, err := c.c.Get(c.srv.URL + path)
	require.NoError(c.t, err)
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

// post sends a form and follows the redirect; it returns the final status, body and path.
func (c *client) post(path string, form url.Values) (int, string, string) {
	c.t.Helper()
	res, err := c.c.PostForm(c.srv.URL+path, form)
	require.NoError(c.t, err)
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b), res.Request.URL.RequestURI()
}

// pages are all pages of the reference with their status.
var pages = map[string]int{
	Home: 200, pathFrame: 200, pathStart: 200, pathCustomers: 200, pathCustomers + "?q=muehl": 200, pathCustomers + "?q=Zeppelin": 200,
	customerURL("K-1001"): 200, customerURL("K-1014"): 200, pathOffers: 200, offerURL("A-2026-001"): 200, offerURL("A-2026-002"): 200,
	pathConcept: 200, pathConceptAssistant: 200, pathRequests: 200, requestURL("R-2026-001"): 200, requestURL("R-2026-003"): 200,
	pathNew: 200, pathStates: 200, pathStates + "?meldung=1": 200, pathBrand: 200, pathAssistant: 200,
	Home + "/gibt-es-nicht": 404, pathStates + "/keine-rechte": 403, pathStates + "/fehler": 500,
}

func TestPagesServe(t *testing.T) {
	c := newClient(t)
	for p, want := range pages {
		got, body := c.get(p)
		assert.Equal(t, want, got, p)
		assert.Contains(t, body, `<html lang="de">`, p)
	}
}

// Every pattern of §11 has its page, and the page names the pattern on top.
func TestEveryPatternHasItsPage(t *testing.T) {
	c := newClient(t)
	for _, p := range patterns {
		got, body := c.get(p.Path)
		assert.Equal(t, http.StatusOK, got, p.ID)
		assert.Contains(t, body, `<span class="ref-explain-id">`+p.ID+`</span>`, p.ID)
		assert.NotEmpty(t, p.Rules, p.ID)
		assert.NotEmpty(t, p.Blocks, p.ID)
	}
	readme, err := os.ReadFile("README.md")
	require.NoError(t, err)
	for _, p := range patterns {
		assert.Contains(t, string(readme), p.ID, "README maps "+p.ID)
	}
}

var external = regexp.MustCompile(`<(?:link|script|img|source|iframe)\b[^>]*\b(?:href|src)="(?:https?:)?//|url\(["']?(?:https?:)?//|@import`)

// GS-09: nothing comes from a CDN: no stylesheet, script, image or font from another host.
func TestGS09_NothingFromOutside(t *testing.T) {
	c := newClient(t)
	for p := range pages {
		_, body := c.get(p)
		assert.False(t, external.MatchString(body), "%s loads from outside: %s", p, external.FindString(body))
	}
	_, css := c.get(c.app.files.URL("npm/@fontsource-variable/geist@5.3.0/index.css"))
	assert.Contains(t, css, "@font-face")
	assert.NotContains(t, css, "//cdn")
}

var groupLabel = regexp.MustCompile(`<div class="kit-sidebar-label">([^<]+)</div>`)

// GS-20: „Start“ on top without a title, then the groups in their order; the reference's own pages last.
func TestGS20_SidebarOrder(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathStart)
	nav := body[strings.Index(body, `class="kit-sidebar-content`):]
	var labels []string
	for _, m := range groupLabel.FindAllStringSubmatch(nav, -1) {
		labels = append(labels, m[1])
	}
	assert.Equal(t, []string{"Tagesgeschäft", "Erledigen", "Referenz"}, labels)
	assert.Less(t, strings.Index(nav, `href="/reference/start"`), strings.Index(nav, `kit-sidebar-label`), "Start stands above the first group")
}

var sidebarLink = regexp.MustCompile(`<a [^>]*class="kit-sidebar-link[^"]*"[^>]*>`)

// GS-21: every entry of the sidebar has a stable anchor.
func TestGS21_SidebarAnchors(t *testing.T) {
	c := newClient(t)
	_, body := c.get(Home)
	links := sidebarLink.FindAllString(body, -1)
	require.NotEmpty(t, links)
	for _, l := range links {
		assert.Contains(t, l, `data-tour="referenz.leiste.`, l)
	}
}

// GS-30: every figure of the start page leads to a filtered list.
func TestGS30_FiguresLeadToFilteredLists(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathStart)
	tiles := regexp.MustCompile(`<a href="([^"]+)" class="ref-tile"`).FindAllStringSubmatch(body, -1)
	require.Len(t, tiles, 4)
	for _, m := range tiles {
		href := strings.ReplaceAll(m[1], "&amp;", "&")
		assert.Contains(t, href, "?", "a tile leads to a filtered list")
		got, list := c.get(href)
		assert.Equal(t, http.StatusOK, got)
		assert.Contains(t, list, "Zurücksetzen", href+" shows the filter as active")
	}
}

var sectionHead = regexp.MustCompile(`(?s)class="kit-section-head[^"]*">(.*?)</h3>(.*?)</div>`)

// GS-41 and GS-06: an island has at most one action, usually „Bearbeiten“.
func TestGS41_IslandHasAtMostOneAction(t *testing.T) {
	c := newClient(t)
	for p := range pages {
		_, body := c.get(p)
		for _, m := range sectionHead.FindAllStringSubmatch(body, -1) {
			n := strings.Count(m[2], "<button") + strings.Count(m[2], "<a ")
			assert.LessOrEqual(t, n, 1, "%s: island %q has %d actions", p, strings.TrimSpace(stripTags(m[1])), n)
		}
	}
}

func stripTags(s string) string { return tagRE.ReplaceAllString(s, "") }

var offcanvasToggle = regexp.MustCompile(`<[^>]*data-bs-toggle="offcanvas"[^>]*>`)

var fieldRE = regexp.MustCompile(`<(?:input|select|textarea)\b[^>]*>`)

// GS-45: no drawers, and no field in the reading view of an island: fields stand only in dialogs, the
// search of a list, the create page and the panel's composer.
func TestGS45_NoDrawersNoFieldsInIslands(t *testing.T) {
	c := newClient(t)
	for p := range pages {
		_, body := c.get(p)
		// The sidebar becomes a menu on small screens (GS-14, Tabler's offcanvas); there is no other drawer.
		assert.NotContains(t, body, `class="offcanvas`, p)
		for _, tag := range offcanvasToggle.FindAllString(body, -1) {
			assert.Contains(t, tag, `data-bs-target="#sidebar"`, p)
		}
		if strings.Contains(body, `class="kit-create`) {
			continue // the create page is a form
		}
		for _, island := range islands(body) {
			if strings.Contains(island[:min(200, len(island))], "data-work") {
				continue // GS-34: on a work page the work itself (here the reply) is edited in the page
			}
			for _, f := range fieldRE.FindAllString(island, -1) {
				assert.Contains(t, f, `type="hidden"`, "%s: a field in an island: %s", p, f)
			}
		}
	}
}

// islands are the HTML of each island up to the next island, dialog or the panel.
func islands(body string) []string {
	var out []string
	parts := strings.Split(body, `kit-section"`)
	for _, part := range parts[1:] {
		end := len(part)
		for _, stop := range []string{`class="modal `, `<aside`, `</body>`} {
			if i := strings.Index(part, stop); i >= 0 && i < end {
				end = i
			}
		}
		out = append(out, part[:end])
	}
	return out
}

// GS-42: the dialog is titled like the island and says below what it is about.
func TestGS42_DialogTitleAndContext(t *testing.T) {
	c := newClient(t)
	_, body := c.get(customerURL("K-1001"))
	assert.Contains(t, body, `<div class="kit-dialog-context">K-1001 · Wohnungsgesellschaft Nordstadt</div>`)
	assert.Contains(t, body, `id="dialog-stammdaten-title">Stammdaten</h2>`)
}

// GS-43: a failed save renders the page again with the dialog's state: it opens again, with the message,
// the values entered and the field marked.
func TestGS43_DialogReopensWithErrorAndValues(t *testing.T) {
	c := newClient(t)
	got, body, _ := c.post(customerURL("K-1001")+"/stammdaten", url.Values{"name": {""}, "street": {"Hafenstraße 12"}, "zip": {"12345"}, "city": {"Musterstadt"}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, `window["kitDialogState"]={"error":"Bitte prüfen Sie die markierten Felder.","id":"dialog-stammdaten"`)
	assert.Contains(t, body, `value="Hafenstraße 12"`)
	assert.Regexp(t, `<input[^>]*id="stammdaten-name"[^>]*class="[^"]*is-invalid|<input[^>]*class="[^"]*is-invalid[^"]*"[^>]*id="stammdaten-name"`, body)
	assert.Contains(t, body, "Bitte geben Sie einen Namen ein.")
	cu, _ := c.app.store.Customer("K-1001")
	assert.Equal(t, "Wohnungsgesellschaft Nordstadt", cu.Name, "nothing saved")
}

// GS-44: a small addition is a dialog of the record's page, opened from its island.
func TestGS44_AdditionIsADialogOfThePage(t *testing.T) {
	c := newClient(t)
	_, body := c.get(customerURL("K-1001"))
	assert.Contains(t, body, `data-bs-target="#dialog-vermerk"`)
	_, body, path := c.post(customerURL("K-1001")+"/vermerke", url.Values{"text": {"Rückruf vereinbart."}})
	assert.Equal(t, customerURL("K-1001"), path)
	assert.Contains(t, body, "Vermerk hinzugefügt.")
	assert.Contains(t, body, "Rückruf vereinbart.")
}

// GS-46: a sub-list shows at most 10 rows and „Alle N“, which leads to the list with a filter.
func TestGS46_SubListAtMostTenRows(t *testing.T) {
	c := newClient(t)
	_, body := c.get(customerURL("K-1001"))
	i := strings.Index(body, "data-sublist")
	require.Positive(t, i)
	sub := body[i:]
	sub = sub[:strings.Index(sub, "</table>")]
	assert.Equal(t, subListMax, strings.Count(sub, "<tr>")-1, "header plus 10 rows")
	assert.Contains(t, body, `<a href="/reference/angebote?q=K-1001" data-all>Alle 14 Angebote</a>`)
	_, list := c.get(pathOffers + "?q=K-1001")
	assert.Contains(t, list, "1–14 von 14")
}

// GS-50: search ignores case and umlauts (ä = ae), every word must occur; one filter; 25 rows.
func TestGS50_SearchFoldsUmlauts(t *testing.T) {
	assert.True(t, matches("muehl", "Mühlhausen"))
	assert.True(t, matches("MÜHL weg", "Hausverwaltung Mühlenweg"))
	assert.True(t, matches("strasse", "Hauptstraße"))
	assert.False(t, matches("mühl nord", "Mühlhausen"))
	c := newClient(t)
	_, body := c.get(pathCustomers + "?q=muehl")
	assert.Contains(t, body, "Mühl")
	assert.NotContains(t, body, "Wohnungsgesellschaft Nordstadt</a>")
}

func TestGS50_PagesOf25AndStateInURL(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathCustomers)
	assert.Contains(t, body, "1–25 von 58")
	assert.Len(t, regexp.MustCompile(`<select[^>]*data-kit-list-submit`).FindAllString(body, -1), 1, "one filter")
	_, body = c.get(pathCustomers + "?posten=gruen&page=2")
	assert.Contains(t, body, "26–")
	assert.Contains(t, body, `href="/reference/kunden?posten=gruen"`, "the link to page 1 keeps the filter")
	assert.Contains(t, body, `<option value="gruen" selected`)
}

// GS-52: numbers right-aligned, status as a pill, the traffic light as a dot with a word.
func TestGS52_Columns(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathCustomers)
	assert.Contains(t, body, `<td class="text-end ref-num">`)
	assert.Contains(t, body, `ref-pill`)
	assert.Regexp(t, `class="status[^"]*ref-signal"><span class="status-dot"></span>überfällig`, body)
}

// GS-53: sorted by what needs action: overdue first.
func TestGS53_OverdueFirst(t *testing.T) {
	cs := NewStore().Customers()
	assert.Equal(t, SignalOverdue, cs[0].Signal)
	assert.Equal(t, SignalNone, cs[len(cs)-1].Signal)
}

// GS-54: wide tables scroll in their island.
func TestGS54_TablesScrollInTheirIsland(t *testing.T) {
	c := newClient(t)
	for _, p := range []string{pathCustomers, offerURL("A-2026-001"), customerURL("K-1001")} {
		_, body := c.get(p)
		assert.Contains(t, body, `class="table-responsive`, p)
	}
}

// GS-55 and GS-56: the list is read on the page; Edit opens the full-page dialog, one form per row.
func TestGS56_EveryRowIsAForm(t *testing.T) {
	c := newClient(t)
	_, body := c.get(offerURL("A-2026-001"))
	assert.Contains(t, body, `data-bs-target="#dialog-positionen"`)
	assert.Contains(t, body, "kit-rowedit-dialog")
	assert.Contains(t, body, "Jede Zeile wird für sich gespeichert.")
	assert.Equal(t, 12, strings.Count(body, `data-form="position"`))
	before, _, _ := strings.Cut(body, `id="dialog-positionen"`)
	for _, island := range islands(before) {
		assert.NotContains(t, island, `name="qty"`, "the page itself only reads")
	}
}

// GS-59: saving a row recomputes amount and total.
func TestGS59_SaveRowRecomputesTotal(t *testing.T) {
	c := newClient(t)
	o, _ := c.app.store.Offer("A-2026-001")
	p := o.Positions[0]
	got, body, path := c.post(offerURL(o.No)+"/positionen/"+p.ID, url.Values{"id": {p.ID}, "no": {p.No}, "text": {p.Text}, "qty": {"60"}, "unit": {p.Unit}, "price": {"89,50"}})
	assert.Equal(t, http.StatusOK, got)
	assert.Equal(t, offerURL(o.No), path)
	assert.Contains(t, body, "Position 01.01.0010 gespeichert.")
	assert.Contains(t, body, Euro(o.Total()+8*8950), "total after 52 → 60")
	assert.Contains(t, body, "5.370,00 €", "amount 60 × 89,50")
}

func TestGS56_FailedRowReopensTheDialogOnThatRow(t *testing.T) {
	c := newClient(t)
	o, _ := c.app.store.Offer("A-2026-001")
	p := o.Positions[1]
	got, body, _ := c.post(offerURL(o.No)+"/positionen/"+p.ID, url.Values{"id": {p.ID}, "no": {p.No}, "text": {p.Text}, "qty": {"abc"}, "price": {"1"}, "_form": {kindPosition}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, `"values":{"_form":["position"],"id":["`+p.ID+`"]}`)
	assert.Contains(t, body, "Position 01.01.0020: Die Menge ist eine Zahl größer als 0")
	assert.Contains(t, body, `value="abc"`)
}

// GS-57: the new row waits at the bottom with the next number proposed.
func TestGS57_NewRowProposesTheNextNumber(t *testing.T) {
	c := newClient(t)
	_, body := c.get(offerURL("A-2026-001"))
	i := strings.Index(body, `data-form="position-neu"`)
	require.Positive(t, i)
	assert.Greater(t, i, strings.LastIndex(body, `data-form="position"`), "below the rows")
	assert.Contains(t, body[i:], `value="02.01.0030"`)
	_, body, _ = c.post(offerURL("A-2026-001")+"/positionen", url.Values{"no": {"02.01.0030"}, "text": {"Neue Leistung"}, "qty": {"2,5"}, "price": {"100"}})
	assert.Contains(t, body, "Position 02.01.0030 hinzugefügt.")
	assert.Contains(t, body, `value="02.01.0040"`)
}

// GS-58: deleting asks first, and the message says what was deleted.
func TestGS58_DeleteAsksAndSaysWhat(t *testing.T) {
	c := newClient(t)
	_, body := c.get(offerURL("A-2026-001"))
	assert.Contains(t, body, `data-kit-confirm="Position 01.01.0010 löschen?"`)
	o, _ := c.app.store.Offer("A-2026-001")
	_, body, _ = c.post(offerURL(o.No)+"/positionen/"+o.Positions[0].ID+"/loeschen", nil)
	assert.Contains(t, body, "Position 01.01.0010 gelöscht.")
}

// GS-59: a list that cannot change shows the reason instead of the buttons, and the server refuses.
func TestGS59_LockedListShowsTheReason(t *testing.T) {
	c := newClient(t)
	_, body := c.get(offerURL("A-2026-002"))
	assert.Contains(t, body, "Die Positionen sind schreibgeschützt.")
	assert.NotContains(t, body, `id="dialog-positionen"`)
	o, _ := c.app.store.Offer("A-2026-002")
	_, err := c.app.store.UpdatePosition(o.No, o.Positions[0].ID, PositionInput{No: "01.01.0010", Text: "x", Qty: "1", Price: "1"})
	assert.ErrorIs(t, err, errLocked)
}

// GS-33 and GS-60: create is a page of its own; afterwards the new detail page with a message and the
// next step. A failed create shows the error at the field.
func TestGS33_CreateIsAPageThenTheDetailPage(t *testing.T) {
	c := newClient(t)
	got, body := c.get(pathNew)
	assert.Equal(t, http.StatusOK, got)
	assert.Contains(t, body, `class="kit-create"`)
	got, body, _ = c.post(pathCustomers, url.Values{"name": {"Stiftung Probe"}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, "Bitte geben Sie einen Ort ein.")
	assert.Contains(t, body, `value="Stiftung Probe"`)
	got, body, path := c.post(pathCustomers, url.Values{"name": {"Stiftung Probe"}, "city": {"Musterstadt"}})
	assert.Equal(t, http.StatusOK, got)
	assert.Equal(t, customerURL("K-1059"), path)
	assert.Contains(t, body, "Kunde Stiftung Probe angelegt.")
	assert.Contains(t, body, `class="alert-action">Ansprechpartner eintragen</a>`)
}

// GS-60: after an action a message with the next step; it shows once.
func TestGS60_MessageWithNextStep(t *testing.T) {
	c := newClient(t)
	_, body, _ := c.post(offerURL("A-2026-001")+"/abgeben", nil)
	assert.Contains(t, body, "Angebot A-2026-001 abgegeben.")
	assert.Contains(t, body, `<a href="/reference/angebote" class="alert-action">Zu den Angeboten</a>`)
	_, body = c.get(offerURL("A-2026-001"))
	assert.NotContains(t, body, "Angebot A-2026-001 abgegeben.", "the message shows once")
}

// GS-61: an empty search says so and offers to reset it.
func TestGS61_EmptySearchOffersReset(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathCustomers + "?q=Zeppelin")
	assert.Contains(t, body, `data-empty="search"`)
	assert.Contains(t, body, "Suche zurücksetzen")
	_, body = c.get(pathStates)
	assert.Contains(t, body, `data-empty="list"`)
}

// GS-62: error pages stand in the frame with the sidebar and say in one sentence what is wrong.
func TestGS62_ErrorPagesInTheFrame(t *testing.T) {
	c := newClient(t)
	for _, p := range []string{Home + "/gibt-es-nicht", pathStates + "/keine-rechte", pathStates + "/fehler", customerURL("K-9999")} {
		_, body := c.get(p)
		assert.Contains(t, body, `class="kit-sidebar`, p)
		assert.Contains(t, body, "data-error-page", p)
		assert.NotContains(t, body, "goroutine", p)
	}
}

// GS-63 and GS-08: read-only hides every edit button, shows the notice, and the server refuses changes.
func TestGS63_ReadOnly(t *testing.T) {
	c := newClient(t)
	c.post(pathStates+"/nur-lesen", nil)
	for _, p := range []string{customerURL("K-1001"), offerURL("A-2026-001"), pathCustomers} {
		_, body := c.get(p)
		assert.Contains(t, body, "data-read-only", p)
		assert.NotRegexp(t, `class="[^"]*kit-section-action`, body, p)
		assert.NotContains(t, body, `data-tour="referenz.kunden.neu"`, p)
		assert.NotContains(t, body, `class="modal `, p)
	}
	got, _, _ := c.post(customerURL("K-1001")+"/stammdaten", url.Values{"name": {"X"}, "city": {"Y"}})
	assert.Equal(t, http.StatusForbidden, got)
	got, _ = c.get(pathNew)
	assert.Equal(t, http.StatusForbidden, got)
}

// GS-64: what takes longer comes after the page: a placeholder first.
func TestGS64_LoadsAfterThePage(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathStates)
	assert.Contains(t, body, `data-ref-load="/reference/zustaende/zahlen"`)
	assert.Contains(t, body, "placeholder-glow")
	_, frag := c.get(pathStates + "/zahlen")
	assert.Contains(t, frag, "Angebote im Entwurf")
	assert.NotContains(t, frag, "<html")
}

var ruleNumber = regexp.MustCompile(`\b(?:GS|GSO|KA|KAO|E|G|NV|UX|RA)-\d+\b`)
var explainBox = regexp.MustCompile(`(?s)<section class="ref-explain".*?</section>|<article class="ref-concept-text.*?</article>`)

// GS-71: the pages show no rule numbers; only the explanation of the reference and the concept's own text
// name them.
func TestGS71_NoRuleNumbersOutsideTheExplanation(t *testing.T) {
	c := newClient(t)
	for p := range pages {
		_, body := c.get(p)
		visible := explainBox.ReplaceAllString(body, "")
		visible = regexp.MustCompile(`(?s)<(script|style)[^>]*>.*?</(script|style)>`).ReplaceAllString(visible, "")
		assert.Empty(t, ruleNumber.FindAllString(stripTags(visible), -1), p)
	}
}

// GS-73: numbers and dates in the German format.
func TestGS73_Formats(t *testing.T) {
	assert.Equal(t, "1.234,50 €", Euro(123450))
	assert.Equal(t, "0,05 €", Euro(5))
	assert.Equal(t, "−12,00 €", Euro(-1200))
	assert.Equal(t, "15.11.2026", Date(Today))
	assert.Equal(t, "01.04.2026 – 31.03.2028", Period(day(4, 1), day(3, 31).AddDate(2, 0, 0)))
	assert.Equal(t, "1.234,5", DecimalOf(1234.5).String())
	for in, want := range map[string]int64{"48,50": 4850, "1.200": 120000, "1.234,56 €": 123456, "12.5": 1250, "7": 700} {
		got, err := ParseCents(in)
		assert.NoError(t, err, in)
		assert.Equal(t, want, got, in)
	}
	_, err := ParseCents("zwölf")
	assert.Error(t, err)
}

// KA-20: the assistant never fills in bank details or e-mail addresses; KA-02: it fills forms that are on
// their page, with the fields named.
func TestKA20_AssistantNeverFillsBankDetails(t *testing.T) {
	_, err := catalog.Check("buero", "kunde-stammdaten", map[string]string{"id": "K-1001", "iban": "DE02120300000000202051"})
	assert.Error(t, err)
	_, err = catalog.Check("buero", "kunde-stammdaten", map[string]string{"id": "K-1001", "street": "DE02120300000000202051"})
	assert.Error(t, err)
	assert.Empty(t, catalog.Problems())
}

func TestKA02_FillActionsFindTheirForm(t *testing.T) {
	c := newClient(t)
	for _, a := range catalog.Actions {
		if a.Form == "" {
			continue
		}
		_, body := c.get(strings.ReplaceAll(a.Path, "{id}", "K-1001"))
		i := strings.Index(body, `data-form="`+a.Form+`"`)
		require.Positive(t, i, a.Name)
		form := body[i : i+strings.Index(body[i:], "</form>")]
		for _, f := range a.Fields {
			assert.Contains(t, form, `name="`+f+`"`, a.Name)
		}
	}
}

func TestAssistantAnswersAndActs(t *testing.T) {
	c := newClient(t)
	res, err := c.c.Post(c.srv.URL+pathAssistant, "application/json", strings.NewReader(`{"message":"Trage für die Wohnungsgesellschaft Nordstadt die Anschrift Hafenstraße 12 ein"}`))
	require.NoError(t, err)
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	assert.Contains(t, string(b), `"form":"kunde-stammdaten"`)
	assert.Contains(t, string(b), `"street":"Hafenstraße 12"`)
	assert.Contains(t, string(b), "speichern Sie selbst")
	_, page := c.get(pathFrame)
	assert.Contains(t, page, "Hafenstraße 12", "the panel shows the conversation again on the next page")

	res, err = c.c.Post(c.srv.URL+pathAssistant, "application/json", strings.NewReader(`{"message":"Ändere die IBAN der Wohnungsgesellschaft Nordstadt"}`))
	require.NoError(t, err)
	b, _ = io.ReadAll(res.Body)
	res.Body.Close()
	assert.Contains(t, string(b), "Bankverbindungen fülle ich nie aus.")
	assert.NotContains(t, string(b), `"form":"kunde-stammdaten"`)
}

// GS-13: the panel has the tabs notes and assistant, comes open from the cookie, and holds no form to
// edit a record.
func TestGS13_PanelFromCookie(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathFrame)
	assert.Contains(t, body, `<aside id="kit-sidepanel" class="kit-sidepanel" data-active="notizen" aria-label="Notizen und KI-Assistent" hidden`)
	u, _ := url.Parse(c.srv.URL)
	c.c.Jar.SetCookies(u, []*http.Cookie{{Name: "kit-sidepanel", Value: tabAssistant, Path: "/"}})
	_, body = c.get(pathFrame)
	assert.Contains(t, body, `data-active="assistent"`)
	assert.Contains(t, body, `kit-sidepanel-open`)
	i := strings.Index(body, `<aside id="kit-sidepanel"`)
	panel := body[i : i+strings.Index(body[i:], "</aside>")]
	assert.NotContains(t, panel, `method="post"`)
}

// Saving in a dialog returns to the record's page, never to an address that is no page (after a failed
// save the address bar shows the form's action).
func TestDialogReturnsToTheRecordsPage(t *testing.T) {
	c := newClient(t)
	o, _ := c.app.store.Offer("A-2026-001")
	_, body, path := c.post(offerURL(o.No)+"/positionen/"+o.Positions[0].ID+"/loeschen", url.Values{"_return": {offerURL(o.No) + "/positionen/" + o.Positions[1].ID}})
	assert.Equal(t, offerURL(o.No), path)
	assert.Contains(t, body, "gelöscht")
	_, _, path = c.post(customerURL("K-1001")+"/vermerke", url.Values{"text": {"x"}, "_return": {"https://example.org/"}})
	assert.Equal(t, customerURL("K-1001"), path)
}

// The concept is shown in the app, and links both ways: the explanation's rules lead to their rows, and the
// examples and rules in the concept lead to the page that shows them.
func TestConceptLinksBothWays(t *testing.T) {
	c := newClient(t)
	bodies := map[string]string{}
	for _, x := range concepts {
		got, body := c.get(x.Path)
		require.Equal(t, http.StatusOK, got)
		bodies[x.Path] = body
		assert.NotRegexp(t, `href="(?:\.\./|[a-z-]+\.md)`, body, "relative links resolved")
	}
	g := bodies[pathConcept]
	assert.Contains(t, g, `<td class="ref-rule-cell" id="gs-41"><a href="/reference/kunden/K-1001" class="ref-backlink"`)
	assert.Contains(t, g, `<a href="/reference/angebote/A-2026-001" class="ref-backlink" title="In der Referenz: B-05 Liste bearbeiten">B-05</a>`)
	assert.Contains(t, g, `id="6-inseln-und-dialoge"`)
	assert.Contains(t, g, `href="/reference/konzept/assistent"`, "the link to the other concept stays in the app")
	sha, _ := source()
	assert.Contains(t, g, "https://github.com/mind-vm/auftragsbutler/blob/"+sha+"/docs/konzepte/navigation.md")

	// Every link of an explanation finds its anchor.
	for _, p := range append(patterns, overviewText) {
		links := []string{p.Link}
		for _, r := range p.Rules {
			links = append(links, conceptLink(r))
		}
		for _, l := range links {
			path, frag, _ := strings.Cut(l, "#")
			require.Contains(t, bodies, path, l)
			assert.Contains(t, bodies[path], `id="`+frag+`"`, "%s: %s has no anchor", p.ID, l)
		}
	}
	_, body := c.get(customerURL("K-1001"))
	assert.Contains(t, body, `<a class="ref-rule" href="/reference/konzept#gs-41"`)
}

func TestSlugIsGitHubs(t *testing.T) {
	assert.Equal(t, "3-application-layout-rahmen", slug("3. Application Layout (Rahmen)"))
	assert.Equal(t, "8-rückmeldung-und-zustände", slug("8. Rückmeldung und Zustände"))
	assert.Equal(t, "71-anzeigen", slug("7.1 Anzeigen"))
}

// GS-15: the theme is chosen in one place; every theme is offered, Graphite by default, and its fonts come
// from the bundle (GS-09 checks that nothing comes from outside).
func TestGS15_ThemeInOnePlace(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathBrand)
	assert.Equal(t, 1, strings.Count(body, `data-ref-theme aria-label`), "one switcher")
	for _, th := range themes.All {
		assert.Contains(t, body, `<option value="`+themeName(th)+`"`, th.Label)
	}
	assert.Contains(t, body, `<option value="graphite" selected>`)
	for _, f := range themeFonts() {
		got, _ := c.get(c.app.files.URL("npm/@fontsource-variable/" + f.Package + "@" + themes.FontsVersion + "/index.css"))
		assert.Equal(t, http.StatusOK, got, f.Package)
	}
	got, _ := c.get(c.app.files.URL("core/css/tabler-themes.min.css"))
	assert.Equal(t, http.StatusOK, got)
}

// GS-44: an action with a reason is a small dialog with a required field and a red button; without a
// reason it opens again with the error at the field.
func TestGS44_ActionWithReason(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathStates)
	i := strings.Index(body, `id="`+dialogWithdraw+`"`)
	require.Positive(t, i)
	assert.Contains(t, body[i:], "btn-danger")
	got, body, _ := c.post(pathStates+"/zuruecknehmen", url.Values{"reason": {" "}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, `window["kitDialogState"]={"error":"Bitte prüfen Sie die markierten Felder.","id":"`+dialogWithdraw+`"`)
	assert.Contains(t, body, "Bitte geben Sie einen Grund an.")
	_, body, _ = c.post(pathStates+"/zuruecknehmen", url.Values{"reason": {"Kunde hat neue Mengen geschickt."}})
	assert.Contains(t, body, "Abgabe von A-2026-002 zurückgenommen")
	o, _ := c.app.store.Offer("A-2026-002")
	assert.Equal(t, OfferDraft, o.Status)
	cu, _ := c.app.store.Customer(o.Customer)
	assert.Contains(t, cu.Notes[0].Text, "Kunde hat neue Mengen geschickt.")
}

// GS-44: an irreversible action asks to type the name; a wrong one changes nothing.
func TestGS44_TypedConfirmation(t *testing.T) {
	c := newClient(t)
	cu, _ := c.app.store.Customer(customerToDelete)
	got, body, _ := c.post(pathStates+"/kunde-loeschen", url.Values{"confirm": {"ja"}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, "Bitte tippen Sie den Namen genau so: "+cu.Name)
	_, ok := c.app.store.Customer(customerToDelete)
	assert.True(t, ok, "nothing deleted")
	_, body, _ = c.post(pathStates+"/kunde-loeschen", url.Values{"confirm": {cu.Name}})
	assert.Contains(t, body, "Kunde "+cu.Name+" endgültig gelöscht.")
	_, ok = c.app.store.Customer(customerToDelete)
	assert.False(t, ok)
}

// GS-63: the notice „Demo-Instanz“ stands where „nur lesen“ stands, above the message line, in info.
func TestGS63_DemoInstance(t *testing.T) {
	c := newClient(t)
	_, body, _ := c.post(pathStates+"/demo", nil)
	i := strings.Index(body, "data-demo")
	require.Positive(t, i)
	assert.Less(t, i, strings.Index(body, "Der Hinweis „Demo-Instanz“ ist an"), "above the message")
	assert.Contains(t, body[strings.LastIndex(body[:i], "<div class=\"alert"):i], "alert-info")
	_, body = c.get(customerURL("K-1001"))
	assert.Contains(t, body, "Demo-Instanz: Alle Daten hier sind erfunden.")
}

// GS-56: the long text is in a closed second line of the row; after an error in it, that row's line is open.
func TestGS56_LongTextOpensAfterAnError(t *testing.T) {
	c := newClient(t)
	_, body := c.get(offerURL("A-2026-001"))
	assert.Contains(t, body, `<details class="kit-rowedit-details"><summary>Langtext (`)
	assert.NotContains(t, body, `<details class="kit-rowedit-details" open>`)
	o, _ := c.app.store.Offer("A-2026-001")
	p := o.Positions[1]
	got, body, _ := c.post(offerURL(o.No)+"/positionen/"+p.ID, url.Values{"id": {p.ID}, "no": {p.No}, "text": {p.Text}, "qty": {"1"}, "price": {"1"},
		"long": {strings.Repeat("x", maxLong+1)}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Equal(t, 1, strings.Count(body, `<details class="kit-rowedit-details" open>`))
	i := strings.Index(body, `value="`+p.ID+`"`)
	require.Positive(t, i)
	assert.Contains(t, body[i:i+strings.Index(body[i:], "</form>")], `<details class="kit-rowedit-details" open>`, "the row that failed")
	assert.Contains(t, body, "Der Langtext hat 601 Zeichen, erlaubt sind 600.")
}

// GS-52: an open offer's deadline is a traffic light whose word is the time left; the detail page names the
// exact time (GS-73).
func TestGS52_DeadlineWithTimeLeft(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathStart)
	assert.Contains(t, body, `data-deadline="near"`)
	assert.Contains(t, body, "noch 3 Tage")
	_, body = c.get(offerURL("A-2026-001"))
	assert.Contains(t, body, "18.11.2026, 12:00 Uhr")
	assert.Contains(t, body, `title="18.11.2026, 12:00 Uhr"`)
}

// GS-54: the offers become cards below 768 px, every cell names its column; the customers scroll instead.
func TestGS54_OffersBecomeCards(t *testing.T) {
	c := newClient(t)
	_, body := c.get(pathOffers)
	assert.Contains(t, body, "table-mobile-md")
	i := strings.Index(body, "<tbody>")
	row := body[i : i+strings.Index(body[i:], "</tr>")]
	assert.Equal(t, strings.Count(row, "<td"), strings.Count(row, "data-label="), "every cell has its label")
	_, body = c.get(pathCustomers)
	assert.NotContains(t, body, "table-mobile-md")
}

// B-10, GS-34, KA-01, KA-02: the AI proposes, the person changes and sends; nothing is sent without them.
func TestB10_ProposalTheyDecide(t *testing.T) {
	c := newClient(t)
	_, body := c.get(requestURL("R-2026-001"))
	assert.Contains(t, body, "Hier schlägt eine KI vor.")
	assert.Contains(t, body, "Sicherheit: hoch")
	assert.Contains(t, body, `href="/reference/angebote/A-2026-002"`, "what it read")
	assert.Contains(t, body, `data-mark="external"`)
	got, body, _ := c.post(requestURL("R-2026-001")+"/senden", url.Values{"reply": {" "}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, "Die Antwort ist leer.")
	_, body, _ = c.post(requestURL("R-2026-001")+"/senden", url.Values{"reply": {"Guten Tag, wir bestätigen den 4. und 11.12."}})
	assert.Contains(t, body, "Antwort an Wohnungsgesellschaft Nordstadt gesendet.")
	assert.Contains(t, body, "Antwort gesendet (geändert):")
	assert.NotContains(t, body, `id="antwort-text"`, "decided: nothing to send any more")
}

// GS-44: discarding needs a reason from the list; without one the dialog opens again.
func TestB10_DiscardNeedsAReason(t *testing.T) {
	c := newClient(t)
	got, body, _ := c.post(requestURL("R-2026-002")+"/verwerfen", url.Values{"reason": {""}})
	assert.Equal(t, http.StatusUnprocessableEntity, got)
	assert.Contains(t, body, `"id":"`+dialogDiscard+`"`)
	assert.Contains(t, body, "Bitte wählen Sie einen Grund.")
	_, body, _ = c.post(requestURL("R-2026-002")+"/verwerfen", url.Values{"reason": {"Erst intern klären"}})
	assert.Contains(t, body, "Vorschlag verworfen; nichts wurde gesendet.")
	assert.Contains(t, body, "Vorschlag verworfen: Erst intern klären")
}
