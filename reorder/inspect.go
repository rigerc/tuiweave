package reorder

import (
	"strconv"
	"strings"

	"github.com/rigerc/tuiweave/action"
	"github.com/rigerc/tuiweave/inspect"
)

// Actions reports stable local intents for focus, navigation, movement, grab
// state, and per-item cursor selection.
func (m Model) Actions() []inspect.Action {
	actions := []inspect.Action{
		{ID: ActionFocus, Label: "Focus reorder list", Enabled: !m.focused},
		{ID: ActionBlur, Label: "Blur reorder list", Enabled: m.focused},
		{ID: ActionNext, Label: "Next item", Enabled: m.cursor >= 0 && m.cursor < len(m.items)-1},
		{ID: ActionPrevious, Label: "Previous item", Enabled: m.cursor > 0},
		{ID: ActionFirst, Label: "First item", Enabled: len(m.items) > 0 && m.cursor != 0},
		{ID: ActionLast, Label: "Last item", Enabled: len(m.items) > 0 && m.cursor != len(m.items)-1},
		{ID: ActionMoveUp, Label: "Move item up", Enabled: m.cursor > 0},
		{ID: ActionMoveDown, Label: "Move item down", Enabled: m.cursor >= 0 && m.cursor < len(m.items)-1},
		{ID: ActionGrab, Label: "Grab item", Enabled: !m.grabbed && len(m.items) > 0},
		{ID: ActionDrop, Label: "Drop item", Enabled: m.grabbed},
	}
	for _, item := range m.items {
		if item.ID == "" {
			continue
		}
		actions = append(actions, inspect.Action{
			ID:      action.SelectID(item.ID),
			Label:   "Focus " + item.Label,
			Enabled: true,
		})
	}
	return actions
}

// Inspect reports the current order, cursor, grab state, focus, and scroll
// position. It intentionally does not include application command details.
func (m Model) Inspect() inspect.Node {
	var selected *inspect.Selection
	if m.cursor >= 0 && m.cursor < len(m.items) {
		selected = &inspect.Selection{
			Index: m.cursor,
			Count: len(m.items),
			Label: m.items[m.cursor].Label,
		}
	}
	attributes := map[string]string{
		"item_count": strconv.Itoa(len(m.items)),
		"grabbed":    strconv.FormatBool(m.grabbed),
	}
	if ids := m.IDs(); len(ids) > 0 {
		attributes["order"] = strings.Join(ids, ",")
	}
	if id := m.SelectedItem().ID; id != "" {
		attributes["selected_id"] = id
	}
	if m.Title != "" {
		attributes["title"] = m.Title
	}
	if m.Inline {
		attributes["inline"] = "true"
	}
	status := "ready"
	if m.grabbed {
		status = "grabbed"
	}
	return inspect.Node{
		Kind:       "reorder",
		Bounds:     inspect.Bounds{Width: m.width, Height: m.height},
		Label:      m.SelectedItem().Label,
		Focused:    m.focused,
		Status:     status,
		Selected:   selected,
		Scroll:     &inspect.Scroll{Total: m.TotalLines(), Visible: m.VisibleLines(), Offset: m.YOffset()},
		Actions:    m.Actions(),
		Attributes: attributes,
	}
}
