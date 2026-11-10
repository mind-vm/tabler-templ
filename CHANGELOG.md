# Changelog

## Unreleased

**Breaking:** the charts are drawn with Apache ECharts 6 (Apache-2.0) instead of ApexCharts. ApexCharts
needs a paid license from USD 2M annual revenue (affiliates included), which reached every app built on
tabler-templ.

- **`chart`:** `chart.Props` stays as it was and draws the same charts (line, area, bar, pie, donut, radial
  bar, radar, polar area, heatmap, treemap, timeline, box plot, bubble, scatter, candlestick; sparkline,
  stacked, funnel, two y-axes, annotations, brush, synced groups, legend toggles). `Library` loads
  `echarts@6.0.0`; `versions.ECharts` replaces `versions.ApexCharts`.
- **`Props.Options` and `RawProps`** take ECharts options now; formatter paths follow them, e.g.
  `yAxis.0.axisLabel.formatter`. `chart.Gauge` builds the option of a single percentage ring.
- **`data-tabler-chart`** holds `{"option": …}` instead of `{"config": …}`; `element.tablerChart` is the
  ECharts instance, and `window.tabler_chart` is gone. `window.tablerCharts` (`init`, `flush`, `destroy`)
  is unchanged.
- **Funnel** charts are drawn as a funnel of trapezoids instead of centered bars.
- **`blocks/metric.Gauge`, `chart.Radial`** use `chart.Gauge`.

## v0.15.1

- **`marketing.FeatureCards`, `marketing.Plans`, `marketing.Team`:** the cards in one row are always equally
  tall; the column is a flex box and the card fills it, so this no longer depends on the marketing sheet or
  on loading it after Tabler. The featured plan's border also holds when the sheet comes before Tabler.

## v0.15.0

A black-and-white landing page in the minimal look of today's component-library sites, built from the
marketing blocks and the Graphite theme.

- **`marketing.Marquee`:** a row of logos, badges or quotes that scrolls endlessly in pure CSS; it pauses on
  hover and focus, can run slow, fast or in reverse, and without motion (`prefers-reduced-motion`) it is a
  centered row. The copy for the loop is `aria-hidden` and `inert`.
- **`marketing.CTABox`:** `Inline` puts the text at the start and the actions at the end of one row;
  `Color: "light"` draws a bordered box without the decorative circle.
- **`marketing.PageProps`:** `Theme` (color mode), `Appearance` and `HTMLAttributes` like `layout.Page`, and
  `themes.Theme.ApplyMarketing` gives a marketing page a ready-made theme.
- **Demo:** `/landing/black-white` (Tasklane: a live issue table in the hero instead of a screenshot, logo
  marquee, feature cards, bento, testimonials, plans with a monthly/yearly switch, inline call to action);
  the landing gallery shows the marquee and the inline call to action.

## v0.14.1

- **`blocks/assistant`:** among several forms of one name, the first form with the action's key is the row
  also when `kitDialog.matchForm` finds none, e.g. because `detail.Script` marks forms with another
  attribute (`data-form`) than `ActionScript` (`WithActionOptions`).

## v0.14.0

The assistant fills in the right form: among several forms of one name it takes the row with the action's
key or the section that shows the fields, focuses the first field once a dialog is open, lets the app
reveal a panel first, and takes another form attribute. Catalogs can require a role for every action.

