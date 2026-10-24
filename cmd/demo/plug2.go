package main

import (
	"math"
	"time"

	"github.com/a-h/templ"
)

// plugin2Pages are the pages of the form and interface plugins.
func plugin2Pages() map[string]templ.Component {
	return map[string]templ.Component{
		"datepicker":   datepickerPage(),
		"dropzone":     dropzonePage(),
		"colorpicker":  colorpickerPage(),
		"wysiwyg":      wysiwygPage(),
		"fullcalendar": fullcalendarPage(),
		"stars-rating": starsRatingPage(),
		"signatures":   signaturesPage(),
		"confetti":     confettiPage(),
		"datatables":   datatablesPage(),
		"sortable":     sortablePage(),
		"tour":         tourPage(),
	}
}

// randomNumber is the preview's pseudo-random number for index x (shared/lib/pseudo-random.ts).
func randomNumber(x, min, max float64) float64 {
	v := math.Mod(x*x*math.Pi*math.E*(max+1)*(math.Sin(x)/math.Cos(x*x)), max+1-min) + min
	v = math.Min(math.Max(v, min), max)
	return math.Floor(v)
}

// randomDate is the preview's pseudo-random date of 2024 for index x.
func randomDate(x float64) time.Time {
	start := float64(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	end := float64(time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC).Unix())
	return time.Unix(int64(randomNumber(x, start, end)), 0).UTC()
}
