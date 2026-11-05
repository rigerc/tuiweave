package textinput

import (
	"strings"
	"testing"

	"github.com/ishan5ain/tuiweave"
	"github.com/ishan5ain/tuiweave/snaptest"
)

func TestFieldLayoutGolden(t *testing.T) {
	ti := New(tuiweave.Dark())
	ti.Title = "Name"
	ti.Description = "Your display name"
	ti.Help = "enter to submit"
	ti.Placeholder = "Ada Lovelace"
	ti.SetSize(36, 4)
	ti.Focus()

	snaptest.Snap(t, ti.View())
	snaptest.SnapCells(t, ti.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineFieldLayoutGolden(t *testing.T) {
	ti := New(tuiweave.Dark())
	ti.Title = "Name"
	ti.Description = "Your display name"
	ti.Placeholder = "Ada Lovelace"
	ti.Inline = true
	ti.Focus()
	ti.SetSize(44, ti.NaturalHeight())

	snaptest.Snap(t, ti.View())
	snaptest.SnapCells(t, ti.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineFieldUsesClackBlock(t *testing.T) {
	ti := New(tuiweave.Dark())
	ti.Title = "Name"
	ti.Description = "Your display name"
	ti.Inline = true
	ti.Focus()
	ti.SetSize(44, ti.NaturalHeight())

	view := ti.View()
	if got := strings.Count(view, "\n") + 1; got != ti.NaturalHeight() {
		t.Fatalf("inline field height = %d, want %d", got, ti.NaturalHeight())
	}
	for _, want := range []string{"◆", "Name", "│", "›", "└"} {
		if !strings.Contains(view, want) {
			t.Fatalf("inline block missing %q: %q", want, view)
		}
	}
}

func TestFieldNaturalHeight(t *testing.T) {
	ti := New(tuiweave.Dark())
	if got := ti.NaturalHeight(); got != 1 {
		t.Fatalf("empty natural height = %d, want 1", got)
	}
	ti.Title = "Name"
	ti.Description = "Your display name"
	ti.Help = "enter"
	if got := ti.NaturalHeight(); got != 4 {
		t.Fatalf("stacked natural height = %d, want 4", got)
	}
	ti.Inline = true
	ti.Focus()
	if got := ti.NaturalHeight(); got != 5 {
		t.Fatalf("inline natural height = %d, want 5", got)
	}
}

func TestInlineHeightCollapsesWhenBlurred(t *testing.T) {
	ti := New(tuiweave.Dark())
	ti.Title = "Name"
	ti.Description = "Your display name"
	ti.Inline = true
	ti.Focus()
	if got := ti.NaturalHeight(); got != 4 {
		t.Fatalf("focused height = %d, want 4", got)
	}
	ti.Blur()
	if got := ti.NaturalHeight(); got != 2 {
		t.Fatalf("blurred height = %d, want 2", got)
	}
}
