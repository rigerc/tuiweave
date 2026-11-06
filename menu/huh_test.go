package menu

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	tea "charm.land/bubbletea/v2"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/layout"
	"github.com/rigerc/tuiweave/snaptest"
)

func huhItems() []Item {
	return []Item{
		{ID: "red", Label: "red"},
		{ID: "green", Label: "green"},
		{ID: "blue", Label: "blue"},
	}
}

func TestFieldLayoutGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Color"
	m.Description = "Pick a color"
	m.SetItems(huhItems()...)
	m.SetSize(36, 5)
	m.Focus()

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineLayoutGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Color"
	m.Description = "Pick a color"
	m.Inline = true
	m.SetItems(huhItems()...)
	m.Focus()
	m.SetSize(36, m.NaturalHeight())

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlinePickerNavigatesWithHorizontalKeys(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Inline = true
	m.SetItems(huhItems()...)
	m.Focus()
	m.SetSize(30, m.NaturalHeight())

	if got := m.SelectedID(); got != "red" {
		t.Fatalf("initial selection = %q, want red", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if got := m.SelectedID(); got != "green" {
		t.Fatalf("right selection = %q, want green", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if got := m.SelectedID(); got != "red" {
		t.Fatalf("left selection = %q, want red", got)
	}

	view := ansi.Strip(m.View())
	for _, want := range []string{"◆", "❯ ● red", "○ green", "│", "└"} {
		if !strings.Contains(view, want) {
			t.Fatalf("inline prompt missing %q: %q", want, view)
		}
	}
}

func TestInlinePickerActivationStillEmitsSelectedMsg(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Inline = true
	m.SetItems(huhItems()...)
	m.SetSize(30, 1)
	m.Focus()

	cmd := m.Activate()
	if cmd == nil {
		t.Fatal("activate returned nil")
	}
	msg, ok := cmd().(SelectedMsg)
	if !ok || msg.ID != "red" {
		t.Fatalf("selected msg = %#v, want red", msg)
	}
}

func TestInlineMenuIsWidthBounded(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Inline = true
	m.SetItems(huhItems()...)
	if got := layout.SizeModeOf(&m); got != layout.SizeWidthBounded {
		t.Fatalf("inline size mode = %d, want SizeWidthBounded", got)
	}
	m.Focus()
	if got := m.NaturalHeight(); got != 6 {
		t.Fatalf("focused height = %d, want 6", got)
	}
	m.Blur()
	if got := m.NaturalHeight(); got != 2 {
		t.Fatalf("blurred height = %d, want 2", got)
	}
}
