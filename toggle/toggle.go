// Package toggle provides a compact, focusable on/off control. It owns only
// the visual state and emits a message when that state changes; applications
// own the setting's meaning and persistence.
//
// Toggles can render as Huh-style stacked fields: Title and Description appear
// above the checkbox, Help appears below, and a thick focus rail marks the
// active field. Inline renders a Clack prompt block with a ◆ header, ❯ ●/○
// control, and └ footer; an unfocused prompt collapses to its submitted label.
// Checked state uses the accent role, not the success role.
package toggle

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/inspect"
	"github.com/rigerc/tuiweave/internal/clackstyle"
	"github.com/rigerc/tuiweave/internal/fieldstyle"
	"github.com/rigerc/tuiweave/layout"
)

const (
	// ActionFocus focuses the toggle.
	ActionFocus = "focus"
	// ActionBlur blurs the toggle.
	ActionBlur = "blur"
	// ActionToggle flips the current value.
	ActionToggle = "toggle"
	// ActionOn turns the value on.
	ActionOn = "on"
	// ActionOff turns the value off.
	ActionOff = "off"
)

// ChangedMsg is emitted when the toggle changes through keyboard or semantic
// action input. ID lets an application route one message type from many
// toggles.
type ChangedMsg struct {
	ID      string
	Checked bool
}

// Model is a focusable toggle with an optional Huh-style field header.
// Create one with New.
type Model struct {
	width, height int
	label         string
	checked       bool
	disabled      bool
	focused       bool

	// ID identifies this toggle in ChangedMsg and application inspection trees.
	ID string

	// Title is an optional field title. It is rendered above the checkbox, or
	// beside it when Inline is true.
	Title string
	// Description is an optional muted line below the title.
	Description string
	// Help is an optional muted hint rendered below the control.
	Help string
	// Inline renders a Clack prompt block. The focused prompt expands; an
	// unfocused prompt collapses to its submitted label.
	Inline bool

	field          fieldstyle.Styles
	clack          clackstyle.Styles
	labelStyle     lipgloss.Style
	disabledLabel  lipgloss.Style
	checkedStyle   lipgloss.Style
	uncheckedStyle lipgloss.Style
	disabledBox    lipgloss.Style
}

// New returns an unchecked toggle styled from the theme's semantic roles.
func New(theme tuiweave.Theme) Model { return NewOn(theme, theme.SurfaceRaised) }

// NewOn returns an unchecked toggle with the given background painted behind
// its rows. Use nil to sit directly on the application surface.
func NewOn(theme tuiweave.Theme, background color.Color) Model {
	base := lipgloss.NewStyle()
	if background != nil {
		base = base.Background(background)
	}
	return Model{
		field:          fieldstyle.NewOn(theme, background),
		clack:          clackstyle.NewOn(theme, background),
		labelStyle:     base.Foreground(theme.Text),
		disabledLabel:  base.Foreground(theme.TextFaint),
		checkedStyle:   lipgloss.NewStyle().Foreground(theme.TextInverted).Background(theme.Accent).Bold(true),
		uncheckedStyle: lipgloss.NewStyle().Foreground(theme.TextMuted).Background(theme.SurfaceSunken).Bold(true),
		disabledBox:    lipgloss.NewStyle().Foreground(theme.TextFaint).Background(theme.SurfaceSunken),
	}
}

// SetSize sets the box the toggle renders in. Stacked mode renders one row plus
// optional field decoration; Inline renders a Clack prompt block whose height
// depends on focus.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
}

// SizeMode reports that the toggle is width-constrained with a natural height
// that grows when field decoration is set.
func (m Model) SizeMode() layout.SizeMode { return layout.SizeWidthBounded }

// SetLabel sets the setting label shown beside the checkbox.
func (m *Model) SetLabel(label string) { m.label = label }

// Label returns the setting label.
func (m Model) Label() string { return m.label }

// NaturalHeight reports the rows View renders at the current settings. Inline
// height depends on focus: two rows when collapsed, the full prompt when
// focused.
func (m Model) NaturalHeight() int {
	if m.Inline {
		if !m.focused {
			return 2
		}
		height := 3 // header, control, footer
		if m.Description != "" {
			height++
		}
		if m.Help != "" {
			height++
		}
		return height
	}
	return fieldstyle.HeaderHeight(m.Title, m.Description) + 1 + boolInt(m.Help != "")
}

// SetChecked sets the on/off value without emitting a message.
func (m *Model) SetChecked(checked bool) { m.checked = checked }

// Checked reports whether the toggle is on.
func (m Model) Checked() bool { return m.checked }

// SetDisabled controls whether the toggle can be focused or changed.
func (m *Model) SetDisabled(disabled bool) {
	m.disabled = disabled
	if disabled {
		m.focused = false
	}
}

// Disabled reports whether the toggle is unavailable.
func (m Model) Disabled() bool { return m.disabled }

// Focus makes the toggle respond to space, x, and enter.
func (m *Model) Focus() {
	if !m.disabled {
		m.focused = true
	}
}

// Blur stops the toggle from responding to keys.
func (m *Model) Blur() { m.focused = false }

