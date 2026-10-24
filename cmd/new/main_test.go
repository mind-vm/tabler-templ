package main

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCreate writes the app and checks that its Go files parse and carry the module path.
func TestCreate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "my-app")
	if err := create(dir, data{Module: "example.com/my-app", Name: "My app", Binary: "my-app", Theme: "Rose"}); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"main.go", "auth.go", "main_test.go", "internal/projects/projects.go"} {
		if _, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, f), nil, 0); err != nil {
			t.Error(err)
		}
	}
	main, _ := os.ReadFile(filepath.Join(dir, "main.go"))
	if !strings.Contains(string(main), `"example.com/my-app/views"`) {
		t.Error("main.go does not import the app's views")
	}
	if _, err := os.Stat(filepath.Join(dir, ".gitignore")); err != nil {
		t.Error(err)
	}
	if err := create(dir, data{}); err == nil {
		t.Error("create into a directory with files: no error")
	}
}

// TestNewApp creates the app against this checkout, generates it and runs its tests. It needs the Go
// toolchain and the module cache (network on first use); set TABLER_NEW_E2E=1 to run it.
func TestNewApp(t *testing.T) {
	if os.Getenv("TABLER_NEW_E2E") == "" {
		t.Skip("set TABLER_NEW_E2E=1")
	}
	dir := filepath.Join(t.TempDir(), "acme")
	d := data{Module: "example.com/acme", Name: "Acme", Binary: "acme", Theme: "Violet"}
	if err := create(dir, d); err != nil {
		t.Fatal(err)
	}
	root, _ := filepath.Abs("../..")
	if err := setup(dir, d, "", root); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}
