# Changelog

## Unreleased

- Validate tuiweave v0.3.0 with Go 1.26.0 and Bubble Tea v2.1.0, and regenerate
  the supported-version matrix, Go migration guidance, and API reference from
  its released commit.
- Resolve audit Go requirements from offline validated tag/commit metadata;
  retain the Go 1.25.8 floor for older supported tags and report unknown pins
  instead of guessing their requirements.
- Add a modern consumer with real view, role-aware, and scenario goldens;
  preserve legacy coverage and fixture immutability in verified audits.
- Exercise exact Go 1.26.0 and stable Go in CI, preserving the required `audit`
  check name. Keep Program Status Protocol reporting optional and app-owned.

## v0.2.0 — 2026-07-21

- Align the skill, playbooks, templates, and generated compatibility reference
  with tuiweave v0.2.0, including owned layout constraints, size modes,
  inspection/action routing, tab hit-testing, and role-aware golden guidance.
- Add compatibility guidance for tuiweave's built-in theme preset catalog.
- Extend the read-only audit with dependency-version reporting, Go compatibility
  guidance, development `replace` detection, `SnapCells` coverage guidance, and
  `go mod verify` support.
- Add a tested, idempotent compatibility-reference generator and update audit
  fixtures to tuiweave v0.2.0.
- Add contribution, conduct, security, support, governance, and issue/PR
  guidance for the public community repository.

## v0.1.0 — 2026-07-11

- Add greenfield, migration, review, task-contract, and quality-gate playbooks.
- Add a portable `tuiweave-app-builder` Agent Skill for Codex, Pi, and compatible hosts.
- Add app, migration, component-map, plan, and review templates.
- Add a read-only tuiweave audit with fixture-based shell tests.
- Record compatibility with tuiweave core commit `a647f53882652763343dd57e31175e35f182427e`.
