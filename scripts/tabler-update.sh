#!/bin/sh
# Prepares the update to a new Tabler release; see "Updating Tabler" in README.md. It changes nothing in
# Tabler's components here: it moves a Tabler checkout to the release, lists what changed there, updates the
# pinned versions and the generated demo data, and compares every rebuilt page with the new preview.
#
#   scripts/tabler-update.sh <tabler checkout> <release tag, e.g. v1.8.0>
#
# Needs git, Go, Node with pnpm (to build Tabler's preview) and a Chromium for Playwright (CHROMIUM=…).
# The reports land in $OUT (default /tmp/tabler-update).
set -eu

[ $# -eq 2 ] || { echo "usage: $0 <tabler checkout> <release tag>" >&2; exit 2; }
tabler=$(cd "$1" && pwd)
tag=$2
root=$(cd "$(dirname "$0")/.." && pwd)
out=${OUT:-/tmp/tabler-update}
port=${PORT:-18391}
mkdir -p "$out"
cd "$root"

old=$(go run ./internal/tools/versions -ref)
echo "== Tabler $old -> $tag"
git -C "$tabler" fetch -q --tags origin || echo "   (fetch failed; using the checkout's refs)"
git -C "$tabler" checkout -q "$tag"

echo "== What changed in Tabler ($out/changes.txt, $out/upgrade-notes.diff)"
git -C "$tabler" diff --stat=200 "$old" "$tag" -- shared/ui shared/components shared/layouts preview/pages \
	shared/data core/scss core/js core/package.json >"$out/changes.txt"
git -C "$tabler" diff "$old" "$tag" -- docs/content/ui/getting-started/upgrade 'docs/content/ui/**/*.mdx' \
	>"$out/upgrade-notes.diff"
tail -n 1 "$out/changes.txt"

echo "== Pinned versions"
go run ./internal/tools/versions -write "$tabler"

echo "== Generated demo data"
go run ./internal/tools/charts "$tabler"
go run ./internal/tools/menu "$tabler"
icons=$(go run ./internal/tools/versions -icons)
go run ./internal/tools/filledicons "$icons"
(cd cmd/demo && go run batch_g_gen.go)
go tool templ generate >/dev/null
go build ./... && go test ./... >"$out/test.txt" 2>&1 || echo "   tests fail, see $out/test.txt"

echo "== Building Tabler's preview"
(cd "$tabler" && pnpm install --frozen-lockfile >/dev/null && pnpm exec turbo build --filter=@tabler/preview >"$out/preview-build.txt" 2>&1) || { echo "   preview build failed, see $out/preview-build.txt"; exit 1; }

echo "== Comparing every page ($out/report.tsv)"
for p in $port $((port + 1)); do
	if curl -s -o /dev/null "http://127.0.0.1:$p/"; then
		echo "   port $p is in use; set PORT to a free pair of ports" >&2
		exit 1
	fi
done
go build -o "$out/demo" ./cmd/demo
"$out/demo" -addr "127.0.0.1:$port" &
demo=$!
python3 -m http.server $((port + 1)) --bind 127.0.0.1 --directory "$tabler/preview/dist" >/dev/null 2>&1 &
orig=$!
trap 'kill $demo $orig 2>/dev/null || true' EXIT
until curl -fs "http://127.0.0.1:$port/pages.json" >/dev/null; do
	kill -0 $demo 2>/dev/null || { echo "   the demo did not start" >&2; exit 1; }
	sleep 0.2
done
(cd scripts/screenshots && node compare-all.mjs "http://127.0.0.1:$port" "http://127.0.0.1:$((port + 1))" "$out" 2000)

cat <<EOF

Next: read $out/upgrade-notes.diff (renamed classes) and $out/changes.txt, fix the pages at the top of
$out/report.tsv (crop.mjs shows where), add new components and pages to internal/catalog, and note
markup changes as "Breaking (look)" in CHANGELOG.md.
EOF
