package sound_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/sound"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var b bytes.Buffer
	if err := c.Render(templ.InitializeContext(context.Background()), &b); err != nil {
		t.Fatal(err)
	}
	return b.String()
}

func TestAttrs(t *testing.T) {
	if a := sound.Attrs(sound.Success); a["data-sound"] != "success" || len(a) != 1 {
		t.Error(a)
	}
}

func TestPlay(t *testing.T) {
	if html := render(t, sound.Play(sound.Message)); html != `<template data-sound-play="message"></template>` {
		t.Error(html)
	}
}

func TestToggle(t *testing.T) {
	html := render(t, templ.Join(sound.Toggle(sound.ToggleProps{Text: "Sound"}), sound.Toggle(sound.ToggleProps{Class: "nav-link", Label: "Töne"})))
	for _, want := range []string{`data-sound-toggle`, `aria-label="Sounds"`, `class="btn btn-icon"`, `class="nav-link"`, `aria-label="Töne"`,
		`ti-volume-off`, `sound-off`, `>Sound</span>`} {
		if !strings.Contains(html, want) {
			t.Errorf("no %s in %s", want, html)
		}
	}
	if n := strings.Count(html, "[data-sound-toggle][aria-pressed=true]"); n != 1 {
		t.Errorf("CSS %d times", n)
	}
}

// TestScript renders the script twice: it comes once, with the configuration of the first.
func TestScript(t *testing.T) {
	html := render(t, templ.Join(
		sound.Script(sound.Config{Volume: 0.7, Off: true, Files: map[string]string{"ding": "/ding.mp3"}}),
		sound.Script(sound.Config{})))
	if n := strings.Count(html, "window.tablerSound = {"); n != 1 {
		t.Errorf("script %d times", n)
	}
	if !strings.Contains(html, `{"files":{"ding":"/ding.mp3"},"off":true,"volume":0.7}`) {
		t.Error(html[len(html)-200:])
	}
	if !strings.Contains(html, `document.addEventListener('sound'`) {
		t.Error("no sound event listener")
	}
}
