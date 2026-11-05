// Package menu provides a focusable vertical menu of application-owned
// actions. It renders the choices and emits a SelectedMsg on activation; the
// application owns the resulting command's behavior.
//
// Menus follow the Huh field grammar: Title and Description render above the
// options, Help below, a thick focus rail marks the active field, and the
// selected option uses a "> " selector. Inline renders a Clack prompt block
// with a ◆ header, ❯ ●/○ options, instructions, and a └ footer; an unfocused
// prompt collapses to the selected value.
package menu

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/action"
	"github.com/ishan5ain/tuiweave/inspect"
	"github.com/ishan5ain/tuiweave/internal/clackstyle"
	"github.com/ishan5ain/tuiweave/internal/fieldstyle"
	"github.com/ishan5ain/tuiweave/layout"
)

const (
	// ActionFocus focuses the menu.
	ActionFocus = "focus"
	// ActionBlur blurs the menu.
	ActionBlur = "blur"
	// ActionNext selects the next enabled item.
	ActionNext = "next"
	// ActionPrevious selects the previous enabled item.
	ActionPrevious = "previous"
	// ActionFirst selects the first enabled item.
	ActionFirst = "first"
	// ActionLast selects the last enabled item.
	ActionLast = "last"
	// ActionActivate activates the selected item.
	ActionActivate = "activate"
	// ActionSelectPrefix prefixes stable item-selection actions. For example,
	// an item with ID "open" exposes the local action "select.open".
	ActionSelectPrefix = action.SelectPrefix
)

// Item is the shared action.Item definition rendered by the menu.
type Item = action.Item

// SelectedMsg is emitted when a menu item is activated with enter or the
// ActionActivate semantic action.
type SelectedMsg struct {
	ID    string
	Index int
	Label string
}

// Model is a focusable, scrolling menu. Create one with New.
type Model struct {
	width, height int
	items         []Item
	selected      int
	off           int
	focused       bool

	// Title is an optional field title. Inline uses it as the Clack prompt
	// message; stacked mode renders it above the options.
	Title string
	// Description is an optional muted line below the title.
	Description string
	// Help is an optional muted hint rendered below the options.
	Help string
	// Inline renders a Clack prompt block. The focused prompt expands its
	// options; an unfocused prompt collapses to the selected value.
	Inline bool

	field fieldstyle.Styles
	clack clackstyle.Styles
}

// New returns an empty menu styled from the theme's roles.
func New(theme tuiweave.Theme) Model { return NewOn(theme, theme.SurfaceRaised) }

// NewOn returns an empty menu with the given background painted behind its
// rows. Use nil for a transparent menu on the application surface.
func NewOn(theme tuiweave.Theme, background color.Color) Model {
	return Model{
		selected: -1,
		field:    fieldstyle.NewOn(theme, background),
		clack:    clackstyle.NewOn(theme, background),
	}
}

// SetSize sets the box the menu renders in.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
	m.scrollIntoView()
}

// SizeMode reports a bounded scrolling menu for stacked options and a
// width-bounded, natural-height Clack prompt block for Inline.
func (m Model) SizeMode() layout.SizeMode {
	if m.Inline {
		return layout.SizeWidthBounded
	}
	return layout.SizeBounded
}

// SetItems replaces the menu items, preserving the selected ID when it is
// still enabled. Otherwise the first enabled item is selected.
func (m *Model) SetItems(items ...Item) {
	selectedID := m.SelectedID()
	m.items = append([]Item(nil), items...)
	m.selected = -1
	if selectedID != "" {
		for i, item := range m.items {
			if item.ID == selectedID && !item.Disabled {
				m.selected = i
				break
			}
		}
	}
	if m.selected < 0 {
		m.selected = m.firstEnabled()
	}
	m.scrollIntoView()
}

// Items returns the configured menu items.
func (m Model) Items() []Item { return append([]Item(nil), m.items...) }

// NaturalHeight reports the rows needed to show every option plus optional
// title, description, and help. Inline height depends on focus: two rows when
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

// Select selects an enabled item by index. Disabled or out-of-range indices
// leave the current selection unchanged.
func (m *Model) Select(index int) {
	if index < 0 || index >= len(m.items) || m.items[index].Disabled {
		return
	}
	m.selected = index
	m.scrollIntoView()
}

// SelectID selects an enabled item by stable ID. It is a no-op when no enabled
// item has that ID.
func (m *Model) SelectID(id string) {
	for i, item := range m.items {
		if item.ID == id && !item.Disabled {
			m.Select(i)
			return
		}
	}
}

// Selected returns the selected item index, or -1 when no enabled item exists.
func (m Model) Selected() int { return m.selected }

// SelectedItem returns the selected item, or a zero Item when there is no
// selection.
func (m Model) SelectedItem() Item {
	if m.selected < 0 || m.selected >= len(m.items) {
		return Item{}
	}
	return m.items[m.selected]
}

// SelectedID returns the selected item's stable ID, or an empty string when
// there is no selection or the item has no ID.
func (m Model) SelectedID() string { return m.SelectedItem().ID }

// Focus makes the menu respond to navigation and activation keys.
func (m *Model) Focus() { m.focused = true }

