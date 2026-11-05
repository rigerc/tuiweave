# tuiweave-agent-kit

Portable workflows, templates, an Agent Skills package, and read-only audit
tooling for applications built with
[`github.com/ishan5ain/tuiweave`](https://github.com/ishan5ain/tuiweave).

tuiweave remains the source of truth for component APIs, design rules, recipes,
and runnable examples. This kit owns operational workflows for greenfield apps,
existing-TUI migrations, reviews, host installation, and repeatable audits.

## Compatibility

The current kit is validated against tuiweave `v0.3.0` at commit
`c05ccf548162f861736844f729c18104949c90e5`, using Go 1.26.0 and Bubble Tea v2.1.0.
It retains Go 1.25.8 audit requirements for validated tuiweave v0.1.0 and v0.2.0
tags. Pin applications to an
explicit tuiweave tag or commit; do not use moving `@latest` instructions. Run
`scripts/sync-tuiweave-reference.sh /path/to/tuiweave` after changing the validated
revision.

The audit reads offline requirements from `scripts/tuiweave-go-versions.tsv`.
Only recorded tags or commit-derived pseudo-versions receive a known Go floor;
unknown pins produce an advisory and fail `--strict` until their contracts are
reviewed and synced. The generator records clean commits and exact release tags
from the source checkout's `go.mod`, preserving validated historical entries.
Run it from a clean checkout of the revision you intend to validate.

## Start here

- New app: [`playbooks/greenfield.md`](playbooks/greenfield.md)
- Existing TUI: [`playbooks/migration.md`](playbooks/migration.md)
- Review or audit: [`playbooks/review.md`](playbooks/review.md)
- Reusable skill: [`skills/tuiweave-app-builder/SKILL.md`](skills/tuiweave-app-builder/SKILL.md)
- Task shape: [`playbooks/task-contract.md`](playbooks/task-contract.md)
- Contributing: [`CONTRIBUTING.md`](CONTRIBUTING.md)
- Support: [`SUPPORT.md`](SUPPORT.md)
- Security: [`SECURITY.md`](SECURITY.md)

Run the audit without modifying the target repository:

```sh
./scripts/tuiweave-audit.sh --repo /path/to/app
./scripts/tuiweave-audit.sh --repo /path/to/app --strict
./scripts/tuiweave-audit.sh --repo /path/to/app --verify
```

`--strict` promotes advisory findings to a failing exit status. `--verify`
also runs `go build ./...`, `go vet ./...`, and `go test ./...`. The audit
never rewrites source or regenerates goldens.
The Go-floor check covers the validated tuiweave pin. Applications can select
newer runtime or indirect dependencies with higher requirements; `--verify`
checks the selected graph. Use the exact intended toolchain with
`GOTOOLCHAIN=local` so automatic toolchain upgrades do not hide an incompatibility.

## Development

```sh
./tests/tuiweave-audit_test.sh
./tests/sync-tuiweave-reference_test.sh
```

See [`install/`](install/) for Codex, Pi, and other Agent Skills-compatible
hosts. The same skill package is used for every host.

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) and [GOVERNANCE.md](GOVERNANCE.md)
for community expectations and project decision-making.
