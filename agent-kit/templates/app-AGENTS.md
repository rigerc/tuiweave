# Application agent instructions

This application uses `github.com/ishan5ain/tuiweave` at: `<PINNED_VERSION_OR_COMMIT>`.
Use the Go minimum declared by that revision and its selected dependency graph;
test the exact minimum with `GOTOOLCHAIN=local`.

API and design truth:

- `AGENT-CONTRACT.md`, `COMPATIBILITY.md`, `AGENT-CATALOG.md`, `AGENTS.md`, and
  the matching runnable examples from the pinned tuiweave checkout.

Route work to the closest runnable tuiweave example before implementing. Keep
layout, focus, routing, product state, persistence, and backend lifecycle in
this application. Reassign every component model returned by `Update`, collect
every command, and verify rendering with snapshots plus interaction scenarios.
The app owns inspection-tree assembly, absolute bounds, privacy, and qualified
action routing; validate the current visible tree before invoking an action.

A local `replace github.com/ishan5ain/tuiweave => ../tuiweave` directive may be used
for development only. Keep the committed dependency pinned. Changes to tuiweave
or a local fork require explicit authorization; report reusable API gaps with a
minimal example and acceptance test.
