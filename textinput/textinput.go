// Package textinput provides a single-line text input with a prompt,
// placeholder, cursor, and grapheme-aware horizontal scrolling.
//
// Inputs can render as Huh-style stacked fields: set Title, Description, and
// Help for a title/description block above the input and a hint below it, with
// a thick focus rail marking the active field. Set Inline to render a Clack
// prompt block instead: a ◆ header, │ body rows, a › input, and a └ footer.
// An unfocused inline input collapses to its submitted value.
package textinput

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/internal/clackstyle"
	"github.com/rigerc/tuiweave/internal/fieldstyle"
	"github.com/rigerc/tuiweave/internal/paste"
	"github.com/rigerc/tuiweave/layout"
)

// Model is a text input component. Create one with New.
type Model struct {
	width, height int
	value         []rune
	pos           int // cursor position in runes
	focused       bool

	// Prompt is rendered before the editable area. Default "> ".
	Prompt string
	// Placeholder is shown, faint, while the value is empty.
	Placeholder string
	// Title is an optional field title. Inline uses it as the Clack prompt
	// message; stacked mode renders it above the input.
	Title string
	// Description is an optional muted line below the title.
	Description string
	// Help is an optional muted hint rendered below the input.
	Help string
	// Inline renders a Clack prompt block. The focused prompt expands; an
	// unfocused prompt collapses to its submitted value.
	Inline bool
	// Mask replaces each visible grapheme with a single mask glyph in the
	// rendered value, as in a password field. It affects display only: Value
	// returns the real text. Use a short, single-cell string such as "•";
	// an empty Mask renders the text normally.
	Mask string

	field            fieldstyle.Styles
	clack            clackstyle.Styles
	promptStyle      lipgloss.Style
	textStyle        lipgloss.Style
	placeholderStyle lipgloss.Style
	cursorStyle      lipgloss.Style
}

const (
	ActionFocus = "focus"
	ActionBlur  = "blur"
	ActionClear = "clear"
)

// New returns an empty text input styled from the theme's roles.
func New(theme tuiweave.Theme) Model { return NewOn(theme, nil) }

// NewOn returns an empty text input that paints background behind its field
// rail, prompt, text, and padding. Use theme.SurfaceRaised when embedding the
// input in a panel or popover so nested style resets cannot leave gaps.
func NewOn(theme tuiweave.Theme, background color.Color) Model {
	base := lipgloss.NewStyle()
	if background != nil {
		base = base.Background(background)
	}
	return Model{
		Prompt:           "> ",
		field:            fieldstyle.NewOn(theme, background),
		clack:            clackstyle.NewOn(theme, background),
		promptStyle:      base.Foreground(theme.Accent),
		textStyle:        base.Foreground(theme.Text),
		placeholderStyle: base.Foreground(theme.TextMuted),
		cursorStyle:      base.Foreground(theme.Text).Reverse(true),
	}
}

// SetSize sets the box the input renders in. The input is width-bounded with a
// natural height: one row normally, plus optional Title/Description rows and a
// Help row. Inline renders a Clack prompt block whose height depends on focus.
func (m *Model) SetSize(width, height int) {
	m.width, m.height = width, height
}

// SizeMode reports that the input is width-constrained with a natural height
// that grows when field decoration is set.
func (m Model) SizeMode() layout.SizeMode { return layout.SizeWidthBounded }

// Focus makes the input accept keys and show its cursor.
func (m *Model) Focus() { m.focused = true }

// Blur stops the input from accepting keys and hides the cursor.
func (m *Model) Blur() { m.focused = false }

// Focused reports whether the input accepts keys.
func (m Model) Focused() bool { return m.focused }

// Value returns the current text.
func (m Model) Value() string { return string(m.value) }

// SetValue replaces the text and moves the cursor to its end.
func (m *Model) SetValue(s string) {
	m.value = []rune(s)
	m.pos = len(m.value)
}

// Reset clears the input.
func (m *Model) Reset() {
	m.value = nil
	m.pos = 0
}

