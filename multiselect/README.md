# multiselect

`multiselect` is a focusable set of application-owned options where any number
of entries can be checked.

```go
m := multiselect.New(theme)
m.Title = "Checks"
m.Description = "Run before commit"
m.SetItems(
    action.Item{ID: "types", Label: "Type checking"},
    action.Item{ID: "tests", Label: "Tests"},
)
m.SetCheckedIDs("types")
m.SetSize(width, m.NaturalHeight())
m.Focus()
m, cmd = m.Update(msg)

// Handle multiselect.ChangedMsg and multiselect.SubmittedMsg in the app.
```

The cursor row uses the shared `> ` selector. Space, `x`, and semantic actions
toggle the cursor item and emit `ChangedMsg`; Enter emits `SubmittedMsg` with
the stable IDs of every checked item. Disabled items stay visible but are
skipped by the cursor and cannot be checked. Set `Inline = true` for a
Clack-style prompt block that collapses to the checked labels when blurred.

See the [inline example](../examples/inline) and [form example](../examples/form),
plus the [multiselect recipe](../AGENTS.md#multiselect).

## API highlights

- `SetItems` preserves checked IDs and the cursor ID that are still present.
- `SetCheckedIDs`, `SetChecked`, `Clear`, and `Checked` manage visible state.
- `SelectedIDs`, `SelectedItems`, and `SelectedCount` report the checked set.
- `ChangedMsg` reports one checkbox flip and the full checked ID list.
- `SubmittedMsg` reports confirmation without claiming persistence.
- `Focus`, `Blur`, `Update`, and `View` implement the local field loop.

## Related documentation

- [API reference](https://pkg.go.dev/github.com/rigerc/tuiweave/multiselect)
- [Inline example](../examples/inline)
- [Form example](../examples/form)
- [Multiselect recipe](../AGENTS.md#multiselect)
