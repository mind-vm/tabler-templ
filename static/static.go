// Package static serves an app's own static files fast: with a version in the path and a year of cache,
// with an ETag for revalidation, and compressed once in memory (brotli, gzip). It is the alternative to
// loading Tabler from jsDelivr (package assets) for apps that serve Tabler themselves: Tabler's CSS and JS,
// a subset of the icon webfont with only the icons the app uses (cmd/static), and the app's own files.
//
//	//go:embed all:web
//	var web embed.FS
//
//	files := static.New("/static/", must(fs.Sub(web, "web")))
//	mux.Handle("/static/", files)
//	… <link rel="stylesheet" href={ files.URL("app.css") }/>
//	… @layout.Page(layout.PageProps{Head: static.Head(files), Scripts: static.Scripts(files)})
//
// Addresses:
//
//   - <prefix>v<version>/<name>, from [Files.URL]: the version is a hash over all files of the bundle; when
//     one changes, all addresses change. Relative references in stylesheets (fonts) stay in the same
//     version. Sent with [CacheVersioned]: the browser keeps the file for a year without asking again.
//   - <prefix>v<other>/<name>: another version (an old page after a deploy, the other instance during a
//     rolling deploy). The current file is sent, but with [CacheOtherVersion], so the browser does not keep
//     it for a year under the wrong address.
//   - <prefix><name>: an address without version (bookmarks, service workers, hand-written links): ETag
//     and a short cache ([CacheUnversioned]). With ?v=<hash of the file> (at least 12 hex digits of the
//     SHA-256 of its content) it is immutable like a versioned address.
//
// Text files (CSS, JS, SVG, JSON …) are compressed once on first use: gzip at once, brotli in the
// background (gzip is sent until brotli is ready). Fonts (woff2) and images are compressed already and
// are sent as they are. Responses already compressed here pass through compress.Handler unchanged.
//
// Measured in an app with Tabler, the icon subset and one font (phone at 4 Mbit/s, 100 ms RTT): the
// first visit of a page went from 1.9 MB in 4.3 s (served uncompressed, without cache headers, the whole
// icon webfont) to 185 KB in 0.7 s, a repeat visit to 12–17 KB (only the HTML) in 0.2 s.
package static

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/andybalholm/brotli"
)

// Cache-Control for the three kinds of address.
const (
	CacheVersioned    = "public, max-age=31536000, immutable"
	CacheOtherVersion = "no-cache"
	CacheUnversioned  = "public, max-age=300"
)

// compressible are the extensions sent compressed.
var compressible = map[string]bool{".css": true, ".js": true, ".mjs": true, ".json": true, ".svg": true,
	".html": true, ".txt": true, ".map": true, ".xml": true, ".webmanifest": true}

// types that mime.TypeByExtension does not know everywhere (Alpine without /etc/mime.types).
var types = map[string]string{
	".css": "text/css; charset=utf-8", ".js": "text/javascript; charset=utf-8", ".mjs": "text/javascript; charset=utf-8",
	".json": "application/json", ".svg": "image/svg+xml", ".woff2": "font/woff2", ".woff": "font/woff",
	".webmanifest": "application/manifest+json", ".png": "image/png", ".txt": "text/plain; charset=utf-8",
	".html": "text/html; charset=utf-8",
}

// Files is a bundle of static files under one prefix, e.g. "/static/". It is an http.Handler.
type Files struct {
	prefix string
	fsys   fs.FS

	once    sync.Once
	version string
	files   map[string]*file
}

type file struct {
	raw  []byte
	gz   []byte                 // nil: not compressed
	br   atomic.Pointer[[]byte] // nil until brotli is ready (or never, when not compressed)
	typ  string
	etag string // without quotes: hash of the content
}

// New serves the files of fsys under prefix, which starts and ends with "/". The files are read and
// compressed on first use and kept in memory; fsys is meant to be an embed.FS (or a part of one).
func New(prefix string, fsys fs.FS) *Files {
	if !strings.HasPrefix(prefix, "/") || !strings.HasSuffix(prefix, "/") {
		panic("static: the prefix needs / at start and end: " + prefix)
	}
	return &Files{prefix: prefix, fsys: fsys}
}

