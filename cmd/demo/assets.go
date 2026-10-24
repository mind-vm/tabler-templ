package main

import "github.com/mind-vm/tabler-templ/kit"

// The preview's sample photos come from the Tabler repository through jsDelivr (pinned to the commit the
// catalog was built from), so this repository holds no images. preview.tabler.io refuses hotlinks.
const staticBase = "https://cdn.jsdelivr.net/gh/tabler/tabler@a23c8417799dfae0a890fac78b03ddcb5e365e87/shared/static/"

func photo(path string) string { return staticBase + path }

func avatarURL(name string) string { return photo("avatars/" + name + ".jpg") }

// kitAttrs merges attribute maps; later maps win.
var kitAttrs = kit.Attrs
