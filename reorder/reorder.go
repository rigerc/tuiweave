// Package reorder provides a focusable list whose items can be moved up and
// down. The component owns the cursor, grab state, and rendered order; it
// emits MovedMsg when an item changes position. Applications own persistence
// and what the order means.
//
// Reorder lists follow the Huh field grammar: Title and Description render
// above the rows, Help below, a thick focus rail marks the active field, and
// the cursor row uses a "> " selector. Inline renders a Clack prompt block with
// a ◆ header, ❯ ≡ rows, instructions, and a └ footer; an unfocused prompt
// collapses to the ordered labels.
package reorder

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/action"
	"github.com/rigerc/tuiweave/inspect"
	"github.com/rigerc/tuiweave/internal/clackstyle"
	"github.com/rigerc/tuiweave/internal/fieldstyle"
	"github.com/rigerc/tuiweave/layout"
)

const (
	// ActionFocus focuses the reorder list.
	ActionFocus = "focus"
	// ActionBlur blurs the reorder list.
	ActionBlur = "blur"
	// ActionNext moves the cursor to the next item.
	ActionNext = "next"
	// ActionPrevious moves the cursor to the previous item.
	ActionPrevious = "previous"
	// ActionFirst moves the cursor to the first item.
	ActionFirst = "first"
	// ActionLast moves the cursor to the last item.
	ActionLast = "last"
	// ActionMoveUp moves the cursor item one position up.
	ActionMoveUp = "move-up"
	// ActionMoveDown moves the cursor item one position down.
	ActionMoveDown = "move-down"
	// ActionGrab starts moving the cursor item with the navigation keys.
	ActionGrab = "grab"
	// ActionDrop releases a grabbed item.
	ActionDrop = "drop"
	// ActionSelectPrefix prefixes stable cursor-selection actions. For example,
	// an item with ID "build" exposes the local action "select.build".
	ActionSelectPrefix = action.SelectPrefix
)

// Item is one reorderable entry. ID should remain stable when the item moves so
// applications and semantic clients can identify it. Label is rendered.
type Item struct {
	ID    string
	Label string
}

// MovedMsg is emitted when an item changes position. From and To are indices
// in the old and new order, and Items is the complete new order.
type MovedMsg struct {
	From  int
	To    int
	Item  Item
	Items []Item
}

// Model is a focusable reorder list. Create one with New.
type Model struct {
	width, height int
	items         []Item
	cursor        int
	off           int
	focused       bool
	grabbed       bool

	// Title is an optional field title. Inline uses it as the Clack prompt
	// message; stacked mode renders it above the rows.
	Title string
	// Description is an optional muted line below the title.
	Description string
	// Help is an optional muted hint rendered below the rows.
	Help string
	// Inline renders a Clack prompt block. The focused prompt expands its
	// rows; an unfocused prompt collapses to the ordered labels.
	Inline bool

	field fieldstyle.Styles
	clack clackstyle.Styles
}

// New returns an empty reorder list styled from the theme's roles.
func New(theme tuiweave.Theme) Model { return NewOn(theme, theme.SurfaceRaised) }

// NewOn returns an empty reorder list with the given background painted behind
// its rows. Use nil for a transparent prompt on the application surface.
func NewOn(theme tuiweave.Theme, background color.Color) Model {
	return Model{
		cursor: -1,
		field:  fieldstyle.NewOn(theme, background),
		clack:  clackstyle.NewOn(theme, background),
	}
}

// SetSize sets the box the reorder list renders in.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
	m.clampCursor()
	m.scrollIntoView()
}

// SizeMode reports a bounded scrolling list for stacked rows and a
// width-bounded, natural-height Clack prompt block for Inline.
func (m Model) SizeMode() layout.SizeMode {
	if m.Inline {
		return layout.SizeWidthBounded
	}
	return layout.SizeBounded
}

// SetItems replaces the items, preserving the cursor item by stable ID when it
// is still present. Otherwise the first item is selected.
func (m *Model) SetItems(items ...Item) {
	cursorID := m.SelectedItem().ID
	m.items = append([]Item(nil), items...)
	m.cursor = 0
	if cursorID != "" {
		for i, item := range m.items {
			if item.ID == cursorID {
				m.cursor = i
				break
			}
		}
	}
	m.clampCursor()
	m.scrollIntoView()
}

