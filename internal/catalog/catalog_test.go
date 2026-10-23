package catalog

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoadmapUpToDate(t *testing.T) {
	var b bytes.Buffer
	Roadmap(&b)
	have, err := os.ReadFile("../../ROADMAP.md")
	require.NoError(t, err)
	assert.Equal(t, b.String(), string(have), "ROADMAP.md is stale: run go run ./internal/catalog/roadmap")
}

func TestUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range Pages {
		assert.False(t, seen[p.Name], p.Name)
		seen[p.Name] = true
	}
	assert.Len(t, Pages, 146)
}
