#!/bin/sh
set -eu

if [ "$#" -ne 1 ]; then
	printf 'usage: %s /path/to/tuiweave\n' "$0" >&2
	exit 2
fi

tuiweave=$(cd "$1" && pwd)
[ -f "$tuiweave/AGENT-CATALOG.md" ] || { printf 'not a tuiweave checkout: %s\n' "$tuiweave" >&2; exit 1; }
[ -f "$tuiweave/COMPATIBILITY.md" ] || { printf 'missing compatibility contract: %s\n' "$tuiweave/COMPATIBILITY.md" >&2; exit 1; }
[ -f "$tuiweave/go.mod" ] || { printf 'missing module contract: %s\n' "$tuiweave/go.mod" >&2; exit 1; }

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
output="$script_dir/../skills/tuiweave-app-builder/references/tuiweave-api-compatibility.md"
revision=$(git -C "$tuiweave" rev-parse HEAD 2>/dev/null || printf unknown)
minimum_go=$(awk '$1 == "go" { print $2; exit }' "$tuiweave/go.mod")
if ! printf '%s\n' "$minimum_go" | grep -Eq '^[0-9]+\.[0-9]+(\.[0-9]+)?$'; then
	printf 'invalid Go minimum in %s\n' "$tuiweave/go.mod" >&2
	exit 1
fi
# Metadata is attributed to a clean commit, never to edited working-tree files.
if [ "$revision" != unknown ] && [ -n "$(git -C "$tuiweave" status --porcelain --untracked-files=no)" ]; then
	printf 'tuiweave checkout has tracked changes; commit them before syncing\n' >&2
	exit 1
fi
tmp=$(mktemp)
metadata_tmp=$(mktemp)
trap 'rm -f "$tmp" "$metadata_tmp"' EXIT HUP INT TERM

{
	printf '# tuiweave API compatibility\n\n'
	printf 'Generated from tuiweave revision `%s`. Verify against the application pinned revision.\n\n' "$revision"
	printf '## Supported versions\n\n'
	sed -n '/^## Supported versions/,/^## Go version and module boundary/p' "$tuiweave/COMPATIBILITY.md" | sed '1d;$d'
	printf '\n## Go version and module boundary\n\n'
	sed -n '/^## Go version and module boundary/,/^## Layout compatibility/p' "$tuiweave/COMPATIBILITY.md" | sed '1d;$d'
	printf '## Relevant package contracts\n\n'
	printf '| Package | Contract |\n|---|---|\n'
	grep -E '^\| (\[)?`(tuiweave|layout|scrollbar|focus|inspect|tabs|snaptest|dialog|autocomplete|textarea)`' "$tuiweave/AGENT-CATALOG.md" |
		sed -E \
			-e 's/^\| (\[[^]]+\]\([^)]*\)|`[^`]+`) \| [^|]+ \| ([^|]+) \|$/| \1 | \2 |/' \
			-e 's/\[`([^`]+)`\]\([^)]*\)/`\1`/'
	printf '\n## Compatibility notes\n\n'
	sed -n '/^## API compatibility notes/,/^## Validation loop/p' "$tuiweave/AGENT-CATALOG.md" | sed '1d;$d'
	printf '\n## Layout constraint compatibility\n\n'
	sed -n '/^### Constraints/,/^## Release and migration policy/p' "$tuiweave/COMPATIBILITY.md" | sed '1d;$d'
	printf '\n## Source of truth\n\n'
	printf 'Read `AGENT-CONTRACT.md`, `COMPATIBILITY.md`, `AGENT-CATALOG.md`, the matching `AGENTS.md` recipe, and the closest runnable example from the pinned tuiweave checkout before editing.\n'
} > "$tmp"

metadata="$script_dir/tuiweave-go-versions.tsv"
if [ "$revision" != unknown ]; then
	# Preserve validated historical releases and add this revision and its
	# exact release tags. LC_ALL makes ordering stable on every supported host.
	{
		awk '!/^#/ && NF == 2' "$metadata"
		printf '%s\t%s\n' "$revision" "$minimum_go"
		git -C "$tuiweave" tag --points-at HEAD | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' |
			while IFS= read -r tag; do printf '%s\t%s\n' "$tag" "$minimum_go"; done
	} | LC_ALL=C sort -u > "$metadata_tmp"
	# A pin must not acquire conflicting floors during regeneration.
	if ! awk 'seen[$1] && seen[$1] != $2 { exit 1 } { seen[$1] = $2 }' "$metadata_tmp"; then
		printf 'conflicting validated Go requirements\n' >&2
		exit 1
	fi
	mv "$metadata_tmp" "$metadata"
fi
mv "$tmp" "$output"
rm -f "$metadata_tmp"
trap - EXIT HUP INT TERM
printf 'updated %s from tuiweave %s\n' "$output" "$revision"
