package main

import (
	"fmt"
	"time"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/blocks/uptime"
	"github.com/mind-vm/tabler-templ/tracking"
)

// appsCPages are the uptime, activity, tracking and task pages.
func appsCPages() map[string]templ.Component {
	return map[string]templ.Component{
		"uptime":     uptimePage(),
		"activity":   activityPage(),
		"tracking":   trackingPage(),
		"tasks":      tasksPage(),
		"tasks-list": tasksListPage(),
	}
}

// cServices draws the history of the uptime page's System status card the way the preview's script does:
// a seeded random state per day (Mulberry32, seed 18), 2% outages and 4% degraded days.
func cServices(names []string, days int, seed uint32) ([]uptime.Service, bool) {
	random := func() float64 {
		seed += 0x6d2b79f5
		t := (seed ^ seed>>15) * (1 | seed)
		t = (t + (t^t>>7)*(61|t)) ^ t
		return float64(t^t>>14) / 4294967296
	}
	today := time.Now()
	allUp := true
	var out []uptime.Service
	for _, name := range names {
		s := uptime.Service{Name: name}
		score := 0.0
		for i := days - 1; i >= 0; i-- {
			roll := random()
			color, label, value := "green", "Operational", 1.0
			if roll < 0.02 {
				color, label, value = "red", "Outage", 0
			} else if roll < 0.06 {
				color, label, value = "yellow", "Degraded performance", 0.5
			}
			date := today.AddDate(0, 0, -i).Format("Jan 2, 2006")
			s.Days = append(s.Days, tracking.Block{Color: color, Title: date + " – " + label})
			score += value
		}
		if score < float64(days) {
			allUp = false
		}
		s.Uptime = fmt.Sprintf("%.2f%%", score/float64(days)*100)
		out = append(out, s)
	}
	return out, allUp
}

var cIncidents = []uptime.Incident{
	{Color: "red", Title: "Email delivery outage", Date: "Aug 24, 2026", Duration: "42 minutes", Resolved: true},
	{Color: "yellow", Title: "Slow webhook deliveries", Date: "Aug 12, 2026", Duration: "1 hour 8 minutes", Resolved: true},
	{Color: "yellow", Title: "Dashboard loading slowly", Date: "Aug 3, 2026", Duration: "17 minutes", Resolved: true},
	{Color: "red", Title: "Email delivery outage", Date: "Jul 21, 2026", Duration: "2 hours 5 minutes", Resolved: true},
}

var cPeriods = []uptime.Period{
	{Period: "Today", Availability: "100%", Downtime: "none", Incidents: "0", Longest: "–", Average: "–"},
	{Period: "Last 7 days", Availability: "100%", Downtime: "none", Incidents: "0", Longest: "–", Average: "–"},
	{Period: "Last 30 days", Availability: "99.9028%", Downtime: "42 minutes", Incidents: "1", Longest: "42 minutes", Average: "42 minutes"},
	{Period: "Last 365 days", Availability: "99.9414%", Downtime: "5 hours 8 minutes", Incidents: "9", Longest: "2 hours 5 minutes", Average: "34 minutes"},
	{Period: "All time", Availability: "99.9533%", Downtime: "11 hours 40 minutes", Incidents: "23", Longest: "2 hours 5 minutes", Average: "30 minutes"},
}

// cUptime is the history of the tracking page (the docs' example).
var cUptime = func() []tracking.Block {
	ok, down, load, none := tracking.Block{Color: "success", Title: "Operational"}, tracking.Block{Color: "danger", Title: "Downtime"},
		tracking.Block{Color: "warning", Title: "Big load"}, tracking.Block{Title: "No data"}
	return []tracking.Block{ok, ok, ok, ok, down, ok, ok, ok, ok, ok, ok, ok, ok, load, ok, ok, ok, down, ok, ok, ok, none, none, ok, ok, ok, ok, ok, ok, ok}
}()
