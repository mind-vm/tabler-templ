#!/bin/sh
# Starts the demo server, screenshots every rebuilt page into docs/screenshots/ and stops the server.
# Needs Go, Node and a Chromium for Playwright (set CHROMIUM=/path/to/chrome to use an installed one).
set -eu
cd "$(dirname "$0")/.."
port=${PORT:-8399}
go build -o /tmp/tabler-templ-demo ./cmd/demo
/tmp/tabler-templ-demo -addr "127.0.0.1:$port" &
pid=$!
trap 'kill $pid' EXIT
until curl -fs "http://127.0.0.1:$port/pages.json" >/dev/null; do sleep 0.2; done
cd scripts/screenshots
[ -d node_modules ] || npm install --no-audit --no-fund
node shoot.mjs "http://127.0.0.1:$port" ../../docs/screenshots