- **`blocks/assistant`, filling in forms** (lessons from the assistant of auftragsbutler):
  - Several forms of one name: the row `kitDialog.matchForm` picks wins only when it carries the action's
    key (a hidden or read-only field the action names, with the action's value) and shows as many of the
    fields as the best form; otherwise the form that shows most of the fields. Sections of one record that
    carry the other sections' fields hidden no longer get the wrong section. Without `kitDialog` the first
    form with the key is the row, so a list's rows work without `detail.Script`.
  - The first field filled in gets the focus once a dialog is shown (`shown.bs.modal`), not after a fixed
    350 ms; a radio button gets it instead of its label.
  - `assistant:reveal` on the form before it is shown, so the app can open a panel or tab the form sits in.
  - A checkbox is also checked by "an" or by its own value.
  - `WithActionOptions(ctx, ActionOptions{FormAttribute: "data-formular"})`: another attribute than
    `data-form`, like `detail.ScriptConfig.FormAttribute`.
- **`Catalog.RequireRoles`:** every action names its roles; `Problems` reports one without, and `Check`,
  `Names` and `Prompt` do not offer it to anyone (otherwise an action without roles is for everyone).

## v0.13.0

An assistant that acts on the page: a catalog of allowed actions (open a page, fill in a form, point at
an element), a copilot panel that stays open from page to page, the conversation kept on the server, and
Markdown without links and images for a model's answers. With the demo app `/copilot`.

- **`markdown`:** `Strict(src)` and `StrictHTML(src)` render without links and images: a link becomes
  its text, an autolink its address as text, an image its alt text. For a model's answer that must not
  lead to another site or load a file from outside (prompt injection). `Markdown` and `HTML` are
  unchanged.
- **`blocks/assistant`, acting on the page:** an assistant beside an app may open pages, fill in forms
  and point at elements, but never save.
  - `Catalog` is the app's allowlist. `Check` turns an agent's request into an `Action` or refuses it
    (`ErrAction`): unknown names, other roles, fields the form lacks, IDs that are not IDs, paths
    outside `Pages`, sensitive field names and values that look like an IBAN or e-mail address.
    `Problems` checks the catalog in the app's tests; `Prompt` describes it for the model.
  - `Stream.Action` sends the action's card and an `action` event. `ActionScript` runs the actions once
    the stream has ended: it finds `form[data-form]`, opens its dialog, fills in and marks the fields
    and puts a banner over the form; `[data-tour]` elements are highlighted. An action for another
    page runs there after loading. A cancelable `assistant:action` event lets the app take over.
  - `ActionCard`, `ActionRefused`, `AnswerStrict` (Markdown without links and images).
  - `Copilot` and `CopilotButton`: the conversation in an offcanvas at the right, without backdrop,
    open from page to page; on wide screens the page makes room for it.
  - `Transcript` with `Stream.Record` and `Stream.Turn`: the server keeps the conversation and
    `Transcript.View` draws it again, e.g. on the next page.
  - The intro hides once a conversation has turns; the conversation starts at the latest turn.
- **Demo app `/copilot`** (`examples/copilot`): orders and customers with a copilot, offline with a
  scripted agent: an overdue list, a new order filled in, an order set to done on its own page, and
  two refusals (bank details, a prompt injection in an order's note).

## v0.12.0

Interactive apps: package `htmx` for components with their own endpoints and a CRUD demo built from them,
an AI assistant with Markdown, enterprise layouts (list beside detail, sticky action bar, workspace
switcher), animations and sounds, less CSS per page, and htmx 2.0.11.

- **`htmx`, after an architecture review:**
  - `Error` sends the real status (404, 409 …) instead of 200. `Setup` swaps these toasts in and does not
    add a generic error toast. Domain errors map to a status and a message with rules: `htmx.Is(err,
    code, message)`.
  - `Redirect` only goes to local paths and refuses anything else with 500, so no `next` value can lead
    to another site; `LocalPath` checks a path, and `RedirectExternal` is for URLs the app chose itself.
  - `Setup` runs once per page, even when rendered twice or swapped in again. It reads the ids from the Go
    constants, looks the elements up when needed, and warns when Tabler's script is missing.
  - `Vary: HX-Request` on every answer (`Error`, `Redirect` too); `htmx.Vary` for pages written by hand.
  - `Respond` logs render errors itself and returns nothing (breaking for callers that used its error).
  - New: `On(events...)` for hx-trigger, `Frame` with `Page(title)` for `FullPage`. `Vals` and `Query`
    panic on an odd number of arguments.
  - htmx 2.0.11 (was 2.0.4).
- **CRUD demo:** the event names live in `ui/events`, so no component imports another to learn them. The
  error messages live once in `ui/uierr`. A version field that is not a number is a 400.
- **Demo:** the assistant, assistant parts, motion and sound pages load htmx with `htmx.Library` and
  `htmx.Setup` instead of their own copy.
- **`blocks/assistant`:** the diff card builds its `hx-vals` with `htmx.Vals`, which escapes as JSON.

- **`htmx` for components with their own endpoints:**
  - `Route` defines an endpoint once: `Pattern` goes to the mux, `URL` to the markup. `CheckRoutes` is the
    contract test.
  - `ID`, `Target`, `Vals` and `Query` name a component instance's elements and requests.
  - `Respond` writes a fragment with its headers first (`Vary: HX-Request`, status, `HX-Trigger` …).
    Requests without htmx get the whole page (`FullPage`). Toasts and other fragments go out of band.
  - `Error` and `Fail` answer failures as a danger toast, without details for 5xx. `Redirect` does
    Post/Redirect/Get, or `HX-Redirect` for htmx.
  - `Setup` is the page's client:
    - toast container and dialog;
    - 422 responses swapped in, `allowEval` off;
    - toasts when offline or on an error;
    - a toast closes once its action succeeded;
    - tooltips and popovers start, and charts are torn down, around swaps.
  - kit texts `ErrorGeneric` and `Offline`.
- **CRUD demo (`examples/crud`, `/crud` in the demo):** contacts built from six co-located components
  (table, form, inline status, delete with undo, live count, detail page). Each has its own contract,
  handlers and tests, and it works without JavaScript too.
- **Tests:** `TestNonceRendered` renders each package in a fresh context. Before, the once-per-page CSS
  made it fail now and then, depending on the map order.

- **`motion`, enter and leave animations** with Motion (motion.dev, formerly Framer Motion), its build
  without React:
  - `Attrs` make an element fade, slide, zoom or pop in when the page loads, when htmx swaps it in, or
    when it scrolls into view (`InView`). `Stagger` brings its children in one after another, `Bounce`
    uses a spring.
  - `Exit` animates an element out when htmx replaces or deletes it; give the swap a delay
    (`hx-swap="outerHTML swap:200ms"`).
  - `Counter` counts a number up, formatted for the page language.
  - Visitors who ask for reduced motion see everything at once. Motion loads from `cdn.Base`
    (`versions.Motion`), and `cmd/vendor` fetches it.
- **`sound`, interface sounds** made by the browser (Web Audio), with no library or sound file:
  - `Message`, `Notify`, `Success`, `Warning`, `Error` and `Click`; `Config.Files` adds an app's own.
  - Play one on click (`Attrs`), when a fragment is swapped in (`Play`), or from the server with
    `HX-Trigger: {"sound":"message"}`.
  - `Toggle` switches the sounds on and off and keeps the choice in the browser.
- Demo: `/motion` and `/sound`.

- **`blocks/assistant`, more parts** (after the examples on beautifului.dev):
  - `Thinking` folds the agent's trace (reasoning, searches, reads, code) under "Thinking…" or "Thought
    for 4s". `ToolChips` folds tool calls into one line with a chip per call.
  - `Cite`, `Sources` and `FollowUps` go with the answer. A follow-up is sent with one click.
  - New block types for `BlockView`: `code` (line numbers, copy, unified diff), `context` (retrieved
    chunks with their source), `screen` (the agent's screen in a browser frame), `tasks` and `insights`
    (paged without a script, with a sparkline and a follow-up question).
  - Interactive blocks: a Go value holds the state, a method applies the form htmx posts and returns
    the next state, and a templ component draws it. `Recommendation.Decide` with `RecommendationCard`
    (confidence, Accept, alternatives), `Approval.Advance` with `ApprovalCard` (questions one at a time,
    Back, Skip, Continue; the answers travel in the form), `Diff.Resolve` with `DiffCard` (proposed edits,
    accepted or rejected per row or all at once), and `TaskListCard`, which polls while a task runs. The
    package has no handlers and the page loads htmx itself.
  - `Conversation`: a model picker (`Models`, posted as `model`), a toolbar for text selected in an
    answer (`SelectionActions`), and the working line counts seconds. Streamed fragments are passed to
    `htmx.process` when htmx is on the page.
  - New texts in `Texts`, English defaults.
  - Demo: `/assistant/parts` shows every part and runs the interactive ones; `/assistant` uses thinking,
    sources, follow-ups, the model picker, selection actions and a recommendation.

- **Less CSS per page:**
  - The blocks `listpage`, `detail`, `create`, `history` and `stats` render their CSS themselves, where
    they are used and once per page.
  - `shell.Styles` no longer carries their CSS. An empty shell page goes from 13.9 KB to 8.6 KB of inline
    CSS.
  - Their `Styles()` still work, e.g. to put the CSS in `<head>`.
  - `sidebar.Sidebar` renders the sidebar's CSS itself as well.
- **`static.WithKitCSS()`:**
  - Puts the kit's own CSS (sidebar, shell, blocks, icon sizing) into the bundle as `kit/kit.css`,
    cached and compressed like Tabler's CSS, and switches the inline CSS off.
  - `static.Head` links the file; `static.Kit(files)` links it on its own.
  - The file holds only the CSS of the packages the program links.

- **Sidebar:**
  - `sidebar.Switcher` puts a workspace (team, organisation) with a dropdown to switch in place of the
    brand.
  - `GroupProps.Collapsible` and `Open` make a group's label a button that folds its items away. With
    the sidebar collapsed to icons every group shows its items.
- **`blocks/split`:**
  - `Split` lays out a list beside its detail and an inspector on the right (mail, tickets, customers).
  - The panes scroll on their own and fill the window.
  - On small screens the page shows the list or, once an item is open, the detail with a back link.
  - `ListItem` is a row of the list; `AsideToggle` shows and hides the inspector.
- **`blocks/actionbar`:**
  - `Bar` keeps a long form's buttons at the bottom of the window.
  - It says "Unsaved changes" once a field changed and, with `Guard`, asks before the reader leaves the
    page.
- **Shell:** `Props.Tabs` puts tabs below the page heading; `Props.Compact` tightens cards, tables,
  forms, buttons and the sidebar for pages dense with data.
- **kit:** new texts `Back`, `ToggleDetails` and `UnsavedChanges` (English and German).
- **Demo:** `/enterprise`, a small help-desk app with all of the above (`?compact=1` for the compact mode).

- **`markdown`:** renders GitHub-flavoured Markdown in Tabler's prose container. Raw HTML and
  `javascript:` links are dropped, so user input and a language model's answers can be shown as they are.
- **`blocks/assistant`:** the blocks of an AI assistant, moved here from thinking-core's
  `examples/blocks-tabler`:
  - the conversation with composer and example prompts;
  - one turn per message, with steps, questions to the user (`Ask`) and result blocks (markdown, table,
    chart, stats);
  - a Markdown answer and the usage line.
  Everything is rendered on the server and streamed as server-sent events (`Open`, `Stream.Send`). A small
  client script puts each fragment in its slot and replaces a block sent again under its id. Texts can be
  changed per request (`WithTexts`). The demo plays back a scripted run at `/assistant`.

## v0.11.0

Landing pages: ten new sections in `blocks/marketing` and a demo with two complete landing pages, plus a
way to tear down charts in swapped content.

- **`window.tablerCharts.destroy(element)`** tears down the charts in an element before it is removed or
  replaced: the ApexCharts instances, their resize listeners and their `window.tabler_chart` entries, and
  charts still queued to be drawn. For content that is swapped repeatedly, e.g. streamed or polled panels.
- **Landing sections in `blocks/marketing`:** `Announcement` (a pill above a hero title, via the new
  `Announcement` field of `Hero` and `SideHero`), `FeatureCards` (2–4 columns, optionally linked), `Bento`
  (tiles of different widths, heights and colors), `Steps` (numbered, joined by a dashed line), `CTABox`
  (a colored box, centered or with media), `Plans` (plan cards with a monthly/yearly switch in pure CSS),
  `Team`, `Quote` (one big testimonial with stars), `Integrations` and `Contact` (ways to reach you and a
  form). The demo shows every landing block at `/landing` and two complete pages at `/landing/saas` and
  `/landing/app`.

## v0.10.1

- **Tabler 1.7.3** (from 1.7.1): fixes only, no markup changes. Datepicker no longer throws when a range
  start is clicked again; `data-bs-dismiss="offcanvas"` closes responsive offcanvas; table variants in dark
  mode; `.card-table` corners in `.table-responsive-*`; `.modal-header` padding; Tom Select styling; Autosize
  no longer moves the scroll position. Icons stay at 3.49.0, as Tabler pins them.
- **templ v0.3.1100** (from v0.3.1070), with upstream security fixes in attribute and JavaScript escaping.
  Rendered `style` attributes now end with a `;`.

## v0.10.0

Serving Tabler yourself, with nothing loaded from a CDN at runtime:
- **Package `static` and `cmd/static`:** Tabler and the icon webfont, cut down to the icons the app uses,
  with versioned addresses, a one-year cache and gzip/brotli.
- **`cdn.Base` and `cmd/vendor`:** the same for charts, plugins and theme fonts.
- **`cmd/new -self-host`:** a starter app set up this way.
- **CSP nonces:** every script and style carries templ's nonce, for a Content-Security-Policy without
  `'unsafe-inline'`.

**Breaking:** the URL constants of the plugin packages are functions now (see below).

- Self-hosting without a CDN:
  - New package `cdn`: `cdn.Base` (default jsDelivr) is where every file loaded in the browser comes from.
    This covers Tabler, its icons, ApexCharts, the plugins, the theme fonts and Plyr's icons.
  - New command `cmd/vendor` downloads the files of the groups an app needs into a folder, along with
    the fonts and images their stylesheets refer to.
  - `cmd/new -self-host` embeds them in the app.
  - The demo takes `-vendor <dir>`.
  - With package `static`, `cmd/vendor` writes into the same bundle and `cdn.Base = files.Dir("npm/")`
    serves the plugins and fonts versioned, cached and compressed like Tabler (new `Files.Dir`).
- **Breaking:** the URL constants are functions now, so they follow `cdn.Base`:
  - `<plugin>.LibraryURL()` and `StylesheetURL()`;
  - `chart.LibraryURL()`, `mapview.LibraryURL()`;
  - `assets.IconsURL()`;
  - `marketing.TypedLibraryURL()`.

  `icon.FilledFontURL` is empty by default; set it only to override `cdn.Base`. Theme fonts load
  Fontsource's `index.css`.
- Content-Security-Policy: every `<script>` and `<style>` the packages render carries templ's nonce
  (`templ.WithNonce`), not only `sidebar` and `blocks/detail`: `modal`, `chart`, `themes`, `icon`, the
  plugins, `confetti`, `assets`, the blocks (layout settings, tables, modals, maps, marketing, shell …) and
  the `cmd/new` app. Without a nonce the output is unchanged. See README, Content-Security-Policy.
- `invoice.PrintButton` prints through a script instead of an inline `onclick`.
- The demo binary (`demo`, 18 MB) is no longer in the repository; `.gitignore` keeps it and `new` out.
- New package `static`: serves an app's own static files fast, for apps that serve Tabler themselves
  instead of loading it from jsDelivr (#29).
  - `static.New(prefix, fsys)` is an `http.Handler`; `URL(name)` gives an address with a hash over the
    bundle in the path, sent with `Cache-Control: public, max-age=31536000, immutable`. Other versions
    get `no-cache`, addresses without version a short cache (or a year with `?v=<hash>`). ETag and 304, HEAD.
  - Text is compressed once in memory: gzip at once, brotli in the background (new dependency
    `github.com/andybalholm/brotli`, templ's already). woff2 is sent as it is.
  - `static.Head`, `Theme` and `Scripts` are `assets.Head`, `Theme` and `Scripts` from the bundle;
    `Scripts` and `Theme` carry templ's nonce.
- New command `cmd/static` (`go get -tool …/cmd/static`, `go tool static`): fetches Tabler's CSS and JS
  for package `static` and cuts the icon webfont down to the icons the app's packages use (in Go, new
  dependency `github.com/tdewolff/font` for the command only): 35 icons for a `cmd/new` app, 5 KB of font
  and 2 KB of CSS instead of 500 KB and 210 KB. `-check` fails when a used icon is missing.
- Measured in an app with Tabler, the icon subset and one font (phone at 4 Mbit/s): the first visit of a
  page went from 1.9 MB in 4.3 s to 185 KB in 0.7 s, a repeat visit to 12–17 KB in 0.2 s.

## v0.9.0

All of Tabler 1.7 in Go and templ. Every component of the Tabler docs is now a package (66 of 67), and the
composed parts of the preview are blocks. 143 of 146 preview pages are rebuilt in the demo and checked against
the original by screenshot.

New in this release:
- **Pages and starter kit:** dashboards, charts and app pages. `cmd/new` creates a running app.
- **Plugins:** packages for Tabler's plugins (date picker, select, editor, maps …).
- **Themes:** Tabler's theme settings and 13 ready-made themes.
- **Tabler updates:** pinned versions in one place and an update script.
- **Speed:** compressed pages, and charts that load quickly on phones.

Breaking changes are marked below; most change the look to Tabler 1.7's markup, a few change types
(`progress` values, `avatar.ShapeCircle`, `shell.Crumb`).

- `chart` draws much faster on phones. Measured at 4G with a 4× slower CPU: the dashboard's visible charts
  appear after 3.3 s instead of 13.3 s, and 36 charts no longer block a page for 51 s.
  - A chart is drawn when it comes near the screen, one at a time, with a frame for the browser between
    charts. Charts in hidden tabs and dialogs are drawn at their real width when shown.
  - The colors are resolved once per chart instead of once per use inside ApexCharts. Each of those uses
    restyled the whole page.
  - Charts follow theme changes (color mode, accent, settings panel) by redrawing.
  - `window.tablerCharts.flush()` draws all charts at once, e.g. before a screenshot; printing calls it.

- New package `compress`: gzip middleware for an app's pages (standard library only; streams, websockets
  and small or binary responses pass through). The demo and `cmd/new` apps use it; the demo's pages are
  5–6 times smaller on the wire (e.g. `/cards` 196 KB → 31 KB).

- New package `themes`: 13 ready-made looks in the spirit of shadcn/ui.
  - Tabler, Zinc, Slate, Ocean, Forest and Paper combine Tabler's own settings.
  - Rose, Violet, Amber, Emerald, Graphite, Sunset and Midnight bring Tailwind CSS palettes and Fontsource
    variable fonts.
  - `Theme.Apply`/`ApplyAuth` give a page a theme; `Head` loads their CSS and fonts and restores a user's
    pick; `Section()` puts them in the settings panel.
  - Your own themes are `themes.Theme` values.
  - The demo has a gallery at `/themes`, and `cmd/new -theme <name>` starts an app with a theme.
- `layout.PageProps.HTMLAttributes` and `auth.PageProps.HTMLAttributes` add attributes to `<html>`.
  `layout.ThemeOption` takes `SVG` and `Values` for any setting, `layout.ThemeSetting` a `Default`.
- `scripts/screenshots/pages.txt` lines take a viewport width (`themes 1600`).

- Theming: new package `theme` (accent color, gray base, font, corner radius as typed values).
  `layout.PageProps.Appearance` and `auth.PageProps.Appearance`/`Theme` render them, and
  `layout.PageProps.OffcanvasFloating` floats side panels. `assets.Props.Themes` loads `tabler-themes.css`.
- New block `layout.ThemeSettings`, Tabler's theme panel (also with only the appearance section, translated,
  or opened from any button via `ThemeSettingsToggle`). The demo shows it on every page, as the preview does.
- `cmd/new`: the app loads Tabler's theme script, so the dark/light toggle sticks (it did nothing before),
  and it has Appearance in the user menu. `views.Appearance` sets the app's default look.

- New package `versions`: Tabler, its icons and the plugin libraries are pinned in one place; every
  package's `Version` constant reads from it. `scripts/tabler-update.sh` prepares a Tabler update. It
  updates the pins (`internal/tools/versions`), regenerates the demo data (now also the menu,
  `internal/tools/menu`), and compares every page with Tabler's preview (`compare-all.mjs`).
- Demo: the sortable list loads SortableJS 1.15.7 like the `sortable` package, not 1.15.6. Two illustration
  images had the wrong size.

- **Breaking (look):** `form.Field`'s `Hint` renders Tabler 1.6+'s `form-text` instead of the deprecated
  `form-hint` (4px closer to the control). A required field's label carries Tabler's hidden " (required)"
  for screen readers by default (`kit.Messages.Required`; `RequiredLabel: "-"` for none).
- Plugin packages, each loading its library and starting itself (also after htmx swaps): `datepicker`,
  `selectbox`, `inputmask`, `rangeslider`, `wysiwyg`, `dropzone`, `colorpicker`, `signature`, `rating`,
  `calendar`, `datatable`, `sortable`, `tour`, `confetti`, `mapview`, `lightbox`, `player`. Form elements,
  form layout, modals and the mail compose dialog now use them; the plugin pages are rebuilt.
- New blocks `blocks/marketing` (the seven marketing pages) and `blocks/maps` (map cards).
  `pricing.CardProps.Content` replaces a card's feature list.

- New command `cmd/new`: creates a running app (frame, sign-in, dashboard with charts, list with search and
  pages, create page, detail page with edit dialog), with its own test. New package `assets`: Tabler's
  stylesheets and scripts from jsDelivr for a page's Head and Scripts.
- Fix: `blocks/detail`'s dialog script found no modal API with Tabler 1.x (it exposes `window.tabler`, not
  `window.bootstrap`), so dialogs did not reopen after a failed save.

Existing components follow Tabler 1.7's markup and cover the variants of Tabler's own components
(`shared/ui/*.astro`). The demo's buttons, alerts and pagination pages now match the preview.

- **Breaking (look):** `icon.Icon` adds Tabler's `.icon` class, so icons get Tabler's size (1.25rem) and
  spacing in buttons, badges and alerts. New: `Size`, `Color`, `Inline`, `icon.Classes`.
- **Breaking:** `alert.Alert` uses Tabler 1.7's markup (`.alert-icon`) and shows the variant's default icon
  when `Icon` is empty (set `NoIcon` for none; `shell.Flash` follows). The close control is a `<button>`.
  New: `Title` (heading with the children as description), `Important`, `Minor`, `Action`, `Link`, `List`.
- **Breaking:** pagination numbers carry `aria-label="Page n"`; disabled links get `tabindex="-1"` and
  `aria-disabled`. New: `Text`, `FirstLast`, `PrevTitle`, `NextTitle`.
- `button`: `Ghost`, `Pill`, `Square`, `Action`, `Block`, `Loading`, `IconEnd`, `Animate`, `SizeXLarge`;
  `List` and `Actions`. Disabled links get `tabindex="-1"`.
- `badge`: `Outline`, `Size`, `Pill`, `Icon`, `IconEnd`, `IconOnly`, `Dot`, `Notification`, `Blink`; `List`.
- `card`: `Size`, `Href`/`Hover`, `Borderless`, `Dashed`, `Transparent`, `Stacked`, `Active`, `Inactive`,
  `Status`/`StatusEdge`, `Light` header, `Scrollable` body; `Subtitle`, `Stamp`, `Progress`, `Image`.
- `kit.Messages`: new texts for pagination. Empty fields of an app's own `Messages` now fall back to English,
  so new texts never render empty.
- New packages: `dropdown`, `nav` (tabs, pills, overflow nav, tab panels, segmented control), `breadcrumb`,
  `offcanvas`, `toast`, `tooltip`, `popover`, `avatar`, `flag`; block `blocks/tabscard`.
- `blocks/shell`: `Crumb` is now `breadcrumb.Crumb` (same fields); the breadcrumbs render through package
  `breadcrumb`.
- New packages for showing data: `listgroup`, `progress`, `steps`, `status`, `timeline`, `datagrid`, `empty`,
  `ribbon`, `tracking`. `breadcrumb.Crumb` gets `Active` and `Disabled` (breadcrumbs as wizard steps).
- `form`: input groups and icons, password toggle, select groups (pills, boxes), image and color checks,
  floating labels, fieldsets, range, number input, OTP input and password strength meter; sizes, validation
  states (also lite) and more variants on `Input`, `Select`, `Textarea`, `Check`; `Field` gets `Description`,
  `Help` and `Valid`. Invalid inputs carry `aria-invalid`. `button.Group`, `button.Check`, `Floating`.
- `blocks/auth`: Tabler's sign-in and account pages (sign in in four variants, sign up, forgot password,
  account locked, 2-step verification, trial ended, onboarding). `blocks/signin` stays for existing apps.
  New packages `divider` and `iconlist`.
