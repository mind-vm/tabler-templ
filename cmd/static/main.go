// Command static fetches Tabler for an app that serves it itself (package static): Tabler's CSS and JS
// and the icon webfont, in the versions this module is built against, cut down to the icons the app uses.
// The whole webfont is 500 KB of woff2 and 210 KB of CSS for 5,000 icons; an app with 50 icons needs about
// 10 KB.
//
//	go get -tool github.com/mind-vm/tabler-templ/cmd/static  # at the app's tabler-templ version
//	go tool static [-out static] [-css tabler,tabler-vendors] [packages]
//	go tool static -check [-out static] [packages]
//
// It writes the layout package static expects (static.CoreDir, static.IconsCSS):
//
//	<out>/core/css/tabler.min.css …, core/js/tabler.min.js, core/js/tabler-theme.min.js
//	<out>/icons/tabler-icons.min.css, icons/fonts/tabler-icons.woff2, icons/LICENSE
//
// The icons are the ti-… classes named in the packages (default ./...) and everything they import,
// tabler-templ's components included: their .go files (without generated ones), .templ files and
// embedded .html and .js files. -dir adds directories of other files (templates, scripts), -keep icons
// that are only put together at run time ("ti-" + name). Run it again after using a new icon; -check
// fails when a used icon is missing from the subset, e.g. in CI or a test. -icons=all keeps the whole
// webfont. Downloads come from jsDelivr.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/mind-vm/tabler-templ/static"
	"github.com/mind-vm/tabler-templ/versions"
)

var (
	coreURL  = "https://cdn.jsdelivr.net/npm/@tabler/core@" + versions.Tabler + "/"
	iconsURL = "https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@" + versions.Icons + "/"
)

// stylesheets are the names -css takes: Tabler's stylesheets that assets.Props switches on.
var stylesheets = []string{"tabler", "tabler-flags", "tabler-payments", "tabler-socials", "tabler-marketing",
	"tabler-themes", "tabler-vendors"}

func main() {
	out := flag.String("out", "static", "directory to write to: the root of the bundle (static.New), which may hold the app's own files too")
	cssList := flag.String("css", "tabler", "Tabler stylesheets to fetch, comma-separated: "+strings.Join(stylesheets, ", "))
	rtl := flag.Bool("rtl", false, "also fetch the right-to-left stylesheets")
	icons := flag.String("icons", "used", `"used": only the icons the packages use; "all": the whole webfont`)
	dirs := flag.String("dir", "", "more directories to scan for ti-… classes, comma-separated (templates, scripts)")
	keep := flag.String("keep", "", "icons to keep in any case, comma-separated, e.g. ti-user,ti-star")
	check := flag.Bool("check", false, "fetch nothing; fail when a used icon is missing from <out>")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: static [flags] [packages]")
		flag.PrintDefaults()
	}
	flag.Parse()
	pkgs := flag.Args()
	if len(pkgs) == 0 {
		pkgs = []string{"./..."}
	}
	if err := run(*out, split(*cssList), *rtl, *icons, pkgs, split(*dirs), split(*keep), *check); err != nil {
		fmt.Fprintln(os.Stderr, "static:", err)
		os.Exit(1)
	}
}

func split(s string) []string {
	var r []string
	for _, x := range strings.Split(s, ",") {
		if x = strings.TrimSpace(x); x != "" {
			r = append(r, x)
		}
	}
	return r
}

func run(out string, css []string, rtl bool, icons string, pkgs, dirs, keep []string, check bool) error {
	if icons != "used" && icons != "all" {
		return fmt.Errorf("-icons is used or all, not %q", icons)
	}
	var used map[string]bool
	if icons == "used" {
		var err error
		if used, err = usedIcons(pkgs, dirs, out); err != nil {
			return err
		}
		for _, k := range keep {
			used[k] = true
		}
	}
	if check {
		if used == nil {
			return nil
		}
		b, err := os.ReadFile(filepath.Join(out, static.IconsCSS))
		if err != nil {
			return err
		}
		if allIcons(string(b)) {
			return nil
		}
		full, err := fetch(iconsURL + "dist/tabler-icons.min.css") // tells icons from words like "ti-…" in texts
		if err != nil {
			return err
		}
		if missing := missingIcons(used, glyphs(string(full)), glyphs(string(b))); len(missing) > 0 {
			return fmt.Errorf("used but not in %s, run cmd/static again: %s", out, strings.Join(missing, " "))
		}
		return nil
	}

	var files []string
	for _, name := range css {
		if !contains(stylesheets, name) {
			return fmt.Errorf("unknown stylesheet %q (one of %s)", name, strings.Join(stylesheets, ", "))
		}
		files = append(files, "css/"+name+".min.css")
		if rtl {
			files = append(files, "css/"+name+".rtl.min.css")
		}
	}
	files = append(files, "js/tabler.min.js", "js/tabler-theme.min.js")
	for _, f := range files {
		b, err := fetch(coreURL + "dist/" + f)
		if err != nil {
			return err
		}
		if err := write(filepath.Join(out, static.CoreDir, f), b); err != nil {
			return err
		}
	}
	// Tabler's own files carry their license in a comment; the icon font's is a file next to it.
	license, err := fetch(iconsURL + "LICENSE")
	if err != nil {
		return err
	}
	if err := write(filepath.Join(out, "icons/LICENSE"), license); err != nil {
		return err
	}

	cssB, err := fetch(iconsURL + "dist/tabler-icons.min.css")
	if err != nil {
		return err
	}
	fontB, err := fetch(iconsURL + "dist/fonts/tabler-icons.woff2")
	if err != nil {
		return err
	}
	var n int
	if used == nil {
		cssB, n = []byte(fontOnlyWOFF2(string(cssB))), len(glyphs(string(cssB)))
	} else {
		var s string
		s, fontB, n, err = subset(string(cssB), fontB, used)
		if err != nil {
			return err
		}
		cssB = []byte(s)
	}
	if err := write(filepath.Join(out, static.IconsCSS), cssB); err != nil {
		return err
	}
	if err := write(filepath.Join(out, "icons/fonts/tabler-icons.woff2"), fontB); err != nil {
		return err
	}
	fmt.Printf("Tabler %s and %d icons (%s font, %s CSS) in %s\n", versions.Tabler, n, kb(len(fontB)), kb(len(cssB)), out)
	return nil
}

// missingIcons are the used icons of the webfont (all) that are not in the subset (have), sorted.
func missingIcons(used map[string]bool, all, have map[string]rune) []string {
	var missing []string
	for k := range used {
		_, icon := all[k]
		if _, ok := have[k]; icon && !ok {
			missing = append(missing, k)
		}
	}
	sort.Strings(missing)
	return missing
}

// allIcons reports whether css is the whole webfont (written with -icons=all); -check passes then.
func allIcons(css string) bool { return !strings.Contains(css, subsetMarker) }

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func kb(n int) string { return fmt.Sprintf("%.1f KB", float64(n)/1024) }

var client = &http.Client{Timeout: time.Minute}

func fetch(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func write(name string, b []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	return os.WriteFile(name, b, 0o644)
}
