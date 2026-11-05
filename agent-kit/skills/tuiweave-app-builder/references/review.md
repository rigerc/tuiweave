# Review mode

Audit:

- theme-role usage and raw colors;
- forbidden Ultraviolet imports;
- tuiweave-owned constraints, `SetSize`, size modes, and application-owned layout;
- MVU model reassignment and command collection;
- fresh focus application after copied model updates;
- mouse-wheel routing by layout bounds without focus gating;
- modal result delivery and app-owned visibility;
- local-coordinate `tabs.IndexAt` mouse routing;
- app-owned inspection bounds, privacy, and current-tree action validation;
- narrow, empty, focus, selection, scrolling, wide-glyph, and combining behavior;
- paired readable/role-aware snapshots and explicit scenario command delivery;
- actual API names/contracts at the pinned tuiweave revision.

Report blocking findings, advisories, and candidate tuiweave gaps separately. Run
the audit read-only and do not regenerate goldens or implement fixes unless the
user requested changes.
