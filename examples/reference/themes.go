package reference

import (
	"encoding/json"
	"strings"

	"github.com/mind-vm/tabler-templ/themes"
)

// The reference shows every theme of the module, chosen in one place (GS-15): Graphite by default, as in
// the products, the others to check that the patterns hold in each (signal colours stay, GS-07). The
// shell sets no attributes on <html>, so a small script in the head sets the theme's attributes (its name,
// its Tabler settings such as corners and grays, and its colour mode) before the page paints. Fonts come
// from the bundle (GS-09).

const defaultTheme = "graphite"

// themeStorage is the browser's key for the theme chosen; ?theme=<name> sets it too.
const themeStorage = "reference-theme"

func themeName(t themes.Theme) string {
	if t.Name == "" {
		return "tabler"
	}
	return t.Name
}

// themeFonts are the fonts of all themes, once each.
func themeFonts() []themes.Font {
	var fs []themes.Font
	seen := map[string]bool{}
	for _, t := range themes.All {
		for _, f := range []*themes.Font{t.Font, t.Mono} {
			if f != nil && !seen[f.Package] {
				seen[f.Package] = true
				fs = append(fs, *f)
			}
		}
	}
	return fs
}

// themeCSS are the rules of all themes; each applies only under its data-bs-theme-pack.
func themeCSS() string {
	var b strings.Builder
	for _, t := range themes.All {
		b.WriteString(t.CSS())
	}
	return b.String()
}

type themeAttrs struct {
	Attrs map[string]any `json:"attrs"`
	Mode  string         `json:"mode,omitempty"`
}

// themeScript sets the chosen theme on <html> and wires the switcher (data-ref-theme).
func themeScript() string {
	m := map[string]themeAttrs{}
	for _, t := range themes.All {
		m[themeName(t)] = themeAttrs{Attrs: t.Attrs(), Mode: t.Mode}
	}
	cfg, _ := json.Marshal(m) // escapes <, > and &
	return `(function(){"use strict";var T=` + string(cfg) + `,K="` + themeStorage + `",D="` + defaultTheme + `",d=document.documentElement;
function apply(n){var t=T[n]||T[D];Array.from(d.attributes).forEach(function(a){if(a.name.indexOf("data-bs-theme")===0)d.removeAttribute(a.name)});
Object.keys(t.attrs).forEach(function(k){d.setAttribute(k,t.attrs[k])});if(t.mode)d.setAttribute("data-bs-theme",t.mode);return T[n]?n:D}
var n=null,p=new URLSearchParams(location.search).get("theme");try{if(p&&T[p])localStorage.setItem(K,p);n=localStorage.getItem(K)}catch(e){}
var cur=apply(n||D);
document.addEventListener("DOMContentLoaded",function(){document.querySelectorAll("[data-ref-theme]").forEach(function(s){s.value=cur;
s.addEventListener("change",function(){cur=apply(s.value);try{localStorage.setItem(K,cur)}catch(e){}})})})})();`
}
