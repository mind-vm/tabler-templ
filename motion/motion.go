// Package motion animates elements as they enter and leave the page with Motion (motion.dev, formerly
// Framer Motion), its vanilla build without React. Attrs make an element (or its children, one after
// another) fade, slide or pop in when the page loads, when htmx swaps it in, or when it scrolls into
// view; Counter counts a number up. Load the library with Library in the head.
//
//	<head> … @motion.Library() </head>
//	<div class="row row-cards" { motion.Attrs(motion.Props{Stagger: 0.06, InView: true})... }> … </div>
//	@motion.Counter(motion.CounterProps{To: 4280, Prefix: "$"})
//
// Elements leave with an animation when htmx swaps them out (Props.Exit); give the swap a delay at least
// as long as the animation, e.g. hx-swap="outerHTML swap:200ms".
//
// Visitors who ask for reduced motion (prefers-reduced-motion) see everything at once, without
// animations; so does a page where the library failed to load.
package motion

import (
	"strconv"

	"github.com/a-h/templ"

	"github.com/mind-vm/tabler-templ/cdn"
	"github.com/mind-vm/tabler-templ/internal/csp"
	"github.com/mind-vm/tabler-templ/internal/plugin"
	"github.com/mind-vm/tabler-templ/versions"
)

// Version is the Motion version Library loads.
const Version = versions.Motion

// LibraryURL is Motion's script build (window.Motion) on cdn.Base.
func LibraryURL() string { return cdn.URL("motion", Version, "dist/motion.js") }

// Preset is how an element enters (or, as Exit, leaves).
type Preset string

const (
	FadeIn     Preset = "fade-in"
	FadeUp     Preset = "fade-up"   // rises a little while it fades in; the default
	FadeDown   Preset = "fade-down" // drops a little while it fades in
	SlideLeft  Preset = "slide-left"
	SlideRight Preset = "slide-right"
	ZoomIn     Preset = "zoom-in" // grows slightly
	Pop        Preset = "pop"     // grows from small with a spring
)

// Props say how an element enters and leaves. The zero value fades it up when it is rendered.
type Props struct {
	Preset   Preset  // how it enters; default FadeUp. Empty with Exit set: it only leaves animated
	Duration float64 // seconds; default 0.4
	Delay    float64 // seconds before it starts
	Stagger  float64 // seconds between the element's children: they enter one after another, not the element
	Bounce   float64 // 0–1: a spring with that much bounce instead of an ease (Pop: 0.4 unless set)
	InView   bool    // enter when it scrolls into view, not when it is rendered
	Exit     Preset  // how it leaves when htmx swaps it (or the element around it) out; needs a swap delay
	// ExitDuration is the exit's length in seconds; default 0.2. The swap delay must be at least as long.
	ExitDuration float64
}

// Attrs animate an element as Props say.
func Attrs(p Props) templ.Attributes {
	a := templ.Attributes{}
	preset := p.Preset
	if preset == "" && p.Exit == "" {
		preset = FadeUp
	}
	if preset != "" {
		a["data-motion"] = string(preset)
		if p.Duration > 0 {
			a["data-motion-duration"] = num(p.Duration)
		}
		if p.Delay > 0 {
			a["data-motion-delay"] = num(p.Delay)
		}
		if p.Stagger > 0 {
			a["data-motion-stagger"] = num(p.Stagger)
		}
		bounce := p.Bounce
		if bounce == 0 && preset == Pop {
			bounce = 0.4
		}
		if bounce > 0 {
			a["data-motion-bounce"] = num(bounce)
		}
		if p.InView {
			a["data-motion-inview"] = true
		}
	}
	if p.Exit != "" {
		a["data-motion-exit"] = string(p.Exit)
		if p.ExitDuration > 0 {
			a["data-motion-exit-duration"] = num(p.ExitDuration)
		}
	}
	return a
}

// CounterProps describe a number that counts up from zero.
type CounterProps struct {
	To       float64
	Decimals int
	Prefix   string  // e.g. "$"
	Suffix   string  // e.g. "%"
	Duration float64 // seconds; default 1.2
	Delay    float64
	InView   bool // count when it scrolls into view, not when it is rendered
	Class    string
}

// counterAttrs are the data attributes Counter puts on its span.
func counterAttrs(p CounterProps) templ.Attributes {
	attrs := templ.Attributes{
		"data-motion-count":    num(p.To),
		"data-motion-decimals": strconv.Itoa(max(p.Decimals, 0)),
	}
	if p.Prefix != "" {
		attrs["data-motion-prefix"] = p.Prefix
	}
	if p.Suffix != "" {
		attrs["data-motion-suffix"] = p.Suffix
	}
	if p.Duration > 0 {
		attrs["data-motion-duration"] = num(p.Duration)
	}
	if p.Delay > 0 {
		attrs["data-motion-delay"] = num(p.Delay)
	}
	if p.InView {
		attrs["data-motion-inview"] = true
	}
	return attrs
}

// counterText is the counter's final text, shown before (or without) the animation.
func counterText(p CounterProps) string {
	return p.Prefix + strconv.FormatFloat(p.To, 'f', max(p.Decimals, 0), 64) + p.Suffix
}