// Update handles editing keys while focused: printable text inserts at the
// cursor; backspace/delete remove; left/right/home/end (and ctrl+a/ctrl+e)
// move; ctrl+u clears before the cursor, ctrl+k after, ctrl+w deletes the
// previous word. Paste is accepted only while focused: CRLF/CR become LF,
// LF and tabs become spaces, and terminal escapes and other C0/C1 controls
// (including DEL) are removed. Paste does not emit a command.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if next, handled := m.applyAction(msg); handled {
		return next, nil
	}
	if !m.focused {
		return m, nil
	}
	if msg, ok := msg.(tea.PasteMsg); ok {
		if content := paste.Normalize(msg.Content, paste.SingleLine); content != "" {
			m.insert([]rune(content))
		}
		return m, nil
	}
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	commandMods := tea.ModCtrl | tea.ModAlt | tea.ModMeta | tea.ModSuper | tea.ModHyper
	if key.Text != "" && key.Mod&commandMods == 0 {
		m.insert([]rune(key.Text))
		return m, nil
	}

	switch key.String() {
	case "backspace":
		if m.pos > 0 {
			start := previousClusterStart(m.value, m.pos)
			end := nextClusterEnd(m.value, start)
			m.value = append(append([]rune{}, m.value[:start]...), m.value[end:]...)
			m.pos = start
		}
	case "delete":
		if m.pos < len(m.value) {
			start := clusterStartAt(m.value, m.pos)
			end := nextClusterEnd(m.value, start)
			m.value = append(append([]rune{}, m.value[:start]...), m.value[end:]...)
			m.pos = start
		}
	case "left":
		m.pos = previousClusterStart(m.value, m.pos)
	case "right":
		m.pos = nextClusterEnd(m.value, m.pos)
	case "home", "ctrl+a":
		m.pos = 0
	case "end", "ctrl+e":
		m.pos = len(m.value)
	case "ctrl+u":
		m.value = append([]rune{}, m.value[m.pos:]...)
		m.pos = 0
	case "ctrl+k":
		m.value = m.value[:m.pos]
	case "ctrl+w":
		m.deleteWordBack()
	}
	return m, nil
}

func (m *Model) insert(runes []rune) {
	value := make([]rune, 0, len(m.value)+len(runes))
	value = append(value, m.value[:m.pos]...)
	value = append(value, runes...)
	value = append(value, m.value[m.pos:]...)
	m.value = value
	m.pos += len(runes)
}

func (m *Model) deleteWordBack() {
	i := m.pos
	for i > 0 && m.value[i-1] == ' ' {
		i--
	}
	for i > 0 && m.value[i-1] != ' ' {
		i--
	}
	m.value = append(append([]rune{}, m.value[:i]...), m.value[m.pos:]...)
	m.pos = i
}

