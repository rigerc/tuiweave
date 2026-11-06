package main

import (
	"testing"

	"github.com/rigerc/tuiweave/snaptest"
)

func TestBannerExampleGolden(t *testing.T) {
	out, err := render()
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	snaptest.Snap(t, out)
}