func (d *Files) load() {
	d.once.Do(func() {
		d.files = map[string]*file{}
		var names []string
		err := fs.WalkDir(d.fsys, ".", func(p string, e fs.DirEntry, err error) error {
			if err != nil || e.IsDir() {
				return err
			}
			b, err := fs.ReadFile(d.fsys, p)
			if err != nil {
				return err
			}
			h := sha256.Sum256(b)
			f := &file{raw: b, typ: contentType(p), etag: hex.EncodeToString(h[:8])}
			if compressible[path.Ext(p)] {
				if gz := gzipped(b); len(gz) < len(b)*9/10 {
					f.gz = gz
				}
			}
			d.files[p] = f
			names = append(names, p)
			return nil
		})
		if err != nil {
			panic("static: " + err.Error())
		}
		sort.Strings(names)
		h := sha256.New()
		for _, n := range names {
			h.Write([]byte(n + "\x00" + d.files[n].etag + "\x00"))
		}
		d.version = hex.EncodeToString(h.Sum(nil))[:12]
		go func() { // brotli takes about a second for Tabler; until then gzip is sent
			for _, n := range names {
				if f := d.files[n]; f.gz != nil {
					br := brotlied(f.raw)
					if len(br) < len(f.gz) {
						f.br.Store(&br)
					}
				}
			}
		}()
	})
}

func contentType(p string) string {
	if t, ok := types[path.Ext(p)]; ok {
		return t
	}
	if t := mime.TypeByExtension(path.Ext(p)); t != "" {
		return t
	}
	return "application/octet-stream"
}

func gzipped(b []byte) []byte {
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, _ = w.Write(b)
	_ = w.Close()
	return buf.Bytes()
}

func brotlied(b []byte) []byte {
	var buf bytes.Buffer
	w := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	_, _ = w.Write(b)
	_ = w.Close()
	return buf.Bytes()
}

// Version is the hash over all files of the bundle (12 hex digits).
func (d *Files) Version() string {
	d.load()
	return d.version
}

// URL is the versioned address of a file, e.g. URL("app.css") = "/static/v3f…/app.css". A file that is not
// in the bundle gets its address without version (and a 404 there); check with [Files.Has] in a test.
func (d *Files) URL(name string) string {
	d.load()
	if _, ok := d.files[name]; !ok {
		return d.prefix + name
	}
	return d.prefix + "v" + d.version + "/" + name
}

// Has reports whether the bundle has the file.
func (d *Files) Has(name string) bool {
	d.load()
	_, ok := d.files[name]
	return ok
}

// isVersion: "v" and 12 hex digits.
func isVersion(s string) bool {
	if len(s) != 13 || s[0] != 'v' {
		return false
	}
	for _, c := range s[1:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

// ServeHTTP sends a file (GET and HEAD); paths without a file (directories too) are 404.
func (d *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.load()
	rest, ok := strings.CutPrefix(r.URL.Path, d.prefix)
	if !ok {
		http.NotFound(w, r)
		return
	}
	cache := CacheUnversioned
	if seg, after, ok := strings.Cut(rest, "/"); ok && isVersion(seg) {
		rest, cache = after, CacheOtherVersion
		if seg[1:] == d.version {
			cache = CacheVersioned
		}
	}
	f := d.files[rest]
	if f == nil {
		http.NotFound(w, r)
		return
	}
	if v := r.URL.Query().Get("v"); cache == CacheUnversioned && len(v) >= 12 && strings.HasPrefix(f.etag, v) {
		cache = CacheVersioned
	}
	h := w.Header()
	h.Set("Cache-Control", cache)
	h.Set("Content-Type", f.typ)
	h.Set("X-Content-Type-Options", "nosniff")
	body, etag := f.raw, f.etag
	if f.gz != nil {
		h.Add("Vary", "Accept-Encoding")
		br := f.br.Load()
		switch {
		case br != nil && accepts(r, "br"):
			body, etag = *br, etag+"-br"
			h.Set("Content-Encoding", "br")
		case accepts(r, "gzip"):
			body, etag = f.gz, etag+"-gz"
			h.Set("Content-Encoding", "gzip")
		}
	}
	h.Set("ETag", `"`+etag+`"`)
	http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(body))
}

// accepts: Accept-Encoding names the encoding (or *) without q=0.
func accepts(r *http.Request, encoding string) bool {
	for _, part := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
		name, params, _ := strings.Cut(strings.TrimSpace(part), ";")
		name = strings.ToLower(strings.TrimSpace(name))
		if name != encoding && name != "*" {
			continue
		}
		if q, ok := strings.CutPrefix(strings.ReplaceAll(params, " ", ""), "q="); ok {
			if v, err := strconv.ParseFloat(q, 64); err == nil && v == 0 {
				return false
			}
		}
		return true
	}
	return false
}
