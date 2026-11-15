#!/bin/sh
# Starts the demo server and takes the screenshots of the reference app (examples/reference) at the desktop and on
# a phone into docs/screenshots/reference/. Fails when a page scrolls sideways (GS-14), logs an error in the browser
# or loads anything from outside (GS-09).
# Needs Go, Node and a Chromium for Playwright (set CHROMIUM=/path/to/chrome to use an installed one).
set -eu
cd "$(dirname "$0")/.."
port=${PORT:-8398}
go build -o /tmp/tabler-templ-demo-reference ./cmd/demo
/tmp/tabler-templ-demo-reference -addr "127.0.0.1:$port" &
pid=$!
trap 'kill $pid' EXIT
until curl -fs "http://127.0.0.1:$port/reference" >/dev/null; do sleep 0.2; done
cd scripts/screenshots
[ -d node_modules ] || npm install --no-audit --no-fund
node reference.mjs "http://127.0.0.1:$port" ../../docs/screenshots/reference
