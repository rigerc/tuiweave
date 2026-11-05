package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/snaptest"
)

func TestModernConsumer(t *testing.T) {
	m := newModel()
	result := snaptest.RunScenario(m,
		snaptest.ScenarioStep{Name: "type", Msg: tea.KeyPressMsg{Code: 'h', Text: "hello"}},
		snaptest.ScenarioStep{Name: "resize", Msg: tea.WindowSizeMsg{Width: 8, Height: 2}},
	)
	final := result.Model.(model)
	if final.input.Value() != "hello" {
		t.Fatalf("input = %q, want hello", final.input.Value())
	}
	snaptest.Snap(t, final.View().Content)
	snaptest.SnapCells(t, final.View().Content, snaptest.WithRoles(tuiweave.Dark()))
	snaptest.SnapScenario(t, result)
}
