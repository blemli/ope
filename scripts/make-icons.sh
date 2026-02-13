#!/bin/bash
# Generate icon files from ope.svg
# Requires: rsvg-convert (librsvg), iconutil (macOS), ImageMagick (for .ico)
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
SVG="$ROOT_DIR/ope.svg"
OUT_DIR="$ROOT_DIR/packaging/icons"

mkdir -p "$OUT_DIR"

# Generate PNG
echo "Generating PNG..."
rsvg-convert -w 1024 -h 1024 "$SVG" -o "$OUT_DIR/icon.png"

# Generate macOS .icns
if [ "$(uname)" = "Darwin" ]; then
    echo "Generating macOS .icns..."
    ICONSET="$OUT_DIR/icon.iconset"
    mkdir -p "$ICONSET"

    for SIZE in 16 32 64 128 256 512; do
        rsvg-convert -w "$SIZE" -h "$SIZE" "$SVG" -o "$ICONSET/icon_${SIZE}x${SIZE}.png"
        DOUBLE=$((SIZE * 2))
        rsvg-convert -w "$DOUBLE" -h "$DOUBLE" "$SVG" -o "$ICONSET/icon_${SIZE}x${SIZE}@2x.png"
    done

    iconutil -c icns -o "$OUT_DIR/icon.icns" "$ICONSET"
    rm -rf "$ICONSET"
    echo "Created $OUT_DIR/icon.icns"
fi

# Generate Windows .ico (requires ImageMagick)
if command -v magick &>/dev/null || command -v convert &>/dev/null; then
    echo "Generating Windows .ico..."
    CONVERT_CMD="convert"
    if command -v magick &>/dev/null; then
        CONVERT_CMD="magick convert"
    fi

    TMPDIR_ICO=$(mktemp -d)
    for SIZE in 16 32 48 64 128 256; do
        rsvg-convert -w "$SIZE" -h "$SIZE" "$SVG" -o "$TMPDIR_ICO/${SIZE}.png"
    done

    $CONVERT_CMD "$TMPDIR_ICO/16.png" "$TMPDIR_ICO/32.png" "$TMPDIR_ICO/48.png" \
        "$TMPDIR_ICO/64.png" "$TMPDIR_ICO/128.png" "$TMPDIR_ICO/256.png" \
        "$OUT_DIR/icon.ico"
    rm -rf "$TMPDIR_ICO"
    echo "Created $OUT_DIR/icon.ico"
else
    echo "Skipping .ico (ImageMagick not found)"
fi

echo "Done! Icons in $OUT_DIR/"
ls -la "$OUT_DIR/"
