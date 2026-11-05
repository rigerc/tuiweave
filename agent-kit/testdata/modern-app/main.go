package main

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/textarea"
)

type model struct{ input textarea.Model }

func main() {}

func newModel() model {
	m := model{input: textarea.New(tuiweave.Dark())}
	m.input.SetSize(24, 3)
	m.input.Focus()
	return m
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.input.SetSize(size.Width, size.Height)
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() tea.View { return tea.NewView(m.input.View()) }