// Items returns the current order.
func (m Model) Items() []Item { return append([]Item(nil), m.items...) }

// IDs returns the current stable item IDs in order.
func (m Model) IDs() []string {
	ids := make([]string, len(m.items))
	for i, item := range m.items {
		ids[i] = item.ID
	}
	return ids
}

// Len returns the number of items.
func (m Model) Len() int { return len(m.items) }

// Select moves the cursor to an item by index, clamped to the available items.
func (m *Model) Select(index int) {
	if len(m.items) == 0 {
		m.cursor = -1
		m.grabbed = false
		return
	}
	m.cursor = max(0, min(index, len(m.items)-1))
	m.scrollIntoView()
}

// SelectID moves the cursor to an item by stable ID. It is a no-op when the ID
// is unknown or empty.
func (m *Model) SelectID(id string) {
	if id == "" {
		return
	}
	for i, item := range m.items {
		if item.ID == id {
			m.Select(i)
			return
		}
	}
}

// Selected returns the cursor index, or -1 when the list is empty.
func (m Model) Selected() int { return m.cursor }

// SelectedItem returns the item under the cursor, or a zero Item when the list
// is empty.
func (m Model) SelectedItem() Item {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return Item{}
	}
	return m.items[m.cursor]
}

// Focus makes the list respond to navigation, grab, and move keys.
func (m *Model) Focus() { m.focused = true }

// Blur stops the list from responding to keys and drops any grabbed item.
func (m *Model) Blur() {
	m.focused = false
	m.grabbed = false
}

// Focused reports whether the list handles keys.
func (m Model) Focused() bool { return m.focused }

// Grabbed reports whether navigation keys currently move the cursor item.
func (m Model) Grabbed() bool { return m.grabbed }

// Grab starts grabbed mode without emitting a message.
func (m *Model) Grab() {
	if len(m.items) > 0 {
		m.grabbed = true
	}
}

// Drop releases grabbed mode without emitting a message.
func (m *Model) Drop() { m.grabbed = false }

// NaturalHeight reports the rows needed to show every item plus optional title,
// description, and help. Inline height depends on focus: two rows when
// collapsed, the full prompt when focused.
func (m Model) NaturalHeight() int {
	if m.Inline {
		if !m.focused {
			return 2
		}
		height := 3 // header, instructions, footer
		if m.Description != "" {
			height++
		}
		height += max(1, len(m.items))
		if m.Help != "" {
			height++
		}
		return height
	}
	return fieldstyle.HeaderHeight(m.Title, m.Description) + boolInt(m.Help != "") + len(m.items)
}

// TotalLines returns the item count for scrollbar.Scrollable.
func (m Model) TotalLines() int { return len(m.items) }

// VisibleLines returns the number of visible rows.
func (m Model) VisibleLines() int {
	if m.Inline {
		return max(1, len(m.items))
	}
	return m.itemArea()
}

// YOffset returns the first visible item index.
func (m Model) YOffset() int { return m.off }

// MoveUp moves the cursor item one position up and returns a command that
// emits MovedMsg. It returns nil at the top of the list.
func (m *Model) MoveUp() tea.Cmd { return m.moveTo(m.cursor - 1) }

// MoveDown moves the cursor item one position down and returns a command that
// emits MovedMsg. It returns nil at the bottom of the list.
func (m *Model) MoveDown() tea.Cmd { return m.moveTo(m.cursor + 1) }

// MoveTo moves the cursor item to target (clamped) and returns a command that
// emits MovedMsg. It returns nil when the item cannot move.
func (m *Model) MoveTo(target int) tea.Cmd { return m.moveTo(target) }

