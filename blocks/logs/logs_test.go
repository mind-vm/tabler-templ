package logs_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/logs"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestLogs(t *testing.T) {
	h := render(t, logs.Pairs(0, []logs.Pair{{Key: "Effective URL", Value: "https://x.org", Href: "https://x.org"}, {Key: "Redirect <count>", Value: "0"}}))
	assert.Equal(t, "<pre>Effective URL" + strings.Repeat("\u00a0", 12) + `<a class="text-reset" target="_blank" href="https://x.org">https://x.org</a><br>Redirect &lt;count&gt;` + strings.Repeat("\u00a0", 9) + "0</pre>", h)
	assert.Equal(t, "<pre><code>a\nb</code></pre>", render(t, logs.Code("a\nb")))
	assert.Equal(t, "<pre>HTTP/1.1 200</pre>", render(t, logs.Text("HTTP/1.1 200")))
	assert.Contains(t, render(t, logs.Section("Checked URL")), `<h2>Checked URL</h2><div>`)
}
