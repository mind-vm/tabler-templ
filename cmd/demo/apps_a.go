package main

import (
	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/progress"
)

// appsAPages are the dashboards: index (also shown by the layout-* pages) and the CRM dashboard.
func appsAPages() map[string]templ.Component {
	return map[string]templ.Component{
		"index":         indexPage(),
		"dashboard-crm": crmPage(),
	}
}

var progressProps75 = progress.Props{Value: 75, Color: "primary"}
