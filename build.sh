#!/bin/bash
# CronPanel Build Script
# Requires: Go 1.18+
# Usage: ./build.sh

set -e

echo "🔨 Building CronPanel..."

# Create output directory
mkdir -p dist

# Build for Linux amd64
echo "📦 Building for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/cronpanel-linux-amd64 .

# Build for Linux arm64
echo "📦 Building for linux/arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/cronpanel-linux-arm64 .

echo ""
echo "✅ Build complete! Binaries in ./dist/"
echo ""
ls -lh dist/
echo ""
echo "Usage:"
echo "  chmod +x dist/cronpanel-linux-amd64"
echo "  ./dist/cronpanel-linux-amd64"
echo "  Then open: http://0.0.0.0:8899"
echo ""
echo "Optional: Set PORT environment variable to change port:"
echo "  PORT=9000 ./dist/cronpanel-linux-amd64"
