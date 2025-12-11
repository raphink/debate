#!/bin/bash
# Generate PWA icons from SVG
# Requires: ImageMagick or rsvg-convert

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PUBLIC_DIR="$SCRIPT_DIR/public"
ICON_SVG="$PUBLIC_DIR/icon.svg"

echo "Generating PWA icons..."

# Check if rsvg-convert is available (preferred for SVG conversion)
if command -v rsvg-convert &> /dev/null; then
    echo "Using rsvg-convert..."
    rsvg-convert -w 192 -h 192 "$ICON_SVG" -o "$PUBLIC_DIR/icon-192.png"
    rsvg-convert -w 512 -h 512 "$ICON_SVG" -o "$PUBLIC_DIR/icon-512.png"
    rsvg-convert -w 32 -h 32 "$ICON_SVG" -o "$PUBLIC_DIR/favicon.ico"
    echo "✓ Generated icon-192.png, icon-512.png, favicon.ico"
elif command -v convert &> /dev/null; then
    echo "Using ImageMagick convert..."
    convert -background none -resize 192x192 "$ICON_SVG" "$PUBLIC_DIR/icon-192.png"
    convert -background none -resize 512x512 "$ICON_SVG" "$PUBLIC_DIR/icon-512.png"
    convert -background none -resize 32x32 "$ICON_SVG" "$PUBLIC_DIR/favicon.ico"
    echo "✓ Generated icon-192.png, icon-512.png, favicon.ico"
else
    echo "Error: Neither rsvg-convert nor ImageMagick convert found."
    echo "Please install one of:"
    echo "  - librsvg (brew install librsvg on macOS)"
    echo "  - ImageMagick (brew install imagemagick on macOS)"
    exit 1
fi

echo "✓ PWA icons generated successfully!"