// Update handles semantic actions and, while focused, navigation, shift/ctrl
// arrows for direct moves, and space/x/enter to grab or drop the cursor item.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if next, cmd, handled := m.applyAction(msg); handled {
		return next, cmd
	}
	if !m.focused {
		return m, nil
	}
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "up", "k":
		if m.grabbed {
			return m, m.MoveUp()
		}
		m.selectRelative(-1)
	case "down", "j":
		if m.grabbed {
			return m, m.MoveDown()
		}
		m.selectRelative(1)
	case "left", "h":
		if m.Inline {
			if m.grabbed {
				return m, m.MoveUp()
			}
			m.selectRelative(-1)
		}
	case "right", "l":
		if m.Inline {
			if m.grabbed {
				return m, m.MoveDown()
			}
			m.selectRelative(1)
		}
	case "shift+up", "ctrl+up", "K":
		return m, m.MoveUp()
	case "shift+down", "ctrl+down", "J":
		return m, m.MoveDown()
	case "pgup":
		if m.grabbed {
			return m, m.MoveTo(m.cursor - max(1, m.itemArea()))
		}
		m.selectRelative(-max(1, m.itemArea()))
	case "pgdown":
		if m.grabbed {
			return m, m.MoveTo(m.cursor + max(1, m.itemArea()))
		}
		m.selectRelative(max(1, m.itemArea()))
	case "home", "g":
		if m.grabbed {
			return m, m.MoveTo(0)
		}
		m.Select(0)
	case "end", "G":
		if m.grabbed {
			return m, m.MoveTo(len(m.items) - 1)
		}
		m.Select(len(m.items) - 1)
	case "space", "x", "enter":
		if m.grabbed {
			m.Drop()
		} else {
			m.Grab()
		}
	}
	return m, nil
}

// View renders the Huh-style rows or, when Inline is set, a Clack-style prompt
// block with a ◆ header, │ rows, ❯ ≡ items, and a └ footer.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	rows := make([]string, 0, m.height)

	if m.Inline {
		rows = append(rows, m.clack.HeaderRow(m.focused, m.Title, m.width))
		if !m.focused {
			rows = append(rows, m.clack.Row(false, m.clack.Inactive.Render(m.summary()), m.width))
			if len(rows) > m.height {
				rows = rows[:m.height]
			}
			return strings.Join(rows, "\n")
		}
		if m.Description != "" {
			rows = append(rows, m.clack.Row(m.focused, m.clack.Desc.Render(m.Description), m.width))
		}
		if len(m.items) == 0 {
			rows = append(rows, m.clack.Row(m.focused, m.clack.Inactive.Render("No items"), m.width))
		} else {
			for i, item := range m.items {
				rows = append(rows, m.clack.ReorderRow(m.focused, i == m.cursor, i == m.cursor && m.grabbed, item.Label, m.width))
			}
		}
		rows = append(rows, m.clack.InstructionsRow(m.focused, m.instructions(), m.width))
		if m.Help != "" {
			rows = append(rows, m.clack.HelpRow(m.focused, m.Help, m.width))
		}
		rows = append(rows, m.clack.FooterRow(m.width))
		if len(rows) > m.height {
			rows = rows[:m.height]
		}
		return strings.Join(rows, "\n")
	}

	rows = append(rows, m.field.HeaderRows(m.focused, m.Title, m.Description, m.width)...)
	for i := 0; i < m.itemArea(); i++ {
		index := m.off + i
		if index >= len(m.items) {
			rows = append(rows, m.field.FrameRow(m.focused, "", m.width))
			continue
		}
		rows = append(rows, m.renderItem(index))
	}
	if help := m.field.HelpRow(m.focused, m.Help, m.width); help != "" {
		rows = append(rows, help)
	}
	for len(rows) < m.height {
		rows = append(rows, m.field.FrameRow(m.focused, "", m.width))
	}
	if len(rows) > m.height {
		rows = rows[:m.height]
	}
	return strings.Join(rows, "\n")
}

