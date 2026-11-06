package multiselect

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	tea "charm.land/bubbletea/v2"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/inspect"
	"github.com/rigerc/tuiweave/layout"
	"github.com/rigerc/tuiweave/snaptest"
)

func demoItems() []Item {
	return []Item{
		{ID: "types", Label: "Type checking"},
		{ID: "format", Label: "Formatting"},
		{ID: "tests", Label: "Tests"},
		{ID: "deploy", Label: "Deploy preview", Disabled: true},
	}
}

func TestMultiSelectFieldGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Checks"
	m.Description = "Run before commit"
	m.Help = "space toggles • enter submits"
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("types", "tests")
	m.SetSize(40, 8)
	m.Focus()

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestMultiSelectEmptyGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Checks"
	m.SetSize(30, 3)
	m.Focus()
	snaptest.Snap(t, m.View())
}

func TestInlineMultiSelectEmptyGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Checks"
	m.Inline = true
	m.Focus()
	m.SetSize(30, m.NaturalHeight())
	snaptest.Snap(t, m.View())
}

func TestInlineMultiSelectGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Checks"
	m.Description = "Run before commit"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("format")
	m.Focus()
	m.SetSize(44, m.NaturalHeight())

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineMultiSelectCollapsedGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Checks"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("types", "tests")
	m.SetSize(44, m.NaturalHeight())

	snaptest.Snap(t, m.View())
}

func TestInlineMultiSelectNarrowGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Checks"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("types")
	m.Focus()
	m.SetSize(20, m.NaturalHeight())

	snaptest.Snap(t, m.View())
}

func TestSpaceTogglesCursorAndEmitsChangedMsg(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)
	m.Focus()

	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := m.Selected(); got != 1 {
		t.Fatalf("cursor after down = %d, want 1", got)
	}
	m, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if cmd == nil {
		t.Fatal("space returned no command")
	}
	msg, ok := cmd().(ChangedMsg)
	if !ok {
		t.Fatalf("command emitted %T, want ChangedMsg", cmd())
	}
	if msg.ID != "format" || !msg.Checked || msg.Index != 1 {
		t.Fatalf("changed msg = %#v", msg)
	}
	if strings.Join(msg.Selected, ",") != "format" {
		t.Fatalf("selected ids = %v, want [format]", msg.Selected)
	}

	_, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	msg = cmd().(ChangedMsg)
	if msg.Checked || len(msg.Selected) != 0 {
		t.Fatalf("second toggle msg = %#v, want unchecked", msg)
	}
}

func TestNavigationSkipsDisabledItems(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)
	m.Focus()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	if got := m.Selected(); got != 2 {
		t.Fatalf("end cursor = %d, want 2 (last enabled)", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if got := m.Selected(); got != 2 {
		t.Fatalf("down past disabled cursor = %d, want 2", got)
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if got := m.Selected(); got != 1 {
		t.Fatalf("up cursor = %d, want 1", got)
	}
	m.Select(3)
	if got := m.Selected(); got != 1 {
		t.Fatalf("direct select of disabled item = %d, want unchanged 1", got)
	}
}

func TestEnterSubmitsCheckedIDs(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("tests", "types")
	m.SetSize(30, 6)
	m.Focus()

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("enter returned no command")
	}
	msg, ok := cmd().(SubmittedMsg)
	if !ok {
		t.Fatalf("command emitted %T, want SubmittedMsg", cmd())
	}
	if strings.Join(msg.Selected, ",") != "types,tests" {
		t.Fatalf("submitted ids = %v, want [types tests] in item order", msg.Selected)
	}
}

func TestSetItemsPreservesCheckedAndCursorByID(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("format")
	m.SelectID("tests")

	m.SetItems(
		Item{ID: "tests", Label: "Tests"},
		Item{ID: "format", Label: "Formatting"},
		Item{ID: "lint", Label: "Linting"},
	)
	if !m.Checked("format") || m.Checked("types") {
		t.Fatalf("checked state not preserved: %v", m.SelectedIDs())
	}
	if m.CursorItem().ID != "tests" {
		t.Fatalf("cursor = %q, want tests", m.CursorItem().ID)
	}
}

func TestSemanticActionsToggleByIDAndSubmit(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)

	next, cmd := m.Update(inspect.Invoke("toggle.tests"))
	m = next
	if cmd == nil {
		t.Fatal("semantic toggle returned no command")
	}
	if msg := cmd().(ChangedMsg); msg.ID != "tests" || !msg.Checked {
		t.Fatalf("semantic toggle msg = %#v", msg)
	}

	_, cmd = m.Update(inspect.Invoke(ActionSubmit))
	if cmd == nil {
		t.Fatal("semantic submit returned no command")
	}
	if msg := cmd().(SubmittedMsg); strings.Join(msg.Selected, ",") != "tests" {
		t.Fatalf("semantic submit = %v, want [tests]", msg.Selected)
	}
}

func TestInspectReportsCheckedIDs(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Checks"
	m.SetItems(demoItems()...)
	m.SetCheckedIDs("types", "tests")
	m.SetSize(30, 6)
	m.Focus()

	node := m.Inspect()
	if node.Kind != "multiselect" {
		t.Fatalf("kind = %q, want multiselect", node.Kind)
	}
	if node.Attributes["selected_ids"] != "types,tests" {
		t.Fatalf("selected_ids = %q", node.Attributes["selected_ids"])
	}
	if node.Attributes["selected_count"] != "2" {
		t.Fatalf("selected_count = %q", node.Attributes["selected_count"])
	}
}

func TestSizeModes(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	if got := layout.SizeModeOf(&m); got != layout.SizeBounded {
		t.Fatalf("stacked size mode = %d, want SizeBounded", got)
	}
	m.Inline = true
	if got := layout.SizeModeOf(&m); got != layout.SizeWidthBounded {
		t.Fatalf("inline size mode = %d, want SizeWidthBounded", got)
	}

	m.Focus()
	if got := m.NaturalHeight(); got != 7 {
		t.Fatalf("focused inline height = %d, want 7", got)
	}
	m.Blur()
	if got := m.NaturalHeight(); got != 2 {
		t.Fatalf("blurred inline height = %d, want 2", got)
	}
}

func TestStaysWithinAssignedBox(t *testing.T) {
	theme := tuiweave.Dark()
	for _, inline := range []bool{false, true} {
		for _, focused := range []bool{false, true} {
			for _, width := range []int{1, 2, 5, 9, 24} {
				for _, height := range []int{1, 2, 5, 8} {
					m := New(theme)
					m.Inline = inline
					m.Title = "Checks"
					m.Description = "Run before commit"
					m.Help = "space toggles"
					m.SetItems(demoItems()...)
					m.SetCheckedIDs("types")
					if focused {
						m.Focus()
					}
					m.SetSize(width, height)
					view := m.View()
					if view == "" {
						t.Fatalf("inline=%v focused=%v size=%dx%d rendered empty", inline, focused, width, height)
					}
					rows := strings.Split(view, "\n")
					if got := len(rows); got > height {
						t.Fatalf("inline=%v focused=%v size=%dx%d rendered %d rows", inline, focused, width, height, got)
					}
					for row, line := range rows {
						if got := ansi.StringWidth(line); got > width {
							t.Fatalf("inline=%v focused=%v size=%dx%d row %d width=%d: %q", inline, focused, width, height, row, got, line)
						}
					}
				}
			}
		}
	}
}
