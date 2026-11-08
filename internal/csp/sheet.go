package csp

import (
	"context"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/a-h/templ"
)

// A sheet is a package's own CSS. Each package registers its sheet once (a package-level var); the
// sheet renders inline at most once per page, where the first component that needs it is drawn, so a
// page carries only the CSS of the blocks it shows. With SetExternal the sheets render nothing: the app
// then serves All as one cached file (package static, WithKitCSS).
type sheet struct {
	name string
	css  func() string
	once *templ.OnceHandle
}

var (
	mu       sync.Mutex
	sheets   []*sheet
	external atomic.Bool
)

// Sheet registers css under name and returns the component that renders it (inline, once per page).
func Sheet(name, css string) templ.Component {
	return SheetFunc(name, func() string { return css })
}

// SheetFunc is Sheet for CSS that depends on settings made at startup (e.g. a font URL under cdn.Base).
func SheetFunc(name string, css func() string) templ.Component {
	s := &sheet{name: name, css: css, once: templ.NewOnceHandle()}
	mu.Lock()
	sheets = append(sheets, s)
	mu.Unlock()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if external.Load() {
			return nil
		}
		return s.once.Once().Render(templ.WithChildren(ctx, Style(s.css())), w)
	})
}

// SetExternal switches the inline sheets off (true) or on again (false).
func SetExternal(on bool) { external.Store(on) }

// External reports whether the sheets are served as a file.
func External() bool { return external.Load() }

// All is the CSS of every registered sheet, in registration order: the sheets of the packages linked
// into the program, and only those.
func All() string {
	mu.Lock()
	defer mu.Unlock()
	var b strings.Builder
	for _, s := range sheets {
		b.WriteString("/* " + s.name + " */\n")
		b.WriteString(s.css())
		b.WriteString("\n")
	}
	return b.String()
}

// Names lists the registered sheets, for tests.
func Names() []string {
	mu.Lock()
	defer mu.Unlock()
	out := make([]string, len(sheets))
	for i, s := range sheets {
		out[i] = s.name
	}
	return out
}
