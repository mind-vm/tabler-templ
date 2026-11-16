package detail

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/kit"
)

// GS-44: a typed confirmation matches regardless of case and spaces, never an empty word.
func TestTypedMatches(t *testing.T) {
	for _, c := range []struct {
		typed, expect string
		want          bool
	}{
		{"Wohnungsgesellschaft Nordstadt", "Wohnungsgesellschaft Nordstadt", true},
		{"  wohnungsgesellschaft   nordstadt ", "Wohnungsgesellschaft Nordstadt", true},
		{"ABSCHLIESSEN", "abschließen", false},
		{"Wohnungsgesellschaft", "Wohnungsgesellschaft Nordstadt", false},
		{"", "", false},
	} {
		assert.Equal(t, c.want, TypedMatches(c.typed, c.expect), "%q / %q", c.typed, c.expect)
	}
}

func TestTypedConfirm(t *testing.T) {
	ctx := kit.WithMessages(templ.InitializeContext(context.Background()), kit.German)
	var b strings.Builder
	require.NoError(t, TypedConfirm(TypedConfirmProps{ID: "x", Expect: "Nordstadt", Error: TypedError(ctx, "Nordstadt")}).Render(ctx, &b))
	h := b.String()
	for _, want := range []string{`data-kit-typed="Nordstadt"`, `name="confirm"`, `autocomplete="off"`, "Zur Bestätigung „Nordstadt“ tippen",
		"Bitte genau so tippen: Nordstadt", "is-invalid"} {
		assert.Contains(t, h, want)
	}
	b.Reset()
	require.NoError(t, Script(ScriptConfig{}).Render(ctx, &b))
	assert.Contains(t, b.String(), "input[data-kit-typed]", "the script disables the buttons until the field matches")
}
