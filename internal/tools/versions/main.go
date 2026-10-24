// Command versions compares the pinned versions in versions/versions.go with a Tabler checkout: the Tabler
// release (core/package.json), its commit, and each library at the version Tabler's pnpm-lock.yaml
// resolves. Run it from the module root:
//
//	go run ./internal/tools/versions <tabler checkout>         # list the differences, exit 1 if any
//	go run ./internal/tools/versions -write <tabler checkout>  # update versions/versions.go
//	go run ./internal/tools/versions -ref                      # print the pinned Tabler commit
//	go run ./internal/tools/versions -icons                    # print the pinned icons version
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mind-vm/tabler-templ/versions"
)

const file = "versions/versions.go"

// constLine matches a pinned version: name, version and, after "npm:", the npm package.
var constLine = regexp.MustCompile(`^(\t(\w+)\s*=\s*")([^"]*)("\s*// npm: (\S+))`)

var commitLine = regexp.MustCompile(`^(\tTablerCommit\s*=\s*")([^"]*)(")`)

// locked reads the versions pnpm-lock.yaml resolves for the importers' dependencies: npm name → version.
// For a package several importers use, the first one wins (core comes before preview and docs).
func locked(path string) map[string]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	out := map[string]string{}
	inImporters := false
	name := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		switch {
		case line == "importers:":
			inImporters = true
		case inImporters && line != "" && !strings.HasPrefix(line, " "):
			return out // the next top-level key ends the importers
		case inImporters && strings.HasPrefix(line, "      ") && !strings.HasPrefix(line, "       ") && strings.HasSuffix(line, ":"):
			name = strings.Trim(strings.TrimSpace(strings.TrimSuffix(line, ":")), "'")
		case inImporters && name != "" && strings.HasPrefix(strings.TrimSpace(line), "version: "):
			v := strings.TrimPrefix(strings.TrimSpace(line), "version: ")
			if i := strings.IndexByte(v, '('); i >= 0 {
				v = v[:i] // peer suffix, e.g. 2.0.0(stylelint@16)
			}
			if _, seen := out[name]; !seen && !strings.HasPrefix(v, "link:") {
				out[name] = v
			}
			name = ""
		}
	}
	return out
}

func main() {
	write := flag.Bool("write", false, "update "+file)
	ref := flag.Bool("ref", false, "print the pinned Tabler commit")
	icons := flag.Bool("icons", false, "print the pinned icons version")
	flag.Parse()
	switch {
	case *ref:
		fmt.Println(versions.TablerCommit)
		return
	case *icons:
		fmt.Println(versions.Icons)
		return
	case flag.NArg() != 1:
		log.Fatal("usage: versions [-write] <tabler checkout> | -ref | -icons")
	}
	tabler := flag.Arg(0)

	want := locked(filepath.Join(tabler, "pnpm-lock.yaml"))
	var core struct{ Version string }
	b, err := os.ReadFile(filepath.Join(tabler, "core/package.json"))
	if err == nil {
		err = json.Unmarshal(b, &core)
	}
	if err != nil {
		log.Fatal(err)
	}
	want["@tabler/core"] = core.Version // a workspace package: not in the lockfile
	commit, err := exec.Command("git", "-C", tabler, "rev-parse", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}

	src, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(src), "\n")
	changed := 0
	for i, l := range lines {
		if m := commitLine.FindStringSubmatch(l); m != nil {
			if c := strings.TrimSpace(string(commit)); c != m[2] {
				fmt.Printf("TablerCommit %s -> %s\n", m[2], c)
				lines[i] = m[1] + c + m[3]
				changed++
			}
			continue
		}
		m := constLine.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		v, ok := want[m[5]]
		switch {
		case !ok:
			fmt.Printf("%s: %s is no longer in Tabler's lockfile\n", m[2], m[5])
			changed++
		case v != m[3]:
			fmt.Printf("%s (%s) %s -> %s\n", m[2], m[5], m[3], v)
			lines[i] = m[1] + v + m[4]
			changed++
		}
	}
	if changed == 0 {
		fmt.Println("versions match Tabler")
		return
	}
	if !*write {
		os.Exit(1)
	}
	if err := os.WriteFile(file, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("updated " + file + "; check the libraries' changelogs for breaking changes")
}
