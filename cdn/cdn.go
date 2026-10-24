// Package cdn is where the browser loads Tabler, its icons, ApexCharts, the plugins and the theme fonts
// from: npm packages on jsDelivr by default, or the app's own copy. Every package builds its URLs with
// URL, so one line moves them all:
//
//	cdn.Base = "/static/npm/" // set before serving pages
//
// cmd/vendor downloads the files into that folder in the same layout (<package>@<version>/<path>), so an
// app can embed them and load nothing from outside at runtime.
//
// Packages say which files they need with Register; cmd/vendor reads the list.
package cdn

import (
	"slices"
	"strings"
	"sync"
)

// JSDelivr is the default Base.
const JSDelivr = "https://cdn.jsdelivr.net/npm/"

// Base is the prefix of every URL: jsDelivr's npm path, or a path or URL of the app's own copy, with a
// trailing slash. Set it once at start-up, before rendering.
var Base = JSDelivr

// URL is a file of an npm package at a version, e.g. URL("apexcharts", "7.0.0", "dist/apexcharts.min.js").
func URL(pkg, version, file string) string {
	base := Base
	if base != "" && !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base + pkg + "@" + version + "/" + strings.TrimPrefix(file, "/")
}

// Files are the files of an npm package one feature of tabler-templ needs.
type Files struct {
	// Group is the feature, as cmd/vendor names it: "tabler", "chart", "datepicker", "theme:rose" …
	Group   string
	Package string
	Version string
	// Paths in the package; "dir/**" stands for everything under dir. Files the CSS files refer to (fonts,
	// images) are found by cmd/vendor.
	Paths []string
}

var (
	mu       sync.Mutex
	registry []Files
)

// Register records the files a group needs; packages call it in init.
func Register(group, pkg, version string, paths ...string) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, Files{Group: group, Package: pkg, Version: version, Paths: paths})
}

// Registered are the files of every group registered, in the order of registration.
func Registered() []Files {
	mu.Lock()
	defer mu.Unlock()
	return slices.Clone(registry)
}

// Groups are the names of the registered groups, sorted.
func Groups() []string {
	var gs []string
	for _, f := range Registered() {
		if !slices.Contains(gs, f.Group) {
			gs = append(gs, f.Group)
		}
	}
	slices.Sort(gs)
	return gs
}
