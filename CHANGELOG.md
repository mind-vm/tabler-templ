# Changelog

## Unreleased

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
- `icon.Styles()` / `icon.CSS`: sizes webfont icons where Tabler's CSS only sizes SVG; `shell`, `signin` and
  `landing` include it. **Breaking:** `avatar.ShapeCircle` is now the round shape (`avatar-rounded`), the
  default is `ShapeDefault`.
- Demo: dropdowns, tabs, offcanvas, toasts, avatars, progress, steps, empty, number-input and password-strength pages, and all authentication pages but license and
  terms-of-service, checked against the
  preview.
- `scripts/screenshots/compare.mjs` screenshots a rebuilt page next to the Tabler original.

## v0.8.0

- **Breaking:** the module moved from `github.com/mind-vm/tablerkit` to `github.com/mind-vm/tabler-templ`
  (the repository was renamed). Replace the import path; the packages and their API are unchanged.
- `cmd/demo`: rebuilds the Tabler preview page by page (first pages: buttons, alerts, badges, pagination,
  modals, sign-in), with an index of all 146 preview pages.
- `ROADMAP.md`: checklist of every Tabler component, block and preview page, generated from `internal/catalog`.
- `scripts/screenshots.sh`: screenshots of the demo's pages for the README.
