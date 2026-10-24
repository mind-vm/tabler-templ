package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tdewolff/font"
)

// subsetMarker starts the comment of a subset stylesheet; without it the stylesheet is the whole webfont.
const subsetMarker = "/* Subset of @tabler/icons-webfont"

var (
	iconClass = regexp.MustCompile(`\bti-[a-z0-9]+(?:-[a-z0-9]+)*`)
	iconRule  = regexp.MustCompile(`\.(ti-[a-z0-9-]+):before\{content:"\\([0-9a-f]+)"\}`)
	fontSrc   = regexp.MustCompile(`src:url\([^}]*`)
	generated = regexp.MustCompile(`(?m)^// Code generated .* DO NOT EDIT\.$`)
)

// glyphs maps the icon classes of the webfont's stylesheet to their code points.
func glyphs(css string) map[string]rune {
	m := map[string]rune{}
	for _, r := range iconRule.FindAllStringSubmatch(css, -1) {
		cp, err := strconv.ParseUint(r[2], 16, 32)
		if err == nil {
			m[r[1]] = rune(cp)
		}
	}
	return m
}

// fontOnlyWOFF2 points the @font-face at the woff2 file only and drops the version query: the version is
// in the path (package static), and every browser that runs Tabler 1.x reads woff2.
func fontOnlyWOFF2(css string) string {
	return fontSrc.ReplaceAllString(css, `src:url("./fonts/tabler-icons.woff2") format("woff2")`)
}

// subset cuts the webfont (its stylesheet and woff2) down to the used icons; names in used that are no
// icon (words like "ti-…" in texts) are left out. It returns the stylesheet, the font and the number of
// icons.
func subset(css string, woff2 []byte, used map[string]bool) (string, []byte, int, error) {
	all := glyphs(css)
	var names []string
	for k := range used {
		if _, ok := all[k]; ok {
			names = append(names, k)
		}
	}
	if len(names) == 0 {
		return "", nil, 0, fmt.Errorf("no icons found")
	}
	sort.Strings(names)
	i := strings.Index(css, ".ti-")
	if i < 0 {
		return "", nil, 0, fmt.Errorf("no icons in the stylesheet")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s (MIT), %d icons, written by tabler-templ's cmd/static; do not edit */\n", subsetMarker, len(names))
	b.WriteString(fontOnlyWOFF2(css[:i]))
	cps := make([]rune, 0, len(names))
	for _, n := range names {
		fmt.Fprintf(&b, `.%s:before{content:"\%x"}`, n, all[n])
		cps = append(cps, all[n])
	}
	b.WriteString("\n")
	f, err := subsetFont(woff2, cps)
	if err != nil {
		return "", nil, 0, err
	}
	return b.String(), f, len(names), nil
}

// subsetFont keeps the glyphs of the code points (and .notdef) in a woff2 font.
func subsetFont(woff2 []byte, cps []rune) ([]byte, error) {
	b, err := font.ToSFNT(woff2)
	if err != nil {
		return nil, err
	}
	f, err := font.ParseSFNT(b, 0)
	if err != nil {
		return nil, err
	}
	ids := []uint16{0}
	seen := map[uint16]bool{0: true}
	for _, cp := range cps {
		id := f.GlyphIndex(cp)
		if id == 0 {
			return nil, fmt.Errorf("no glyph for U+%X", cp)
		}
		if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	s, err := f.Subset(ids, font.SubsetOptions{Tables: font.KeepMinTables})
	if err != nil {
		return nil, err
	}
	return s.WriteWOFF2()
}

// usedIcons collects the ti-… classes of the packages (go list patterns) with everything they import, and
// of the files in dirs. The output directory is skipped: its stylesheet names the icons of the last run.
func usedIcons(pkgs, dirs []string, out string) (map[string]bool, error) {
	used := map[string]bool{}
	skip, _ := filepath.Abs(out)
	scan := func(name string) error {
		if a, _ := filepath.Abs(name); a == skip || strings.HasPrefix(a, skip+string(filepath.Separator)) {
			return nil
		}
		b, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		if strings.HasSuffix(name, ".go") && generated.Match(b) {
			return nil // templ's _templ.go (the .templ is scanned) and generated data
		}
		for _, m := range iconClass.FindAll(b, -1) {
			used[string(m)] = true
		}
		return nil
	}
	cmd := exec.Command("go", append([]string{"list", "-deps", "-json=Dir,Standard,GoFiles,EmbedFiles"}, pkgs...)...)
	cmd.Stderr = os.Stderr
	listed, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(listed))
	for dec.More() {
		var p struct {
			Dir                 string
			Standard            bool
			GoFiles, EmbedFiles []string
		}
		if err := dec.Decode(&p); err != nil {
			return nil, err
		}
		if p.Standard {
			continue
		}
		names := p.GoFiles
		templs, _ := filepath.Glob(filepath.Join(p.Dir, "*.templ"))
		for _, t := range templs {
			names = append(names, filepath.Base(t))
		}
		for _, e := range p.EmbedFiles {
			if ext := filepath.Ext(e); ext == ".html" || ext == ".js" || ext == ".mjs" || ext == ".templ" {
				names = append(names, e)
			}
		}
		for _, n := range names {
			if err := scan(filepath.Join(p.Dir, n)); err != nil {
				return nil, err
			}
		}
	}
	for _, d := range dirs {
		err := filepath.WalkDir(d, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if e.IsDir() {
				if n := e.Name(); n == "node_modules" || n == "vendor" || strings.HasPrefix(n, ".") && p != d {
					return filepath.SkipDir
				}
				return nil
			}
			switch filepath.Ext(p) {
			case ".html", ".js", ".mjs", ".templ", ".go", ".tmpl", ".gohtml":
				if !strings.HasSuffix(p, "_test.go") {
					return scan(p)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return used, nil
}
