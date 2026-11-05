// Package fieldstyle holds the shared Huh-inspired field grammar used by
// form-like components: a thick left focus rail, title and description rows,
// a right-pointing selector, and action chips.
//
// It is an internal composition helper. Components still derive every color
// from tuiweave.Theme and remain responsible for their own layout and
// interaction behavior.
package fieldstyle

import (
	"image/color"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ishan5ain/tuiweave"
)

// RailWidth is the standard number of cells the focus rail and its inset
// occupy on normal-width boxes.
const RailWidth = 2

// RailWidthFor returns the rail budget for a box. Narrow boxes reduce the rail
// to a single cell so content is not starved; very narrow boxes drop it.
func RailWidthFor(width int) int {
	switch {
	case width <= 0:
		return 0
	case width < 6:
		return 0
	case width < 7:
		return 1
	default:
		return RailWidth
	}
}

// BodyWidth returns the cells available after the focus rail.
func BodyWidth(width int) int {
	return max(0, width-RailWidthFor(width))
}

// Styles is the shared field vocabulary derived from a Theme.
type Styles struct {
	RailFocus  lipgloss.Style
	Title      lipgloss.Style
	Desc       lipgloss.Style
	Help       lipgloss.Style
	SelectMark lipgloss.Style
	Selected   lipgloss.Style
	Normal     lipgloss.Style
	Disabled   lipgloss.Style
	Pad        lipgloss.Style

	buttonFocused  lipgloss.Style
	buttonBlurred  lipgloss.Style
	buttonDisabled lipgloss.Style
}

// New returns field styles with no background. Use it for controls that sit
// directly on the application surface.
func New(theme tuiweave.Theme) Styles { return NewOn(theme, nil) }

// NewOn returns field styles that paint background behind every run. Raised
// panels, menus, popovers, dialogs, and chips pass theme.SurfaceRaised so
// nested style resets cannot leave unpainted gaps inside their surface.
func NewOn(theme tuiweave.Theme, background color.Color) Styles {
	base := lipgloss.NewStyle()
	if background != nil {
		base = base.Background(background)
	}
	fg := func(c color.Color) lipgloss.Style { return base.Foreground(c) }
	return Styles{
		RailFocus:  base.Foreground(theme.Accent).Bold(true),
		Title:      fg(theme.Accent).Bold(true),
		Desc:       fg(theme.TextMuted),
		Help:       fg(theme.TextMuted),
		SelectMark: fg(theme.Accent).Bold(true),
		Selected:   fg(theme.Accent).Bold(true),
		Normal:     fg(theme.Text),
		Disabled:   fg(theme.TextFaint),
		Pad:        base,
		buttonFocused: lipgloss.NewStyle().
			Foreground(theme.TextInverted).
			Background(theme.Accent).
			Bold(true),
		buttonBlurred: lipgloss.NewStyle().
			Foreground(theme.Text).
			Background(theme.SurfaceSunken),
		buttonDisabled: lipgloss.NewStyle().
			Foreground(theme.TextFaint).
			Background(theme.SurfaceSunken),
	}
}

func (s Styles) chipStyle(focused, disabled bool) lipgloss.Style {
	switch {
	case disabled:
		return s.buttonDisabled
	case focused:
		return s.buttonFocused
	default:
		return s.buttonBlurred
	}
}

// rail returns the focus rail for a box width. A focused field gets a thick
// left border; a blurred field gets the same inset without the border so
// content does not jump between focus states.
func (s Styles) rail(width int, focused bool) string {
	switch RailWidthFor(width) {
	case 0:
		return ""
	case 1:
		if focused {
			return s.RailFocus.Render("┃")
		}
		return s.Pad.Render(" ")
	default:
		if focused {
			return s.RailFocus.Render("┃") + s.Pad.Render(" ")
		}
		return s.Pad.Render("  ")
	}
}

// FrameRow places body inside the field rail and pads or truncates it to
// exactly width cells. Narrow boxes degrade to a one-cell rail or no rail.
func (s Styles) FrameRow(focused bool, body string, width int) string {
	if width <= 0 {
		return ""
	}

	inner := BodyWidth(width)
	body = ansi.Truncate(body, inner, "…")
	if pad := inner - ansi.StringWidth(body); pad > 0 {
		body += s.Pad.Render(strings.Repeat(" ", pad))
	}
	return s.rail(width, focused) + body
}

// HeaderRows renders the optional title and description inside the field rail.
// The caller owns placement; rows are already exact-width.
func (s Styles) HeaderRows(focused bool, title, description string, width int) []string {
	if width <= 0 {
		return nil
	}
	rows := make([]string, 0, 2)
	if title != "" {
		rows = append(rows, s.FrameRow(focused, s.Title.Render(title), width))
	}
	if description != "" {
		rows = append(rows, s.FrameRow(focused, s.Desc.Render(description), width))
	}
	return rows
}

// HeaderHeight reports how many rows a title/description pair occupies.
func HeaderHeight(title, description string) int {
	height := 0
	if title != "" {
		height++
	}
	if description != "" {
		height++
	}
	return height
}

// InlineHeader renders title and description as one label followed by two
// spaces, truncated to maxWidth. It returns "" when both are empty. Callers
// subtract the returned width from the control area.
func (s Styles) InlineHeader(title, description string, maxWidth int) string {
	if maxWidth <= 2 {
		return ""
	}
	var b strings.Builder
	if title != "" {
		b.WriteString(s.Title.Render(title))
	}
	if description != "" {
		if b.Len() > 0 {
			b.WriteString(s.Pad.Render("  "))
		}
		b.WriteString(s.Desc.Render(description))
	}
	if b.Len() == 0 {
		return ""
	}
	label := ansi.Truncate(b.String(), maxWidth-2, "…")
	return label + s.Pad.Render("  ")
}

// Selector returns the Huh-style selected-option marker.
func (s Styles) Selector(selected bool) string {
	if selected {
		return s.SelectMark.Render("> ")
	}
	return s.Pad.Render("  ")
}

// Label renders an option label according to its state.
func (s Styles) Label(label string, selected, disabled bool) string {
	switch {
	case disabled:
		return s.Disabled.Render(label)
	case selected:
		return s.Selected.Render(label)
	default:
		return s.Normal.Render(label)
	}
}

// Chip renders one button label as a padded action chip.
func (s Styles) Chip(label string, focused, disabled bool) string {
	return s.chipStyle(focused, disabled).Padding(0, 2).Render(label)
}

// ChipNarrow renders a chip without horizontal padding for boxes too narrow
// for the standard 2-cell insets.
func (s Styles) ChipNarrow(label string, focused, disabled bool, width int) string {
	if width <= 0 {
		return ""
	}
	return s.chipStyle(focused, disabled).Render(ansi.Truncate(label, width, "…"))
}

// Arrow renders an inline-picker arrow. Disabled arrows are faint.
func (s Styles) Arrow(right, disabled bool) string {
	glyph := "←"
	if right {
		glyph = "→"
	}
	if disabled {
		return s.Disabled.Render(glyph)
	}
	return s.SelectMark.Render(glyph)
}

// HelpRow renders an optional hint row inside the field rail.
func (s Styles) HelpRow(focused bool, help string, width int) string {
	if help == "" || width <= 0 {
		return ""
	}
	return s.FrameRow(focused, s.Help.Render(help), width)
}
