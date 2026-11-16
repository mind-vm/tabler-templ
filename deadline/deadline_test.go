package deadline_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/deadline"
	"github.com/mind-vm/tabler-templ/kit"
)

var now = time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)

// GS-52: a deadline is a traffic light whose word is the time left.
func TestGS52_DeadlineWord(t *testing.T) {
	de := kit.WithMessages(context.Background(), kit.German)
	for due, want := range map[time.Time]string{
		now.Add(58 * time.Minute):            "noch 58 Min.",
		now.Add(20 * time.Second):            "noch 1 Min.",
		now.Add(5 * time.Hour):               "noch 5 Std.",
		now.Add(-20 * time.Hour):             "überfällig seit 20 Std.",
		now.Add(3*24*time.Hour + time.Hour):  "noch 3 Tage",
		now.Add(-2*24*time.Hour - time.Hour): "überfällig seit 2 Tagen",
	} {
		assert.Equal(t, want, deadline.Word(de, due, now))
	}
	assert.Equal(t, "overdue by 3 days", deadline.Word(context.Background(), now.Add(-72*time.Hour), now))
}

func TestGS52_DeadlineState(t *testing.T) {
	assert.Equal(t, deadline.Overdue, deadline.StateOf(now.Add(-time.Minute), now, 0))
	assert.Equal(t, deadline.Near, deadline.StateOf(now.Add(23*time.Hour), now, 0))
	assert.Equal(t, deadline.Open, deadline.StateOf(now.Add(25*time.Hour), now, 0))
	assert.Equal(t, deadline.Near, deadline.StateOf(now.Add(3*24*time.Hour), now, 4*24*time.Hour))

	var b strings.Builder
	require.NoError(t, deadline.Deadline(deadline.Props{Due: now.Add(-20 * time.Hour), Now: now}).Render(kit.WithMessages(context.Background(), kit.German), &b))
	h := b.String()
	assert.Contains(t, h, "status-red")
	assert.Contains(t, h, `title="14.11.2026, 13:00 Uhr"`)
	assert.Contains(t, h, `data-deadline="overdue"`)
	assert.Contains(t, h, "überfällig seit 20 Std.")
}
