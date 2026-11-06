package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave/snaptest"
)

func update(t *testing.T, m model, msg tea.Msg) (model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	got, ok := next.(model)
	if !ok {
		t.Fatalf("Update returned %T, want model", next)
	}
	return got, cmd
}

func sized(t *testing.T) model {
	t.Helper()
	m, _ := update(t, newModel(), tea.WindowSizeMsg{Width: 80, Height: 36})
	return m
}

func TestFormScreenGolden(t *testing.T) {
	m := sized(t)
	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}

func TestFormInlineInteractionGolden(t *testing.T) {
	m := sized(t)

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // search
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'l', Text: "l"})
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'o', Text: "o"})
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'g', Text: "g"})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // environment
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyRight})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // checks
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})
	m, cmd := update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})
	if cmd == nil {
		t.Fatal("multiselect toggle did not emit a command")
	}
	m, _ = update(t, m, cmd())

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // steps
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})
	m, cmd = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd == nil {
		t.Fatal("reorder move did not emit a command")
	}
	m, _ = update(t, m, cmd())
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab})   // auto-refresh
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace}) // toggle

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // save
	m, cmd = update(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("save did not emit a command")
	}
	m, _ = update(t, m, cmd())

	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}
