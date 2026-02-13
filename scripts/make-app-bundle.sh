#!/bin/bash
# Create Ope.app bundle from a compiled binary
# Usage: ./make-app-bundle.sh <binary> <version> [output-dir]
set -e

BINARY="${1:?Usage: make-app-bundle.sh <binary> <version> [output-dir]}"
VERSION="${2:?Usage: make-app-bundle.sh <binary> <version> [output-dir]}"
OUTPUT_DIR="${3:-.}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

APP_DIR="$OUTPUT_DIR/Ope.app"
CONTENTS="$APP_DIR/Contents"
MACOS="$CONTENTS/MacOS"
RESOURCES="$CONTENTS/Resources"

# Clean previous bundle
rm -rf "$APP_DIR"

# Create structure
mkdir -p "$MACOS" "$RESOURCES"

# Copy binary
cp "$BINARY" "$MACOS/ope"
chmod +x "$MACOS/ope"

# Write Info.plist with version
sed "s/VERSION/$VERSION/g" "$ROOT_DIR/packaging/macos/Info.plist" > "$CONTENTS/Info.plist"

# Copy icon if available
ICON="$ROOT_DIR/packaging/icons/icon.icns"
if [ -f "$ICON" ]; then
    cp "$ICON" "$RESOURCES/icon.icns"
fi

echo "Created $APP_DIR"
