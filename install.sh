#!/bin/bash

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)          ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

URL="https://github.com/insigmo/jetreset/releases/latest/download/jetreset-${OS}-${ARCH}"

echo "Downloading: $URL"
curl -L "$URL" -o jetreset || { echo "Download failed"; exit 1; }
chmod +x jetreset

# Снять карантин на macOS (Gatekeeper блокирует неподписанные бинарники)
if [ "$OS" = "darwin" ]; then
  xattr -d com.apple.quarantine ./jetreset 2>/dev/null
  echo "Quarantine attribute removed (macOS)"
fi

echo 'Run exec file for reseting: ./jetreset'
