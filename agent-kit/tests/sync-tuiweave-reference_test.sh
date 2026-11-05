#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM

mkdir -p "$tmp/kit/scripts" "$tmp/kit/skills/tuiweave-app-builder/references" "$tmp/core"
cp "$root/scripts/sync-tuiweave-reference.sh" "$tmp/kit/scripts/"
cp "$root/scripts/tuiweave-go-versions.tsv" "$tmp/kit/scripts/"

cat > "$tmp/core/AGENT-CATALOG.md" <<'EOF'
| Package | Use when | First API / contract |
|---|---|---|
| `tuiweave` | Choosing themes | `Presets()` |
| [`snaptest`](snaptest/README.md) | Testing | `Snap`, `SnapCells`, `RunScenario` |
| `inspect` | Inspection | `BindAt` and qualified actions |
| `tabs` | Navigation | `IndexAt` accepts local coordinates |

## API compatibility notes

- Compatibility sentinel.

## Validation loop
EOF

cat > "$tmp/core/COMPATIBILITY.md" <<'EOF'
## Supported versions

| Go | 1.26.0 minimum | exact-floor CI |
| Bubble Tea | v2.1.0 | full tests |

## Go version and module boundary

Consumers must install Go 1.26.0 before adopting v0.3.0.

## Layout compatibility

### Constraints

Use tuiweave-owned constraint constructors.

## Release and migration policy
EOF
cat > "$tmp/core/go.mod" <<'EOF'
module github.com/ishan5ain/tuiweave

go 1.26.0
EOF
# Use a real clean commit/tag so generation exercises provenance metadata.
git -C "$tmp/core" init -q
git -C "$tmp/core" add .
git -C "$tmp/core" -c user.name=fixture -c user.email=fixture@example.com commit -qm fixture
git -C "$tmp/core" tag v9.8.7
revision=$(git -C "$tmp/core" rev-parse HEAD)

"$tmp/kit/scripts/sync-tuiweave-reference.sh" "$tmp/core" > "$tmp/run.out"
output="$tmp/kit/skills/tuiweave-app-builder/references/tuiweave-api-compatibility.md"

grep -q '| `snaptest` | `Snap`, `SnapCells`, `RunScenario` |' "$output"
grep -q '| `inspect` | `BindAt` and qualified actions |' "$output"
grep -q '| `tabs` | `IndexAt` accepts local coordinates |' "$output"
grep -q 'Compatibility sentinel' "$output"
grep -q 'Use tuiweave-owned constraint constructors' "$output"
grep -q 'Go | 1.26.0 minimum' "$output"
grep -q 'Consumers must install Go 1.26.0' "$output"
awk '$1 == "v9.8.7" && $2 == "1.26.0" { found=1 } END { exit !found }' "$tmp/kit/scripts/tuiweave-go-versions.tsv"
awk -v revision="$revision" '$1 == revision && $2 == "1.26.0" { found=1 } END { exit !found }' "$tmp/kit/scripts/tuiweave-go-versions.tsv"

before=$(cksum "$output" "$tmp/kit/scripts/tuiweave-go-versions.tsv")
"$tmp/kit/scripts/sync-tuiweave-reference.sh" "$tmp/core" > "$tmp/rerun.out"
after=$(cksum "$output" "$tmp/kit/scripts/tuiweave-go-versions.tsv")
[ "$before" = "$after" ]

# Editing the module contract without committing it must not attribute the new
# requirement to an old commit, or overwrite either generated artifact.
printf '\n// dirty fixture\n' >> "$tmp/core/go.mod"
if "$tmp/kit/scripts/sync-tuiweave-reference.sh" "$tmp/core" > "$tmp/dirty.out" 2>&1; then
	printf 'FAIL: dirty core checkout accepted\n' >&2
	exit 1
fi
[ "$after" = "$(cksum "$output" "$tmp/kit/scripts/tuiweave-go-versions.tsv")" ]

# A moved tag with a changed minimum conflicts with its recorded contract.
# Reject it without losing any previously generated data.
sed 's/go 1.26.0/go 1.27.0/' "$tmp/core/go.mod" > "$tmp/new.mod"
mv "$tmp/new.mod" "$tmp/core/go.mod"
git -C "$tmp/core" add go.mod
git -C "$tmp/core" -c user.name=fixture -c user.email=fixture@example.com commit -qm changed-floor
git -C "$tmp/core" tag -f v9.8.7 > /dev/null
if "$tmp/kit/scripts/sync-tuiweave-reference.sh" "$tmp/core" > "$tmp/conflict.out" 2>&1; then
	printf 'FAIL: conflicting tag minimum accepted\n' >&2
	exit 1
fi
[ "$after" = "$(cksum "$output" "$tmp/kit/scripts/tuiweave-go-versions.tsv")" ]

printf 'PASS: tuiweave reference sync\n'
