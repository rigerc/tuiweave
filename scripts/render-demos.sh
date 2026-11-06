#!/usr/bin/env bash
# Re-render every README demo locally. The GIFs are committed assets; MP4s
# stay in demos/vhs/outputs/ for inspection.
set -euo pipefail
cd "$(dirname "$0")/.."

mkdir -p docs/demos demos/vhs/outputs
for tape in demos/vhs/*.tape; do
  name=$(basename "$tape" .tape)
  echo "rendering $name"
  vhs "$tape"
  test -s "docs/demos/$name.gif"
  test -s "demos/vhs/outputs/$name.mp4"
  ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=p=0 \
    "demos/vhs/outputs/$name.mp4"
done