// Blur stops the menu from responding to keys.
func (m *Model) Blur() { m.focused = false }

// Focused reports whether the menu handles keys.
func (m Model) Focused() bool { return m.focused }

// TotalLines returns the configured item count for scrollbar.Scrollable.
func (m Model) TotalLines() int { return len(m.items) }

// VisibleLines returns the number of visible option rows.
func (m Model) VisibleLines() int {
	if m.Inline {
		return max(1, len(m.items))
	}
	return m.itemArea()
}

// YOffset returns the first visible item index.
func (m Model) YOffset() int { return m.off }

// Activate returns a command that emits SelectedMsg for the current enabled
// selection, or nil when the menu has no selection.
func (m Model) Activate() tea.Cmd {
	item := m.SelectedItem()
	if m.selected < 0 || item.Disabled {
		return nil
	}
	return func() tea.Msg {
		return SelectedMsg{ID: item.ID, Index: m.selected, Label: item.Label}
	}
}

// Update handles semantic actions and, while focused, up/down arrows, j/k,
// page navigation, home/g, end/G, and enter activation.
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
		m.selectRelative(-1)
	case "down", "j":
		m.selectRelative(1)
	case "left", "h":
		if m.Inline {
			m.selectRelative(-1)
		}
	case "right", "l":
		if m.Inline {
			m.selectRelative(1)
		}
	case "pgup":
		m.selectRelative(-max(1, m.itemArea()))
	case "pgdown":
		m.selectRelative(max(1, m.itemArea()))
	case "home", "g":
		m.Select(m.firstEnabled())
	case "end", "G":
		m.Select(m.lastEnabled())
	case "enter":
		return m, m.Activate()
	}
	return m, nil
}

// View renders the Huh-style option list or, when Inline is set, a Clack-style
// prompt block with a ◆ header, │ rows, ❯ ●/○ options, and a └ footer.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	rows := make([]string, 0, m.height)

	if m.Inline {
		rows = append(rows, m.clack.HeaderRow(m.focused, m.Title, m.width))
		if !m.focused {
			value := "No selection"
			if item := m.SelectedItem(); item.Label != "" {
				value = item.Label
			}
			rows = append(rows, m.clack.Row(false, m.clack.Inactive.Render(value), m.width))
			if len(rows) > m.height {
				rows = rows[:m.height]
			}
			return strings.Join(rows, "\n")
		}
		if m.Description != "" {
			rows = append(rows, m.clack.Row(m.focused, m.clack.Desc.Render(m.Description), m.width))
		}
		if len(m.items) == 0 {
			rows = append(rows, m.clack.Row(m.focused, m.clack.Inactive.Render("No options"), m.width))
		} else {
			for i, item := range m.items {
				rows = append(rows, m.clack.OptionRow(m.focused, i == m.selected, item.Disabled, item.Label, m.width))
			}
		}
		rows = append(rows, m.clack.InstructionsRow(m.focused, "↑/↓ to navigate • Enter: confirm", m.width))
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
	item := m.items[index]
	bodyWidth := fieldstyle.BodyWidth(m.width)
	if bodyWidth <= 0 {
		return m.field.FrameRow(m.focused, "", m.width)
	}

	selected := index == m.selected && !item.Disabled
	selector := m.field.Selector(selected)
	labelWidth := max(0, bodyWidth-lipgloss.Width(selector))
	label := ansi.Truncate(item.Label, labelWidth, "…")
	body := selector + m.field.Label(label, selected, item.Disabled)
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func (m Model) firstEnabled() int {
	for i, item := range m.items {
		if !item.Disabled {
			return i
		}
	}
	return -1
}

func (m Model) lastEnabled() int {
	for i := len(m.items) - 1; i >= 0; i-- {
		if !m.items[i].Disabled {
			return i
		}
	}
	return -1
}

func (m *Model) selectRelative(delta int) {
	if len(m.items) == 0 {
		return
	}
	if m.selected < 0 {
		if delta < 0 {
			m.Select(m.lastEnabled())
		} else {
			m.Select(m.firstEnabled())
		}
		return
	}
	step := 1
	if delta < 0 {
		step = -1
	}
	remaining := abs(delta)
	index := m.selected
	for remaining > 0 {
		index += step
		for index >= 0 && index < len(m.items) && m.items[index].Disabled {
			index += step
		}
		if index < 0 || index >= len(m.items) {
			return
		}
		remaining--
	}
	m.Select(index)
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
	if m.selected < 0 {
		m.off = 0
		return
	}
	if m.selected < m.off {
		m.off = m.selected
	}
	if m.selected >= m.off+area {
		m.off = m.selected - area + 1
	}
	m.off = max(0, min(m.off, max(0, len(m.items)-area)))
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
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
		m.Select(m.firstEnabled())
	case ActionLast:
		m.Select(m.lastEnabled())
	case ActionActivate:
		return m, m.Activate(), true
	default:
		id, ok := action.ParseSelectID(event.ID)
		if !ok {
			return m, nil, false
		}
		for i, item := range m.items {
			if item.ID == id && !item.Disabled {
				m.Select(i)
				return m, nil, true
			}
		}
		return m, nil, false
	}
	return m, nil, true
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
