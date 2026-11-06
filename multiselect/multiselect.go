// Package multiselect provides a focusable list of application-owned options
// where any number of entries can be checked. The component owns the cursor,
// checked state, and rendering; it emits ChangedMsg when a box flips and
// SubmittedMsg when the user confirms the field. Applications own validation,
// persistence, and what the selection means.
//
// Multiselects follow the Huh field grammar: Title and Description render above
// the options, Help below, a thick focus rail marks the active field, and the
// cursor row uses a "> " selector. Inline renders a Clack prompt block with a
// ◆ header, ❯ ◼/◻ options, instructions, and a └ footer; an unfocused prompt
// collapses to the checked labels.
package multiselect

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
	// ActionFocus focuses the multiselect.
	ActionFocus = "focus"
	// ActionBlur blurs the multiselect.
	ActionBlur = "blur"
	// ActionNext moves the cursor to the next enabled item.
	ActionNext = "next"
	// ActionPrevious moves the cursor to the previous enabled item.
	ActionPrevious = "previous"
	// ActionFirst moves the cursor to the first enabled item.
	ActionFirst = "first"
	// ActionLast moves the cursor to the last enabled item.
	ActionLast = "last"
	// ActionToggle flips the checked state of the cursor item.
	ActionToggle = "toggle"
	// ActionCheck checks the cursor item.
	ActionCheck = "check"
	// ActionUncheck unchecks the cursor item.
	ActionUncheck = "uncheck"
	// ActionClear unchecks every item.
	ActionClear = "clear"
	// ActionSubmit emits SubmittedMsg for the current checked items.
	ActionSubmit = "submit"
	// ActionSelectPrefix prefixes stable cursor-selection actions. For example,
	// an item with ID "tests" exposes the local action "select.tests".
	ActionSelectPrefix = action.SelectPrefix
	// ActionTogglePrefix prefixes stable checkbox actions. For example, an
	// item with ID "tests" exposes the local action "toggle.tests".
	ActionTogglePrefix = "toggle."
)

// Item is the shared action.Item definition rendered by the multiselect.
type Item = action.Item

// ChangedMsg is emitted when a checkbox flips through keyboard or semantic
// action input. Selected contains the stable IDs of every checked item in item
// order, so an application can route the whole field state from one message.
type ChangedMsg struct {
	ID       string
	Index    int
	Label    string
	Checked  bool
	Selected []string
}

// SubmittedMsg is emitted when the user confirms the field with enter or the
// ActionSubmit semantic action. Selected contains the stable IDs of every
// checked item in item order.
type SubmittedMsg struct {
	Selected []string
}

// Model is a focusable multiselect. Create one with New.
type Model struct {
	width, height int
	items         []Item
	checked       []bool
	cursor        int
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
	// options; an unfocused prompt collapses to the checked labels.
	Inline bool

	field fieldstyle.Styles
	clack clackstyle.Styles
}

// New returns an empty multiselect styled from the theme's roles.
func New(theme tuiweave.Theme) Model { return NewOn(theme, theme.SurfaceRaised) }

// NewOn returns an empty multiselect with the given background painted behind
// its rows. Use nil for a transparent prompt on the application surface.
func NewOn(theme tuiweave.Theme, background color.Color) Model {
	return Model{
		cursor: -1,
		field:  fieldstyle.NewOn(theme, background),
		clack:  clackstyle.NewOn(theme, background),
	}
}

// SetSize sets the box the multiselect renders in.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
	m.clampCursor()
	m.scrollIntoView()
}

// SizeMode reports a bounded scrolling field for stacked options and a
// width-bounded, natural-height Clack prompt block for Inline.
func (m Model) SizeMode() layout.SizeMode {
	if m.Inline {
		return layout.SizeWidthBounded
	}
	return layout.SizeBounded
}

