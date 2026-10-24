// Command demo rebuilds the Tabler preview (preview.tabler.io) page by page with this module's
// components and blocks. The index lists every preview page; rebuilt pages are linked, open ones grey.
//
//	go run ./cmd/demo -addr :8080
package main

import (
	"encoding/json"
	"flag"
	"log"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/catalog"
)

// Tabler's CSS and JS and the icon font come from a CDN; the module itself ships none of them.
const (
	tablerVersion = "1.7.1"
	iconsVersion  = "3.49.0"
)

// pages maps a preview page name to its rebuild. Keep catalog.Pages' Status in step with it.
// It is filled in init: the pages' menus link to the rebuilt pages, so they refer to the map.
var pages map[string]templ.Component

func init() {
	pages = map[string]templ.Component{
		"error-404":                error404Page(),
		"error-500":                error500Page(),
		"error-maintenance":        maintenancePage(),
		"page-loader":              pageLoaderPage(),
		"sign-in-cover":            signInCoverPage(),
		"sign-in-illustration":     signInIllustrationPage(),
		"sign-in-link":             signInLinkPage(),
		"sign-up":                  signUpPage(),
		"forgot-password":          forgotPasswordPage(),
		"auth-lock":                authLockPage(),
		"2-step-verification":      twoStepPage(),
		"2-step-verification-code": twoStepCodePage(),
		"trial-ended":              trialEndedPage(),
		"onboarding":               onboardingPage(),
		"alerts":                   alertsPage(),
		"avatars":                  avatarsPage(),
		"badges":                   badgesPage(),
		"buttons":                  buttonsPage(),
		"dropdowns":                dropdownsPage(),
		"empty":                    emptyPage(),
		"number-input":             numberInputPage(),
		"offcanvas":                offcanvasPage(),
		"tabs":                     tabsPage(),
		"toasts":                   toastsPage(),
		"modals":                   modalsPage(),
		"pagination":               paginationPage(),
		"password-strength":        passwordStrengthPage(),
		"progress":                 progressPage(),
		"sign-in":                  signInPage(),
		"steps":                    stepsPage(),
	}
	for _, batch := range []map[string]templ.Component{batchAPages(), batchBPages(), batchCPages(), batchDPages(), batchEPages(), batchFPages(), batchGPages(), batchHPages()} {
		for k, v := range batch {
			pages[k] = v
		}
	}
	for name := range layoutVariants {
		pages[name] = layoutPage(name)
	}
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(index()))
	// pages.json lists the rebuilt pages, for the screenshot script.
	mux.HandleFunc("GET /pages.json", func(w http.ResponseWriter, r *http.Request) {
		names := slices.Sorted(maps.Keys(pages))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(names)
	})
	mux.HandleFunc("GET /{name...}", func(w http.ResponseWriter, r *http.Request) {
		c, ok := pages[strings.TrimSuffix(r.PathValue("name"), ".html")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		templ.Handler(c).ServeHTTP(w, r)
	})
	return mux
}

func main() {
	addr := flag.String("addr", "localhost:8080", "listen address")
	flag.Parse()
	log.Printf("tabler-templ demo on http://%s (%d of %d preview pages)", *addr, len(pages), len(catalog.Pages))
	log.Fatal(http.ListenAndServe(*addr, handler()))
}