- `blocks/layout`: Tabler's page frame with navbar, sidebar, data-driven menu, navbar actions, page header and
  footer; the layout variants of the preview as props (data-bs-* on <html>). `blocks/errorpage`: error,
  maintenance and loader pages. The demo runs in Tabler's frame with the preview's own menu.
- `modal`: sizes xl and full-width, blur, status line, no header/close, sections, inline display, native
  `Dialog`; the modal gets `role="dialog"` and `aria-modal`. `blocks/modals`: Status, Confirm, ConfirmDelete,
  ChangePassword. `form.SelectItem` boxes with title and description.
- `icon`: filled icons next to outline ones (`Filled`, or names ending in `-filled`): the glyph comes from
  the filled webfont, which `icon.Styles()` loads (`icon.FilledFontURL`).
- New packages: `accordion`, `carousel`, `chat`, `hovergallery`, `imagelabel`, `legend`, `mention`, `mockup`,
  `placeholder`, `pricing`, `shape`, `shimmer`, `signal`, `sparkline`, `spinner`, `switchicon`, `tag`, `trending`,
  `typography`. `avatar.Props.Empty`, `listgroup` as `<ul>` or `<nav>`. `card`:
  gradients, rotation, image places, overlay, cover, card buttons, meta, code, dropdown; `tabscard.Tabs`.
