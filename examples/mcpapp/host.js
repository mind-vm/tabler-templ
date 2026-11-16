(function () {
	'use strict';
	// A minimal MCP Apps host for the demo: it plays what Claude or ChatGPT do around a view.
	var base = location.pathname.replace(/\/$/, '');
	var frame = document.getElementById('view');
	var logList = document.getElementById('log');
	var chat = document.getElementById('chat');
	var theme = 'light', hostStyles = false, initialized = false, args = {};

	var log = function (dir, m) {
		var li = document.createElement('li');
		li.className = 'list-group-item py-2';
		var what = m.method || (m.error ? 'error' : 'result') + ' #' + m.id;
		var body = JSON.stringify(m.params || m.result || m.error || {});
		if (body.length > 160) body = body.slice(0, 160) + '…';
		li.innerHTML = '<span class="badge ' + (dir === 'in' ? 'bg-blue-lt' : 'bg-green-lt') + ' me-1"></span><strong></strong> <span class="text-secondary"></span>';
		li.children[0].textContent = dir === 'in' ? 'view → host' : 'host → view';
		li.children[1].textContent = what;
		li.children[2].textContent = body;
		li.dataset.method = what;
		logList.appendChild(li);
		logList.scrollTop = logList.scrollHeight;
	};
	var post = function (m) { m.jsonrpc = '2.0'; log('out', m); frame.contentWindow.postMessage(m, '*'); };
	var reply = function (id, result) { post({ id: id, result: result }); };
	var fail = function (id, message) { post({ id: id, error: { code: -32000, message: message } }); };
	var notify = function (method, params) { if (initialized) post({ method: method, params: params }); };

	var callTool = function (name, a) {
		return fetch(base + '/call', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name: name, arguments: a || {} }) })
			.then(function (r) { if (!r.ok) throw new Error(r.statusText); return r.json(); });
	};

	var variables = function () {
		return hostStyles ? {
			'--color-background-primary': 'light-dark(#faf9f5, #262624)',
			'--color-background-secondary': 'light-dark(#f0eee6, #1f1e1d)',
			'--color-text-primary': 'light-dark(#141413, #faf9f5)',
			'--color-text-secondary': 'light-dark(#5e5d59, #c2c0b6)',
			'--color-border-primary': 'light-dark(#e3dfd2, #3e3e38)',
			'--font-sans': 'Georgia, serif',
			'--border-radius-md': '10px'
		} : {};
	};
	var hostContext = function () {
		return { theme: theme, styles: { variables: variables() }, displayMode: 'inline', availableDisplayModes: ['inline'],
			containerDimensions: { maxHeight: 4000 }, locale: 'en-US', platform: 'web' };
	};

	// The host builds the sandbox's CSP from the resource's _meta.ui.csp (default: nothing outside).
	var policy = function (meta) {
		var c = (meta && meta.ui && meta.ui.csp) || {};
		var res = (c.resourceDomains || []).join(' ');
		var conn = (c.connectDomains || []).join(' ') || "'none'";
		return "default-src 'none'; script-src 'self' 'unsafe-inline' " + res + "; style-src 'self' 'unsafe-inline' " + res +
			"; img-src 'self' data: " + res + "; font-src 'self' " + res + "; media-src 'self' data: " + res +
			'; connect-src ' + conn + "; frame-src 'none'; object-src 'none'; base-uri 'self'";
	};

	var open = function (toolArgs) {
		args = toolArgs || {};
		initialized = false;
		fetch(base + '/tools').then(function (r) { return r.json(); }).then(function (tools) {
			var tool = tools.find(function (t) { return t.name === 'list_orders'; });
			var uri = tool._meta.ui.resourceUri;
			return fetch(base + '/resource?uri=' + encodeURIComponent(uri)).then(function (r) { return r.json(); });
		}).then(function (res) {
			var csp = '<meta http-equiv="Content-Security-Policy" content="' + policy(res._meta) + '">';
			frame.srcdoc = res.text.replace(/<head>/i, '<head>' + csp);
			var border = res._meta && res._meta.ui && res._meta.ui.prefersBorder;
			document.getElementById('frame-box').classList.toggle('border', border !== false);
		});
	};

	window.addEventListener('message', function (e) {
		if (e.source !== frame.contentWindow) return;
		var m = e.data;
		if (!m || m.jsonrpc !== '2.0') return;
		log('in', m);
		var p = m.params || {};
		switch (m.method) {
			case 'ui/initialize':
				reply(m.id, { protocolVersion: '2026-01-26', hostInfo: { name: 'tabler-templ demo host', version: '1.0.0' },
					hostCapabilities: { openLinks: {}, serverTools: {}, logging: {} }, hostContext: hostContext() });
				return;
			case 'ui/notifications/initialized':
				initialized = true;
				notify('ui/notifications/tool-input', { arguments: args });
				callTool('list_orders', args).then(function (r) { notify('ui/notifications/tool-result', r); });
				return;
			case 'ui/notifications/size-changed':
				frame.style.height = Math.min(p.height, 4000) + 'px';
				return;
			case 'tools/call':
				callTool(p.name, p.arguments).then(function (r) { reply(m.id, r); }, function (err) { fail(m.id, err.message); });
				return;
			case 'ui/message':
				var text = (p.content || []).map(function (c) { return c.text || ''; }).join('\n');
				var b = document.createElement('div');
				b.className = 'chat-bubble chat-bubble-me ms-auto mt-2';
				b.style.maxWidth = '24rem';
				b.textContent = text;
				chat.appendChild(b);
				reply(m.id, {});
				return;
			case 'ui/open-link':
				if (!/^https?:\/\//.test(p.url)) { fail(m.id, 'Invalid URL'); return; }
				window.open(p.url, '_blank', 'noopener');
				reply(m.id, {});
				return;
			case 'notifications/message':
				return;
		}
		if (m.id != null && m.method) post({ id: m.id, error: { code: -32601, message: 'Method not found' } });
	});

	document.querySelectorAll('[data-theme]').forEach(function (b) {
		b.addEventListener('click', function () {
			theme = b.dataset.theme;
			document.querySelectorAll('[data-theme]').forEach(function (x) { x.classList.toggle('active', x === b); });
			document.documentElement.setAttribute('data-bs-theme', theme);
			notify('ui/notifications/host-context-changed', { theme: theme, styles: { variables: variables() } });
		});
	});
	document.getElementById('host-styles').addEventListener('click', function (e) {
		hostStyles = !hostStyles;
		e.currentTarget.classList.toggle('active', hostStyles);
		e.currentTarget.setAttribute('aria-pressed', String(hostStyles));
		// Views take them with mcpapp.Props.HostStyles, as the demo's does.
		notify('ui/notifications/host-context-changed', { styles: { variables: variables() } });
	});

	open({});
})();
