// Command inline demonstrates the Clack-style inline prompts used by the
// Vercel Skills CLI: a ◆ header, │ body, › input, ❯ ●/○ or ◼/◻ options, and a
// └ footer. The focused prompt expands; the others collapse to submitted
// values. The form covers text, password, select, multiselect, reorder,
// validated number, and confirm fields.
package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/action"
	"github.com/ishan5ain/tuiweave/focus"
	"github.com/ishan5ain/tuiweave/menu"
	"github.com/ishan5ain/tuiweave/multiselect"
	"github.com/ishan5ain/tuiweave/reorder"
	"github.com/ishan5ain/tuiweave/textinput"
	"github.com/ishan5ain/tuiweave/toggle"
)

type model struct {
	width, height int
	theme         tuiweave.Theme

	name      textinput.Model
	email     textinput.Model
	password  textinput.Model
	env       menu.Model
	features  multiselect.Model
	steps     reorder.Model
	port      textinput.Model
	auto      toggle.Model
	themePick menu.Model

	notice string
	fm     focus.Manager
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

	password := textinput.New(theme)
	password.Title = "Password"
	password.Placeholder = "stored locally"
	password.Mask = "•"
	password.Inline = true

	env := menu.NewOn(theme, nil)
	env.Title = "Environment"
	env.Inline = true
	env.SetItems(
		action.Item{ID: "staging", Label: "Staging"},
		action.Item{ID: "production", Label: "Production"},
		action.Item{ID: "development", Label: "Development"},
	)

	features := multiselect.NewOn(theme, nil)
	features.Title = "Features"
	features.Description = "Applied to new workspaces"
	features.Inline = true
	features.SetItems(
		action.Item{ID: "types", Label: "Type checking"},
		action.Item{ID: "format", Label: "Formatting"},
		action.Item{ID: "tests", Label: "Tests"},
		action.Item{ID: "preview", Label: "Deploy preview", Disabled: true},
	)
	features.SetCheckedIDs("types", "tests")

	steps := reorder.NewOn(theme, nil)
	steps.Title = "Steps"
	steps.Description = "Space grabs; arrows move"
	steps.Inline = true
	steps.SetItems(
		reorder.Item{ID: "design", Label: "Design"},
		reorder.Item{ID: "build", Label: "Build"},
		reorder.Item{ID: "test", Label: "Test"},
		reorder.Item{ID: "ship", Label: "Ship"},
	)

	port := textinput.New(theme)
	port.Title = "Port"
	port.Description = "Digits only; the app validates"
	port.Placeholder = "3000"
	port.Inline = true

	auto := toggle.NewOn(theme, nil)
	auto.ID = "auto-refresh"
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
		password:  password,
		env:       env,
		features:  features,
		steps:     steps,
		port:      port,
		auto:      auto,
		themePick: themePick,
		fm:        focus.NewManager(9),
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

	case multiselect.SubmittedMsg:
		m.notice = "features: " + selectedLabels(m.features, msg.Selected)

	case reorder.MovedMsg:
		m.notice = "steps: " + orderedLabels(m.steps)

	case toggle.ChangedMsg:
		state := "off"
		if msg.Checked {
			state = "on"
		}
		m.notice = "auto-refresh: " + state

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
	m.password, cmd = m.password.Update(msg)
	cmds = append(cmds, cmd)
	m.env, cmd = m.env.Update(msg)
	cmds = append(cmds, cmd)
	m.features, cmd = m.features.Update(msg)
	cmds = append(cmds, cmd)
	m.steps, cmd = m.steps.Update(msg)
	cmds = append(cmds, cmd)
	m.port, cmd = m.port.Update(msg)
	cmds = append(cmds, cmd)
	m.auto, cmd = m.auto.Update(msg)
	cmds = append(cmds, cmd)
	m.themePick, cmd = m.themePick.Update(msg)
	cmds = append(cmds, cmd)

	// Application-owned validation: the port field accepts digits only.
	if value := m.port.Value(); value != digitsOnly(value) {
		m.port.SetValue(digitsOnly(value))
	}
	return m, tea.Batch(cmds...)
}

func (m *model) resize() {
	m.name.SetSize(m.width, m.name.NaturalHeight())
	m.email.SetSize(m.width, m.email.NaturalHeight())
	m.password.SetSize(m.width, m.password.NaturalHeight())
	m.env.SetSize(m.width, m.env.NaturalHeight())
	m.features.SetSize(m.width, m.features.NaturalHeight())
	m.steps.SetSize(m.width, m.steps.NaturalHeight())
	m.port.SetSize(m.width, m.port.NaturalHeight())
	m.auto.SetSize(m.width, m.auto.NaturalHeight())
	m.themePick.SetSize(m.width, m.themePick.NaturalHeight())
}

func (m *model) applyFocus() {
	m.fm.Apply(
		&m.name,
		&m.email,
		&m.password,
		&m.env,
		&m.features,
		&m.steps,
		&m.port,
		&m.auto,
		&m.themePick,
	)
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
		m.password.View(),
		m.env.View(),
		m.features.View(),
		m.steps.View(),
		m.port.View(),
		m.auto.View(),
		m.themePick.View(),
		"",
	}
	if m.notice != "" {
		lines = append(lines, m.noticeRow(), "")
	}
	lines = append(lines, m.outro())

	view := fitHeight(strings.Join(lines, "\n"), m.width, m.height)
	v := tea.NewView(view)
	v.AltScreen = true
	return v
}

func (m model) intro() string {
	mark := lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("┌")
	title := lipgloss.NewStyle().Foreground(m.theme.Text).Bold(true).Render("Inline fields")
	return mark + "  " + title
}

func (m model) bar() string {
	return lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("│")
}

func (m model) noticeRow() string {
	mark := lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("│")
	text := lipgloss.NewStyle().Foreground(m.theme.TextMuted).Render(m.notice)
	return mark + "  " + text
}

func (m model) outro() string {
	mark := lipgloss.NewStyle().Foreground(m.theme.TextFaint).Render("└")
	help := lipgloss.NewStyle().Foreground(m.theme.TextMuted).
		Render("tab focus • ↑/↓ change • space toggle • enter submit • esc quit")
	return mark + "  " + help
}

// digitsOnly keeps the characters the example accepts for a port field.
func digitsOnly(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// selectedLabels maps submitted stable IDs back to display labels.
func selectedLabels(field multiselect.Model, ids []string) string {
	labels := make([]string, 0, len(ids))
	for _, id := range ids {
		for _, item := range field.Items() {
			if item.ID == id {
				labels = append(labels, item.Label)
				break
			}
		}
	}
	if len(labels) == 0 {
		return "none"
	}
	return strings.Join(labels, ", ")
}

// orderedLabels renders the current reorder sequence for the example notice.
func orderedLabels(field reorder.Model) string {
	labels := make([]string, 0, field.Len())
	for _, item := range field.Items() {
		labels = append(labels, item.Label)
	}
	return strings.Join(labels, " → ")
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
