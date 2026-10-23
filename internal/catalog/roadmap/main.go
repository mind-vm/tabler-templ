// Command roadmap writes ROADMAP.md from the catalog. Run it from the module root:
//
//	go run ./internal/catalog/roadmap
package main

import (
	"bytes"
	"log"
	"os"

	"github.com/mind-vm/tabler-templ/internal/catalog"
)

func main() {
	var b bytes.Buffer
	catalog.Roadmap(&b)
	if err := os.WriteFile("ROADMAP.md", b.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
