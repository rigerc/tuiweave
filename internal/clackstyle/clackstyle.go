// Package clackstyle holds the Clack-inspired prompt grammar used by inline
// fields: a ◆/◇ step marker, │ body rows, a › input prompt, ❯ ●/○ selection
// rows, and a └ footer.
//
// It mirrors the visual language of the Vercel Skills CLI, which is built on
// @clack/prompts. It is an internal composition helper; components still
// derive every color from tuiweave.Theme.
package clackstyle

import (
	"image/color"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ishan5ain/tuiweave"
)

// PrefixWidth is the standard width of the bar, gutter, and content start.
const PrefixWidth = 3 // "│  "

// Styles is the Clack prompt vocabulary derived from a Theme.
type Styles struct {
	Step       lipgloss.Style
	StepIdle   lipgloss.Style
	Bar        lipgloss.Style
	Footer     lipgloss.Style
	Title      lipgloss.Style
	TitleIdle  lipgloss.Style
	Desc       lipgloss.Style
	Help       lipgloss.Style
	Prompt     lipgloss.Style
	PromptIdle lipgloss.Style
	Cursor     lipgloss.Style
	Selected   lipgloss.Style
	Inactive   lipgloss.Style
	Disabled   lipgloss.Style
	Pad        lipgloss.Style
}

// New returns Clack styles with no background. Use it for prompts that sit
// directly on the application surface.
func New(theme tuiweave.Theme) Styles { return NewOn(theme, nil) }

// NewOn returns Clack styles that paint background behind every run. Use it
// when a prompt is embedded in a panel so nested resets cannot leave gaps.
func NewOn(theme tuiweave.Theme, background color.Color) Styles {
	base := lipgloss.NewStyle()
	if background != nil {
		base = base.Background(background)
	}
	fg := func(c color.Color) lipgloss.Style { return base.Foreground(c) }
	return Styles{
		Step:       fg(theme.Info).Bold(true),
		StepIdle:   fg(theme.Success),
		Bar:        fg(theme.TextFaint),
		Footer:     fg(theme.TextFaint),
		Title:      fg(theme.Text).Bold(true),
		TitleIdle:  fg(theme.TextMuted),
		Desc:       fg(theme.TextMuted),
		Help:       fg(theme.TextMuted),
		Prompt:     fg(theme.Info).Bold(true),
		PromptIdle: fg(theme.TextMuted),
		Cursor:     fg(theme.Info).Bold(true),
		Selected:   fg(theme.Success).Bold(true),
		Inactive:   fg(theme.TextMuted),
		Disabled:   fg(theme.TextFaint),
		Pad:        base,
	}
}

// PrefixWidthFor returns the bar budget for a box width.
func PrefixWidthFor(width int) int {
	switch {
	case width <= 0:
		return 0
	case width < 4:
		return 1
	default:
		return PrefixWidth
	}
}

// BodyWidth returns the cells available after the bar and its gutter.
func BodyWidth(width int) int {
	return max(0, width-PrefixWidthFor(width))
}

func (s Styles) prefix(width int, _ bool) string {
	switch PrefixWidthFor(width) {
	case 0:
		return ""
	case 1:
		return s.Bar.Render("│")
	default:
		return s.Bar.Render("│") + s.Pad.Render("  ")
	}
}

func (s Styles) fit(body string, width int) string {
	body = ansi.Truncate(body, width, "…")
	if pad := width - ansi.StringWidth(body); pad > 0 {
		body += s.Pad.Render(strings.Repeat(" ", pad))
	}
	return body
}

// HeaderRow renders "◆  Title" when focused and "◇  Title" otherwise.
func (s Styles) HeaderRow(focused bool, title string, width int) string {
	if width <= 0 {
		return ""
	}
	prefixWidth := PrefixWidthFor(width)
	if prefixWidth == 0 {
		return ""
	}
	step := s.Step.Render("◆")
	titleStyle := s.Title
	if !focused {
		step = s.StepIdle.Render("◇")
		titleStyle = s.TitleIdle
	}
	if prefixWidth == 1 {
		if width == 1 {
			return step
		}
		return step + s.Pad.Render(" ")
	}
	return step + s.Pad.Render("  ") + s.fit(titleStyle.Render(title), width-PrefixWidth)
}

// Row renders a body line under the bar: "│  body".
func (s Styles) Row(focused bool, body string, width int) string {
	if width <= 0 {
		return ""
	}
	prefix := s.prefix(width, focused)
	return prefix + s.fit(body, BodyWidth(width))
}

// FooterRow renders the closing "└" line.
func (s Styles) FooterRow(width int) string {
	if width <= 0 {
		return ""
	}
	return s.Footer.Render("└") + s.Pad.Render(strings.Repeat(" ", max(0, width-1)))
}

// InputRow renders a Clack text prompt line: "│  › value".
func (s Styles) InputRow(focused bool, prompt, value string, width int) string {
	style := s.Prompt
	if !focused {
		style = s.PromptIdle
	}
	return s.Row(focused, style.Render(prompt)+value, width)
}

// OptionRow renders one Clack selection row. The selected row gets the cursor
// and a filled radio; other rows get an empty radio and muted text.
func (s Styles) OptionRow(focused, selected, disabled bool, label string, width int) string {
	cursor := s.Pad.Render("  ")
	mark := s.Inactive.Render("○")
	labelStyle := s.Inactive
	if disabled {
		labelStyle = s.Disabled
	} else if selected {
		if focused {
			cursor = s.Cursor.Render("❯ ")
			labelStyle = s.Title
		} else {
			labelStyle = s.TitleIdle
		}
		mark = s.Selected.Render("●")
	}
	body := cursor + mark + s.Pad.Render(" ") + labelStyle.Render(label)
	return s.Row(focused, body, width)
}

// HelpRow renders a muted hint under the bar.
func (s Styles) HelpRow(focused bool, text string, width int) string {
	if text == "" {
		return ""
	}
	return s.Row(focused, s.Help.Render(text), width)
}

// InstructionsRow renders muted prompt instructions under the bar.
func (s Styles) InstructionsRow(focused bool, text string, width int) string {
	if text == "" {
		return ""
	}
	return s.Row(focused, s.Help.Render(text), width)
}
