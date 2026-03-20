#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)

rm -rf "$SCRIPT_DIR/dist"
mkdir -p "$SCRIPT_DIR"
cp -R "$ROOT_DIR/client-vue/dist" "$SCRIPT_DIR/dist"

cd "$SCRIPT_DIR"
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o "$ROOT_DIR/ai-kline-web.exe" .
