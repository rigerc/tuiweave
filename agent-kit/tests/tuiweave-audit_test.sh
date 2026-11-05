#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
audit="$root/scripts/tuiweave-audit.sh"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM

fail() {
	printf 'FAIL: %s\n' "$1" >&2
	exit 1
}

expect_pass() {
	name=$1
	shift
	if ! "$audit" "$@" > "$tmp/$name.out" 2>&1; then
		sed -n '1,160p' "$tmp/$name.out" >&2
		fail "$name should pass"
	fi
}

expect_fail() {
	name=$1
	needle=$2
	shift 2
	if "$audit" "$@" > "$tmp/$name.out" 2>&1; then
		fail "$name should fail"
	fi
	grep -q "$needle" "$tmp/$name.out" || {
		sed -n '1,160p' "$tmp/$name.out" >&2
		fail "$name did not report $needle"
	}
}

before=$(find "$root/testdata/clean-app" -type f -exec cksum {} + | sort)
expect_pass clean --repo "$root/testdata/clean-app" --strict
after=$(find "$root/testdata/clean-app" -type f -exec cksum {} + | sort)
[ "$before" = "$after" ] || fail "audit rewrote the clean fixture"

expect_pass clean-verify --repo "$root/testdata/clean-app" --strict --verify
verified=$(find "$root/testdata/clean-app" -type f -exec cksum {} + | sort)
[ "$before" = "$verified" ] || fail "verified audit rewrote the clean fixture"

expect_fail raw-color 'raw color' --repo "$root/testdata/raw-color-violation"
expect_fail ultraviolet 'Ultraviolet' --repo "$root/testdata/ultraviolet-violation"
expect_pass missing-sizing-advisory --repo "$root/testdata/missing-sizing"
expect_fail missing-sizing-strict 'WindowSizeMsg' --repo "$root/testdata/missing-sizing" --strict
expect_fail replace-directive 'replace directive' --repo "$root/testdata/replace-directive" --strict
expect_fail low-go-version 'supported 1.25.8 floor' --repo "$root/testdata/low-go-version" --strict
expect_fail missing-snapcells 'SnapCells' --repo "$root/testdata/missing-snapcells" --strict

# The modern consumer has real snapshot/scenario tests and remains immutable
# even when the audit runs module verification, build, vet, and tests.
modern_before=$(find "$root/testdata/modern-app" -type f -exec cksum {} + | sort)
expect_pass modern --repo "$root/testdata/modern-app" --strict --verify
modern_after=$(find "$root/testdata/modern-app" -type f -exec cksum {} + | sort)
[ "$modern_before" = "$modern_after" ] || fail "verified audit rewrote the modern fixture"

make_mod_fixture() {
	name=$1
	version=$2
	go_floor=$3
	mkdir -p "$tmp/$name"
	printf '# Application instructions\n' > "$tmp/$name/AGENTS.md"
	printf 'module example.com/%s\n\ngo %s\n\nrequire github.com/ishan5ain/tuiweave %s\n' "$name" "$go_floor" "$version" > "$tmp/$name/go.mod"
}

make_mod_fixture modern-low v0.3.0 1.25.8
expect_pass modern-low-advisory --repo "$tmp/modern-low"
expect_fail modern-low-strict 'supported 1.26.0 floor' --repo "$tmp/modern-low" --strict
make_mod_fixture modern-exact v0.3.0 1.26.0
expect_pass modern-exact --repo "$tmp/modern-exact" --strict
make_mod_fixture modern-above v0.3.0 1.27.2
expect_pass modern-above --repo "$tmp/modern-above" --strict
make_mod_fixture legacy-exact v0.2.0 1.25.8
expect_pass legacy-exact --repo "$tmp/legacy-exact" --strict
make_mod_fixture unknown v9.9.9 1.27.2
expect_pass unknown-advisory --repo "$tmp/unknown"
expect_fail unknown-strict 'Go compatibility is unknown' --repo "$tmp/unknown" --strict

revision=$(awk '$1 ~ /^[0-9a-f]+$/ && length($1) == 40 && $2 == "1.26.0" { print substr($1, 1, 12); exit }' "$root/scripts/tuiweave-go-versions.tsv")
[ -n "$revision" ] || fail "missing modern commit metadata"
make_mod_fixture pseudo-exact "v0.2.1-0.20261010000000-$revision" 1.26.0
expect_pass pseudo-exact --repo "$tmp/pseudo-exact" --strict
make_mod_fixture pseudo-low "v0.2.1-0.20261010000000-$revision" 1.25.8
expect_fail pseudo-low 'supported 1.26.0 floor' --repo "$tmp/pseudo-low" --strict
make_mod_fixture pseudo-unknown v0.2.1-0.20261010000000-deadbeef1234 1.27.2
expect_fail pseudo-unknown 'Go compatibility is unknown' --repo "$tmp/pseudo-unknown" --strict

printf 'PASS: tuiweave audit fixtures\n'
