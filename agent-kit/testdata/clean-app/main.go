package main

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/textarea"
)

type model struct{ input textarea.Model }

func main() {}

func newModel() model { return model{input: textarea.New(tuiweave.Dark())} }

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.input.SetSize(size.Width, size.Height)
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
