package motion_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/motion"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestAttrs(t *testing.T) {
	for name, tc := range map[string]struct {
		p    motion.Props
		want templ.Attributes
	}{
		"zero fades up": {motion.Props{}, templ.Attributes{"data-motion": "fade-up"}},
		"all": {
			motion.Props{Preset: motion.SlideLeft, Duration: 0.5, Delay: 0.1, Stagger: 0.05, Bounce: 0.2, InView: true, Exit: motion.FadeIn, ExitDuration: 0.3},
			templ.Attributes{"data-motion": "slide-left", "data-motion-duration": "0.5", "data-motion-delay": "0.1",
				"data-motion-stagger": "0.05", "data-motion-bounce": "0.2", "data-motion-inview": true,
				"data-motion-exit": "fade-in", "data-motion-exit-duration": "0.3"},
		},
		"pop springs": {motion.Props{Preset: motion.Pop}, templ.Attributes{"data-motion": "pop", "data-motion-bounce": "0.4"}},
		"exit only":   {motion.Props{Exit: motion.ZoomIn}, templ.Attributes{"data-motion-exit": "zoom-in"}},
	} {
		got := motion.Attrs(tc.p)
		if len(got) != len(tc.want) {
			t.Errorf("%s: %v, want %v", name, got, tc.want)
			continue
		}
		for k, v := range tc.want {
			if got[k] != v {
				t.Errorf("%s: %s = %v, want %v", name, k, got[k], v)
			}
		}
	}
}

func TestCounter(t *testing.T) {
	html := render(t, motion.Counter(motion.CounterProps{To: 1234.5, Decimals: 1, Prefix: "<$", Suffix: "%", InView: true, Class: "h1"}))
	for _, want := range []string{`class="h1"`, `data-motion-count="1234.5"`, `data-motion-decimals="1"`, `data-motion-prefix="&lt;$"`,
		`data-motion-inview`, `>&lt;$1234.5%</span>`} {
		if !strings.Contains(html, want) {
			t.Errorf("no %s in %s", want, html)
		}
	}
}

// TestLibraryOnce renders the library twice on a page: the script and its CSS come once.
func TestLibraryOnce(t *testing.T) {
	html := render(t, templ.Join(motion.Library(), motion.Library(), motion.Script()))
	if n := strings.Count(html, "window.tablerMotion = function"); n != 1 {
		t.Errorf("setup %d times", n)
	}
	if n := strings.Count(html, ".motion-on [data-motion]"); n != 1 {
		t.Errorf("CSS %d times", n)
	}
	if !strings.Contains(html, `src="https://cdn.jsdelivr.net/npm/motion@`+motion.Version+`/dist/motion.js"`) {
		t.Error(html)
	}
}
