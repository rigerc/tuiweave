#!/bin/sh
set -eu

usage() {
	printf '%s\n' "usage: $0 --repo /path/to/app [--strict] [--verify]"
}

repo=""
strict=0
verify=0
while [ "$#" -gt 0 ]; do
	case "$1" in
	--repo)
		[ "$#" -ge 2 ] || { usage >&2; exit 2; }
		repo=$2
		shift 2
		;;
	--strict)
		strict=1
		shift
		;;
	--verify)
		verify=1
		shift
		;;
	-h|--help)
		usage
		exit 0
		;;
	*)
		usage >&2
		exit 2
		;;
	esac
done

[ -n "$repo" ] || { usage >&2; exit 2; }
[ -d "$repo" ] || { printf 'BLOCK: repository does not exist: %s\n' "$repo" >&2; exit 1; }
repo=$(cd "$repo" && pwd)

blocking=0
advisory=0
script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
tuiweave_version=""

block() {
	blocking=$((blocking + 1))
	printf 'BLOCK: %s\n' "$1"
}

advise() {
	advisory=$((advisory + 1))
	printf 'ADVISORY: %s\n' "$1"
}

if [ ! -f "$repo/go.mod" ] || ! grep -q 'github.com/ishan5ain/tuiweave' "$repo/go.mod"; then
	block "missing github.com/ishan5ain/tuiweave dependency in go.mod"
