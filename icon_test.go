package tablertempl_test

import (
	"context"
	"github.com/mind-vm/tabler-templ/icon"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilledIcon(t *testing.T) {
	h := render(t, context.Background(), icon.Icon("ti-star-filled", icon.Props{Color: "yellow"}))
	assert.Equal(t, "<i class=\"ti ti-filled icon text-yellow\" data-icon=\"ti-star\" aria-hidden=\"true\"></i>", h)
	assert.Equal(t, render(t, context.Background(), icon.Icon("ti-star", icon.Props{Filled: true})), "<i class=\"ti ti-filled icon\" data-icon=\"ti-star\" aria-hidden=\"true\"></i>")
	assert.Contains(t, render(t, context.Background(), icon.Icon("ti-no-such-filled")), "ti-no-such-filled")
}
