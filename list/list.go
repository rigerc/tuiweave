// Package list provides a scrolling list of single-line items with a
// selection cursor and optional filtering.
//
// Lists use the Huh selection grammar: a thick focus rail, a "> " selector,
// and accent text for the selected row. Title, Description, and Help add
// optional field decoration above and below the rows.
//
// Filtering is display-only state: SetFilter narrows what is shown and
// navigated, but Selected always reports the index into the original items,
// so app logic never deals with filtered indices.
package list

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/internal/fieldstyle"
)

// Model is a list component. Create one with New.
type Model struct {
	width, height int
	items         []string
	filter        string
	matches       []int // original indices matching filter; nil when no filter
	pos           int   // cursor position in display space
	off           int   // first visible display position
	focused       bool

	// Title is an optional field title rendered above the rows.
	Title string
	// Description is an optional muted line below the title.
	Description string
	// Help is an optional muted hint rendered below the rows.
	Help string

	field fieldstyle.Styles
}

const (
	ActionFocus    = "focus"
	ActionBlur     = "blur"
	ActionNext     = "next"
	ActionPrevious = "previous"
	ActionFirst    = "first"
	ActionLast     = "last"
)

// New returns an empty list styled from the theme's roles.
func New(theme tuiweave.Theme) Model {
	return Model{field: fieldstyle.NewOn(theme, theme.SurfaceRaised)}
}

// SetSize sets the box the list renders in.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
	m.scrollIntoView()
}

// SetItems replaces the items, reapplying any active filter and clamping the
// selection.
func (m *Model) SetItems(items ...string) {
	m.items = append([]string(nil), items...)
	m.applyFilter(m.Selected())
}

// Items returns the original, unfiltered items.
func (m Model) Items() []string { return append([]string(nil), m.items...) }

// Len returns the number of original items.
func (m Model) Len() int { return len(m.items) }

// NaturalHeight reports the rows needed to show every displayed item plus
// optional title, description, and help.
func (m Model) NaturalHeight() int {
	return fieldstyle.HeaderHeight(m.Title, m.Description) + boolInt(m.Help != "") + m.displayCount()
}

// SetFilter narrows the list to items containing query, case-insensitively.
// An empty query clears the filter. The current selection is kept when it
// still matches; otherwise the first match is selected.
func (m *Model) SetFilter(query string) {
	keep := m.Selected()
	m.filter = query
	m.applyFilter(keep)
}

// Filter returns the active filter query.
func (m Model) Filter() string { return m.filter }

// FilteredLen returns how many items are currently displayed.
func (m Model) FilteredLen() int { return m.displayCount() }

// applyFilter rebuilds the match set, trying to keep the item with original
// index keep selected.
func (m *Model) applyFilter(keep int) {
	if m.filter == "" {
		m.matches = nil
	} else {
		q := strings.ToLower(m.filter)
		m.matches = make([]int, 0, len(m.items))
		for i, item := range m.items {
			if strings.Contains(strings.ToLower(item), q) {
				m.matches = append(m.matches, i)
			}
		}
	}
	m.pos = 0
	for p := range m.displayCount() {
		if m.origIndex(p) == keep {
			m.pos = p
			break
		}
	}
	m.clampPos()
	m.scrollIntoView()
}

// displayCount is how many items are shown under the current filter.
func (m Model) displayCount() int {
	if m.matches == nil {
		return len(m.items)
	}
	return len(m.matches)
}

// origIndex maps a display position to an original item index.
func (m Model) origIndex(pos int) int {
	if m.matches == nil {
		return pos
	}
	return m.matches[pos]
}

// Select moves the selection to the item with original index i. Without a
// filter the index is clamped; with a filter active, non-matching indices
// leave the selection unchanged.
func (m *Model) Select(i int) {
	if m.matches == nil {
		m.pos = max(0, min(i, len(m.items)-1))
		m.scrollIntoView()
		return
	}
	for p, orig := range m.matches {
		if orig == i {
			m.pos = p
			m.scrollIntoView()
			return
		}
	}
}

func (m *Model) selectPos(p int) {
	m.pos = max(0, min(p, m.displayCount()-1))
	m.clampPos()
	m.scrollIntoView()
}

func (m *Model) clampPos() {
	m.pos = max(0, min(m.pos, max(0, m.displayCount()-1)))
}

// Selected returns the original index of the selected item, or -1 when
// nothing is displayed.
func (m Model) Selected() int {
	if m.displayCount() == 0 {
		return -1
	}
	return m.origIndex(m.pos)
}

// SelectedItem returns the selected item, or "" when nothing is displayed.
func (m Model) SelectedItem() string {
	if m.displayCount() == 0 {
		return ""
	}
	return m.items[m.origIndex(m.pos)]
}

// Focus makes the list respond to navigation keys and highlights the
// selection with the selection roles.
func (m *Model) Focus() { m.focused = true }

// Blur stops the list from responding to keys.
func (m *Model) Blur() { m.focused = false }

// Focused reports whether the list handles keys.
func (m Model) Focused() bool { return m.focused }

// TotalLines returns the displayed item count (scrollbar.Scrollable).
func (m Model) TotalLines() int { return m.displayCount() }

// VisibleLines returns how many item rows are shown at once
// (scrollbar.Scrollable).
func (m Model) VisibleLines() int { return m.itemArea() }

// YOffset returns the first visible display position (scrollbar.Scrollable).
func (m Model) YOffset() int { return m.off }

func (m Model) itemArea() int {
	area := m.height - fieldstyle.HeaderHeight(m.Title, m.Description)
	if m.Help != "" {
		area--
	}
	return max(0, area)
}

func (m *Model) scrollIntoView() {
	area := m.itemArea()
	if area <= 0 {
		return
	}
	if m.pos < m.off {
		m.off = m.pos
	}
	if m.pos >= m.off+area {
		m.off = m.pos - area + 1
	}
	m.off = max(0, min(m.off, max(0, m.displayCount()-area)))
}

// Update handles navigation while focused: up/k, down/j, pgup, pgdown,
// g/home, G/end — all within the displayed (filtered) items.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if next, handled := m.applyAction(msg); handled {
		return next, nil
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
		m.selectPos(m.pos - 1)
	case "down", "j":
		m.selectPos(m.pos + 1)
	case "pgup":
		m.selectPos(m.pos - max(1, m.itemArea()))
	case "pgdown":
		m.selectPos(m.pos + max(1, m.itemArea()))
	case "g", "home":
		m.selectPos(0)
	case "G", "end":
		m.selectPos(m.displayCount() - 1)
	}
	return m, nil
}

// View renders the optional title/description, the visible item window, and
// an optional hint row. Selection is shown with the Huh selector and accent
// text; the left rail marks focus. The box is always filled exactly.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	rows := make([]string, 0, m.height)
	rows = append(rows, m.field.HeaderRows(m.focused, m.Title, m.Description, m.width)...)
	for i := 0; i < m.itemArea(); i++ {
		p := m.off + i
		if p >= m.displayCount() {
			rows = append(rows, m.field.FrameRow(m.focused, "", m.width))
			continue
		}
		rows = append(rows, m.renderItem(p))
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

func (m Model) renderItem(pos int) string {
	bodyWidth := fieldstyle.BodyWidth(m.width)
	if bodyWidth <= 0 {
		return m.field.FrameRow(m.focused, "", m.width)
	}

	selected := pos == m.pos
	selector := m.field.Selector(selected)
	textWidth := max(0, bodyWidth-lipgloss.Width(selector))
	text := ansi.Truncate(m.items[m.origIndex(pos)], textWidth, "…")
	body := selector + m.field.Label(text, selected, false)
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
