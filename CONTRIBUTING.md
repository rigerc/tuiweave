# Contributing

Issues and pull requests are welcome. For a new component or public API, open an
issue first so its ownership, contract, and alternatives can be agreed before a
large implementation. Small fixes may go directly to a pull request.

Follow [AGENTS.md](AGENTS.md) and [DESIGN.md](DESIGN.md). Changes must be formatted,
build cleanly, pass vet and tests, and include behavior coverage. Every rendering
change requires updated plain and role-aware goldens whose diff has been reviewed.
New components also require a runnable example and an AGENTS.md recipe.

Pull requests should be focused and explain API and visual decisions. No CLA or
DCO sign-off is required. Participation is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## Vendored agent kit

[`agent-kit/`](agent-kit/README.md) is a plain snapshot of
[`tuiweave-agent-kit`](https://github.com/ishan5ain/tuiweave-agent-kit) at commit
`bfe0f1dc15f442da86346a575dde062671736a39`. It carries the app playbooks, the
Agent Skills package, templates, and the read-only audit, so a tuiweave checkout
can use them without a second clone.

- Keep the copy byte-identical to upstream. Do not edit files under
  `agent-kit/` here; fix the kit upstream and re-vendor.
- `agent-kit/.github/` is inert, because GitHub only reads workflows, CODEOWNERS,
  and Dependabot configuration from the repository root. The kit's fixture
  suites run from the root [CI workflow](.github/workflows/ci.yml) instead.
- Update by cloning the intended upstream commit, copying it over the directory
  with `.git` excluded, confirming `diff -r` reports nothing else, then updating
  the commit in this section and in `CHANGELOG.md`.
- Preserve the executable bit on `agent-kit/scripts/*.sh` and
  `agent-kit/tests/*.sh`. This checkout has `core.filemode=false`, so a copied
  script can land as `100644`; restore it with
  `git update-index --chmod=+x <path>`. The kit documents running those files
  directly, and `scripts/agent-kit-verify.sh` fails when the bit is missing.
- [`scripts/agent-kit-verify.sh`](scripts/agent-kit-verify.sh) runs the vendored
  suites; run it after any kit, Go-minimum, or dependency change. The throwaway
  fixture modules under `agent-kit/testdata/` pin their dependencies on purpose
  and stay outside Dependabot's scope.
