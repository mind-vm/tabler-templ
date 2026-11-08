package markdown

import (
	"context"
	"strings"
	"testing"
)

func render(t *testing.T, src string) string {
	t.Helper()
	var b strings.Builder
	if err := Markdown(src).Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestMarkdown(t *testing.T) {
	out := render(t, "# Title\n\n**bold** and a [link](https://example.com)\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n- [x] done")
	for _, want := range []string{`<div class="prose">`, "<h1>Title</h1>", "<strong>bold</strong>",
		`<a href="https://example.com">link</a>`, "<table>", "<td>1</td>", `type="checkbox"`} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in %s", want, out)
		}
	}
}

func TestMarkdownIsSafe(t *testing.T) {
	out := render(t, "<script>alert(1)</script>\n\n<img src=x onerror=alert(1)>\n\n[x](javascript:alert(1))")
	for _, bad := range []string{"<script>", "onerror", "javascript:"} {
		if strings.Contains(out, bad) {
			t.Errorf("unsafe %q in %s", bad, out)
		}
	}
}
