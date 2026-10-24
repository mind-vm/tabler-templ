package icon

// CSS sizes webfont icons where Tabler's CSS sizes icons through width and height only, which works for
// its inline SVG icons but not for a font glyph: avatars (and avatar uploads), card stamps, shapes, empty
// states and page titles. Render Styles (or serve CSS) after Tabler's stylesheet.
const CSS = `.avatar .ti.icon{font-size:var(--tblr-avatar-icon-size)}` +
	`.card-stamp-icon .ti.icon{font-size:calc(var(--tblr-stamp-size) * .75)}` +
	`.shape .ti.icon{font-size:var(--tblr-shape-icon-size)}` +
	`.empty-icon .ti.icon{width:100%;height:100%;font-size:var(--tblr-empty-icon-size)}` +
	`.page-title .ti.icon{width:1.5rem;height:1.5rem;font-size:1.5rem;margin-inline-end:.25rem}`
