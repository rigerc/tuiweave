// Package palette provides a bounded command-palette foundation: a query
// input, filtered action window, stable IDs, and semantic activation. The
// application owns visibility and the meaning of each action.
//
// Palettes use the Huh field grammar: the query input supplies the focus rail,
// results use a "> " selector and accent text, and descriptions stay muted.
package palette

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/action"
	"github.com/rigerc/tuiweave/inspect"
	"github.com/rigerc/tuiweave/internal/fieldstyle"
	"github.com/rigerc/tuiweave/textinput"
)

const (
	// ActionFocus focuses the palette query and action window.
	ActionFocus = "focus"
	// ActionBlur blurs the palette.
	ActionBlur = "blur"
	// ActionNext selects the next enabled matching action.
	ActionNext = "next"
	// ActionPrevious selects the previous enabled matching action.
	ActionPrevious = "previous"
	// ActionFirst selects the first enabled matching action.
	ActionFirst = "first"
	// ActionLast selects the last enabled matching action.
	ActionLast = "last"
	// ActionClear clears the query.
	ActionClear = "clear"
	// ActionActivate activates the selected action.
	ActionActivate = "activate"
	// ActionSelectPrefix prefixes stable action-selection IDs.
	ActionSelectPrefix = action.SelectPrefix
)

// Item is the shared action.Item definition rendered by the palette. The
// description is included in its search window.
type Item = action.Item

// SelectedMsg is emitted when an action is activated with enter, space, or a
// semantic activation action.
type SelectedMsg struct {
	ID    string
	Label string
}

// Model is a bounded, focusable command palette. Create one with New.
type Model struct {
	width, height int
	items         []Item
	matches       []int // original item indices in the current query window
	pos           int   // selected position in matches; -1 when none enabled
	off           int   // first visible position in matches
	focused       bool
	input         textinput.Model

	field fieldstyle.Styles
}

// New returns an empty command palette styled from the theme's semantic roles.
func New(theme tuiweave.Theme) Model {
	input := textinput.NewOn(theme, theme.SurfaceRaised)
	input.Prompt = "› "
	return Model{
		pos:   -1,
		input: input,
		field: fieldstyle.NewOn(theme, theme.SurfaceRaised),
	}
}

// SetSize sets the palette's bounded box. The first row is the query; the
// remaining rows are the filtered action window.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
	m.input.SetSize(width, 1)
	m.scrollIntoView()
}

// SetItems replaces the actions, preserving the selected stable ID when it
// remains enabled and visible under the current query.
func (m *Model) SetItems(items ...Item) {
	keepIndex, keepID := m.Selected(), m.SelectedID()
	m.items = append([]Item(nil), items...)
	m.applyFilter(keepIndex, keepID)
}

// Items returns the original, unfiltered actions.
func (m Model) Items() []Item { return append([]Item(nil), m.items...) }

// SetQuery changes the search query and keeps the current action selected when
// possible. Matching is case-insensitive over ID, label, and description.
func (m *Model) SetQuery(query string) {
	keepIndex, keepID := m.Selected(), m.SelectedID()
	m.input.SetValue(query)
	m.applyFilter(keepIndex, keepID)
}

// Query returns the current search query.
func (m Model) Query() string { return m.input.Value() }

// FilteredLen returns the number of actions matching the current query.
func (m Model) FilteredLen() int { return len(m.matches) }

// Select selects an enabled action by original item index when it matches the
// current query. Non-matching, disabled, and out-of-range indices are ignored.
func (m *Model) Select(index int) {
	for p, original := range m.matches {
		if original == index && !m.items[original].Disabled {
			m.pos = p
			m.scrollIntoView()
			return
		}
	}
}

// SelectID selects an enabled action by stable ID. If the action is hidden by
// the current query, the query is cleared first so semantic selection remains
// deterministic.
func (m *Model) SelectID(id string) {
	original := -1
	for i, item := range m.items {
		if item.ID == id && !item.Disabled {
			original = i
			break
		}
	}
	if original < 0 {
		return
	}
	for _, index := range m.matches {
		if index == original {
			m.Select(original)
			return
		}
	}
	m.SetQuery("")
	m.Select(original)
}

