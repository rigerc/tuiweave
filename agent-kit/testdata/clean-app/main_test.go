package main

import "github.com/ishan5ain/tuiweave/snaptest"

func auditCoverageNames() {
	_ = snaptest.Snap
	_ = snaptest.SnapCells
	_ = snaptest.RunScenario
}
