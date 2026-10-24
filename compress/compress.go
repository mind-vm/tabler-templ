// Package compress gzips an app's responses: the HTML of templ pages compresses to a sixth or less. Tabler,
// its icons and the plugins come from jsDelivr already minified and compressed; the pages themselves are
// the app's to compress.
//
//	http.ListenAndServe(addr, compress.Handler(mux))
//
// Behind a proxy that compresses already (Caddy, Traefik, nginx with gzip on), leave it out.
package compress

import (
	"bufio"
	"compress/gzip"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
)

// MinSize is the smallest response worth compressing; smaller ones are sent as they are.
const MinSize = 1024

var writers = sync.Pool{New: func() any {
	w, _ := gzip.NewWriterLevel(io.Discard, gzip.DefaultCompression)
	return w
}}

// compressible are the content types that shrink: text, scripts, styles, JSON and SVG.
func compressible(contentType string) bool {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	return strings.HasPrefix(ct, "text/") || ct == "application/javascript" || ct == "application/json" ||
		ct == "image/svg+xml" || ct == "application/xml" || strings.HasSuffix(ct, "+json")
}

// Handler gzips h's responses for clients that accept gzip, when they are text of at least MinSize bytes
// and not encoded already. Responses that stream (Flush) are compressed as they go.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding")
		if r.Method == http.MethodHead || !acceptsGzip(r.Header.Get("Accept-Encoding")) {
			h.ServeHTTP(w, r)
			return
		}
		cw := &writer{ResponseWriter: w}
		defer cw.close()
		h.ServeHTTP(cw, r)
	})
}

func acceptsGzip(header string) bool {
	for _, part := range strings.Split(header, ",") {
		name, q, _ := strings.Cut(strings.TrimSpace(part), ";")
		if strings.EqualFold(strings.TrimSpace(name), "gzip") {
			return strings.ReplaceAll(strings.TrimSpace(q), " ", "") != "q=0"
		}
	}
	return false
}

// writer holds back the start of the body until it knows whether to compress: the status is written when
// MinSize bytes have come, the handler flushes, or the response ends.
type writer struct {
	http.ResponseWriter
	status  int
	buf     []byte
	gz      *gzip.Writer
	decided bool
}

func (w *writer) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}

func (w *writer) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	if !w.decided {
		w.buf = append(w.buf, p...)
		if len(w.buf) < MinSize {
			return len(p), nil
		}
		if err := w.decide(); err != nil {
			return 0, err
		}
		return len(p), nil
	}
	if w.gz != nil {
		return w.gz.Write(p)
	}
	return w.ResponseWriter.Write(p)
}

// decide starts the response, compressed or not, and writes what was held back.
func (w *writer) decide() error {
	w.decided = true
	hdr := w.Header()
	if hdr.Get("Content-Type") == "" && len(w.buf) > 0 {
		hdr.Set("Content-Type", http.DetectContentType(w.buf))
	}
	if len(w.buf) >= MinSize && hdr.Get("Content-Encoding") == "" && compressible(hdr.Get("Content-Type")) &&
		w.status != http.StatusNoContent && w.status != http.StatusNotModified && w.status >= http.StatusOK {
		hdr.Set("Content-Encoding", "gzip")
		hdr.Del("Content-Length")
		w.gz = writers.Get().(*gzip.Writer)
		w.gz.Reset(w.ResponseWriter)
	}
	w.ResponseWriter.WriteHeader(w.status)
	buf := w.buf
	w.buf = nil
	if len(buf) == 0 {
		return nil
	}
	var err error
	if w.gz != nil {
		_, err = w.gz.Write(buf)
	} else {
		_, err = w.ResponseWriter.Write(buf)
	}
	return err
}

func (w *writer) close() {
	if !w.decided {
		if w.status == 0 {
			return // nothing written: net/http sends its own empty 200
		}
		w.decide()
	}
	if w.gz != nil {
		w.gz.Close()
		writers.Put(w.gz)
		w.gz = nil
	}
}

// Flush sends what is there, e.g. for server-sent events or a page streamed in parts.
func (w *writer) Flush() {
	if !w.decided {
		if w.status == 0 {
			w.status = http.StatusOK
		}
		w.decide()
	}
	if w.gz != nil {
		w.gz.Flush()
	}
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack lets websockets through.
func (w *writer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := w.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// Unwrap gives http.ResponseController the underlying writer.
func (w *writer) Unwrap() http.ResponseWriter { return w.ResponseWriter }