// Focused reports whether the toggle accepts keyboard input.
func (m Model) Focused() bool { return m.focused }

// Toggle flips the value without emitting a message. Use Update when the app
// needs a ChangedMsg command for the transition.
func (m *Model) Toggle() {
	if !m.disabled {
		m.checked = !m.checked
	}
}

// Update handles semantic actions and, while focused, space, x, and enter.
// State changes return a command that emits ChangedMsg.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if next, cmd, handled := m.applyAction(msg); handled {
		return next, cmd
	}
	if !m.focused || m.disabled {
		return m, nil
	}
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "space", "x", "enter":
		return m, m.setChecked(!m.checked)
	}
	return m, nil
}

// View renders the optional title/description, the checkbox row, and an
// optional hint row. The focus rail is the only focus indicator; checked state
// uses the accent role rather than the success role.
//
// The stacked view is width-bounded with a natural height: one row when no
// decoration is set, plus one row each for Title, Description, and Help.
// Inline renders a Clack prompt block instead.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	box := "[ ]"
	if m.checked {
		box = "[x]"
	}
	boxStyle, labelStyle := m.styles()
	var rows []string
	if m.Inline {
		rows = append(rows, m.clack.HeaderRow(m.focused, m.Title, m.width))
		if !m.focused {
			style := m.clack.Inactive
			mark := "○"
			if m.disabled {
				style = m.clack.Disabled
			} else if m.checked {
				mark = "●"
			}
			rows = append(rows, m.clack.Row(false, style.Render(mark+" "+m.label), m.width))
			if len(rows) > m.height {
				rows = rows[:m.height]
			}
			return strings.Join(rows, "\n")
		}
		if m.Description != "" {
			rows = append(rows, m.clack.Row(m.focused, m.clack.Desc.Render(m.Description), m.width))
		}
		rows = append(rows, m.clack.Row(m.focused, m.clackControl(labelStyle), m.width))
		if m.Help != "" {
			rows = append(rows, m.clack.HelpRow(m.focused, m.Help, m.width))
		}
		rows = append(rows, m.clack.FooterRow(m.width))
	} else {
		rows = append(rows, m.field.HeaderRows(m.focused, m.Title, m.Description, m.width)...)
		control := m.controlRow(box, boxStyle, labelStyle)
		if control != "" {
			rows = append(rows, control)
		}
		if help := m.field.HelpRow(m.focused, m.Help, m.width); help != "" {
			rows = append(rows, help)
		}
	}
	if len(rows) > m.height {
		rows = rows[:m.height]
	}
	return strings.Join(rows, "\n")
}

func (m Model) clackControl(labelStyle lipgloss.Style) string {
	cursor := m.clack.Pad.Render("  ")
	if m.focused {
		cursor = m.clack.Cursor.Render("❯ ")
	}
	mark := m.clack.Inactive.Render("○")
	if m.disabled {
		mark = m.clack.Disabled.Render("○")
	} else if m.checked {
		mark = m.clack.Selected.Render("●")
	}
	return cursor + mark + m.clack.Pad.Render(" ") + labelStyle.Render(m.label)
}

func (m Model) controlRow(box string, boxStyle, labelStyle lipgloss.Style) string {
	bodyWidth := fieldstyle.BodyWidth(m.width)
	if bodyWidth <= 0 {
		return m.field.FrameRow(m.focused, "", m.width)
	}
	if bodyWidth < 4 {
		return m.field.FrameRow(m.focused, boxStyle.Render(ansi.Truncate(box, bodyWidth, "…")), m.width)
	}

	labelWidth := bodyWidth - 4
	label := ansi.Truncate(m.label, labelWidth, "…")
	body := boxStyle.Render(box) + labelStyle.Render(" "+label)
	if pad := bodyWidth - lipgloss.Width(body); pad > 0 {
		body += m.field.Pad.Render(strings.Repeat(" ", pad))
	}
	return m.field.FrameRow(m.focused, body, m.width)
}

func (m Model) styles() (lipgloss.Style, lipgloss.Style) {
	if m.disabled {
		return m.disabledBox, m.disabledLabel
	}
	return m.checkedStyleOrUnchecked(), m.labelStyle
}

func (m Model) checkedStyleOrUnchecked() lipgloss.Style {
	if m.checked {
		return m.checkedStyle
	}
	return m.uncheckedStyle
}

func (m *Model) setChecked(checked bool) tea.Cmd {
	if m.disabled || m.checked == checked {
		return nil
	}
	m.checked = checked
	id := m.ID
	return func() tea.Msg { return ChangedMsg{ID: id, Checked: checked} }
}

func (m Model) applyAction(msg tea.Msg) (Model, tea.Cmd, bool) {
	action, ok := msg.(inspect.ActionMsg)
	if !ok {
		return m, nil, false
	}
	switch action.ID {
	case ActionFocus:
		m.Focus()
	case ActionBlur:
		m.Blur()
	case ActionToggle:
		return m, m.setChecked(!m.checked), true
	case ActionOn:
		return m, m.setChecked(true), true
	case ActionOff:
		return m, m.setChecked(false), true
	default:
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
