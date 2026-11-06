package reorder

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
		{ID: "design", Label: "Design"},
		{ID: "build", Label: "Build"},
		{ID: "test", Label: "Test"},
		{ID: "ship", Label: "Ship"},
	}
}

func TestReorderFieldGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Steps"
	m.Description = "Order the release pipeline"
	m.Help = "space grabs • shift+↑/↓ moves"
	m.SetItems(demoItems()...)
	m.Select(1)
	m.SetSize(40, 8)
	m.Focus()

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestReorderGrabbedGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Steps"
	m.SetItems(demoItems()...)
	m.Select(2)
	m.SetSize(40, 5)
	m.Focus()
	m.Grab()

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineReorderGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Steps"
	m.Description = "Order the release pipeline"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.Select(1)
	m.Focus()
	m.SetSize(48, m.NaturalHeight())

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineReorderGrabbedGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Steps"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.Focus()
	m.Grab()
	m.SetSize(48, m.NaturalHeight())

	snaptest.Snap(t, m.View())
	snaptest.SnapCells(t, m.View(), snaptest.WithRoles(tuiweave.Dark()))
}

func TestInlineReorderCollapsedGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Steps"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.SetSize(48, m.NaturalHeight())

	snaptest.Snap(t, m.View())
}

func TestInlineReorderNarrowGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Steps"
	m.Inline = true
	m.SetItems(demoItems()...)
	m.Select(2)
	m.Focus()
	m.SetSize(20, m.NaturalHeight())

	snaptest.Snap(t, m.View())
}

func TestMoveDownEmitsMovedMsg(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)
	m.Focus()

	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModShift})
	if cmd == nil {
		t.Fatal("shift+down returned no command")
	}
	msg, ok := cmd().(MovedMsg)
	if !ok {
		t.Fatalf("command emitted %T, want MovedMsg", cmd())
	}
	if msg.From != 0 || msg.To != 1 || msg.Item.ID != "design" {
		t.Fatalf("moved msg = %#v", msg)
	}
	if strings.Join(m.IDs(), ",") != "build,design,test,ship" {
		t.Fatalf("order = %v", m.IDs())
	}
	if got := m.Selected(); got != 1 {
		t.Fatalf("cursor = %d, want 1 (follows the item)", got)
	}
	if strings.Join(idsOf(msg.Items), ",") != "build,design,test,ship" {
		t.Fatalf("message order = %v", idsOf(msg.Items))
	}
}

func TestGrabMovesCursorItemAndDrops(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)
	m.Focus()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if !m.Grabbed() {
		t.Fatal("space did not grab the item")
	}
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd == nil {
		t.Fatal("down while grabbed returned no command")
	}
	if strings.Join(m.IDs(), ",") != "build,design,test,ship" {
		t.Fatalf("order after move = %v", m.IDs())
	}
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace})
	if m.Grabbed() {
		t.Fatal("space did not drop the item")
	}
	m, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if cmd != nil {
		t.Fatal("navigation after drop emitted a move command")
	}
	if got := m.Selected(); got != 2 {
		t.Fatalf("cursor after drop + down = %d, want 2", got)
	}
}

func TestMoveToClampsAndEdgesAreNoOps(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.Select(2)

	if cmd := m.MoveUp(); cmd == nil {
		t.Fatal("move up emitted no command")
	}
	if cmd := m.MoveUp(); cmd == nil {
		t.Fatal("second move up emitted no command")
	}
	if cmd := m.MoveUp(); cmd != nil {
		t.Fatal("move up at the top emitted a command")
	}
	if strings.Join(m.IDs(), ",") != "test,design,build,ship" {
		t.Fatalf("order = %v", m.IDs())
	}

	m.Select(1)
	m.MoveTo(99)
	if strings.Join(m.IDs(), ",") != "test,build,ship,design" {
		t.Fatalf("order after clamped move = %v", m.IDs())
	}
	if got := m.Selected(); got != 3 {
		t.Fatalf("cursor after clamped move = %d, want 3", got)
	}
}

