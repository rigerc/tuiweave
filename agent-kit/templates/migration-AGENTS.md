# TUI migration agent instructions

Target tuiweave revision: `<PINNED_VERSION_OR_COMMIT>`

Use the tuiweave source documentation at https://github.com/ishan5ain/tuiweave and
the migration playbook in tuiweave-agent-kit. Inventory existing rendering,
input, scrolling, focus, mouse, modal, inspection/action, and domain behavior
before editing.
Classify every concern as reusable tuiweave behavior, application-owned behavior,
or a candidate tuiweave API gap.

Work in vertical slices using `migration-plan.md`. Preserve behavior before
intentional visual changes and review every golden diff. Do not parallelize
shared model/render/update files without isolated worktrees. Tuiweave changes and
local forks require explicit authorization. A local `replace` directive is
development-only; committed configuration remains pinned.
