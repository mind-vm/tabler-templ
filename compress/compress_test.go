package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serve(h http.HandlerFunc, accept string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", "/", nil)
	if accept != "" {
		r.Header.Set("Accept-Encoding", accept)
	}
	w := httptest.NewRecorder()
	Handler(h).ServeHTTP(w, r)
	return w
}

func TestHandler(t *testing.T) {
	page := strings.Repeat("<div class=\"card\">Hello</div>", 200)
	html := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusTeapot)
		io.WriteString(w, page[:100])
		io.WriteString(w, page[100:])
	}

	w := serve(html, "br, gzip;q=0.8")
	if w.Code != http.StatusTeapot || w.Header().Get("Content-Encoding") != "gzip" || w.Header().Get("Vary") != "Accept-Encoding" {
		t.Fatalf("%d %v", w.Code, w.Header())
	}
	if w.Body.Len() >= len(page)/5 {
		t.Errorf("%d bytes for %d", w.Body.Len(), len(page))
	}
	zr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if b, _ := io.ReadAll(zr); string(b) != page {
		t.Error("body changed")
	}

	if w := serve(html, ""); w.Header().Get("Content-Encoding") != "" || w.Body.String() != page {
		t.Error("compressed without Accept-Encoding")
	}
	if w := serve(html, "gzip;q=0"); w.Header().Get("Content-Encoding") != "" {
		t.Error("compressed though gzip is refused")
	}
	small := func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "<p>hi</p>") }
	if w := serve(small, "gzip"); w.Header().Get("Content-Encoding") != "" || w.Body.String() != "<p>hi</p>" {
		t.Error("small response compressed")
	}
	png := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(make([]byte, 4096))
	}
	if w := serve(png, "gzip"); w.Header().Get("Content-Encoding") != "" {
		t.Error("image compressed")
	}
	redirect := func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/x", http.StatusSeeOther) }
	if w := serve(redirect, "gzip"); w.Code != http.StatusSeeOther || w.Header().Get("Location") != "/x" {
		t.Errorf("redirect: %d", w.Code)
	}
	empty := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }
	if w := serve(empty, "gzip"); w.Code != http.StatusNoContent || w.Header().Get("Content-Encoding") != "" {
		t.Errorf("no content: %d", w.Code)
	}
}
