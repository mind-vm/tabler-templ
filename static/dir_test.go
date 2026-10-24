package static

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/mind-vm/tabler-templ/cdn"
)

// TestDir serves files cmd/vendor wrote into the bundle at cdn.Base = Dir("npm/"), cached like the others.
func TestDir(t *testing.T) {
	f := New("/static/", fstest.MapFS{"npm/apexcharts@7.0.0/dist/apexcharts.min.js": {Data: []byte(strings.Repeat("x", 2000))}})
	base := f.Dir("npm/")
	if !strings.HasPrefix(base, "/static/v") || !strings.HasSuffix(base, "/npm/") {
		t.Fatal(base)
	}
	cdn.Base = base
	defer func() { cdn.Base = cdn.JSDelivr }()
	w := httptest.NewRecorder()
	f.ServeHTTP(w, httptest.NewRequest("GET", cdn.URL("apexcharts", "7.0.0", "dist/apexcharts.min.js"), nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Header().Get("Cache-Control"), "immutable") {
		t.Errorf("%d %v", w.Code, w.Header())
	}
}
