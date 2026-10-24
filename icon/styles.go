package icon

import "github.com/mind-vm/tabler-templ/cdn"

// CSS sizes webfont icons where Tabler's CSS sizes icons through width and height only, which works for
// its inline SVG icons but not for a font glyph: avatars (and avatar uploads), card stamps, shapes, empty
// states, page titles and segmented controls. Render Styles (or serve CSS) after Tabler's stylesheet.
const CSS = `.avatar .ti.icon{font-size:var(--tblr-avatar-icon-size)}` +
	`.card-stamp-icon .ti.icon{font-size:calc(var(--tblr-stamp-size) * .75)}` +
	`.shape .ti.icon{font-size:var(--tblr-shape-icon-size)}` +
	`.empty-icon .ti.icon{width:100%;height:100%;font-size:var(--tblr-empty-icon-size)}` +
	`.page-title .ti.icon{width:1.5rem;height:1.5rem;font-size:1.5rem;margin-inline-end:.25rem}` +
	`.nav-segmented .ti.nav-link-icon{font-size:var(--tblr-nav-link-icon-size)}`

// FilledFontURL is where the filled webfont is loaded from: on cdn.Base, unless set to a URL of its own.
var FilledFontURL = ""

func filledFontURL() string {
	if FilledFontURL != "" {
		return FilledFontURL
	}
	return cdn.URL("@tabler/icons-webfont", FilledVersion, "dist/fonts/tabler-icons-filled.woff2")
}

// FilledCSS loads the filled webfont for icons with Filled. The webfont's own stylesheet cannot be used
// next to the outline one: it redefines every .ti icon.
func FilledCSS() string {
	return `@font-face{font-family:"tabler-icons-filled";font-style:normal;font-weight:400;font-display:block;src:url("` + filledFontURL() + `") format("woff2")}` +
		`.ti.ti-filled{font-family:"tabler-icons-filled"!important;font-style:normal}`
}
