package static

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/mind-vm/tabler-templ/compress"
)

var styles = strings.Repeat(".a{color:red}\n", 500)

func bundle() *Files {
	return New("/x/", fstest.MapFS{
		"a.css":          {Data: []byte(styles)},
		"fonts/f.woff2":  {Data: []byte("wOF2" + strings.Repeat("\x00", 500))},
		"sub/b.js":       {Data: []byte(strings.Repeat("var b=1;\n", 300))},
		"small.txt":      {Data: []byte("x")},
		"vendor/LICENSE": {Data: []byte("MIT")},
	})
}

func get(t *testing.T, h http.Handler, path string, header ...string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, path, nil)
	for i := 0; i+1 < len(header); i += 2 {
		r.Header.Set(header[i], header[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// With the version in the path: a year, immutable; without version: short; with another version: revalidate.
func TestCachePerAddress(t *testing.T) {
	d := bundle()
	url := d.URL("a.css")
	if url != "/x/v"+d.Version()+"/a.css" || len(d.Version()) != 12 {
		t.Fatalf("URL %q", url)
	}
	for path, want := range map[string]string{
		url:                      CacheVersioned,
		"/x/a.css":               CacheUnversioned,
		"/x/v000000000000/a.css": CacheOtherVersion,
		d.URL("fonts/f.woff2"):   CacheVersioned,
		d.URL("sub/b.js"):        CacheVersioned,
		"/x/a.css?v=" + d.files["a.css"].etag[:12]: CacheVersioned,
		"/x/a.css?v=000000000000":                  CacheUnversioned,
		"/x/a.css?v=" + d.files["a.css"].etag[:4]:  CacheUnversioned,
	} {
		w := get(t, d, path)
		if w.Code != http.StatusOK || w.Header().Get("Cache-Control") != want || w.Header().Get("ETag") == "" {
			t.Errorf("%s: %d %v", path, w.Code, w.Header())
		}
	}
	if got := d.URL("does/not/exist.css"); got != "/x/does/not/exist.css" {
		t.Errorf("unknown file: %q", got)
	}
	for _, path := range []string{"/x/", "/x/sub/", "/x/sub", "/x/missing.css", "/x/v" + d.Version() + "/", "/y/a.css"} {
		if w := get(t, d, path); w.Code != http.StatusNotFound {
			t.Errorf("%s: %d instead of 404", path, w.Code)
		}
	}
}

// The version changes with the content of any file, and with it the addresses of all files.
func TestVersionFollowsContent(t *testing.T) {
	a := New("/x/", fstest.MapFS{"a.css": {Data: []byte("a")}, "b.css": {Data: []byte("b")}})
	b := New("/x/", fstest.MapFS{"a.css": {Data: []byte("a")}, "b.css": {Data: []byte("B")}})
	if a.Version() == b.Version() || a.URL("a.css") == b.URL("a.css") {
		t.Error("same version for other content")
	}
}

// Text is sent gzipped for Accept-Encoding gzip (with Vary and its own ETag), fonts never, without
// Accept-Encoding as it is; brotli once it is ready in the background.
func TestCompression(t *testing.T) {
	d := bundle()
	w := get(t, d, d.URL("a.css"), "Accept-Encoding", "gzip")
	if w.Header().Get("Content-Encoding") != "gzip" || w.Header().Get("Vary") != "Accept-Encoding" ||
		!strings.HasSuffix(w.Header().Get("ETag"), `-gz"`) || w.Header().Get("Content-Type") != "text/css; charset=utf-8" {
		t.Fatalf("gzip: %v", w.Header())
	}
	if w.Body.Len() >= len(styles) {
		t.Error("gzip not smaller")
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := io.ReadAll(zr); string(b) != styles {
		t.Error("gzip content wrong")
	}
	w = get(t, d, d.URL("a.css"))
	if w.Header().Get("Content-Encoding") != "" || w.Body.String() != styles {
		t.Errorf("without Accept-Encoding: %v", w.Header())
	}
	w = get(t, d, d.URL("a.css"), "Accept-Encoding", "gzip;q=0, deflate")
	if w.Header().Get("Content-Encoding") != "" {
		t.Errorf("gzip;q=0: %v", w.Header())
	}
	w = get(t, d, d.URL("fonts/f.woff2"), "Accept-Encoding", "gzip, br")
	if w.Header().Get("Content-Encoding") != "" || w.Header().Get("Content-Type") != "font/woff2" || w.Header().Get("Vary") != "" {
		t.Errorf("woff2: %v", w.Header())
	}
	w = get(t, d, d.URL("small.txt"), "Accept-Encoding", "gzip")
	if w.Header().Get("Content-Encoding") != "" || w.Body.String() != "x" {
		t.Errorf("too small to compress: %v", w.Header())
	}
	deadline := time.Now().Add(10 * time.Second)
	for {
		w = get(t, d, d.URL("a.css"), "Accept-Encoding", "gzip, deflate, br")
		if w.Header().Get("Content-Encoding") == "br" || time.Now().After(deadline) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if w.Header().Get("Content-Encoding") != "br" || !strings.HasSuffix(w.Header().Get("ETag"), `-br"`) {
		t.Fatalf("brotli: %v", w.Header())
	}
	if b, _ := io.ReadAll(brotli.NewReader(bytes.NewReader(w.Body.Bytes()))); string(b) != styles {
		t.Error("brotli content wrong")
	}
}

// With the ETag the browser revalidates and gets a 304 without body, per encoding.
func TestETag304(t *testing.T) {
	d := bundle()
	for _, enc := range []string{"", "gzip"} {
		for _, path := range []string{d.URL("a.css"), "/x/a.css"} {
			w := get(t, d, path, "Accept-Encoding", enc)
			etag := w.Header().Get("ETag")
			w = get(t, d, path, "Accept-Encoding", enc, "If-None-Match", etag)
			if w.Code != http.StatusNotModified || w.Body.Len() != 0 {
				t.Errorf("%s %q: %d", path, enc, w.Code)
			}
		}
	}
	if w := get(t, d, "/x/a.css", "If-None-Match", `"anders"`); w.Code != http.StatusOK {
		t.Errorf("other ETag: %d", w.Code)
	}
}

func TestHeadRequest(t *testing.T) {
	d := bundle()
	r := httptest.NewRequest(http.MethodHead, d.URL("a.css"), nil)
	w := httptest.NewRecorder()
	d.ServeHTTP(w, r)
	if w.Code != http.StatusOK || w.Body.Len() != 0 || w.Header().Get("Cache-Control") != CacheVersioned {
		t.Errorf("HEAD: %d %v", w.Code, w.Header())
	}
}

// Behind compress.Handler the files keep their own encoding (brotli or gzip), fonts stay as they are.
func TestBehindCompress(t *testing.T) {
	d := bundle()
	h := compress.Handler(d)
	w := get(t, h, d.URL("a.css"), "Accept-Encoding", "gzip")
	zr, err := gzip.NewReader(w.Body)
	if err != nil || w.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("%v %v", err, w.Header())
	}
	if b, _ := io.ReadAll(zr); string(b) != styles {
		t.Error("compressed twice")
	}
	w = get(t, h, d.URL("fonts/f.woff2"), "Accept-Encoding", "gzip")
	if w.Header().Get("Content-Encoding") != "" || w.Header().Get("Cache-Control") != CacheVersioned {
		t.Errorf("woff2: %v", w.Header())
	}
}
