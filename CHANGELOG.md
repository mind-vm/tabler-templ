# Changelog

## v0.8.0

- **Breaking:** the module moved from `github.com/mind-vm/tablerkit` to `github.com/mind-vm/tabler-templ`
  (the repository was renamed). Replace the import path; the packages and their API are unchanged.
- `cmd/demo`: rebuilds the Tabler preview page by page (first pages: buttons, alerts, badges, pagination,
  modals, sign-in), with an index of all 146 preview pages.
- `ROADMAP.md`: checklist of every Tabler component, block and preview page, generated from `internal/catalog`.
- `scripts/screenshots.sh`: screenshots of the demo's pages for the README.
