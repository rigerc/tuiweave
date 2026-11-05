// Command inline demonstrates the Clack-style inline prompts used by the
// Vercel Skills CLI: a ◆ header, │ body, › input, ❯ ●/○ options, and a └
// footer. The focused prompt expands; the others collapse to submitted values.
package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/action"
	"github.com/ishan5ain/tuiweave/focus"
	"github.com/ishan5ain/tuiweave/menu"
	"github.com/ishan5ain/tuiweave/textinput"
	"github.com/ishan5ain/tuiweave/toggle"
)

type model struct {
	width, height int
	theme         tuiweave.Theme

	name      textinput.Model
	email     textinput.Model
	env       menu.Model
	auto      toggle.Model
	themePick menu.Model

	fm focus.Manager
}

func newModel() model {
	theme := tuiweave.Dark()

	name := textinput.New(theme)
	name.Title = "Name"
	name.Description = "Shown on your profile"
	name.Placeholder = "Ada Lovelace"
	name.Inline = true

	email := textinput.New(theme)
	email.Title = "Email"
	email.Placeholder = "ada@example.com"
	email.Inline = true

	env := menu.NewOn(theme, nil)
	env.Title = "Environment"
	env.Inline = true
	env.SetItems(
		action.Item{ID: "staging", Label: "Staging"},
		action.Item{ID: "production", Label: "Production"},
		action.Item{ID: "development", Label: "Development"},
	)

	auto := toggle.NewOn(theme, nil)
	auto.Title = "Updates"
	auto.SetLabel("Auto-refresh")
	auto.SetChecked(true)
	auto.Inline = true

	themePick := menu.NewOn(theme, nil)
	themePick.Title = "Theme"
	themePick.Inline = true
	themePick.SetItems(
		action.Item{ID: "dark", Label: "Dark"},
		action.Item{ID: "light", Label: "Light"},
		action.Item{ID: "system", Label: "System"},
	)

	return model{
		theme:     theme,
		name:      name,
		email:     email,
		env:       env,
		auto:      auto,
		themePick: themePick,
		fm:        focus.NewManager(5),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.applyFocus()
		m.resize()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.fm.Next()
			m.applyFocus()
			m.resize()
			return m, nil
		case "shift+tab":
			m.fm.Prev()
			m.applyFocus()
			m.resize()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)
	m.email, cmd = m.email.Update(msg)
	cmds = append(cmds, cmd)
	m.env, cmd = m.env.Update(msg)
	cmds = append(cmds, cmd)
	m.auto, cmd = m.auto.Update(msg)
	cmds = append(cmds, cmd)
	m.themePick, cmd = m.themePick.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) resize() {
	m.name.SetSize(m.width, m.name.NaturalHeight())
	m.email.SetSize(m.width, m.email.NaturalHeight())
	m.env.SetSize(m.width, m.env.NaturalHeight())
	m.auto.SetSize(m.width, m.auto.NaturalHeight())
	m.themePick.SetSize(m.width, m.themePick.NaturalHeight())
}

func (m *model) applyFocus() {
	m.fm.Apply(&m.name, &m.email, &m.env, &m.auto, &m.themePick)
}

func (m model) View() tea.View {
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}

	lines := []string{
		m.intro(),
		m.bar(),
		m.name.View(),
		m.email.View(),
		m.env.View(),
		m.auto.View(),
		m.themePick.View(),
		"",
		m.outro(),
	}
	view := fitHeight(strings.Join(lines, "\n"), m.width, m.height)
	return tea.NewView(view)
}

func (m model) intro() string {
	mark := lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("┌")
	title := lipgloss.NewStyle().Foreground(m.theme.Text).Bold(true).Render("Inline fields")
	return mark + "  " + title
}

func (m model) bar() string {
	return lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("│")
}

func (m model) outro() string {
	mark := lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("└")
	help := lipgloss.NewStyle().Foreground(m.theme.TextMuted).
		Render("tab focus • ↑/↓ change • space toggle • esc quit")
	return mark + "  " + help
}

// fitHeight pads or clips a rendered block to exactly height rows.
func fitHeight(view string, width, height int) string {
	rows := strings.Split(view, "\n")
	if len(rows) > height {
		rows = rows[:height]
	}
	for len(rows) < height {
		rows = append(rows, strings.Repeat(" ", width))
	}
	return strings.Join(rows, "\n")
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		panic(err)
	}
}
