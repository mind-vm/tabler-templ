// Command vendor downloads the files tabler-templ loads in the browser (Tabler, its icons, ApexCharts, the
// plugins, the theme fonts) into a folder, in the layout cdn.URL builds: <package>@<version>/<path>. Serve
// the folder (or embed it) and point cdn.Base at it, and the app loads nothing from outside at runtime:
//
//	go run github.com/mind-vm/tabler-templ/cmd/vendor -out static/npm tabler chart theme:graphite
//	cdn.Base = "/static/npm/" // in the app
//
// The arguments are groups (-list shows them); "all" is every group. Files the stylesheets refer to (fonts,
// flags, images) are downloaded as well. Run it again after updating tabler-templ.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/mind-vm/tabler-templ/cdn"

	// Every package that registers files.
	_ "github.com/mind-vm/tabler-templ/assets"
	_ "github.com/mind-vm/tabler-templ/blocks/marketing"
	_ "github.com/mind-vm/tabler-templ/calendar"
	_ "github.com/mind-vm/tabler-templ/chart"
	_ "github.com/mind-vm/tabler-templ/colorpicker"
	_ "github.com/mind-vm/tabler-templ/datatable"
	_ "github.com/mind-vm/tabler-templ/datepicker"
	_ "github.com/mind-vm/tabler-templ/dropzone"
	_ "github.com/mind-vm/tabler-templ/inputmask"
	_ "github.com/mind-vm/tabler-templ/lightbox"
	_ "github.com/mind-vm/tabler-templ/mapview"
	_ "github.com/mind-vm/tabler-templ/player"
	_ "github.com/mind-vm/tabler-templ/rangeslider"
	_ "github.com/mind-vm/tabler-templ/rating"
	_ "github.com/mind-vm/tabler-templ/selectbox"
	_ "github.com/mind-vm/tabler-templ/signature"
	_ "github.com/mind-vm/tabler-templ/sortable"
	_ "github.com/mind-vm/tabler-templ/themes"
	_ "github.com/mind-vm/tabler-templ/tour"
	_ "github.com/mind-vm/tabler-templ/vectormap"
	_ "github.com/mind-vm/tabler-templ/wysiwyg"
)

var client = &http.Client{Timeout: 60 * time.Second}

func main() {
	out := flag.String("out", "static/npm", "folder to write the files to")
	from := flag.String("from", cdn.JSDelivr, "where to download from (an npm CDN with jsDelivr's layout)")
	list := flag.Bool("list", false, "list the groups and their files")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: vendor [-out dir] group … (tabler, chart, theme:rose, … or all; -list shows them)")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *list {
		for _, g := range cdn.Groups() {
			fmt.Println(g)
			for _, f := range cdn.Registered() {
				if f.Group == g {
					fmt.Printf("\t%s@%s: %s\n", f.Package, f.Version, strings.Join(f.Paths, " "))
				}
			}
		}
		return
	}
	groups := flag.Args()
	if len(groups) == 0 {
		groups = []string{"tabler"}
	}
	all := slices.Contains(groups, "all")
	for _, g := range groups {
		if g != "all" && !slices.Contains(cdn.Groups(), g) {
			log.Fatalf("no group %q; vendor -list shows them", g)
		}
	}

	v := &vendor{out: *out, from: strings.TrimSuffix(*from, "/") + "/", done: map[string]bool{}, listings: map[string][]string{}}
	for _, f := range cdn.Registered() {
		if !all && !slices.Contains(groups, f.Group) {
			continue
		}
		for _, p := range f.Paths {
			for _, file := range v.expand(f.Package, f.Version, p) {
				v.fetch(f.Package+"@"+f.Version, file)
			}
		}
	}
	if v.failed > 0 {
		log.Fatalf("%d files failed", v.failed)
	}
	log.Printf("%d files in %s; set cdn.Base to the URL path it is served at", len(v.done), *out)
}

type vendor struct {
	out, from string
	done      map[string]bool
	listings  map[string][]string
	failed    int
}

// expand turns "dir/**" into the files under dir, from jsDelivr's listing of the package.
func (v *vendor) expand(pkg, version, p string) []string {
	if !strings.HasSuffix(p, "/**") {
		return []string{p}
	}
	key := pkg + "@" + version
	files, ok := v.listings[key]
	if !ok {
		res, err := client.Get("https://data.jsdelivr.com/v1/packages/npm/" + key + "?structure=flat")
		if err != nil {
			log.Fatal(err)
		}
		var listing struct{ Files []struct{ Name string } }
		err = json.NewDecoder(res.Body).Decode(&listing)
		res.Body.Close()
		if err != nil {
			log.Fatalf("listing %s: %v", key, err)
		}
		for _, f := range listing.Files {
			files = append(files, strings.TrimPrefix(f.Name, "/"))
		}
		v.listings[key] = files
	}
	dir := strings.TrimSuffix(p, "**")
	var out []string
	for _, f := range files {
		if strings.HasPrefix(f, dir) {
			out = append(out, f)
		}
	}
	return out
}

// cssURL finds url(…) references in a stylesheet.
var cssURL = regexp.MustCompile(`url\(\s*['"]?([^'")]+)['"]?\s*\)`)

// fetch downloads one file of a package and, for a stylesheet, the files it refers to.
func (v *vendor) fetch(pkgVersion, file string) {
	file = path.Clean(file)
	key := pkgVersion + "/" + file
	if v.done[key] {
		return
	}
	v.done[key] = true
	res, err := client.Get(v.from + key)
	if err != nil {
		log.Printf("%s: %v", key, err)
		v.failed++
		return
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("%s: %s %v", key, res.Status, err)
		v.failed++
		return
	}
	dst := filepath.Join(v.out, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(dst, body, 0o644); err != nil {
		log.Fatal(err)
	}
	if !strings.HasSuffix(file, ".css") {
		return
	}
	for _, m := range cssURL.FindAllSubmatch(body, -1) {
		ref := string(m[1])
		if strings.HasPrefix(ref, "data:") || strings.HasPrefix(ref, "#") || strings.Contains(ref, "//") {
			continue
		}
		u, err := url.Parse(ref)
		if err != nil || u.Path == "" {
			continue
		}
		v.fetch(pkgVersion, path.Join(path.Dir(file), u.Path))
	}
}
