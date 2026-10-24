# Changelog

## Unreleased

- Self-hosting without a CDN:
  - New package `cdn`: `cdn.Base` (default jsDelivr) is where every file loaded in the browser comes from.
    This covers Tabler, its icons, ApexCharts, the plugins, the theme fonts and Plyr's icons.
  - New command `cmd/vendor` downloads the files of the groups an app needs into a folder, along with
    the fonts and images their stylesheets refer to.
  - `cmd/new -self-host` embeds them in the app.
  - The demo takes `-vendor <dir>`.
- **Breaking:** the URL constants are functions now, so they follow `cdn.Base`:
  - `<plugin>.LibraryURL()` and `StylesheetURL()`;
  - `chart.LibraryURL()`, `mapview.LibraryURL()`;
  - `assets.IconsURL()`;
  - `marketing.TypedLibraryURL()`.

  `icon.FilledFontURL` is empty by default; set it only to override `cdn.Base`. Theme fonts load
  Fontsource's `index.css`.

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
