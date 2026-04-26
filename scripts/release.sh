#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="$ROOT_DIR/dist"
BIN_DIR="$DIST_DIR/bin"
COMP_DIR="$DIST_DIR/completions"

mkdir -p "$BIN_DIR" "$COMP_DIR"

build() {
  local goos="$1"
  local goarch="$2"
  local ext=""
  if [[ "$goos" == "windows" ]]; then
    ext=".exe"
  fi
  local out="$BIN_DIR/gotree_${goos}_${goarch}${ext}"
  echo "Building $out"
  GOOS="$goos" GOARCH="$goarch" go build -o "$out" ./src/cmd/gotree
}

build linux amd64
build darwin amd64
build darwin arm64
build windows amd64

echo "Generating completions..."
go run ./src/cmd/gotree --completion bash > "$COMP_DIR/gotree.bash"
go run ./src/cmd/gotree --completion zsh > "$COMP_DIR/_gotree"
go run ./src/cmd/gotree --completion fish > "$COMP_DIR/gotree.fish"
go run ./src/cmd/gotree --completion powershell > "$COMP_DIR/gotree.ps1"

echo "Done. Artifacts in $DIST_DIR"
