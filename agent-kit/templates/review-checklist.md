# tuiweave review checklist

- [ ] Tuiweave dependency is pinned and API names match that revision.
- [ ] The Go directive and toolchain meet that pin's requirement and the selected dependency graph; exact-minimum checks use `GOTOOLCHAIN=local`.
- [ ] Colors come from semantic theme roles.
- [ ] App code has no Ultraviolet imports.
- [ ] Layout constraints use tuiweave constructors; no Ultraviolet values leak into app state.
- [ ] `WindowSizeMsg` owns layout, bounded components receive `SetSize`, and natural heights follow `layout.SizeModeOf`.
- [ ] Returned models are reassigned and commands collected.
- [ ] Focus managers receive fresh addresses after value-model updates.
- [ ] Wheel events are routed by bounds and are not focus-gated.
- [ ] Modal results are delivered and modal visibility is app-owned.
- [ ] Tab clicks are bounds-checked and translated to local x before `tabs.IndexAt`.
- [ ] Inspection IDs, absolute bounds, privacy, and qualified-action routing are app-owned.
- [ ] Actions are validated against the current visible tree before invocation.
- [ ] Narrow, empty, focus, selection, scroll, wide-glyph, and combining states are covered.
- [ ] Themed rendering pairs readable and role-aware goldens; transitions have named scenarios with explicit command delivery.
- [ ] Golden changes were intentional and reviewed.
- [ ] Reusable gaps are reported instead of implemented as local forks.
