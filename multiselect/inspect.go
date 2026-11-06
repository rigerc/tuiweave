package multiselect

import (
	"strconv"
	"strings"

	"github.com/ishan5ain/tuiweave/action"
	"github.com/ishan5ain/tuiweave/inspect"
)

// Actions reports stable local intents for focus, cursor navigation, checkbox
// changes, clearing, submission, and per-item selection. Disabled items remain
// discoverable with Enabled=false.
func (m Model) Actions() []inspect.Action {
	cursor := m.CursorItem()
	actions := []inspect.Action{
		{ID: ActionFocus, Label: "Focus multiselect", Enabled: !m.focused},
		{ID: ActionBlur, Label: "Blur multiselect", Enabled: m.focused},
		{ID: ActionNext, Label: "Next option", Enabled: m.nextEnabled() >= 0},
		{ID: ActionPrevious, Label: "Previous option", Enabled: m.previousEnabled() >= 0},
		{ID: ActionFirst, Label: "First option", Enabled: m.firstEnabled() >= 0},
		{ID: ActionLast, Label: "Last option", Enabled: m.lastEnabled() >= 0},
		{ID: ActionToggle, Label: "Toggle " + cursor.Label, Enabled: m.cursor >= 0 && !cursor.Disabled},
		{ID: ActionCheck, Label: "Check " + cursor.Label, Enabled: m.cursor >= 0 && !cursor.Disabled && !m.isChecked(m.cursor)},
		{ID: ActionUncheck, Label: "Uncheck " + cursor.Label, Enabled: m.cursor >= 0 && !cursor.Disabled && m.isChecked(m.cursor)},
		{ID: ActionClear, Label: "Clear selection", Enabled: m.SelectedCount() > 0},
		{ID: ActionSubmit, Label: "Submit selection", Enabled: true},
	}
	for _, item := range m.items {
		if item.ID == "" {
			continue
		}
		actions = append(actions,
			inspect.Action{
				ID:          action.SelectID(item.ID),
				Label:       "Focus " + item.Label,
				Description: item.Description,
				Enabled:     !item.Disabled,
			},
			inspect.Action{
				ID:          ActionTogglePrefix + item.ID,
				Label:       "Toggle " + item.Label,
				Description: item.Description,
				Enabled:     !item.Disabled,
			},
		)
	}
	return actions
}

// Inspect reports the cursor, checked IDs, focus, and scroll position. It
// intentionally does not include application command details.
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
		"item_count":     strconv.Itoa(len(m.items)),
		"enabled_count":  strconv.Itoa(action.EnabledCount(m.items)),
		"selected_count": strconv.Itoa(m.SelectedCount()),
	}
	if ids := m.SelectedIDs(); len(ids) > 0 {
		attributes["selected_ids"] = strings.Join(ids, ",")
	}
	if m.Title != "" {
		attributes["title"] = m.Title
	}
	if m.Inline {
		attributes["inline"] = "true"
	}
	return inspect.Node{
		Kind:       "multiselect",
		Bounds:     inspect.Bounds{Width: m.width, Height: m.height},
		Label:      m.CursorItem().Label,
		Focused:    m.focused,
		Selected:   selected,
		Scroll:     &inspect.Scroll{Total: m.TotalLines(), Visible: m.VisibleLines(), Offset: m.YOffset()},
		Actions:    m.Actions(),
		Attributes: attributes,
	}
}

func (m Model) nextEnabled() int {
	if m.cursor < 0 {
		return m.firstEnabled()
	}
	for i := m.cursor + 1; i < len(m.items); i++ {
		if !m.items[i].Disabled {
			return i
		}
	}
	return -1
}

func (m Model) previousEnabled() int {
	if m.cursor < 0 {
		return m.lastEnabled()
	}
	for i := m.cursor - 1; i >= 0; i-- {
		if !m.items[i].Disabled {
			return i
		}
	}
	return -1
}

func (m Model) isChecked(index int) bool {
	return index >= 0 && index < len(m.checked) && m.checked[index]
}
