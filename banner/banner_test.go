package banner

import (
	"strings"
	"testing"

	"github.com/ishan5ain/tuiweave/snaptest"
)

func TestRenderPlain(t *testing.T) {
	got, err := Render(Config{
		Text:   "tuiweave",
		Font:   "standard",
		Width:  80,
		Parser: "terminal",
		Color:  "red", // must not add ANSI when the plain parser is explicit
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if got == "" {
		t.Fatal("Render returned empty output")
	}
	if strings.Contains(got, "\x1b") {
		t.Fatalf("plain parser emitted ANSI: %q", got)
	}
	if strings.Count(got, "\n") < 2 {
		t.Fatalf("banner is not multi-line: %q", got)
	}
}

func TestRenderGolden(t *testing.T) {
	got, err := Render(Config{
		Text:   "tuiweave",
		Font:   "standard",
		Width:  80,
		Parser: "terminal",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	snaptest.Snap(t, got)
}

func TestRenderColor(t *testing.T) {
	got, err := Render(Config{
		Text:   "hi",
		Font:   "standard",
		Width:  40,
		Parser: "terminal-color",
		Color:  "#FF0000",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(got, "\x1b[38;2;255;0;0m") {
		t.Fatalf("truecolor prefix missing: %q", got)
	}
}

func TestRenderGradient(t *testing.T) {
	got, err := Render(Config{
		Text:     "hi",
		Font:     "standard",
		Width:    40,
		Parser:   "terminal-color",
		Gradient: &GradientOcean,
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(got, "\x1b[38;2;") {
		t.Fatalf("gradient truecolor prefix missing: %q", got)
	}
}

func TestRenderHTML(t *testing.T) {
	got, err := Render(Config{
		Text:     "hi",
		Font:     "standard",
		Width:    40,
		Parser:   "html",
		Gradient: &GradientNeon,
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(got, "<code>") || !strings.Contains(got, "<span style='color: rgb(") {
		t.Fatalf("html output missing wrapper or color span: %q", got)
	}
}

func TestMutuallyExclusiveColors(t *testing.T) {
	_, err := Render(Config{
		Text:           "hi",
		Font:           "standard",
		Color:          "red",
		Gradient:       &GradientOcean,
		RandomGradient: true,
	})
	if err == nil {
		t.Fatal("expected mutually exclusive error")
	}
}

func TestInvalidColor(t *testing.T) {
	_, err := Render(Config{Text: "hi", Font: "standard", Color: "not-a-color"})
	if err == nil {
		t.Fatal("expected invalid color error")
	}
}

func TestInvalidGradient(t *testing.T) {
	bad := Gradient{Name: "bad", Colors: []string{"not-hex"}}
	_, err := Render(Config{Text: "hi", Font: "standard", Gradient: &bad})
	if err == nil {
		t.Fatal("expected invalid gradient error")
	}
}

func TestInvalidParser(t *testing.T) {
	_, err := Render(Config{Text: "hi", Font: "standard", Parser: "nope"})
	if err == nil {
		t.Fatal("expected invalid parser error")
	}
}

func TestGradientCatalog(t *testing.T) {
	gradients := Gradients()
	if len(gradients) == 0 {
		t.Fatal("no predefined gradients")
	}
	gradients[0].Name = "changed"
	gradients[0].Colors[0] = "000000"
	if got, _ := GradientForName("sunset"); got.Name != "sunset" || got.Colors[0] == "000000" {
		t.Fatal("Gradients returned shared mutable data")
	}
	if got, ok := GradientForName("OCEAN"); !ok || got.Name != "ocean" {
		t.Fatalf("GradientForName(OCEAN) = (%#v, %v)", got, ok)
	}
	if _, ok := GradientForName("missing"); ok {
		t.Fatal("GradientForName accepted an unknown name")
	}
}

func TestRandomFontIsNotEmpty(t *testing.T) {
	if got := RandomFont(); got == "" {
		t.Fatal("RandomFont returned empty string")
	}
}

func TestRandomGradientIsKnown(t *testing.T) {
	got := RandomGradient()
	if _, ok := GradientForName(got.Name); !ok {
		t.Fatalf("RandomGradient returned unknown gradient %q", got.Name)
	}
}

func TestEmptyText(t *testing.T) {
	_, err := Render(Config{Font: "standard", Parser: "terminal"})
	if err == nil {
		t.Fatal("expected empty text error")
	}
}

func TestEmptyGradient(t *testing.T) {
	bad := Gradient{Name: "empty"}
	_, err := Render(Config{Text: "hi", Font: "standard", Gradient: &bad})
	if err == nil {
		t.Fatal("expected empty gradient error")
	}
}
