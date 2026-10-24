package main

import (
	"github.com/mind-vm/tabler-templ/kit"
	"github.com/mind-vm/tabler-templ/progress"
	"github.com/mind-vm/tabler-templ/versions"
)

// The preview's sample photos come from the Tabler repository from raw.githubusercontent.com (pinned to the commit the
// catalog was built from), so this repository holds no images. preview.tabler.io refuses hotlinks.
const staticBase = "https://raw.githubusercontent.com/tabler/tabler/" + versions.TablerCommit + "/shared/static/"

func photo(path string) string { return staticBase + path }

func avatarURL(name string) string { return photo("avatars/" + name + ".jpg") }

// kitAttrs merges attribute maps; later maps win.
var kitAttrs = kit.Attrs

// kitIf returns class when cond holds.
var kitIf = kit.If

func progressSteps5() progress.StepsProps { return progress.StepsProps{Count: 5, Active: 1} }
