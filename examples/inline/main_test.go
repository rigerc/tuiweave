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

func sized(t *testing.T, width, height int) model {
	t.Helper()
	m, _ := update(t, newModel(), tea.WindowSizeMsg{Width: width, Height: height})
	return m
}

func TestInlineScreenGolden(t *testing.T) {
	m := sized(t, 88, 26)
	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}

func TestInlineNarrowGolden(t *testing.T) {
	m := sized(t, 60, 20)
	snaptest.Snap(t, m.View().Content)
}

func TestInlineInteractionGolden(t *testing.T) {
	m := sized(t, 88, 26)

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // email
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'a', Text: "a"})
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'd', Text: "d"})
	m, _ = update(t, m, tea.KeyPressMsg{Code: 'a', Text: "a"})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // environment
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyRight})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // auto-refresh
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})

	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab}) // theme
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})

	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}