// SetItems replaces the items, preserving the checked IDs and cursor ID that
// are still present and enabled. New items start unchecked.
func (m *Model) SetItems(items ...Item) {
	checkedIDs := make(map[string]bool, len(m.checked))
	for i, item := range m.items {
		if i < len(m.checked) && m.checked[i] && item.ID != "" {
			checkedIDs[item.ID] = true
		}
	}
	cursorID := ""
	if item := m.CursorItem(); item.ID != "" {
		cursorID = item.ID
	}

	m.items = append([]Item(nil), items...)
	m.checked = make([]bool, len(m.items))
	for i, item := range m.items {
		m.checked[i] = !item.Disabled && item.ID != "" && checkedIDs[item.ID]
	}
	m.cursor = -1
	if cursorID != "" {
		for i, item := range m.items {
			if item.ID == cursorID && !item.Disabled {
				m.cursor = i
				break
			}
		}
	}
	m.clampCursor()
	m.scrollIntoView()
}

// Items returns the configured items.
func (m Model) Items() []Item { return append([]Item(nil), m.items...) }

// Len returns the number of items.
func (m Model) Len() int { return len(m.items) }

// SetChecked sets an enabled item's checked state by stable ID without
// emitting a message. It is a no-op when the ID is unknown or disabled.
func (m *Model) SetChecked(id string, checked bool) {
	if id == "" {
		return
	}
	for i, item := range m.items {
		if item.ID == id && !item.Disabled {
			m.checked[i] = checked
			return
		}
	}
}

// SetCheckedIDs replaces the checked set with the given stable IDs. Unknown
// and disabled IDs are ignored.
func (m *Model) SetCheckedIDs(ids ...string) {
	want := make(map[string]bool, len(ids))
	for _, id := range ids {
		want[id] = true
	}
	for i, item := range m.items {
		m.checked[i] = !item.Disabled && want[item.ID]
	}
}

// Clear unchecks every item without emitting a message.
func (m *Model) Clear() {
	for i := range m.checked {
		m.checked[i] = false
	}
}

// Checked reports whether the item with the given stable ID is checked.
func (m Model) Checked(id string) bool {
	if id == "" {
		return false
	}
	for i, item := range m.items {
		if item.ID == id {
			return i < len(m.checked) && m.checked[i]
		}
	}
	return false
}

// SelectedIDs returns the stable IDs of every checked item in item order.
func (m Model) SelectedIDs() []string {
	ids := make([]string, 0, len(m.items))
	for i, item := range m.items {
		if i < len(m.checked) && m.checked[i] && item.ID != "" {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

// SelectedItems returns every checked item in item order.
func (m Model) SelectedItems() []Item {
	items := make([]Item, 0, len(m.items))
	for i, item := range m.items {
		if i < len(m.checked) && m.checked[i] {
			items = append(items, item)
		}
	}
	return items
}

// SelectedCount reports how many items are checked.
func (m Model) SelectedCount() int {
	count := 0
	for _, checked := range m.checked {
		if checked {
			count++
		}
	}
	return count
}

// Select moves the cursor to an enabled item by index. Disabled or
// out-of-range indices leave the cursor unchanged.
func (m *Model) Select(index int) {
	if index < 0 || index >= len(m.items) || m.items[index].Disabled {
		return
	}
	m.cursor = index
	m.scrollIntoView()
}

// SelectID moves the cursor to an enabled item by stable ID.
func (m *Model) SelectID(id string) {
	if id == "" {
		return
	}
	for i, item := range m.items {
		if item.ID == id && !item.Disabled {
			m.Select(i)
			return
		}
	}
}

// Selected returns the cursor index, or -1 when no enabled item exists.
func (m Model) Selected() int { return m.cursor }

// CursorItem returns the item under the cursor, or a zero Item when there is
// no cursor.
func (m Model) CursorItem() Item {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return Item{}
	}
	return m.items[m.cursor]
}

// Focus makes the multiselect respond to navigation, toggle, and submit keys.
func (m *Model) Focus() { m.focused = true }

// Blur stops the multiselect from responding to keys.
func (m *Model) Blur() { m.focused = false }

// Focused reports whether the field handles keys.
func (m Model) Focused() bool { return m.focused }

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

// ToggleCursor flips the cursor item and returns a command that emits
// ChangedMsg. It returns nil when the cursor has no enabled item.
func (m Model) ToggleCursor() tea.Cmd {
	index := m.cursor
	if index < 0 || index >= len(m.items) || m.items[index].Disabled {
		return nil
	}
	return m.setChecked(index, !m.checked[index])
}

// Submit returns a command that emits SubmittedMsg for the checked items.
func (m Model) Submit() tea.Cmd {
	selected := m.SelectedIDs()
	return func() tea.Msg { return SubmittedMsg{Selected: selected} }
}

// Update handles semantic actions and, while focused, up/down, j/k, left/right
// when Inline, page navigation, home/g, end/G, space/x to toggle the cursor
// item, and enter to submit.
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
	case "space", "x":
		return m, m.ToggleCursor()
	case "enter":
		return m, m.Submit()
	}
	return m, nil
}