- App and content pages: blocks `settings`, `plans`, `pricing`, `invoice`, `billing`, `legal`, `faq`, `profile`,
  `users`, `gallery`, `search`, `jobs`, `logs`, `cookie`, `cards`, `tables`; packages `payment`, `clipboard`, `prose`,
  `pattern`, `social`, `colorswatch`. `table` covers Tabler's table docs; `form` gets horizontal fields,
  fieldsets as fields, wrapped controls, check descriptions, datalists and more; `nav` overflow helpers and
  link segments; `sparkline` stroke and spot sizes.
- **Breaking:** `progress` values are `float64` (decimals such as 82.54); a `form.Field` without `For` renders its
  label as `<div class="form-label">` like Tabler; `.table-responsive` gets `tabindex="0"`.
- `card.Title` and `pricing.Card` take a heading level; `badge` as a link; `timeline` strong marker colors;
  `accordion` body class; `layout.Navbar` attributes and `layout.Brand` as a link; `typography.Subheader` as
  legend; `mention` as anchor; `form.ColorCheck` with any CSS color.
- `icon.Styles()` / `icon.CSS`: sizes webfont icons where Tabler's CSS only sizes SVG; `shell`, `signin` and
  `landing` include it. **Breaking:** `avatar.ShapeCircle` is now the round shape (`avatar-rounded`), the
  default is `ShapeDefault`.
