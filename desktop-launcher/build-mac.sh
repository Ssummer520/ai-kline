#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
APP_NAME="AI KLine Web"
APP_DIR="$ROOT_DIR/$APP_NAME.app"
DMG_PATH="$ROOT_DIR/ai-kline-macos-arm64.dmg"
PKG_PATH="$ROOT_DIR/ai-kline-macos-arm64.pkg"
BIN_DIR="$APP_DIR/Contents/MacOS"
RES_DIR="$APP_DIR/Contents/Resources"
MAC_BIN="$BIN_DIR/$APP_NAME"

rm -rf "$SCRIPT_DIR/dist" "$APP_DIR" "$DMG_PATH" "$PKG_PATH"
mkdir -p "$SCRIPT_DIR" "$BIN_DIR" "$RES_DIR"

cp -R "$ROOT_DIR/client-vue/dist" "$SCRIPT_DIR/dist"
cp "$SCRIPT_DIR/Info.plist" "$APP_DIR/Contents/Info.plist"

cd "$SCRIPT_DIR"
GOCACHE="${GOCACHE:-/tmp/ai-kline-go-cache}" GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o "$MAC_BIN" .
chmod +x "$MAC_BIN"

pkgbuild \
  --component "$APP_DIR" \
  --install-location /Applications \
  "$PKG_PATH"

hdiutil create \
  -volname "$APP_NAME" \
  -srcfolder "$APP_DIR" \
  -ov \
  -format UDZO \
  "$DMG_PATH"