else
	tuiweave_version=$(awk '
		$1 == "require" && $2 == "github.com/ishan5ain/tuiweave" { print $3; exit }
		$1 == "github.com/ishan5ain/tuiweave" { print $2; exit }
	' "$repo/go.mod")
	if [ -n "$tuiweave_version" ]; then
		printf 'INFO: tuiweave dependency %s\n' "$tuiweave_version"
	else
		advise "could not determine pinned tuiweave version from go.mod"
	fi

	replace_matches=$(grep -nE '^[[:space:]]*replace[[:space:]]+github[.]com/ishan5ain/tuiweave([[:space:]]|$)|^[[:space:]]*github[.]com/ishan5ain/tuiweave[[:space:]]+=>[[:space:]]+' "$repo/go.mod" || true)
	if [ -n "$replace_matches" ]; then
		advise "tuiweave replace directive is development-only and should not be committed"
		printf '%s\n' "$replace_matches"
	fi
fi

go_version=$(awk '$1 == "go" { print $2; exit }' "$repo/go.mod" 2>/dev/null || true)
# Resolve only validated pins. A pseudo-version identifies a commit by its
# trailing hash; do not infer a future release's floor from its version number.
minimum_go=$(awk -v pin="$tuiweave_version" '
	$1 == pin && pin != "" { print $2; exit }
	pin ~ /^v[0-9]+\.[0-9]+\.[0-9]+-.*-[0-9a-f]+$/ && $1 ~ /^[0-9a-f]+$/ && length($1) == 40 {
		n = split(pin, parts, "-"); hash = parts[n]
		if (length(hash) >= 12 && length(hash) <= 40 && substr($1, 1, length(hash)) == hash) {
			print $2; exit
		}
	}
' "$script_dir/tuiweave-go-versions.tsv")
if [ -n "$tuiweave_version" ] && [ -z "$minimum_go" ]; then
	advise "Go compatibility is unknown for tuiweave $tuiweave_version; inspect its pinned go.mod and COMPATIBILITY.md"
elif [ -n "$minimum_go" ]; then
	printf 'INFO: validated tuiweave Go minimum %s\n' "$minimum_go"
	if [ -z "$go_version" ]; then
		advise "missing go directive; tuiweave $tuiweave_version requires Go $minimum_go"
	elif ! awk -v actual="$go_version" -v minimum="$minimum_go" 'BEGIN {
		split(actual, a, "."); split(minimum, m, ".")
		for (i = 1; i <= 3; i++) {
			a[i] += 0; m[i] += 0
			if (a[i] > m[i]) exit 0
			if (a[i] < m[i]) exit 1
		}
		exit 0
	}'; then
		advise "go directive $go_version is below tuiweave $tuiweave_version supported $minimum_go floor"
	fi
fi

raw_matches=$(find "$repo" -type f -name '*.go' ! -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -HnE '#[[:xdigit:]]{3,8}([^[:xdigit:]]|$)|lipgloss[.]Color[[:space:]]*[(]' {} + 2>/dev/null | grep -Ev '/(theme|themes|color|colors)(/|[^/]*[.]go:)|:[0-9]+:[[:space:]]*//' || true)
if [ -n "$raw_matches" ]; then
	block "raw color literal or lipgloss.Color outside theme definitions"
	printf '%s\n' "$raw_matches"
fi

uv_matches=$(find "$repo" -type f -name '*.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -Hn 'charmbracelet/ultraviolet' {} + 2>/dev/null || true)
if [ -n "$uv_matches" ]; then
	block "forbidden Ultraviolet import in application code"
	printf '%s\n' "$uv_matches"
fi

if [ ! -f "$repo/AGENTS.md" ]; then
	advise "missing app-level AGENTS.md"
fi

ui_sources=$(find "$repo" -type f -name '*.go' ! -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -l 'github.com/ishan5ain/tuiweave' {} + 2>/dev/null || true)
if [ -n "$ui_sources" ]; then
	snapshot_tests=$(find "$repo" -type f -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -l 'snaptest[.]Snap' {} + 2>/dev/null || true)
	role_snapshot_tests=$(find "$repo" -type f -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -l 'snaptest[.]SnapCells' {} + 2>/dev/null || true)
	scenario_tests=$(find "$repo" -type f -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -l 'snaptest[.]RunScenario' {} + 2>/dev/null || true)
	[ -n "$snapshot_tests" ] || advise "UI packages have no snaptest.Snap coverage"
	[ -n "$role_snapshot_tests" ] || advise "UI packages have no snaptest.SnapCells role-aware coverage"
	[ -n "$scenario_tests" ] || advise "UI packages have no snaptest.RunScenario coverage"
fi

discarded_updates=$(find "$repo" -type f -name '*.go' ! -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -HnE '^[[:space:]]*[[:alnum:]_.]+[.]Update[(]|^[[:space:]]*_[[:space:]]*=[^=]*[.]Update[(]|^[[:space:]]*_,[[:space:]]*_[[:space:]]*=[^=]*[.]Update[(]' {} + 2>/dev/null || true)
if [ -n "$discarded_updates" ]; then
	advise "suspicious component Update call may discard the returned model or command"
	printf '%s\n' "$discarded_updates"
fi

window_files=$(find "$repo" -type f -name '*.go' ! -name '*_test.go' ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -l 'WindowSizeMsg' {} + 2>/dev/null || true)
old_ifs=$IFS
IFS='
'
for file in $window_files; do
	if ! grep -Eq 'SetSize[[:space:]]*[(]|[.]Apply[[:space:]]*[(]' "$file"; then
		advise "WindowSizeMsg without SetSize or layout.Apply in $file"
	fi
done
IFS=$old_ifs

latest_matches=$(find "$repo" -type f \( -name '*.md' -o -name '*.sh' -o -name 'go.mod' \) ! -path '*/vendor/*' ! -path '*/.git/*' -exec grep -Hn '@latest' {} + 2>/dev/null | grep -Eiv 'do not|never|avoid|unpinned|moving' || true)
if [ -n "$latest_matches" ]; then
	advise "unpinned @latest dependency instruction"
	printf '%s\n' "$latest_matches"
fi

run_verify() {
	label=$1
	shift
	printf 'VERIFY: %s\n' "$label"
	if ! (cd "$repo" && "$@"); then
		block "$label failed"
	fi
}

if [ "$verify" -eq 1 ]; then
	run_verify "go mod verify" go mod verify
	verify_build_dir=$(mktemp -d)
	run_verify "go build ./..." go build -mod=readonly -o "$verify_build_dir/" ./...
	rm -rf "$verify_build_dir"
	run_verify "go vet ./..." go vet -mod=readonly ./...
	run_verify "go test ./..." go test -mod=readonly ./...
fi

printf 'SUMMARY: %d blocking, %d advisory\n' "$blocking" "$advisory"
if [ "$blocking" -gt 0 ] || { [ "$strict" -eq 1 ] && [ "$advisory" -gt 0 ]; }; then
	exit 1
fi
