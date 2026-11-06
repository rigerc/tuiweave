package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/rigerc/tuiweave/snaptest"
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

func tab(t *testing.T, m model, times int) model {
	t.Helper()
	for range times {
		m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyTab})
	}
	return m
}

func typeText(t *testing.T, m model, text string) model {
	t.Helper()
	for _, r := range text {
		m, _ = update(t, m, tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return m
}

func TestInlineScreenGolden(t *testing.T) {
	m := sized(t, 88, 30)
	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}

func TestInlineNarrowGolden(t *testing.T) {
	m := sized(t, 60, 20)
	snaptest.Snap(t, m.View().Content)
}

func TestInlineNarrowMultiselectGolden(t *testing.T) {
	m := sized(t, 44, 24)
	m = tab(t, m, 4) // features
	snaptest.Snap(t, m.View().Content)
}

func TestInlineNarrowReorderGolden(t *testing.T) {
	m := sized(t, 44, 26)
	m = tab(t, m, 5) // steps
	snaptest.Snap(t, m.View().Content)
}

func TestInlineInteractionGolden(t *testing.T) {
	m := sized(t, 88, 30)

	m = typeText(t, m, "Ada")
	m = tab(t, m, 1) // email
	m = typeText(t, m, "ada@example.com")
	m = tab(t, m, 1) // password
	m = typeText(t, m, "hunter2")

	m = tab(t, m, 1) // environment
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyRight})

	m = tab(t, m, 1) // features
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})
	var cmd tea.Cmd
	m, cmd = update(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		m, _ = update(t, m, cmd())
	}

	m = tab(t, m, 1) // steps
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})
	m, cmd = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd != nil {
		m, _ = update(t, m, cmd())
	}
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})

	m = tab(t, m, 1) // port
	m = typeText(t, m, "8000x")

	m = tab(t, m, 1) // auto-refresh
	m, cmd = update(t, m, tea.KeyPressMsg{Code: tea.KeySpace})
	if cmd != nil {
		m, _ = update(t, m, cmd())
	}

	m = tab(t, m, 1) // theme
	m, _ = update(t, m, tea.KeyPressMsg{Code: tea.KeyDown})

	snaptest.Snap(t, m.View().Content)
	snaptest.SnapCells(t, m.View().Content, snaptest.WithRoles(m.theme))
}