func num(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// hidden keeps elements that are still to enter invisible, only while the animations run (the script
// sets motion-on on <html>), so the page never flashes them before they animate.
var hidden = csp.Sheet("motion", `.motion-on [data-motion]:not([data-motion-done]):not([data-motion-stagger]),`+
	`.motion-on [data-motion-stagger]:not([data-motion-done])>*{opacity:0}`)

// setup is the plugin script: it marks <html> (when the animations run), registers the exit animations
// once, and returns the function that starts the elements below a root.
const setup = `(function () {
	if (window.tablerMotion) return window.tablerMotion;
	var M = window.Motion;
	var reduce = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
	var on = !!(M && M.animate) && !reduce;
	if (on) document.documentElement.classList.add('motion-on');
	var presets = {
		'fade-in': { opacity: [0, 1] },
		'fade-up': { opacity: [0, 1], y: [16, 0] },
		'fade-down': { opacity: [0, 1], y: [-16, 0] },
		'slide-left': { opacity: [0, 1], x: [32, 0] },
		'slide-right': { opacity: [0, 1], x: [-32, 0] },
		'zoom-in': { opacity: [0, 1], scale: [0.94, 1] },
		'pop': { opacity: [0, 1], scale: [0.6, 1] }
	};
	var ease = [0.22, 1, 0.36, 1];
	function num(el, name, d) { var v = parseFloat(el.getAttribute(name)); return isNaN(v) ? d : v; }
	function keyframes(name, back) {
		var k = presets[name] || presets['fade-up'], r = {};
		for (var p in k) r[p] = back ? k[p].slice().reverse() : k[p];
		return r;
	}
	function options(el, delay) {
		var d = num(el, 'data-motion-duration', 0.4), b = num(el, 'data-motion-bounce', 0);
		return b > 0 ? { type: 'spring', bounce: b, visualDuration: d, delay: delay } : { duration: d, ease: ease, delay: delay };
	}
	function enter(el) {
		var k = keyframes(el.getAttribute('data-motion')), delay = num(el, 'data-motion-delay', 0), s = num(el, 'data-motion-stagger', 0);
		if (s > 0) M.animate(Array.prototype.slice.call(el.children), k, options(el, M.stagger(s, { startDelay: delay })));
		else M.animate(el, k, options(el, delay));
		el.setAttribute('data-motion-done', '');
	}
	function formatter(el) {
		var d = parseInt(el.getAttribute('data-motion-decimals') || '0', 10);
		var f = new Intl.NumberFormat(document.documentElement.lang || undefined, { minimumFractionDigits: d, maximumFractionDigits: d });
		var pre = el.getAttribute('data-motion-prefix') || '', suf = el.getAttribute('data-motion-suffix') || '';
		return function (v) { return pre + f.format(v) + suf; };
	}
	function count(el) {
		var f = formatter(el), to = parseFloat(el.getAttribute('data-motion-count')) || 0;
		M.animate(0, to, { duration: num(el, 'data-motion-duration', 1.2), delay: num(el, 'data-motion-delay', 0), ease: ease,
			onUpdate: function (v) { el.textContent = f(v); } });
		el.setAttribute('data-motion-done', '');
	}
	function start(el, run) {
		if (el.hasAttribute('data-motion-done') || el.hasAttribute('data-motion-wait')) return;
		el.setAttribute('data-motion-wait', '');
		if (el.hasAttribute('data-motion-inview')) M.inView(el, function () { run(el); }, { amount: 0.2 });
		else run(el);
	}
	function each(root, sel, fn) {
		if (root.matches && root.matches(sel)) fn(root);
		if (root.querySelectorAll) root.querySelectorAll(sel).forEach(fn);
	}
	if (on) document.addEventListener('htmx:beforeSwap', function (e) {
		if (!e.detail.shouldSwap || !e.detail.target) return;
		// Only swaps that replace content take elements away; beforeend, afterbegin … only add.
		var src = e.detail.requestConfig && e.detail.requestConfig.elt, holder = src && src.closest && src.closest('[hx-swap]');
		var style = (e.detail.swapOverride || (holder && holder.getAttribute('hx-swap')) || (window.htmx && htmx.config.defaultSwapStyle) || 'innerHTML').split(/\s+/)[0];
		if (style !== 'innerHTML' && style !== 'outerHTML' && style !== 'delete') return;
		var target = e.detail.target;
		each(target, '[data-motion-exit]', function (el) {
			if (el === target && style === 'innerHTML') return;
			M.animate(el, keyframes(el.getAttribute('data-motion-exit'), true), { duration: num(el, 'data-motion-exit-duration', 0.2), ease: 'easeIn' });
		});
	});
	window.tablerMotion = function (root) {
		if (!on) return;
		each(root, '[data-motion]', function (el) { start(el, enter); });
		each(root, '[data-motion-count]', function (el) {
			if (!el.hasAttribute('data-motion-done') && !el.hasAttribute('data-motion-wait')) el.textContent = formatter(el)(0);
			start(el, count);
		});
	};
	return window.tablerMotion;
})()`

// Script starts the animated elements on the page and in content swapped in by htmx (Library renders it).
func Script() templ.Component {
	return templ.Join(hidden, plugin.Start("motion", setup))
}

// Library loads Motion and starts the animations; put it in the head.
func Library() templ.Component {
	return templ.Join(plugin.Script(LibraryURL()), Script())
}

func init() {
	cdn.Register("motion", "motion", Version, "dist/motion.js")
}