// Selected returns the selected action's original item index, or -1 when no
// enabled matching action is selected.
func (m Model) Selected() int {
	if m.pos < 0 || m.pos >= len(m.matches) {
		return -1
	}
	original := m.matches[m.pos]
	if m.items[original].Disabled {
		return -1
	}
	return original
}

// SelectedItem returns the selected action, or a zero Item when none is
// selected.
func (m Model) SelectedItem() Item {
	selected := m.Selected()
	if selected < 0 {
		return Item{}
	}
	return m.items[selected]
}

// SelectedID returns the selected action's stable ID, or an empty string.
func (m Model) SelectedID() string { return m.SelectedItem().ID }

// Focus makes the palette accept query, navigation, and activation keys.
func (m *Model) Focus() {
	m.focused = true
	m.input.Focus()
}

// Blur stops the palette from accepting keys and hides the query cursor.
func (m *Model) Blur() {
	m.focused = false
	m.input.Blur()
}

// Focused reports whether the palette accepts keyboard input.
func (m Model) Focused() bool { return m.focused }

// TotalLines returns the number of matching actions for scrollbar.Scrollable.
func (m Model) TotalLines() int { return len(m.matches) }

// VisibleLines returns the number of result rows below the query.
func (m Model) VisibleLines() int { return max(0, m.height-1) }

// YOffset returns the first visible matching-action position.
func (m Model) YOffset() int { return m.off }

// Activate returns a command for the selected action, or nil when there is no
// enabled selection.
func (m Model) Activate() tea.Cmd {
	item := m.SelectedItem()
	if item.ID == "" && item.Label == "" {
		return nil
	}
	return func() tea.Msg { return SelectedMsg{ID: item.ID, Label: item.Label} }
}

// Update handles semantic actions and, while focused, query editing,
// navigation, and enter/space activation. The application still owns palette
// visibility and the command's side effects.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if next, cmd, handled := m.applyAction(msg); handled {
		return next, cmd
	}
	if !m.focused {
		return m, nil
	}
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "up", "k":
			m.selectRelative(-1)
			return m, nil
		case "down", "j":
			m.selectRelative(1)
			return m, nil
		case "pgup":
			m.selectRelative(-max(1, m.VisibleLines()))
			return m, nil
		case "pgdown":
			m.selectRelative(max(1, m.VisibleLines()))
			return m, nil
		case "home", "g":
			m.selectFirst()
			return m, nil
		case "end", "G":
			m.selectLast()
			return m, nil
		case "enter", "space":
			return m, m.Activate()
		}
	}
	keepIndex, keepID := m.Selected(), m.SelectedID()
	before := m.Query()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.Query() != before {
		m.applyFilter(keepIndex, keepID)
	}
	return m, cmd
}

// View renders one exact-width query row followed by the visible result rows.
// The query input supplies the focus rail; result rows use the same rail and a
// Huh-style selector. A no-match state is explicit so a palette does not look
// broken while a user is searching.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	rows := []string{m.input.View()}
	visible := m.VisibleLines()
	if visible == 0 {
		return strings.Join(rows, "\n")
	}
	if len(m.matches) == 0 {
		empty := m.field.FrameRow(m.focused, m.field.Label("  No matching commands", false, false), m.width)
		rows = append(rows, empty)
		for len(rows) < m.height {
			rows = append(rows, m.field.FrameRow(m.focused, "", m.width))
		}
		return strings.Join(rows, "\n")
	}
	for i := 0; i < visible; i++ {
		position := m.off + i
		if position >= len(m.matches) {
			rows = append(rows, m.field.FrameRow(m.focused, "", m.width))
			continue
		}
		rows = append(rows, m.renderItem(position))
	}
	return strings.Join(rows, "\n")
}