func (m Model) renderItem(index int) string {
	bodyWidth := fieldstyle.BodyWidth(m.width)
	if bodyWidth <= 0 {
		return m.field.FrameRow(m.focused, "", m.width)
	}

	cursor := index == m.cursor
	grabbed := cursor && m.grabbed
	selector := m.field.Selector(cursor)
	if grabbed {
		selector = m.field.SelectMark.Render("↕ ")
	}
	handle := m.field.Normal.Render("≡")
	if cursor {
		handle = m.field.SelectMark.Render("≡")
	}
	if grabbed {
		handle = m.field.Selected.Render("≡")
	}
	selectorWidth := lipgloss.Width(selector)
	handleWidth := lipgloss.Width(handle)

	var body string
	switch {
	case bodyWidth >= selectorWidth+handleWidth+1:
		labelWidth := bodyWidth - selectorWidth - handleWidth - 1
		label := ansi.Truncate(m.items[index].Label, labelWidth, "…")
		body = selector + handle + m.field.Pad.Render(" ") + m.field.Label(label, cursor, false)
	case bodyWidth >= selectorWidth+handleWidth:
		body = selector + handle
	default:
		body = handle
	}
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func (m Model) summary() string {
	if len(m.items) == 0 {
		return "No items"
	}
	labels := make([]string, 0, len(m.items))
	for _, item := range m.items {
		labels = append(labels, item.Label)
	}
	return strings.Join(labels, ", ")
}

func (m Model) instructions() string {
	if m.grabbed {
		return "↑/↓ move • Space drop"
	}
	return "↑/↓ • Space grab • Shift+↑/↓ move"
}

func (m *Model) clampCursor() {
	if len(m.items) == 0 {
		m.cursor = -1
		m.grabbed = false
		return
	}
	m.cursor = max(0, min(m.cursor, len(m.items)-1))
}

func (m *Model) selectRelative(delta int) {
	if len(m.items) == 0 {
		m.cursor = -1
		return
	}
	index := m.cursor
	if index < 0 {
		if delta < 0 {
			index = len(m.items)
		} else {
			index = -1
		}
	}
	index = max(0, min(index+delta, len(m.items)-1))
	m.Select(index)
}

func (m *Model) moveTo(target int) tea.Cmd {
	from := m.cursor
	if from < 0 || from >= len(m.items) || len(m.items) < 2 {
		return nil
	}
	target = max(0, min(target, len(m.items)-1))
	if target == from {
		return nil
	}

	item := m.items[from]
	next := make([]Item, 0, len(m.items))
	next = append(next, m.items[:from]...)
	next = append(next, m.items[from+1:]...)
	next = append(next, Item{})
	copy(next[target+1:], next[target:])
	next[target] = item

	m.items = next
	m.cursor = target
	m.scrollIntoView()

	items := append([]Item(nil), m.items...)
	return func() tea.Msg {
		return MovedMsg{From: from, To: target, Item: item, Items: items}
	}
}

func (m Model) itemArea() int {
	area := m.height - fieldstyle.HeaderHeight(m.Title, m.Description)
	if m.Help != "" {
		area--
	}
	return max(0, area)
}

func (m *Model) scrollIntoView() {
	if m.Inline || len(m.items) == 0 {
		m.off = 0
		return
	}
	area := m.itemArea()
	if area <= 0 {
		m.off = 0
		return
	}
	if m.cursor < 0 {
		m.off = 0
		return
	}
	if m.cursor < m.off {
		m.off = m.cursor
	}
	if m.cursor >= m.off+area {
		m.off = m.cursor - area + 1
	}
	m.off = max(0, min(m.off, max(0, len(m.items)-area)))
}

func (m Model) applyAction(msg tea.Msg) (Model, tea.Cmd, bool) {
	event, ok := msg.(inspect.ActionMsg)
	if !ok {
		return m, nil, false
	}
	switch event.ID {
	case ActionFocus:
		m.Focus()
	case ActionBlur:
		m.Blur()
	case ActionNext:
		m.selectRelative(1)
	case ActionPrevious:
		m.selectRelative(-1)
	case ActionFirst:
		m.Select(0)
	case ActionLast:
		m.Select(len(m.items) - 1)
	case ActionMoveUp:
		return m, m.MoveUp(), true
	case ActionMoveDown:
		return m, m.MoveDown(), true
	case ActionGrab:
		m.Grab()
	case ActionDrop:
		m.Drop()
	default:
		id, ok := action.ParseSelectID(event.ID)
		if !ok {
			return m, nil, false
		}
		m.SelectID(id)
		return m, nil, true
	}
	return m, nil, true
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