// View renders the Huh-style option list or, when Inline is set, a Clack-style
// prompt block with a ◆ header, │ rows, ❯ ◼/◻ options, and a └ footer.
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
			rows = append(rows, m.clack.Row(m.focused, m.clack.Inactive.Render("No options"), m.width))
		} else {
			for i, item := range m.items {
				checked := i < len(m.checked) && m.checked[i]
				rows = append(rows, m.clack.MultiOptionRow(m.focused, i == m.cursor, checked, item.Disabled, item.Label, m.width))
			}
		}
		rows = append(rows, m.clack.InstructionsRow(m.focused, "↑/↓ to navigate • Space: select • Enter: submit", m.width))
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

	cursor := index == m.cursor
	checked := index < len(m.checked) && m.checked[index]
	selector := m.field.Selector(cursor)
	box := m.box(item.Disabled, checked)
	selectorWidth := lipgloss.Width(selector)
	boxWidth := lipgloss.Width(box)

	var body string
	switch {
	case bodyWidth >= selectorWidth+boxWidth+1:
		labelWidth := bodyWidth - selectorWidth - boxWidth - 1
		label := ansi.Truncate(item.Label, labelWidth, "…")
		body = selector + box + m.field.Pad.Render(" ") + m.field.Label(label, cursor, item.Disabled)
	case bodyWidth >= selectorWidth+boxWidth:
		body = selector + box
	default:
		body = box
	}
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func (m Model) box(disabled, checked bool) string {
	switch {
	case disabled:
		return m.field.Disabled.Render("[ ]")
	case checked:
		return m.field.SelectMark.Render("[x]")
	default:
		return m.field.Normal.Render("[ ]")
	}
}

func (m Model) summary() string {
	items := m.SelectedItems()
	if len(items) == 0 {
		return "None"
	}
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	return strings.Join(labels, ", ")
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

func (m *Model) clampCursor() {
	if len(m.items) == 0 {
		m.cursor = -1
		return
	}
	if m.cursor >= 0 && m.cursor < len(m.items) && !m.items[m.cursor].Disabled {
		return
	}
	m.cursor = m.firstEnabled()
}

func (m *Model) selectRelative(delta int) {
	if len(m.items) == 0 {
		m.clampCursor()
		return
	}
	if m.cursor < 0 {
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
	index := m.cursor
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

func (m Model) setChecked(index int, checked bool) tea.Cmd {
	if index < 0 || index >= len(m.items) || m.items[index].Disabled {
		return nil
	}
	if m.checked[index] == checked {
		return nil
	}
	m.checked[index] = checked
	item := m.items[index]
	selected := m.SelectedIDs()
	return func() tea.Msg {
		return ChangedMsg{
			ID:       item.ID,
			Index:    index,
			Label:    item.Label,
			Checked:  checked,
			Selected: selected,
		}
	}
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
	case ActionToggle:
		return m, m.ToggleCursor(), true
	case ActionCheck:
		return m, m.setChecked(m.cursor, true), true
	case ActionUncheck:
		return m, m.setChecked(m.cursor, false), true
	case ActionClear:
		m.Clear()
	case ActionSubmit:
		return m, m.Submit(), true
	default:
		if id, ok := parseToggleID(event.ID); ok {
			for i, item := range m.items {
				if item.ID == id {
					return m, m.setChecked(i, !m.checked[i]), true
				}
			}
			return m, nil, false
		}
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

func parseToggleID(id string) (string, bool) {
	if !strings.HasPrefix(id, ActionTogglePrefix) {
		return "", false
	}
	itemID := strings.TrimPrefix(id, ActionTogglePrefix)
	return itemID, itemID != ""
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
