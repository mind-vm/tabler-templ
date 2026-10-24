// Command new creates a web app with tabler-templ: Tabler's page frame, sign-in, a dashboard with charts
// and a list of records with search, filter, pages, a create page and a detail page with an edit dialog.
// It writes the files, adds the modules and generates the templ code, so `go run .` starts the app.
//
//	go run github.com/mind-vm/tabler-templ/cmd/new@latest [-module example.com/myapp] [-name "My app"] myapp
package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"
	"unicode"
)

//go:embed all:template
var files embed.FS

const module = "github.com/mind-vm/tabler-templ"

// templVersion is the templ version tabler-templ is built with; used when the build has no module info.
const templVersion = "v0.3.1070"

type data struct {
	Module string // the new app's module path
	Name   string // the app's name, e.g. "My app"
	Binary string // the directory's name
}

func main() {
	modPath := flag.String("module", "", "module path of the new app (default: the directory's name)")
	name := flag.String("name", "", "the app's name, shown in the navigation (default: from the directory's name)")
	version := flag.String("version", "", "tabler-templ version to use (default: the version of this command, or latest)")
	replace := flag.String("replace", "", "use a local copy of tabler-templ (a replace directive), e.g. for developing tabler-templ")
	noSetup := flag.Bool("no-setup", false, "only write the files; skip go get and templ generate")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: new [flags] <directory>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	dir := flag.Arg(0)
	base := filepath.Base(filepath.Clean(dir))
	d := data{Module: *modPath, Name: *name, Binary: base}
	if d.Module == "" {
		d.Module = base
	}
	if d.Name == "" {
		d.Name = title(base)
	}
	if err := create(dir, d); err != nil {
		fail(err)
	}
	if !*noSetup {
		if err := setup(dir, d, *version, *replace); err != nil {
			fail(err)
		}
	}
	fmt.Printf("\nCreated %s in %s. Start it with\n\n\tcd %s && go run .\n\nand sign in with admin@example.com / admin.\n", d.Name, dir, dir)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "new:", err)
	os.Exit(1)
}

// title makes "my-app" "My app".
func title(s string) string {
	s = strings.NewReplacer("-", " ", "_", " ").Replace(s)
	r := []rune(s)
	if len(r) > 0 {
		r[0] = unicode.ToUpper(r[0])
	}
	return string(r)
}

// create writes the template's files into dir, which must not exist or be empty.
func create(dir string, d data) error {
	if entries, err := os.ReadDir(dir); err == nil && len(entries) > 0 {
		return fmt.Errorf("%s is not empty", dir)
	}
	return fs.WalkDir(files, "template", func(p string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() {
			return err
		}
		src, err := files.ReadFile(p)
		if err != nil {
			return err
		}
		t, err := template.New(p).Delims("[[", "]]").Parse(string(src))
		if err != nil {
			return err
		}
		var b bytes.Buffer
		if err := t.Execute(&b, d); err != nil {
			return err
		}
		rel := strings.TrimSuffix(strings.TrimPrefix(p, "template/"), ".tmpl")
		if path.Base(rel) == "gitignore" {
			rel = path.Join(path.Dir(rel), ".gitignore")
		}
		out := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return err
		}
		fmt.Println("  " + rel)
		return os.WriteFile(out, b.Bytes(), 0o644)
	})
}

// versions are the tabler-templ and templ versions to require: the ones this command was built with
// (go run …@v0.9.0), else latest.
func versions(override string) (kit, templ string) {
	kit, templ = "latest", templVersion
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; info.Main.Path == module && v != "" && v != "(devel)" {
			kit = v
		}
		for _, m := range info.Deps {
			if m.Path == "github.com/a-h/templ" {
				templ = m.Version
			}
		}
	}
	if override != "" {
		kit = override
	}
	return kit, templ
}

// setup creates go.mod, adds tabler-templ and the templ tool, and generates the templ code.
func setup(dir string, d data, version, replace string) error {
	kit, templ := versions(version)
	steps := [][]string{{"go", "mod", "init", d.Module}}
	if replace != "" {
		abs, err := filepath.Abs(replace)
		if err != nil {
			return err
		}
		steps = append(steps, []string{"go", "mod", "edit", "-require=" + module + "@v0.0.0", "-replace=" + module + "=" + abs})
	} else {
		steps = append(steps, []string{"go", "get", module + "@" + kit})
	}
	steps = append(steps,
		[]string{"go", "get", "-tool", "github.com/a-h/templ/cmd/templ@" + templ},
		[]string{"go", "tool", "templ", "generate"},
		[]string{"go", "mod", "tidy"},
	)
	for _, s := range steps {
		fmt.Println("$ " + strings.Join(s, " "))
		cmd := exec.Command(s[0], s[1:]...)
		cmd.Dir, cmd.Stdout, cmd.Stderr = dir, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", strings.Join(s, " "), err)
		}
	}
	return nil
}