func (m Model) renderItem(position int) string {
	item := m.items[m.matches[position]]
	bodyWidth := fieldstyle.BodyWidth(m.width)
	if bodyWidth <= 0 {
		return m.field.FrameRow(m.focused, "", m.width)
	}

	selected := position == m.pos && !item.Disabled
	selector := m.field.Selector(selected)
	available := max(0, bodyWidth-lipgloss.Width(selector))
	label := ansi.Truncate(item.Label, available, "…")
	description := ""
	gap := 0
	if item.Description != "" && available >= 4 {
		label = ansi.Truncate(item.Label, available-3, "…")
		gap = 2
		descriptionWidth := available - ansi.StringWidth(label) - gap
		if descriptionWidth > 0 {
			description = ansi.Truncate(item.Description, descriptionWidth, "…")
		} else {
			gap = 0
		}
	}
	body := selector + m.field.Label(label, selected, item.Disabled)
	if gap > 0 {
		descStyle := m.field.Desc
		if item.Disabled {
			descStyle = m.field.Disabled
		}
		body += m.field.Pad.Render(strings.Repeat(" ", gap)) + descStyle.Render(description)
	}
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func (m *Model) selectFirst() {
	m.pos = m.firstEnabledPosition()
	m.scrollIntoView()
}

func (m *Model) selectLast() {
	m.pos = m.lastEnabledPosition()
	m.scrollIntoView()
}

func (m *Model) selectRelative(delta int) {
	if len(m.matches) == 0 {
		m.pos = -1
		return
	}
	if m.Selected() < 0 {
		if delta < 0 {
			m.selectLast()
		} else {
			m.selectFirst()
		}
		return
	}
	step := 1
	if delta < 0 {
		step = -1
	}
	remaining := abs(delta)
	position := m.pos
	for remaining > 0 {
		position += step
		for position >= 0 && position < len(m.matches) && m.items[m.matches[position]].Disabled {
			position += step
		}
		if position < 0 || position >= len(m.matches) {
			return
		}
		remaining--
	}
	m.pos = position
	m.scrollIntoView()
}

func (m Model) firstEnabledPosition() int {
	for position, original := range m.matches {
		if !m.items[original].Disabled {
			return position
		}
	}
	return -1
}

func (m Model) lastEnabledPosition() int {
	for position := len(m.matches) - 1; position >= 0; position-- {
		if !m.items[m.matches[position]].Disabled {
			return position
		}
	}
	return -1
}

func (m *Model) scrollIntoView() {
	visible := m.VisibleLines()
	if visible <= 0 || m.pos < 0 {
		m.off = 0
		return
	}
	if m.pos < m.off {
		m.off = m.pos
	}
	if m.pos >= m.off+visible {
		m.off = m.pos - visible + 1
	}
	m.off = max(0, min(m.off, max(0, len(m.matches)-visible)))
}

func (m *Model) applyFilter(keepIndex int, keepID string) {
	query := strings.ToLower(m.Query())
	m.matches = m.matches[:0]
	for i, item := range m.items {
		searchable := strings.ToLower(item.ID + " " + item.Label + " " + item.Description)
		if query == "" || strings.Contains(searchable, query) {
			m.matches = append(m.matches, i)
		}
	}

	m.pos = -1
	if keepID != "" {
		for position, original := range m.matches {
			if m.items[original].ID == keepID && !m.items[original].Disabled {
				m.pos = position
				break
			}
		}
	}
	if m.pos < 0 && keepIndex >= 0 {
		for position, original := range m.matches {
			if original == keepIndex && !m.items[original].Disabled {
				m.pos = position
				break
			}
		}
	}
	if m.pos < 0 {
		m.pos = m.firstEnabledPosition()
	}
	if m.pos < 0 && len(m.matches) > 0 {
		m.pos = 0
	}
	m.scrollIntoView()
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
		m.selectFirst()
	case ActionLast:
		m.selectLast()
	case ActionClear:
		m.SetQuery("")
	case ActionActivate:
		return m, m.Activate(), true
	default:
		id, ok := action.ParseSelectID(event.ID)
		if !ok {
			return m, nil, false
		}
		for _, item := range m.items {
			if item.ID == id && !item.Disabled {
				m.SelectID(id)
				return m, nil, true
			}
		}
		return m, nil, false
	}
	return m, nil, true
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