func TestSetItemsPreservesCursorByID(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SelectID("test")

	m.SetItems(
		Item{ID: "ship", Label: "Ship"},
		Item{ID: "test", Label: "Test"},
		Item{ID: "design", Label: "Design"},
	)
	if got := m.SelectedItem().ID; got != "test" {
		t.Fatalf("cursor = %q, want test", got)
	}
}

func TestInspectReportsOrderAndGrabState(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Steps"
	m.SetItems(demoItems()...)
	m.Select(1)
	m.SetSize(30, 6)
	m.Focus()
	m.Grab()

	node := m.Inspect()
	if node.Kind != "reorder" {
		t.Fatalf("kind = %q, want reorder", node.Kind)
	}
	if node.Status != "grabbed" {
		t.Fatalf("status = %q, want grabbed", node.Status)
	}
	if node.Attributes["order"] != "design,build,test,ship" {
		t.Fatalf("order = %q", node.Attributes["order"])
	}
	if node.Attributes["selected_id"] != "build" {
		t.Fatalf("selected_id = %q", node.Attributes["selected_id"])
	}
}

func TestSemanticActionsMoveAndGrab(t *testing.T) {
	m := New(tuiweave.Dark())
	m.SetItems(demoItems()...)
	m.SetSize(30, 6)

	next, cmd := m.Update(inspect.Invoke(ActionGrab))
	m = next
	if cmd != nil || !m.Grabbed() {
		t.Fatalf("grab: cmd=%v grabbed=%v", cmd, m.Grabbed())
	}

	next, cmd = m.Update(inspect.Invoke(ActionMoveDown))
	m = next
	if cmd == nil {
		t.Fatal("semantic move-down returned no command")
	}
	if strings.Join(m.IDs(), ",") != "build,design,test,ship" {
		t.Fatalf("order = %v", m.IDs())
	}

	next, _ = m.Update(inspect.Invoke(ActionDrop))
	m = next
	if m.Grabbed() {
		t.Fatal("drop left the item grabbed")
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

func TestEmptyGolden(t *testing.T) {
	m := New(tuiweave.Dark())
	m.Title = "Steps"
	m.SetSize(30, 3)
	m.Focus()
	snaptest.Snap(t, m.View())
}

func TestInlineEmptyGolden(t *testing.T) {
	m := NewOn(tuiweave.Dark(), nil)
	m.Title = "Steps"
	m.Inline = true
	m.Focus()
	m.SetSize(30, m.NaturalHeight())
	snaptest.Snap(t, m.View())
}

func TestStaysWithinAssignedBox(t *testing.T) {
	theme := tuiweave.Dark()
	for _, inline := range []bool{false, true} {
		for _, grabbed := range []bool{false, true} {
			for _, width := range []int{1, 2, 5, 9, 24} {
				for _, height := range []int{1, 2, 5, 8} {
					m := New(theme)
					m.Inline = inline
					m.Title = "Steps"
					m.Description = "Order the release pipeline"
					m.Help = "space grabs"
					m.SetItems(demoItems()...)
					m.Focus()
					if grabbed {
						m.Grab()
					}
					m.SetSize(width, height)
					view := m.View()
					if view == "" {
						t.Fatalf("inline=%v grabbed=%v size=%dx%d rendered empty", inline, grabbed, width, height)
					}
					rows := strings.Split(view, "\n")
					if got := len(rows); got > height {
						t.Fatalf("inline=%v grabbed=%v size=%dx%d rendered %d rows", inline, grabbed, width, height, got)
					}
					for row, line := range rows {
						if got := ansi.StringWidth(line); got > width {
							t.Fatalf("inline=%v grabbed=%v size=%dx%d row %d width=%d: %q", inline, grabbed, width, height, row, got, line)
						}
					}
				}
			}
		}
	}
}

func idsOf(items []Item) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
