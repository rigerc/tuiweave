# reorder

`reorder` is a focusable list whose items can be moved up and down.

```go
r := reorder.New(theme)
r.Title = "Steps"
r.Description = "Release pipeline order"
r.SetItems(
    reorder.Item{ID: "design", Label: "Design"},
    reorder.Item{ID: "build", Label: "Build"},
    reorder.Item{ID: "ship", Label: "Ship"},
)
r.SetSize(width, r.NaturalHeight())
r.Focus()
r, cmd = r.Update(msg)

// Handle reorder.MovedMsg in the application.
```

Arrow keys move the cursor. Shift+↑/↓ (or ctrl+↑/↓, `K`/`J`) move the item
directly; Space grabs the cursor item, arrow keys move it, and Space drops it.
Moves emit `MovedMsg` with the old and new indices plus the complete new order.
Set `Inline = true` for a Clack-style prompt block that collapses to the ordered
labels when blurred.

See the [inline example](../examples/inline) and [form example](../examples/form),
plus the [reorder recipe](../AGENTS.md#reorder).

## API highlights

- `SetItems` preserves the cursor item by stable ID when it is still present.
- `MoveUp`, `MoveDown`, and `MoveTo` mutate the order and return a command that
  emits `MovedMsg`, or `nil` when the item cannot move.
- `Grab`, `Drop`, and `Grabbed` manage cursor-driven movement.
- `Items`, `IDs`, `Len`, `Selected`, and `SelectedItem` report the current order.
- `Focus`, `Blur`, `Update`, and `View` implement the local field loop.

## Related documentation

- [API reference](https://pkg.go.dev/github.com/rigerc/tuiweave/reorder)
- [Inline example](../examples/inline)
- [Form example](../examples/form)
- [Reorder recipe](../AGENTS.md#reorder)
