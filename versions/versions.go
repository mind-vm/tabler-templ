// Package versions pins what this module is built against: the Tabler release, its icons and the
// libraries of Tabler's plugins, all loaded from jsDelivr in these versions. The packages read them from
// here (chart.Version is versions.ApexCharts …), so a Tabler update changes this file only.
//
// `go run ./internal/tools/versions <tabler checkout>` compares them with the versions Tabler's lockfile
// pins (the npm package after each "npm:"); with -write it updates them. scripts/tabler-update.sh runs it.
package versions

const (
	Tabler = "1.7.1" // npm: @tabler/core
	// TablerCommit is the commit of the Tabler release; the demo reads Tabler's images and data at it.
	TablerCommit = "a23c8417799dfae0a890fac78b03ddcb5e365e87"
	Icons        = "3.49.0" // npm: @tabler/icons

	ApexCharts      = "7.0.0"  // npm: apexcharts
	MapLibre        = "6.11.2" // npm: maplibre-gl
	VanillaCalendar = "3.4.0"  // npm: vanilla-calendar-pro
	TomSelect       = "2.6.2"  // npm: tom-select
	IMask           = "7.6.1"  // npm: imask
	NoUiSlider      = "15.8.1" // npm: nouislider
	HugeRTE         = "1.0.14" // npm: hugerte
	Dropzone        = "6.3.5"  // npm: dropzone
	Coloris         = "0.25.0" // npm: @melloware/coloris
	SignaturePad    = "5.1.4"  // npm: signature_pad
	StarRating      = "4.3.1"  // npm: star-rating.js
	FullCalendar    = "6.1.22" // npm: fullcalendar
	ListJS          = "2.3.1"  // npm: list.js
	SortableJS      = "1.15.7" // npm: sortablejs
	Driver          = "1.10.2" // npm: driver.js
	Fslightbox      = "3.8.0"  // npm: fslightbox
	Plyr            = "3.8.5"  // npm: plyr
	Typed           = "3.0.0"  // npm: typed.js
)
