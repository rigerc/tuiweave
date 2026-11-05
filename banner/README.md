# banner

`banner` renders FIGlet ASCII art with optional ANSI, truecolor gradients, and
HTML output. It is a pure string utility for splash screens and CLI headers,
not a Bubble Tea component.

```go
out, err := banner.Render(banner.Config{
    Text:     "tuiweave",
    Font:     "standard",
    Width:    80,
    Parser:   "terminal-color",
    Gradient: &banner.GradientOcean,
})
```

- `Parser` supports `terminal-color` (default), `terminal` (plain text), and
  `html`. An explicit `terminal` parser suppresses ANSI even when a color or
  gradient is set.
- `Color` accepts ANSI names or hex values. `Gradient`, `RandomGradient`, and
  `RandomColor` are mutually exclusive with each other and with `Color`.
- `Gradients()` and `GradientForName(name)` expose the named presets.
- For deterministic output, set `Font` and one color source, or use the
  `terminal`/`html` parser. Empty `Font` and omitted colors select random
  values.
- `Width` defaults to 80; values below 20 clamp to 20.

Run the example with `go run ./examples/banner`.

## API highlights

- `Render(Config)` produces plain, ANSI/truecolor, or HTML banner output.
- `Config` controls text, font, width, justification, direction, colors, and parser.
- `Gradient` and the `Gradient*` presets provide truecolor stops.
- `Gradients()` and `GradientForName` discover the named palette.
- `RandomFont`, `RandomGradient`, and `RandomColor` provide playful defaults.

## Related documentation

- [API reference](https://pkg.go.dev/github.com/ishan5ain/tuiweave/banner)
- [Banner recipe](../AGENTS.md#banner)
- [Figlet-go module](https://github.com/lsferreira42/figlet-go)

The figlet fonts are embedded by
[`figlet-go`](https://github.com/lsferreira42/figlet-go); see that module for
the FIGlet font licensing terms.
