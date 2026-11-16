// Package deadline shows a deadline as a traffic light whose word is the time left (GS-52 of the design
// concept „Gestaltung“): „noch 58 Min.“ in yellow when it is near, „überfällig seit 20 Std.“ in red when it
// has passed, green otherwise. The exact time goes into the title.
//
//	@deadline.Deadline(deadline.Props{Due: o.Due, Now: time.Now()})
package deadline

import (
	"context"
	"fmt"
	"time"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/kit"
)

// DefaultSoon is how near a deadline turns yellow when Props.Soon is 0.
const DefaultSoon = 24 * time.Hour

type Props struct {
	Due, Now time.Time
	Soon     time.Duration // from when on the deadline is near (yellow); default DefaultSoon
	// Layout of the exact time in the title, default kit.T(ctx).DateTime.
	Layout     string
	Class      string
	Attributes templ.Attributes
}

// State of a deadline.
type State string

const (
	Open    State = "open"    // green
	Near    State = "near"    // yellow
	Overdue State = "overdue" // red
)

// StateOf is the deadline's state at now.
func StateOf(due, now time.Time, soon time.Duration) State {
	if soon <= 0 {
		soon = DefaultSoon
	}
	switch left := due.Sub(now); {
	case left < 0:
		return Overdue
	case left <= soon:
		return Near
	}
	return Open
}

func (s State) color() string {
	switch s {
	case Overdue:
		return "red"
	case Near:
		return "yellow"
	}
	return "green"
}

// Word is the time left or overdue, e.g. „noch 58 Min.“ or „überfällig seit 2 Tagen“.
func Word(ctx context.Context, due, now time.Time) string {
	m := kit.T(ctx)
	d := due.Sub(now)
	if d < 0 {
		return fmt.Sprintf(m.DeadlineOverdue, span(m, -d, true))
	}
	return fmt.Sprintf(m.DeadlineLeft, span(m, d, false))
}

// span is a duration in the largest whole unit: minutes below an hour, hours below two days, else days.
// since: the dative German needs after „seit“ („2 Tagen“).
func span(m kit.Messages, d time.Duration, since bool) string {
	switch {
	case d < time.Hour:
		return fmt.Sprintf(m.Minutes, max(1, int(d.Minutes())))
	case d < 48*time.Hour:
		return fmt.Sprintf(m.Hours, int(d.Hours()))
	}
	days := int(d.Hours() / 24)
	if since && m.DaysSince != "" {
		return fmt.Sprintf(m.DaysSince, days)
	}
	return fmt.Sprintf(m.Days, days)
}

func layout(ctx context.Context, p Props) string {
	return kit.Or(p.Layout, kit.T(ctx).DateTime)
}