// NaturalHeight reports the number of rows View renders at the current
// settings, before any caller-imposed height cap. Inline height depends on
// focus: an unfocused prompt is two rows; a focused prompt includes its field
// rows and footer.
func (m Model) NaturalHeight() int {
	if m.Inline {
		if !m.focused {
			return 2 // submitted header and value
		}
		height := 3 // header, input, footer
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

// View renders the optional title/description, the input row, and an optional
// hint row. The focus rail is the only focus indicator.
func (m Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	var rows []string
	if m.Inline {
		rows = append(rows, m.clack.HeaderRow(m.focused, m.Title, m.width))
		if !m.focused {
			value := m.displayString()
			if value == "" {
				value = m.Placeholder
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
		bodyWidth := clackstyle.BodyWidth(m.width)
		rows = append(rows, m.clack.Row(m.focused, m.renderControlWithStyle("› ", m.clack.Prompt, bodyWidth), m.width))
		if m.Help != "" {
			rows = append(rows, m.clack.HelpRow(m.focused, m.Help, m.width))
		}
		rows = append(rows, m.clack.FooterRow(m.width))
	} else {
		rows = append(rows, m.field.HeaderRows(m.focused, m.Title, m.Description, m.width)...)
		bodyWidth := fieldstyle.BodyWidth(m.width)
		if m.width <= 1 {
			bodyWidth = max(0, m.width)
		}
		rows = append(rows, m.field.FrameRow(m.focused, m.renderControl(bodyWidth), m.width))
		if help := m.field.HelpRow(m.focused, m.Help, m.width); help != "" {
			rows = append(rows, help)
		}
	}
	if len(rows) > m.height {
		rows = rows[:m.height]
	}
	return strings.Join(rows, "\n")
}

// renderControl renders the prompt, visible text window, and cursor within
// avail cells. It returns an unpadded fragment; FrameRow supplies the field
// rail and exact-width padding.
func (m Model) renderControl(avail int) string {
	return m.renderControlWithPrompt(m.Prompt, avail)
}

func (m Model) renderControlWithPrompt(rawPrompt string, avail int) string {
	return m.renderControlWithStyle(rawPrompt, m.promptStyle, avail)
}

func (m Model) renderControlWithStyle(rawPrompt string, promptStyle lipgloss.Style, avail int) string {
	if avail <= 0 {
		return ""
	}
	prompt := ansi.Truncate(rawPrompt, avail, "")
	promptWidth := ansi.StringWidth(prompt)
	textWidth := avail - promptWidth
	if textWidth <= 0 {
		return promptStyle.Render(prompt)
	}
	prompt = promptStyle.Render(prompt)

	value, pos := m.displayValue()
	if len(value) == 0 {
		return prompt + m.renderPlaceholder(textWidth)
	}

	start, end, focusStart, focusEnd, focusVisible := windowForCursor(value, pos, textWidth, m.focused)
	var b strings.Builder
	b.WriteString(prompt)
	if m.focused {
		if focusVisible && focusStart < focusEnd {
			b.WriteString(m.textStyle.Render(string(value[start:focusStart])))
			b.WriteString(m.cursorStyle.Render(string(value[focusStart:focusEnd])))
			b.WriteString(m.textStyle.Render(string(value[focusEnd:end])))
		} else {
			b.WriteString(m.textStyle.Render(string(value[start:end])))
			if pos >= len(value) || !focusVisible {
				b.WriteString(m.cursorStyle.Render(" "))
			}
		}
	} else {
		b.WriteString(m.textStyle.Render(string(value[start:end])))
	}
	return b.String()
}

// displayValue returns the rendered runes and the cursor position mapped into
// them. With Mask set, every grapheme cluster becomes one mask sequence so
// horizontal windowing and the cursor stay aligned without exposing the value.
func (m Model) displayValue() ([]rune, int) {
	if m.Mask == "" {
		return m.value, m.pos
	}
	mask := []rune(m.Mask)
	if len(mask) == 0 {
		return m.value, m.pos
	}
	clusters := clustersOf(m.value)
	if len(clusters) == 0 {
		return nil, 0
	}
	display := make([]rune, 0, len(clusters)*len(mask))
	pos := 0
	for _, cluster := range clusters {
		if cluster.end <= m.pos {
			pos += len(mask)
		}
		display = append(display, mask...)
	}
	return display, pos
}

func (m Model) displayString() string {
	value, _ := m.displayValue()
	return string(value)
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func (m Model) renderPlaceholder(avail int) string {
	ph := []rune(m.Placeholder)
	if !m.focused {
		return m.placeholderStyle.Render(ansi.Truncate(string(ph), avail, ""))
	}
	if len(ph) == 0 || avail <= 0 {
		return m.cursorStyle.Render(" ")
	}
	clusters := clustersOf(ph)
	if len(clusters) == 0 || clusters[0].width > avail {
		return m.cursorStyle.Render(" ")
	}
	first := clusters[0]
	remaining := ansi.Truncate(string(ph[first.end:]), avail-first.width, "")
	return m.cursorStyle.Render(string(ph[first.start:first.end])) + m.placeholderStyle.Render(remaining)
}
