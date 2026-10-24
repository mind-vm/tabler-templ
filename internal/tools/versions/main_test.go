package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEveryVersionNamesItsPackage keeps versions.go checkable: each pinned version names its npm package.
func TestEveryVersionNamesItsPackage(t *testing.T) {
	src, err := os.ReadFile("../../../" + file)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, l := range strings.Split(string(src), "\n") {
		if !strings.HasPrefix(l, "\t") || !strings.Contains(l, `= "`) || commitLine.MatchString(l) {
			continue
		}
		if !constLine.MatchString(l) {
			t.Errorf("no // npm: package: %s", strings.TrimSpace(l))
		}
		n++
	}
	if n < 10 {
		t.Errorf("found %d versions", n)
	}
}

func TestLocked(t *testing.T) {
	lock := `lockfileVersion: '9.0'

importers:

  .:
    devDependencies:
      '@tabler/icons':
        specifier: ^3.49.0
        version: 3.49.0
  core:
    devDependencies:
      apexcharts:
        specifier: 7.0.0
        version: 7.0.0
      stylelint-config:
        specifier: ^2.0.0
        version: 2.0.0(stylelint@16.26.1)
  preview:
    dependencies:
      '@tabler/core':
        specifier: workspace:*
        version: link:../core
      apexcharts:
        specifier: 6.0.0
        version: 6.0.0

packages:

  apexcharts@7.0.0:
    resolution: {integrity: sha512-x}
`
	path := filepath.Join(t.TempDir(), "pnpm-lock.yaml")
	os.WriteFile(path, []byte(lock), 0o644)
	got := locked(path)
	want := map[string]string{"@tabler/icons": "3.49.0", "apexcharts": "7.0.0", "stylelint-config": "2.0.0"}
	if len(got) != len(want) {
		t.Errorf("got %v", got)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
}
