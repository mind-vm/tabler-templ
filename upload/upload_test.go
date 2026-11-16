package upload_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/upload"
)

func render(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(templ.InitializeContext(ctx), &b))
	return b.String()
}

var de = kit.WithMessages(context.Background(), kit.German)

func TestButton(t *testing.T) {
	h := render(t, de, upload.Button(upload.ButtonProps{ID: "f", Name: "file", Accept: ".pdf", Required: true}))
	assert.Contains(t, h, `<label class="btn mb-0" for="f">`)
	assert.Contains(t, h, "Datei wählen")
	assert.Regexp(t, `<input class="kit-upload-input" type="file" id="f" name="file" accept=".pdf" required>`, h)
	assert.Contains(t, h, `data-none="Keine Datei gewählt" aria-live="polite">Keine Datei gewählt</span>`)

	h = render(t, de, upload.Button(upload.ButtonProps{ID: "p", Name: "file", Camera: true, Submit: true, Primary: true}))
	assert.Contains(t, h, "Foto aufnehmen")
	assert.Contains(t, h, `capture="environment"`)
	assert.Contains(t, h, "data-kit-upload-submit")
	assert.Contains(t, h, "btn btn-primary")
	assert.NotContains(t, h, "data-kit-upload-name", "sent at once: no name to show")
}

// GS-44: a file dropped on the record goes into the dialog's form, which opens to ask for the rest.
func TestDrop(t *testing.T) {
	h := render(t, de, upload.Drop(upload.DropProps{Form: "file-add", Open: true}))
	assert.Contains(t, h, `data-kit-drop="file-add" data-kit-drop-open data-kit-drop-hint="Zum Hochladen loslassen"`)
	s := render(t, de, templ.Join(upload.Script(), upload.Script()))
	assert.Equal(t, 1, strings.Count(s, "<script"))
	assert.Contains(t, s, "field.files = ev.dataTransfer.files")
}
