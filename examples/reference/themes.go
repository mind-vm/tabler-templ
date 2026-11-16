package reference

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mind-vm/tabler-templ/themes"
)

// The reference shows every theme of the module, chosen in one place (GS-15): Graphite by default, as in
// the products, the others to check that the patterns hold in each (signal colours stay, GS-07). The
// shell sets no attributes on <html>, so a small script in the head sets the theme's attributes (its name,
// its Tabler settings such as corners and grays, and its colour mode) before the page paints. Fonts come
// from the bundle (GS-09).

const defaultTheme = "graphite"

// themeStorage is the cookie of the theme chosen; ?theme=<name> sets it too.
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
		for _, f := range fontsOf(t) {
			if !seen[f.Package] {
				seen[f.Package] = true
				fs = append(fs, f)
			}
		}
	}
	return fs
}

func fontsOf(t themes.Theme) []themes.Font {
	var fs []themes.Font
	for _, f := range []*themes.Font{t.Font, t.Mono} {
		if f != nil {
			fs = append(fs, *f)
		}
	}
	return fs
}

func themeNamed(name string) (themes.Theme, bool) {
	for _, t := range themes.All {
		if themeName(t) == name {
			return t, true
		}
	}
	return themes.Theme{}, false
}

// themeOf is the theme of the request: ?theme=<name> (which it also remembers), else the cookie, else
// Graphite. The server knows it, so the page loads only that theme's fonts (GS-09).
func themeOf(w http.ResponseWriter, r *http.Request) string {
	if q := r.URL.Query().Get("theme"); q != "" {
		if _, ok := themeNamed(q); ok {
			http.SetCookie(w, &http.Cookie{Name: themeStorage, Value: q, Path: Home, MaxAge: 31536000, SameSite: http.SameSiteLaxMode})
			return q
		}
	}
	if c, err := r.Cookie(themeStorage); err == nil {
		if _, ok := themeNamed(c.Value); ok {
			return c.Value
		}
	}
	return defaultTheme
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
	Fonts []string       `json:"fonts,omitempty"`
}

// themeScript sets the theme on <html> before the page paints and wires the switcher (data-ref-theme): a
// new choice goes into the cookie, and the new theme's fonts are added to the page when it has none yet.
func themeScript(current string, fontURLs map[string][]string) string {
	m := map[string]themeAttrs{}
	for _, t := range themes.All {
		m[themeName(t)] = themeAttrs{Attrs: t.Attrs(), Mode: t.Mode, Fonts: fontURLs[themeName(t)]}
	}
	cfg, _ := json.Marshal(m) // escapes <, > and &
	cur, _ := json.Marshal(current)
	return `(function(){"use strict";var T=` + string(cfg) + `,K="` + themeStorage + `",d=document.documentElement,cur=` + string(cur) + `;
function apply(n){var t=T[n];if(!t)return;Array.from(d.attributes).forEach(function(a){if(a.name.indexOf("data-bs-theme")===0)d.removeAttribute(a.name)});
Object.keys(t.attrs).forEach(function(k){d.setAttribute(k,t.attrs[k])});if(t.mode)d.setAttribute("data-bs-theme",t.mode)}
function fonts(n){(T[n].fonts||[]).forEach(function(u){if(!document.querySelector('link[data-ref-font][href="'+u+'"]')){var l=document.createElement("link");l.rel="stylesheet";l.href=u;l.setAttribute("data-ref-font","");document.head.appendChild(l)}})}
apply(cur);
document.addEventListener("DOMContentLoaded",function(){document.querySelectorAll("[data-ref-theme]").forEach(function(s){s.value=cur;
s.addEventListener("change",function(){if(!T[s.value])return;cur=s.value;fonts(cur);apply(cur);document.cookie=K+"="+cur+"; path=/reference; max-age=31536000; samesite=lax"})})})})();`
}
