// Package sound plays short interface sounds: a new message, a notification, success, a warning, an
// error. The browser makes them itself with the Web Audio API, so there is no library and no sound file
// to load; an app can add its own files (Config.Files).
//
//	<body> … @sound.Script(sound.Config{}) </body>
//	@button.Button(button.Props{Attributes: sound.Attrs(sound.Success)}) { Save }  // plays on click
//	@sound.Play(sound.Message)                                                     // plays when swapped in
//	w.Header().Set("HX-Trigger", `{"sound":"message"}`)                            // plays on the htmx response
//	@sound.Toggle(sound.ToggleProps{})                                             // sounds on or off
//
// Browsers play sound only after the visitor has clicked or typed on the page once; sounds asked for
// before that are dropped. The visitor's choice of Toggle is kept in the browser (localStorage).
//
// In the browser, window.tablerSound.play(name) plays a sound, tablerSound.mute(true|false) switches
// them off or on, and the document gets a "sound" event (detail: the name, or {value: name}) to play one.
package sound

import (
	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/internal/plugin"
)

// Sound is the name of a sound: one of the built-in ones, or one of Config.Files.
type Sound string

const (
	Click   Sound = "click"   // a short tick, for toggles
	Message Sound = "message" // two rising notes: a chat message came in
	Notify  Sound = "notify"  // a soft chime: something needs attention
	Success Sound = "success" // a rising chord: done, saved, sent
	Warning Sound = "warning" // two low beeps
	Error   Sound = "error"   // two falling notes
)

// Event is the DOM event that plays a sound; with htmx, answer with HX-Trigger: {"sound":"message"}.
const Event = "sound"

// Config sets the sounds up for the page.
type Config struct {
	Volume float64 // 0–1; default 0.4
	// Off starts with the sounds off until the visitor switches them on (Toggle); their choice wins
	// once they made one.
	Off bool
	// Files are sounds of the app's own, name → URL (mp3, ogg, wav); a name of a built-in sound
	// replaces it.
	Files map[string]string
}

// Attrs make an element play s when it is clicked.
func Attrs(s Sound) templ.Attributes { return templ.Attributes{"data-sound": string(s)} }

// toggleStyles show the icon of Toggle's state.
var toggleStyles = csp.Sheet("sound", `[data-sound-toggle][aria-pressed=true] .sound-off,`+
	`[data-sound-toggle][aria-pressed=false] .sound-on{display:none}`)

// Script sets the sounds up: it plays them on clicks, on swapped-in Play and on "sound" events, and keeps
// Toggle buttons in step. Render it once in the body; later renders on the page are dropped.
func Script(c Config) templ.Component {
	if c.Volume <= 0 || c.Volume > 1 {
		c.Volume = 0.4
	}
	if c.Files == nil {
		c.Files = map[string]string{}
	}
	cfg := plugin.JSON(map[string]any{"volume": c.Volume, "off": c.Off, "files": c.Files})
	return plugin.Start("sound", `(function (cfg) {
	if (window.tablerSound) return window.tablerSound.setup;
	var AC = window.AudioContext || window.webkitAudioContext, ctx = null, buffers = {}, KEY = 'tabler-sound';
	var saved = null;
	try { saved = localStorage.getItem(KEY); } catch (e) {}
	var muted = saved ? saved === 'off' : !!cfg.off;
	// Each tone: frequency (Hz), start and length (s), wave, loudness (0–1).
	var tones = {
		click: [[1400, 0, 0.03, 'sine', 0.3]],
		message: [[880, 0, 0.12, 'sine', 0.9], [1320, 0.1, 0.2, 'sine', 0.9]],
		notify: [[660, 0, 0.4, 'triangle', 1], [990, 0.12, 0.5, 'triangle', 0.8]],
		success: [[523.25, 0, 0.14, 'triangle', 1], [659.25, 0.09, 0.14, 'triangle', 1], [783.99, 0.18, 0.32, 'triangle', 1]],
		warning: [[440, 0, 0.12, 'square', 0.3], [440, 0.18, 0.12, 'square', 0.3]],
		error: [[311, 0, 0.16, 'sawtooth', 0.3], [233, 0.14, 0.3, 'sawtooth', 0.3]]
	};
	function load() {
		Object.keys(cfg.files).forEach(function (name) {
			fetch(cfg.files[name]).then(function (r) { return r.arrayBuffer(); })
				.then(function (b) { return ctx.decodeAudioData(b); })
				.then(function (buf) { buffers[name] = buf; }).catch(function () {});
		});
	}
	function unlock() {
		if (!AC) return;
		if (!ctx) { ctx = new AC(); load(); }
		if (ctx.state === 'suspended') ctx.resume();
	}
	['pointerdown', 'keydown', 'touchstart'].forEach(function (t) { document.addEventListener(t, unlock, { capture: true, passive: true }); });
	function play(name) {
		if (muted || !name || !ctx || ctx.state !== 'running') return false;
		var out = ctx.createGain();
		out.gain.value = cfg.volume;
		out.connect(ctx.destination);
		if (buffers[name]) {
			var src = ctx.createBufferSource();
			src.buffer = buffers[name];
			src.connect(out);
			src.start();
			return true;
		}
		var t = tones[name];
		if (!t) return false;
		var now = ctx.currentTime;
		t.forEach(function (n) {
			var o = ctx.createOscillator(), g = ctx.createGain(), at = now + n[1];
			o.type = n[3];
			o.frequency.value = n[0];
			g.gain.setValueAtTime(0.0001, at);
			g.gain.exponentialRampToValueAtTime(n[4], at + 0.01);
			g.gain.exponentialRampToValueAtTime(0.0001, at + n[2]);
			o.connect(g);
			g.connect(out);
			o.start(at);
			o.stop(at + n[2] + 0.05);
		});
		return true;
	}
	function each(root, sel, fn) {
		if (root.matches && root.matches(sel)) fn(root);
		if (root.querySelectorAll) root.querySelectorAll(sel).forEach(fn);
	}
	function sync(root) {
		each(root, '[data-sound-toggle]', function (b) { b.setAttribute('aria-pressed', String(!muted)); });
	}
	function mute(m) {
		muted = !!m;
		try { localStorage.setItem(KEY, muted ? 'off' : 'on'); } catch (e) {}
		sync(document);
		document.dispatchEvent(new CustomEvent('sound:change', { detail: { muted: muted } }));
	}
	document.addEventListener('click', function (e) {
		if (!e.target.closest) return;
		var t = e.target.closest('[data-sound-toggle]');
		if (t) {
			mute(!muted);
			play('click');
			return;
		}
		var s = e.target.closest('[data-sound]');
		if (s) play(s.getAttribute('data-sound'));
	});
	document.addEventListener('`+Event+`', function (e) {
		var d = e.detail;
		play(typeof d === 'string' ? d : d && d.value);
	});
	function setup(root) {
		sync(root);
		each(root, 'template[data-sound-play]', function (el) {
			play(el.getAttribute('data-sound-play'));
			el.remove();
		});
	}
	window.tablerSound = { play: play, mute: mute, muted: function () { return muted; }, setup: setup };
	return setup;
})(`+cfg+`)`)
}

// Styles is the CSS of Toggle; Toggle renders it itself.
func Styles() templ.Component { return toggleStyles }
