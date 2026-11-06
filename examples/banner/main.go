// Command banner renders a figlet banner with a gradient. It is a static
// example of the banner utility; interactive apps can place Render's output in
// any view.
package main

import (
	"fmt"
	"os"

	"github.com/rigerc/tuiweave/banner"
)

func render() (string, error) {
	return banner.Render(banner.Config{
		Text:     "tuiweave",
		Font:     "standard",
		Width:    80,
		Parser:   "terminal-color",
		Gradient: &banner.GradientOcean,
	})
}

func main() {
	out, err := render()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}