- Demo: dropdowns, tabs, offcanvas, toasts, avatars, progress, steps, empty, number-input and password-strength pages, the authentication pages and some 40 app,
  content and component pages, checked against the
  preview.
- `scripts/screenshots/compare.mjs` screenshots a rebuilt page next to the Tabler original.
- New package `chart`: Tabler's ApexCharts setup ported to Go (all chart types of the preview, sparklines,
  annotations, brush, synced groups, legend toggles); one script per page starts every chart, also after
  htmx swaps. New block `blocks/metric` with the dashboards' metric, overview and chart cards. The demo
  loads `tabler-vendors.css` (ApexCharts tooltips) and rebuilds the charts pages and the legend charts.
- App pages with charts and apps: the dashboard (also inside every layout-* page), CRM and crypto
  dashboards, widgets, uptime, activity, tracking, tasks, task list, chat, mail inbox, emails, music, wizard.
  New blocks `blocks/dashboard`, `blocks/crm`, `blocks/crypto`, `blocks/activity`, `blocks/uptime`,
  `blocks/tasks`, `blocks/projects`, `blocks/messenger`, `blocks/mail`, `blocks/music`, `blocks/wizard`; more
  cards in `blocks/metric`, `blocks/users` (`Link`, `Profile`, `Cover`, `Follow`) and `blocks/cards` (`Code`,
  `Definitions`). New component `vectormap`. `chat`: `MessageProps.Col`, `Typing`. `chart`: `RawProps`,
  `Props.SerieColor`.

## v0.8.0

- **Breaking:** the module moved from `github.com/mind-vm/tablerkit` to `github.com/mind-vm/tabler-templ`
  (the repository was renamed). Replace the import path; the packages and their API are unchanged.
- `cmd/demo`: rebuilds the Tabler preview page by page (first pages: buttons, alerts, badges, pagination,
  modals, sign-in), with an index of all 146 preview pages.
- `ROADMAP.md`: checklist of every Tabler component, block and preview page, generated from `internal/catalog`.
- `scripts/screenshots.sh`: screenshots of the demo's pages for the README.
