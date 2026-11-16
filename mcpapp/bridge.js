(function () {
	'use strict';
	// The bridge between a tabler-templ page in an MCP Apps sandbox and its host: JSON-RPC 2.0 over
	// postMessage. It shakes hands (ui/initialize, ui/notifications/initialized), applies the host's theme,
	// reports the page's size, swaps the HTML of tool results in and turns data-mcp-* attributes into
	// requests (tools/call, ui/message, ui/open-link). Package mcpapp documents the attributes.
	if (window.mcpApp) return;
	var me = document.currentScript;
	var KEY = me.dataset.key, VERSION = me.dataset.protocol, hostStyles = me.hasAttribute('data-host-styles');
	var host = window.parent;
	var standalone = host === window;
	var seq = 0, pending = {}, context = {};

	var emit = function (el, name, detail) {
		(el || document).dispatchEvent(new CustomEvent('mcp:' + name, { bubbles: true, detail: detail }));
	};
	var post = function (msg) { msg.jsonrpc = '2.0'; host.postMessage(msg, '*'); };
	var send = function (method, params) {
		if (standalone) return Promise.reject(new Error('mcpapp: not inside an MCP Apps host'));
		return new Promise(function (resolve, reject) {
			var id = ++seq;
			pending[id] = { resolve: resolve, reject: reject };
			post({ id: id, method: method, params: params || {} });
		});
	};
	var notify = function (method, params) { if (!standalone) post({ method: method, params: params || {} }); };

	// Theme: Tabler's color mode from the host's theme; with HostStyles the host's colors, font and radius
	// for Tabler's variables (the host's own variables are set too, for the app's CSS).
	var tablerVars = {
		'--color-background-primary': ['--tblr-body-bg', '--tblr-bg-surface'],
		'--color-background-secondary': ['--tblr-bg-surface-secondary'],
		'--color-text-primary': ['--tblr-body-color'],
		'--color-text-secondary': ['--tblr-secondary', '--tblr-secondary-color'],
		'--color-border-primary': ['--tblr-border-color'],
		'--font-sans': ['--tblr-font-sans-serif', '--tblr-body-font-family'],
		'--font-mono': ['--tblr-font-monospace'],
		'--border-radius-md': ['--tblr-border-radius']
	};
	var fontsStyle, applied = [];
	var applyContext = function (c) {
		if (!c) return;
		for (var k in c) context[k] = c[k];
		var root = document.documentElement;
		if (c.theme === 'light' || c.theme === 'dark') {
			root.setAttribute('data-bs-theme', c.theme);
			root.style.colorScheme = c.theme;
		}
		if (c.locale && !root.getAttribute('lang')) root.setAttribute('lang', c.locale);
		if (!hostStyles || !c.styles) return;
		// A new set of variables replaces the last one.
		applied.forEach(function (name) { root.style.removeProperty(name); });
		applied = [];
		var vars = c.styles.variables || {};
		var set = function (name, value) { root.style.setProperty(name, value); applied.push(name); };
		Object.keys(vars).forEach(function (name) {
			if (vars[name] == null) return;
			set(name, vars[name]);
			(tablerVars[name] || []).forEach(function (t) { set(t, 'var(' + name + ')'); });
		});
		var fonts = c.styles.css && c.styles.css.fonts;
		if (fonts) {
			if (!fontsStyle) { fontsStyle = document.createElement('style'); document.head.appendChild(fontsStyle); }
			fontsStyle.textContent = fonts;
		}
	};

	// Size: the host sizes the iframe from these notifications (flexible containers).
	var last = '';
	var measure = function () {
		var r = document.querySelector('[data-mcp-root]') || document.body;
		var b = r.getBoundingClientRect();
		var cs = getComputedStyle(document.body);
		var w = Math.ceil(Math.max(b.width, r.scrollWidth)), h = Math.ceil(b.height + parseFloat(cs.marginTop) + parseFloat(cs.marginBottom));
		var key = w + 'x' + h;
		if (key === last) return;
		last = key;
		notify('ui/notifications/size-changed', { width: w, height: h });
	};
	var queued = false;
	var resized = function () { if (!queued) { queued = true; requestAnimationFrame(function () { queued = false; measure(); }); } };

	// Results: the HTML from _meta (or structuredContent) under KEY, or an alert with the error's text.
	var esc = function (s) { var d = document.createElement('div'); d.textContent = s; return d.innerHTML; };
	var textOf = function (r) {
		return ((r && r.content) || []).filter(function (c) { return c.type === 'text'; }).map(function (c) { return c.text; }).join('\n');
	};
	var htmlOf = function (r) {
		if (!r) return null;
		if (r.isError) return '<div class="alert alert-danger" role="alert">' + esc(textOf(r) || 'Error') + '</div>';
		var h = (r._meta && r._meta[KEY]) || (r.structuredContent && r.structuredContent[KEY]);
		return typeof h === 'string' ? h : null;
	};
	var swap = function (el, html, outer) {
		if (!el || html == null) return;
		var parent = el.parentNode;
		if (outer) {
			var t = document.createElement('template');
			t.innerHTML = html;
			var nodes = Array.prototype.slice.call(t.content.childNodes);
			el.replaceWith(t.content);
			nodes.forEach(function (n) { if (n.nodeType === 1) emit(n, 'swapped'); });
		} else {
			el.innerHTML = html;
			emit(el, 'swapped');
		}
		if (parent) resized();
	};
	var showResult = function (r) {
		swap(document.querySelector('[data-mcp-result]'), htmlOf(r), false);
	};

	// Messages from the host: responses to our requests, requests (teardown, ping) and notifications.
	window.addEventListener('message', function (e) {
		if (e.source !== host) return;
		var m = e.data;
		if (!m || m.jsonrpc !== '2.0') return;
		if (m.method == null) {
			var p = pending[m.id];
			if (!p) return;
			delete pending[m.id];
			if (m.error) p.reject(Object.assign(new Error(m.error.message || 'error'), { code: m.error.code, data: m.error.data }));
			else p.resolve(m.result);
			return;
		}
		var params = m.params || {};
		switch (m.method) {
			case 'ui/notifications/tool-input': emit(null, 'tool-input', params.arguments || {}); break;
			case 'ui/notifications/tool-input-partial': emit(null, 'tool-input-partial', params.arguments || {}); break;
			case 'ui/notifications/tool-result': showResult(params); emit(null, 'tool-result', params); break;
			case 'ui/notifications/tool-cancelled': emit(null, 'tool-cancelled', params); break;
			case 'ui/notifications/host-context-changed': applyContext(params); emit(null, 'host-context', context); break;
		}
		if (m.id != null) {
			if (m.method === 'ui/resource-teardown') emit(null, 'teardown', params);
			if (m.method === 'ui/resource-teardown' || m.method === 'ping') post({ id: m.id, result: {} });
			else post({ id: m.id, error: { code: -32601, message: 'Method not found' } });
		}
	});

	var ready = send('ui/initialize', {
		protocolVersion: VERSION,
		appInfo: { name: me.dataset.name || document.title || 'tabler-templ', version: me.dataset.version || '0' },
		appCapabilities: {}
	}).then(function (r) {
		applyContext(r && r.hostContext);
		notify('ui/notifications/initialized', {});
		if (window.ResizeObserver) new ResizeObserver(resized).observe(document.querySelector('[data-mcp-root]') || document.body);
		measure();
		emit(null, 'initialized', r);
		return r;
	});
	if (standalone) {
		ready.catch(function () {});
		console.info('mcpapp: opened outside an MCP Apps host; tools, prompts and links do nothing');
	}

	// The calls the attributes make; also window.mcpApp for the app's own scripts.
	var call = function (name, args) { return ready.then(function () { return send('tools/call', { name: name, arguments: args || {} }); }); };
	var message = function (text) { return ready.then(function () { return send('ui/message', { role: 'user', content: [{ type: 'text', text: text }] }); }); };
	var openLink = function (url) { return ready.then(function () { return send('ui/open-link', { url: url }); }); };

	var busy = function (el, on) {
		if (!el) return;
		el.toggleAttribute('disabled', on);
		el.setAttribute('aria-busy', on ? 'true' : 'false');
		if (el.classList.contains('btn')) el.classList.toggle('btn-loading', on);
	};
	var json = function (s) { try { return s ? JSON.parse(s) : {}; } catch (_) { return {}; } };
	var fields = function (form, submitter) {
		var out = {};
		var fd;
		try { fd = new FormData(form, submitter); } catch (_) { fd = new FormData(form); }
		fd.forEach(function (v, k) {
			var el = form.elements[k];
			var input = el && el.length && !el.tagName ? el[0] : el;
			if (input && (input.type === 'number' || input.type === 'range') && v !== '') v = Number(v);
			else if (input && input.type === 'checkbox' && v === 'on' && !input.hasAttribute('value')) v = true;
			if (k in out) out[k] = [].concat(out[k], v);
			else out[k] = v;
		});
		return out;
	};
	var run = function (el, button, args) {
		var sel = el.getAttribute('data-mcp-target');
		var name = el.getAttribute('data-mcp-tool');
		busy(button, true);
		return call(name, args).then(function (r) {
			var target = sel ? document.querySelector(sel) : null;
			swap(target, htmlOf(r), el.getAttribute('data-mcp-swap') === 'outer');
			emit(target && target.isConnected ? target : document, 'tool-result', { name: name, result: r });
			return r;
		}, function (err) {
			var target = sel ? document.querySelector(sel) : null;
			if (target && !standalone) swap(target, '<div class="alert alert-danger" role="alert">' + esc(err.message) + '</div>', false);
			emit(el, 'error', { name: name, error: err });
		}).then(function (r) { if (button && button.isConnected) busy(button, false); return r; });
	};

	document.addEventListener('click', function (e) {
		var el = e.target.closest('[data-mcp-tool], [data-mcp-prompt], a[href], button, input[type=submit]');
		if (!el || el.tagName === 'FORM') return;
		var form = el.form;
		if (form && form.hasAttribute('data-mcp-tool') && el.type === 'submit') {
			// The form's own submit is blocked in a sandbox without allow-forms (before any submit event).
			e.preventDefault();
			submit(form, el);
			return;
		}
		if (el.hasAttribute('data-mcp-tool')) {
			e.preventDefault();
			run(el, el, json(el.getAttribute('data-mcp-args')));
		} else if (el.hasAttribute('data-mcp-prompt')) {
			e.preventDefault();
			busy(el, true);
			message(el.getAttribute('data-mcp-prompt')).catch(function (err) { emit(el, 'error', { error: err }); })
				.then(function () { busy(el, false); });
		} else if (/^https?:/.test(el.href) && el.getAttribute('href').charAt(0) !== '#') {
			e.preventDefault();
			openLink(el.href).catch(function (err) { emit(el, 'error', { error: err }); });
		}
	});
	var submit = function (form, submitter) {
		if (!form.checkValidity()) { form.reportValidity(); return; }
		var button = submitter || form.querySelector('button:not([type]), [type=submit]');
		var args = json(form.getAttribute('data-mcp-args'));
		var f = fields(form, submitter);
		for (var k in f) args[k] = f[k];
		run(form, button, args);
	};
	// Enter in a field submits the form, as the browser would.
	document.addEventListener('keydown', function (e) {
		var el = e.target;
		if (e.key !== 'Enter' || e.isComposing || !el.form || !el.form.hasAttribute('data-mcp-tool')) return;
		if (el.tagName === 'TEXTAREA' || el.tagName === 'BUTTON' || el.tagName === 'SELECT') return;
		e.preventDefault();
		submit(el.form, null);
	});
	document.addEventListener('submit', function (e) {
		if (!e.target.hasAttribute('data-mcp-tool')) return;
		e.preventDefault();
		submit(e.target, e.submitter);
	});

	window.mcpApp = {
		ready: ready, context: context,
		call: call, message: message, openLink: openLink,
		request: send, notify: notify,
		log: function (level, data) { notify('notifications/message', { level: level || 'info', data: data }); }
	};
})();
