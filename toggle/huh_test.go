package toggle

import (
	"strings"
	"testing"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/snaptest"
)

func TestFieldLayoutGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Refresh"
	m.Description = "Keep the list current"
	m.SetLabel("Auto-refresh")
	m.SetChecked(true)
	m.SetSize(36, 3)
	m.Focus()

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineFieldLayoutGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Refresh"
	m.Description = "Keep the list current"
	m.SetLabel("Auto-refresh")
	m.Inline = true
	m.Focus()
	m.SetSize(44, m.NaturalHeight())

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineFieldUsesClackBlock(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Refresh"
	m.Description = "Keep the list current"
	m.SetLabel("Auto-refresh")
	m.Inline = true
	m.Focus()
	m.SetSize(44, m.NaturalHeight())

	view := m.View()
	if got := strings.Count(view, "\n") + 1; got != m.NaturalHeight() {
		t.Fatalf("inline field height = %d, want %d", got, m.NaturalHeight())
	}
	for _, want := range []string{"◆", "│", "❯", "○", "└"} {
		if !strings.Contains(view, want) {
			t.Fatalf("inline block missing %q: %q", want, view)
		}
	}
}

func TestInlineHeightCollapsesWhenBlurred(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Refresh"
	m.Description = "Keep the list current"
	m.SetLabel("Auto-refresh")
	m.Inline = true
	m.Focus()
	if got := m.NaturalHeight(); got != 4 {
		t.Fatalf("focused height = %d, want 4", got)
	}
	m.Blur()
	if got := m.NaturalHeight(); got != 2 {
		t.Fatalf("blurred height = %d, want 2", got)
	}
}
