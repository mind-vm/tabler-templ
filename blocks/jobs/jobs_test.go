package jobs_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mind-vm/tabler-templ/blocks/jobs"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b strings.Builder
	require.NoError(t, c.Render(context.Background(), &b))
	return b.String()
}

func TestJob(t *testing.T) {
	h := render(t, jobs.Job(jobs.JobProps{Logo: "/l.png", Title: "Web Developer", Salary: "$1", Tags: []string{"PHP"},
		Facts: []jobs.Fact{{Icon: "ti-license", Text: "Full-time"}}}))
	assert.True(t, strings.HasPrefix(h, `<div class="card"><div class="row g-0"><div class="col-auto"><div class="card-body"><span class="avatar avatar-md" style="background-image: url(/l.png);"></span>`))
	assert.Contains(t, h, `<h2 class="h3 mb-0"><a href="#">Web Developer</a></h2>`)
	assert.Contains(t, h, `<div class="col-auto fs-3 text-green">$1</div>`)
	assert.Contains(t, h, `list-inline list-inline-dots`)
	assert.Contains(t, h, `<div class="mt-3 list mb-0 text-secondary d-block d-sm-none">`)
	assert.Equal(t, 2, strings.Count(h, "Full-time"))
	assert.Contains(t, h, `<div class="badge-list mt-3">`)
	assert.Contains(t, h, `>PHP</span>`)
}
