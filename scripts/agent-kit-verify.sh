#!/bin/sh
# Verify the vendored agent kit under agent-kit/ against its own fixture suites.
#
# agent-kit/ is a snapshot of https://github.com/ishan5ain/tuiweave-agent-kit, so
# these suites catch a broken or partial copy as well as a Go toolchain or audit
# regression. They never modify the target repositories; see the kit's README.
# Updating the snapshot is documented in CONTRIBUTING.md.
set -eu

root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
kit="$root/agent-kit"

if [ ! -f "$kit/scripts/tuiweave-audit.sh" ]; then
	printf 'missing vendored agent kit at %s\n' "$kit" >&2
	exit 1
fi

status=0

# The kit documents direct invocation (./scripts/…) where it is installed, so a
# copy that loses the executable bit breaks its instructions even though the
# suites below run through sh.
for path in "$kit"/scripts/*.sh "$kit"/tests/*.sh; do
	if [ ! -x "$path" ]; then
		printf 'not executable: %s\n  fix with: git update-index --chmod=+x %s\n' \
			"${path#"$root"/}" "${path#"$root"/}" >&2
		status=1
	fi
done

for suite in tuiweave-audit_test.sh sync-tuiweave-reference_test.sh; do
	if ! sh "$kit/tests/$suite"; then
		printf 'agent-kit fixture suite failed: %s\n' "$suite" >&2
		status=1
	fi
done

if [ "$status" -eq 0 ]; then
	printf 'agent-kit fixtures pass\n'
fi
exit "$status"
