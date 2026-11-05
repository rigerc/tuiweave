# Greenfield playbook

1. Inspect the repository, `go.mod`, Go version, current test conventions, and
   the pinned tuiweave tag or commit. Read tuiweave's `AGENT-CONTRACT.md`,
   `COMPATIBILITY.md`, `AGENT-CATALOG.md`, matching `AGENTS.md` recipes, and the
   closest runnable example.
2. Write a component map, layout rectangles, size modes, focus order,
   message/result flow, inspection/action ownership, and application-owned state
   map. Keep backend clients, persistence, and product orchestration outside
   tuiweave components.
3. Create tasks using the standard task contract. Start with the smallest
   working shell: root model, `WindowSizeMsg` layout, sizing, focus, delegation,
   and one composed `tea.View`.
4. Add behavior in vertical slices. Always reassign component models returned
   by `Update` and collect every `tea.Cmd`.
5. Add readable snapshots for default, focused, empty/full, narrow, wide-glyph,
   and combining-grapheme states. Pair `Snap` and `SnapCells` when rendering
   assigns theme roles, and add named scenarios with explicit command delivery.
6. Run build, vet, tests, intentional golden update/review, and
   `tuiweave-audit.sh --strict --verify`.

Pin tuiweave in committed configuration. A local `replace` directive may be used
temporarily during development and should not become the release dependency.
Tuiweave changes require explicit authorization.
