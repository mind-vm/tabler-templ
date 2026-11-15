#!/bin/sh
# Copies the concepts the reference app shows (Gestaltung, Assistent) from mind-vm/auftragsbutler into
# examples/reference/konzept/, with the commit they come from in QUELLE. The concept stays the source of truth
# (GS-03, GSO-01: it lives in the Auftragsbutler until v1.0); run this after it changes and commit the copy.
# Needs the GitHub CLI with access to the repository.
set -eu
cd "$(dirname "$0")/../examples/reference/konzept"
repo=mind-vm/auftragsbutler
sha=$(gh api "repos/$repo/commits?path=docs/konzepte&per_page=1" --jq '.[0].sha')
date=$(gh api "repos/$repo/commits/$sha" --jq '.commit.committer.date[:10]')
for f in gestaltung assistent; do
	gh api "repos/$repo/contents/docs/konzepte/$f.md?ref=$sha" -H "Accept: application/vnd.github.raw" >"$f.md"
done
printf '%s %s\n' "$sha" "$date" >QUELLE
echo "konzept: $repo@$(echo "$sha" | cut -c1-8) ($date)"
