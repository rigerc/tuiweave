# VHS demos

Each `.tape` records one runnable example for the README gallery. The committed
GIFs live in [`docs/demos`](../../docs/demos); MP4s are local/CI evidence and
are ignored.

Render every demo locally:

```sh
./scripts/render-demos.sh
```

Render one demo:

```sh
vhs demos/vhs/ops.tape
```

Tapes run from the repository root so `go run ./examples/...` resolves. Keep
the terminal dimensions, theme, input sequence, and waits deterministic. CI
renders the same tapes with `charmbracelet/vhs-action@v2.1.1` and uploads the
GIF/MP4 as artifacts; the checked-in GIFs are the README assets.
