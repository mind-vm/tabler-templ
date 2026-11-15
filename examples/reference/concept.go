package reference

import (
	"embed"
	"html"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/mind-vm/tabler-templ/markdown"
)

// The concepts the reference shows, as an embedded copy of mind-vm/auftragsbutler (docs/konzepte). The
// original stays the source of truth (GS-03, GSO-01); scripts/reference-konzept.sh copies it again and
// writes the commit into konzept/QUELLE. The pages link both ways: the explanation of a pattern leads to
// its rules here, and every rule number and example here that a page of the reference shows leads to it.

//go:embed konzept/*.md konzept/QUELLE
var conceptFiles embed.FS

const (
	pathConcept          = Home + "/konzept"
	pathConceptAssistant = pathConcept + "/assistent"
	conceptRepo          = "https://github.com/mind-vm/auftragsbutler/blob/"
)

type concept struct {
	File  string // in konzept/ and in docs/konzepte of the source
	Path  string // its page here
	Title string
}

var concepts = []concept{
	{"gestaltung.md", pathConcept, "Gestaltung"},
	{"assistent.md", pathConceptAssistant, "Assistent"},
}

// source is the commit and date of the copy (konzept/QUELLE).
func source() (sha, date string) {
	b, _ := conceptFiles.ReadFile("konzept/QUELLE")
	sha, date, _ = strings.Cut(strings.TrimSpace(string(b)), " ")
	return sha, date
}

func sourceURL(file string) string {
	sha, _ := source()
	return conceptRepo + sha + "/docs/konzepte/" + file
}

// conceptLink is where a rule's explanation leads: the row of the rule on its concept page.
func conceptLink(rule string) string {
	if strings.HasPrefix(rule, "KA-") || strings.HasPrefix(rule, "KAO-") {
		return pathConceptAssistant + "#" + strings.ToLower(rule)
	}
	return pathConcept + "#" + strings.ToLower(rule)
}

var (
	rendered     = map[string]string{}
	renderedOnce sync.Once
)

// conceptHTML is the concept rendered once, with anchors and links.
func conceptHTML(c concept) string {
	renderedOnce.Do(func() {
		for _, c := range concepts {
			src, _ := conceptFiles.ReadFile("konzept/" + c.File)
			rendered[c.File] = enrich(markdown.HTML(string(src)), c)
		}
	})
	return rendered[c.File]
}

var (
	headingRE = regexp.MustCompile(`<h([1-6])>(.*?)</h[1-6]>`)
	hrefRE    = regexp.MustCompile(`href="([^"]*)"`)
	ruleRowRE = regexp.MustCompile(`<td>((?:GS|GSO|KA|KAO)-\d+)</td>`)
	tableRE   = regexp.MustCompile(`<table>`)
	tagRE     = regexp.MustCompile(`<[^>]+>`)
	refRE     = regexp.MustCompile(`\b(?:B-0\d|(?:GS|KA)-\d+)\b`)
)

func enrich(h string, c concept) string {
	// Headings get GitHub's anchors, so links into the original work here too.
	h = headingRE.ReplaceAllStringFunc(h, func(m string) string {
		p := headingRE.FindStringSubmatch(m)
		return `<h` + p[1] + ` id="` + slug(html.UnescapeString(tagRE.ReplaceAllString(p[2], ""))) + `">` + p[2] + `</h` + p[1] + `>`
	})
	// The row of a rule gets the rule's anchor (#gs-41).
	h = ruleRowRE.ReplaceAllString(h, `<td id="${1}">${1}</td>`)
	h = regexp.MustCompile(`<td id="([A-Z]+-\d+)">`).ReplaceAllStringFunc(h, strings.ToLower)
	h = strings.ReplaceAll(h, `<td id="`, `<td class="ref-rule-cell" id="`)
	// Relative links: the other concept shown here, everything else on GitHub at the copy's commit.
	h = hrefRE.ReplaceAllStringFunc(h, func(m string) string {
		return `href="` + html.EscapeString(resolve(html.UnescapeString(hrefRE.FindStringSubmatch(m)[1]))) + `"`
	})
	// Wide tables scroll in their box (GS-14, GS-54).
	h = tableRE.ReplaceAllString(h, `<div class="table-responsive"><table class="table table-vcenter ref-concept-table">`)
	h = strings.ReplaceAll(h, `</table>`, `</table></div>`)
	return backlinks(h)
}

func resolve(href string) string {
	if href == "" || strings.HasPrefix(href, "#") || strings.Contains(href, "://") || strings.HasPrefix(href, "mailto:") {
		return href
	}
	path, frag, _ := strings.Cut(href, "#")
	if frag != "" {
		frag = "#" + frag
	}
	for _, c := range concepts {
		if path == c.File {
			return c.Path + frag
		}
	}
	base, _ := url.Parse(sourceURL("x"))
	ref, err := url.Parse(path)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String() + frag
}

// backlinks: in the text (not in tags, links, code or headings), the examples B-01 … and every rule a page
// of the reference shows link to that page.
func backlinks(h string) string {
	shows := ruleTargets()
	var b strings.Builder
	skip := 0
	for len(h) > 0 {
		i := strings.IndexByte(h, '<')
		if i < 0 {
			i = len(h)
		}
		text := h[:i]
		if skip == 0 {
			text = refRE.ReplaceAllStringFunc(text, func(r string) string {
				if p, ok := shows[r]; ok {
					return `<a href="` + p.Path + `" class="ref-backlink" title="In der Referenz: ` + p.ID + ` ` + p.Name + `">` + r + `</a>`
				}
				return r
			})
		}
		b.WriteString(text)
		h = h[i:]
		if h == "" {
			break
		}
		j := strings.IndexByte(h, '>')
		if j < 0 {
			b.WriteString(h)
			break
		}
		tag := h[:j+1]
		name := strings.ToLower(strings.TrimLeft(strings.Fields(strings.Trim(tag, "<>/"))[0], "/"))
		if slices.Contains([]string{"a", "code", "pre", "h1", "h2", "h3", "h4", "h5", "h6"}, name) {
			if strings.HasPrefix(tag, "</") {
				skip--
			} else {
				skip++
			}
		}
		b.WriteString(tag)
		h = h[j+1:]
	}
	return b.String()
}

// ruleTargets maps an example (B-04) and each rule a page names to that page.
func ruleTargets() map[string]pattern {
	m := map[string]pattern{}
	for _, p := range patterns {
		m[p.ID] = p
		for _, r := range p.Rules {
			if _, ok := m[r]; !ok {
				m[r] = p
			}
		}
	}
	return m
}

// slug is GitHub's anchor of a heading: lower case, letters, digits, '-' and '_' kept, spaces to '-'.
func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}
