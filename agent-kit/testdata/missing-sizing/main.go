package main

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/textarea"
)

type model struct{ input textarea.Model }

func newModel() model { return model{input: textarea.New(tuiweave.Dark())} }

func (m model) Update(msg tea.Msg) (model, tea.Cmd) {
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		// Intentionally missing SetSize for the audit fixture.
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
