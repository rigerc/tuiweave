package textarea

import (
	"testing"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/snaptest"
)

func TestFieldLayoutGolden(t *testing.T) {
	m := newFocused(40, 8)
	m.Title = "Instructions"
	m.Description = "Optional details for the agent"
	m.Help = "ctrl+enter sends"
	m.SetValue("first line\nsecond line")
	m.SetSize(40, 8)

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestContentHeightIncludesFieldDecoration(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Instructions"
	m.Description = "Optional"
	m.Help = "ctrl+enter sends"
	if got := m.ContentHeight(); got != 4 {
		t.Fatalf("ContentHeight = %d, want 4 (content + title + description + help)", got)
	}
}
