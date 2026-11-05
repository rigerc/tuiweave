# nice-llama migration dogfood

Target reviewed read-only:
`nice-llama-server/docs/tuiweave-migration-plan.md`.

The migration workflow found the following:

- The plan's corrected-assumptions section now rejects the original guessed
  APIs (`SetSuggestions`, `dialog.Message`, `scrollbar.Model`,
  `ForViewport`, and `/` delegated to `tabs.Update`). The skill's compatibility
  reference independently enforces those names.
- The plan explicitly serializes slices that share `model.go`, `render.go`,
  `update.go`, and tests, with isolated worktrees as the only parallel option.
- The textarea section is stale after this rollout: it proposes `Cursor()` and
  an old `ReplaceRange(row, start, end, text)` shape. It must use
  `CursorPosition() textarea.Position` and
  `ReplaceRange(start, end textarea.Position, value string)`. Completion must
  derive the logical token range at the cursor, not mirror `Value()` alone.
- Snapshot requirements are recognized globally, but each slice's task
  contract should name focused/narrow goldens and command-delivery scenarios;
  a generic “update tests” acceptance line is insufficient.
- Mouse routing now retains rectangles and uses `mouse.InBounds`, satisfying
  the multi-pane design requirement. Wheel delegation must remain independent
  of keyboard focus.
- Focus risk is recognized, but “conditional registration” and manager
  “capacity” are not tuiweave APIs. The implementation plan must define concrete
  manager/scope/stack counts, fresh-address `Apply` calls, visibility routing,
  and modal isolation before Slice 10.
- The prerequisites still use `go get ...@latest` and developer-specific
  absolute paths. Generated app instructions must pin tuiweave and use portable
  repository links; local `replace` paths are development-only.

No nice-llama files were modified during this dogfood review.
