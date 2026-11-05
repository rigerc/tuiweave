# Quality gates

Before handoff:

Resolve the minimum Go version from the application's pinned tuiweave revision
and selected dependency graph. Test with that exact toolchain and
`GOTOOLCHAIN=local`, then with current stable Go. This prevents automatic
toolchain upgrades from disguising a minimum-version failure. tuiweave v0.3.0
requires Go 1.26.0; validated v0.1.0 and v0.2.0 tags retain Go 1.25.8.
Keep terminal Program Status Protocol reporting optional and application-owned;
the dependency upgrade alone does not require using `tea.View.ProgramStatus`.

```sh
go mod tidy -diff
go mod verify
go build ./...
go vet ./...
go test ./...
go test ./path/to/changed/package -update
git diff -- '**/testdata/**'
go test ./...
./path/to/tuiweave-agent-kit/scripts/tuiweave-audit.sh --repo . --strict --verify
```

Only run `-update` after an intentional rendering change. Read every changed
golden and explain the visual delta. Update only affected snapshot-bearing
packages; never use repository-wide regeneration as a way to silence failures.
Pair readable `Snap` and role-aware `SnapCells` goldens when rendering assigns
theme roles.

Record the pinned tuiweave revision, commands run, failures, deferred work, and
candidate library gaps in the handoff.
