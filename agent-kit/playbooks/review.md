# Review playbook

Audit the pinned tuiweave version and actual contracts before reviewing behavior.
Check:

- all colors use `tuiweave.Theme` roles;
- app code does not import Ultraviolet;
- layout rectangles are app-owned, constraints come from tuiweave constructors,
  bounded components receive `SetSize`, and natural heights follow
  `layout.SizeModeOf`;
- each component `Update` return model is reassigned and each command collected;
- fresh component addresses are applied after copied focus-model updates;
- wheel messages reach scrollables even while blurred and use layout bounds
  when multiple panes exist;
- modal command results are delivered and visibility remains app-owned;
- tab clicks are bounds-checked and translated to local coordinates before
  `tabs.IndexAt`;
- inspection trees use stable IDs and app-owned absolute bounds and privacy;
  qualified actions are revalidated for visibility and `Enabled` state before
  their local suffix is routed;
- narrow, empty, focused, scrolled, wide-glyph, and combining-grapheme states
  behave intentionally;
- snapshots cover rendering, themed renderers pair readable and role-aware
  goldens, and scenarios explicitly deliver command-produced results;
- referenced tuiweave API names exist at the pinned revision.

Classify findings as blocking correctness defects, advisory coverage/design
risks, or candidate tuiweave API gaps. Never regenerate goldens during an audit.
