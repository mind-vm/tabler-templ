package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mind-vm/tabler-templ/internal/catalog"
)

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	handler().ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

// The catalog's status and the demo's pages must agree: a page counts as rebuilt only if it is served.
func TestPagesMatchCatalog(t *testing.T) {
	for _, p := range catalog.Pages {
		_, served := pages[p.Name]
		assert.Equal(t, p.Status != catalog.Open, served, p.Name)
	}
	for name := range pages {
		_, ok := catalog.Find(name)
		assert.True(t, ok, name)
	}
}

func TestServe(t *testing.T) {
	w := get(t, "/")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `href="/buttons"`)
	for name := range pages {
		w := get(t, "/"+name)
		assert.Equal(t, http.StatusOK, w.Code, name)
	}
	assert.Equal(t, http.StatusNotFound, get(t, "/dropdowns").Code)
}
