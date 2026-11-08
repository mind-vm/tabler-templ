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

func TestStrict(t *testing.T) {
	src := "See [the docs](https://evil.example/x) and <https://auto.example> or www.bare.example.\n\n" +
		"Mail me@example.com. ![Logo](https://tracker.example/p.png) [*bold* link][ref]\n\n" +
		"[ref]: https://ref.example\n\n| a |\n|---|\n| [x](https://cell.example) |"
	out := StrictHTML(src)
	for _, bad := range []string{"<a", "<img", "href", "src=", "evil.example/x", "tracker.example", "ref.example"} {
		if strings.Contains(out, bad) {
			t.Errorf("strict output has %q: %s", bad, out)
		}
	}
	for _, want := range []string{"the docs", "https://auto.example", "www.bare.example", "me@example.com",
		"Logo", "<em>bold</em> link", "<td>x</td>"} {
		if !strings.Contains(out, want) {
			t.Errorf("strict output misses %q: %s", want, out)
		}
	}
}

func TestStrictIsSafe(t *testing.T) {
	out := StrictHTML("<script>alert(1)</script>\n\n<img src=x onerror=alert(1)>\n\n[<b>x</b>](javascript:alert(1)) <https://a.example/\"><x>")
	for _, bad := range []string{"<script>", "onerror", "javascript:", "<b>", "<x>"} {
		if strings.Contains(out, bad) {
			t.Errorf("unsafe %q in %s", bad, out)
		}
	}
}

func TestStrictComponent(t *testing.T) {
	var b strings.Builder
	if err := Strict("**a** [b](https://b.example)").Render(context.Background(), &b); err != nil {
		t.Fatal(err)
	}
	if out := b.String(); !strings.Contains(out, `<div class="prose">`) || !strings.Contains(out, "<strong>a</strong> b") {
		t.Errorf("unexpected %s", out)
	}
}

// Markdown keeps links: Strict must not change the default renderer.
func TestMarkdownKeepsLinks(t *testing.T) {
	if out := HTML("[a](https://a.example) ![i](https://i.example/i.png)"); !strings.Contains(out, `<a href="https://a.example">`) || !strings.Contains(out, "<img") {
		t.Errorf("HTML lost links or images: %s", out)
	}
}
