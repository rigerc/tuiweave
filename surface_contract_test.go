package tuiweave_test

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"

	"github.com/rigerc/tuiweave"
	"github.com/rigerc/tuiweave/action"
	"github.com/rigerc/tuiweave/agentic/permission"
	"github.com/rigerc/tuiweave/autocomplete"
	"github.com/rigerc/tuiweave/button"
	"github.com/rigerc/tuiweave/dialog"
	"github.com/rigerc/tuiweave/frame"
	"github.com/rigerc/tuiweave/palette"
	"github.com/rigerc/tuiweave/progress"
	"github.com/rigerc/tuiweave/snaptest"
	"github.com/rigerc/tuiweave/statusbar"
	"github.com/rigerc/tuiweave/tabs"
	"github.com/rigerc/tuiweave/toggle"
	"github.com/rigerc/tuiweave/toolbar"
)

func TestInheritedRaisedSurfaceGolden(t *testing.T) {
	theme := tuiweave.Dark()
	theme.SurfaceRaised = lipgloss.NoColor{}

	nav := tabs.New(theme)
	nav.SetTabs(
		tabs.Tab{ID: "services", Label: "Services"},
		tabs.Tab{ID: "jobs", Label: "Jobs"},
	)
	nav.SetSize(28, 1)
	nav.Focus()

	bar := statusbar.New(theme)
	bar.SetLeft(
		statusbar.Segment{Text: "tuiweave", Kind: statusbar.KindAccent},
		statusbar.Segment{Text: " ready", Kind: statusbar.KindNormal},
	)
	bar.SetRight(statusbar.Segment{Text: "healthy", Kind: statusbar.KindSuccess})
	bar.SetSize(28, 1)

	completions := autocomplete.New(theme)
	completions.SetItems(
		autocomplete.Item{ID: "status", Value: "git status", Label: "git status", Description: "inspect"},
		autocomplete.Item{ID: "test", Value: "go test ./...", Label: "go test ./...", Description: "verify"},
	)
	completions.SetSize(28, 2)
	completions.Focus()

	actions := []action.Item{
		{ID: "open", Label: "Open", Description: "open workspace"},
		{ID: "refresh", Label: "Refresh", Description: "refresh state"},
	}

	tools := toolbar.New(theme)
	tools.SetItems(actions...)
	tools.SetSize(28, 1)
	tools.Focus()

	setting := toggle.New(theme)
	setting.SetLabel("Auto-refresh")
	setting.SetChecked(true)
	setting.SetSize(28, 1)

	task := progress.New(theme)
	task.SetLabel("Indexing")
	task.SetPercent(0.72)
	task.SetStatus(progress.StatusSuccess)
	task.SetSize(28, 1)

	commands := palette.New(theme)
	commands.SetItems(actions...)
	commands.SetSize(28, 3)
	commands.Focus()

	open := button.New(theme)
	open.SetLabel("Open workspace")
	open.SetSize(28, 1)

	confirm := dialog.New(theme)
	confirm.Title = "Delete workspace?"
	confirm.Body = "This cannot be undone."
	confirm.SetSize(32, 8)

	approval := permission.New(theme)
	approval.Title = "Run tests?"
	approval.Body = "The command reads this repository."
	approval.SetProvenance(permission.Provenance{Tool: "shell", Operation: "execute"})
	approval.SetSize(36, 12)

	view := strings.Join([]string{
		"panel",
		frame.Panel(theme, "Application-owned body", 28, frame.PanelOptions{Title: "Services", Padding: 1}),
		"badge",
		frame.Badge(theme, "healthy", frame.BadgeSuccess),
		"tabs",
		nav.View(),
		"statusbar",
		bar.View(),
		"autocomplete",
		completions.View(),
		"toolbar",
		tools.View(),
		"toggle",
		setting.View(),
		"progress",
		task.View(),
		"palette",
		commands.View(),
		"button",
		open.View(),
		"dialog",
		confirm.View(),
		"permission",
		approval.View(),
	}, "\n")

	snaptest.Snap(t, view)
	snaptest.SnapCells(t, view, snaptest.WithRoles(theme))
}
