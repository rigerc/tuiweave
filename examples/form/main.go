// Command form demonstrates Huh-style stacked fields: a focus rail, title and
// description rows, a vertical select, a checkbox, and an action chip.
package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/action"
	"github.com/ishan5ain/tuiweave/button"
	"github.com/ishan5ain/tuiweave/focus"
	"github.com/ishan5ain/tuiweave/frame"
	"github.com/ishan5ain/tuiweave/help"
	"github.com/ishan5ain/tuiweave/menu"
	"github.com/ishan5ain/tuiweave/stack"
	"github.com/ishan5ain/tuiweave/statusbar"
	"github.com/ishan5ain/tuiweave/textinput"
	"github.com/ishan5ain/tuiweave/toggle"
)

type model struct {
	width, height int
	theme         tuiweave.Theme

	header statusbar.Model
	help   help.Model
	name   textinput.Model
	search textinput.Model
	env    menu.Model
	auto   toggle.Model
	save   button.Model

	fm         focus.Manager
	noticeKind statusbar.Kind
}

func newModel() model {
	theme := tuiweave.Dark()

	header := statusbar.New(theme)
	header.SetLeft(
		statusbar.Segment{Text: "FORM", Kind: statusbar.KindAccent},
		statusbar.Segment{Text: "field demo", Kind: statusbar.KindNormal},
	)
	header.SetRight(statusbar.Segment{Text: "ready", Kind: statusbar.KindMuted})

	footer := help.New(theme)
	footer.SetBindings(
		help.Binding{Key: "tab", Desc: "focus"},
		help.Binding{Key: "enter", Desc: "activate"},
		help.Binding{Key: "←/→", Desc: "change"},
		help.Binding{Key: "esc", Desc: "quit"},
	)

	name := textinput.NewOn(theme, theme.SurfaceRaised)
	name.Title = "Name"
	name.Description = "Shown on your profile"
	name.Placeholder = "Ada Lovelace"
	name.SetValue("Ada")

	search := textinput.NewOn(theme, theme.SurfaceRaised)
	search.Title = "Search"
	search.Placeholder = "type to filter"

	env := menu.New(theme)
	env.Title = "Environment"
	env.Description = "Choose a target"
	env.SetItems(
		action.Item{ID: "staging", Label: "Staging"},
		action.Item{ID: "production", Label: "Production"},
		action.Item{ID: "development", Label: "Development"},
	)

	auto := toggle.New(theme)
	auto.Title = "Updates"
	auto.SetLabel("Auto-refresh")
	auto.SetChecked(true)

	save := button.New(theme)
	save.ID = "save"
	save.SetLabel("Save changes")

	return model{
		theme:      theme,
		header:     header,
		help:       footer,
		name:       name,
		search:     search,
		env:        env,
		auto:       auto,
		save:       save,
		fm:         focus.NewManager(5),
		noticeKind: statusbar.KindMuted,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		contentWidth := max(0, msg.Width-4)
		m.header.SetSize(msg.Width, 1)
		m.help.SetSize(msg.Width, 1)
		m.name.SetSize(contentWidth, m.name.NaturalHeight())
		m.search.SetSize(contentWidth, m.search.NaturalHeight())
		m.env.SetSize(contentWidth, m.env.NaturalHeight())
		m.auto.SetSize(contentWidth, m.auto.NaturalHeight())
		m.save.SetSize(contentWidth, 1)
		m.applyFocus()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab":
			m.fm.Next()
			m.applyFocus()
			return m, nil
		case "shift+tab":
			m.fm.Prev()
			m.applyFocus()
			return m, nil
		}

	case button.PressedMsg:
		m.noticeKind = statusbar.KindSuccess
		m.setNotice("saved")
	case toggle.ChangedMsg:
		m.noticeKind = statusbar.KindInfo
		m.setNotice(fmt.Sprintf("auto-refresh %v", msg.Checked))
	case menu.SelectedMsg:
		m.noticeKind = statusbar.KindInfo
		m.setNotice("environment: " + msg.Label)
	}

	var cmd tea.Cmd
	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)
	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)
	m.env, cmd = m.env.Update(msg)
	cmds = append(cmds, cmd)
	m.auto, cmd = m.auto.Update(msg)
	cmds = append(cmds, cmd)
	m.save, cmd = m.save.Update(msg)
	cmds = append(cmds, cmd)

	if msg, ok := msg.(tea.KeyPressMsg); ok && msg.String() == "enter" {
		if m.fm.Index() == 0 {
			m.noticeKind = statusbar.KindInfo
			m.setNotice("name: " + m.name.Value())
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *model) applyFocus() {
	m.fm.Apply(&m.name, &m.search, &m.env, &m.auto, &m.save)
}

func (m *model) setNotice(text string) {
	m.header.SetRight(statusbar.Segment{Text: text, Kind: m.noticeKind})
}

func (m model) View() tea.View {
	if m.width <= 0 || m.height <= 0 {
		return tea.NewView("")
	}
	bodyHeight := max(0, m.height-2)
	contentWidth := max(0, m.width-4)
	contentHeight := max(0, bodyHeight-4)
	content := stack.Vertical(m.theme, contentWidth, stack.Options{Gap: 1},
		func(int) string { return m.name.View() },
		func(int) string { return m.search.View() },
		func(int) string { return m.env.View() },
		func(int) string { return m.auto.View() },
		func(int) string { return m.save.View() },
	)
	content = fitHeight(content, contentWidth, contentHeight)
	panel := frame.Panel(m.theme, content, m.width, frame.PanelOptions{
		Title:   "Preferences",
		Padding: 1,
		Focused: true,
	})
	parts := []string{m.header.View()}
	if bodyHeight > 0 {
		parts = append(parts, fitHeight(panel, m.width, bodyHeight))
	}
	if m.height > 1 {
		parts = append(parts, m.help.View())
	}
	v := tea.NewView(strings.Join(parts, "\n"))
	v.AltScreen = true
	return v
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
