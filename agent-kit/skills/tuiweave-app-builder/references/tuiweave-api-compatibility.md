# tuiweave API compatibility

Generated from tuiweave revision `c05ccf548162f861736844f729c18104949c90e5`. Verify against the application pinned revision.

## Supported versions


| Surface | Supported version | Evidence |
|---|---|---|
| Go | 1.26.0 minimum | `go.mod`; Linux CI at the exact floor |
| Go | current stable | Linux, macOS, and Windows CI |
| Bubble Tea | `charm.land/bubbletea/v2 v2.1.0` | `go.mod`; full CI suite and consumer runtime test |
| Lip Gloss | `charm.land/lipgloss/v2 v2.0.6` | `go.mod`; full CI suite |
| Glamour | `charm.land/glamour/v2 v2.0.1` | `go.mod`; agentic Markdown tests |
| Ultraviolet | `v0.0.0-20261008173134-6b8d4baf91b4` | `go.mod`; layout, overlay, and SnapCells coverage |
| ANSI utilities | `github.com/charmbracelet/x/ansi v0.11.9` | `go.mod`; rendering and cell-width coverage |

The supported dependency graph is the graph selected by `go.mod`. tuiweave
does not promise compatibility with arbitrary newer upstream versions merely
because Go's module solver can select them. Dependency upgrades are deliberate
compatibility events and include changelog review, the full test suite, and
golden review when rendering changes.

The matrix is reproduced in CI by testing the exact minimum Go version on
Linux and current stable Go on Linux, macOS, and Windows. Locally, validate the
selected graph with:

```sh
go mod tidy -diff
go mod verify
go build ./...
go vet ./...
go test ./...
go test -race ./...
```


## Go version and module boundary


tuiweave v1 will remain one Go module and will require Go 1.26.0 or newer.
Bubble Tea v2.1.0 and the selected Ultraviolet revision establish that floor.
Splitting `agentic` into another module would add coordinated
release tags, dependency wiring, and example complexity without demonstrated
application demand.

The decision can be revisited in a future major version if applications need a
lower Go floor or the agentic packages require an independent release cadence.
Until then, all top-level and `agentic/...` packages share one module version.

### Go 1.26 decision (2026-10-10)

tuiweave v0.3.0 raises the minimum from Go 1.25.8 to Go 1.26.0 to
adopt [Bubble Tea v2.1.0](https://github.com/charmbracelet/bubbletea/releases/tag/v2.1.0).
The upgrade brings fixes for shutdown signal deadlocks, zero-size terminal
startup and resizing, and terminal queries issued with input disabled. Its
required Ultraviolet update also fixes renderer repainting and input handling.
These runtime fixes justify revising the earlier Go floor; the review found
no documented consumer requirement to retain Go 1.25.8. This does not imply
that every external consumer has been audited.

Consumers must install Go 1.26.0 or newer before adopting v0.3.0.
With `GOTOOLCHAIN=local`, Go 1.25 rejects the upgraded module. Existing tags
retain their original requirements. No tuiweave API migration is required,
and terminal Program Status Protocol reporting remains application-owned and
opt-in through Bubble Tea's `tea.View`.

## Relevant package contracts

| Package | Contract |
|---|---|
| `tuiweave` | `tuiweave.Presets()` + `ThemeForPreset`, or a named constructor; never raw colors |
| `layout` | `layout.Vertical(...).Apply(area, &components...)`; `SizeModeOf` for sizing exceptions |
| `snaptest` | `Snap`, grapheme-preserving `SnapCells`, `RunScenario`, `SnapScenario`; follow the compile-checked tutorial |
| `inspect` | `Inspect()`, `Bind`, `BindAt`, `Group`, `Marshal`; follow the field/action contract in `inspect/README.md` |
| `focus` | `NewManager(n)`, `NewScope(n)`, or `NewStack(n)`; apply fresh addresses after changes |
| `scrollbar` | `scrollbar.For(theme, component)` in a 1-cell layout segment |
| `textarea` | Logical-rune selection, cell-aware wrapping, bounded undo/redo, word movement/deletion, bounded kill/yank, `CursorPosition`, and one-edit `ReplaceRange`; Enter inserts a newline |
| `dialog` | App owns visibility; width-bounded titles, bodies, and buttons; result arrives as `ResultMsg` |
| `tabs` | `SetTabs`, `Focus`, `SelectedID`; `IndexAt` maps a local cell x-coordinate to a rendered tab |
| `autocomplete` | `SetItems`, `SetQuery`, `SelectedMsg`; app owns query and insertion, disabled candidates are skipped |

## Compatibility notes


- `dialog` and `agentic/permission` are width-bounded and answer with a typed
  `ResultMsg`; the application owns visibility and modal routing. There is no
  `dialog.Message`.
- `autocomplete` is a suggestion window beside an app-owned input. Call
  `SetItems` and `SetQuery`; accept with `SelectedMsg` and apply insertion in
  the app. It does not expose `SetSuggestions`.
- `scrollbar.For(theme, component)` is the standalone renderer for scrolling
  components. There is no `scrollbar.Model` or `scrollbar.ForViewport`.
- `focus.Manager`, `focus.Scope`, and `focus.Stack` store indices/state, not
  component pointers. Apply fresh addresses after copied model updates.
- `textarea.Position` is logical rune space. `CursorPosition` reports the
  cursor, and `ReplaceRange` clamps/normalizes a range, records one undoable
  edit, and moves the cursor to the inserted value's end.
- `layout` owns geometry through `layout.Rect`; applications call `SetSize`
  from their `WindowSizeMsg` layout pass and do not import Ultraviolet.


## Layout constraint compatibility


The stable application API is the `layout.Constraint` name plus the `Len`,
`Min`, `Max`, `Percent`, `Ratio`, and `Fill` constructors. tuiweave owns the
sealed constraint representation and translates it to the Ultraviolet solver
only when constructing a `Vertical` or `Horizontal` layout. Ultraviolet
constraint identity is not part of the v1 contract.

Applications following the documented API require no migration:

```go
parts := []layout.Constraint{
	layout.Len(3),
	layout.Fill(1),
}
view := layout.Vertical(parts...)
```

Direct construction of, conversion to, or storage as
`github.com/charmbracelet/ultraviolet/layout.Constraint` is unsupported and no
longer type-compatible. An application that used that pre-v1 implementation
detail must replace Ultraviolet constraint values with tuiweave's constructors.
This boundary was completed in
[Issue #28](https://github.com/ishan5ain/tuiweave/issues/28).


## Source of truth

Read `AGENT-CONTRACT.md`, `COMPATIBILITY.md`, `AGENT-CATALOG.md`, the matching `AGENTS.md` recipe, and the closest runnable example from the pinned tuiweave checkout before editing.
