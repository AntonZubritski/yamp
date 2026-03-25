#!/bin/sh
set -e

REPO="AntonZubritski/yamp"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
    linux|darwin) ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

BINARY="yamp-${OS}-${ARCH}"
if [ "$OS" = "windows" ]; then
    BINARY="${BINARY}.exe"
fi

URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo ""
echo "  yamp installer"
echo "  OS: ${OS}  ARCH: ${ARCH}"
echo "  Downloading ${BINARY}..."
echo ""

TMP=$(mktemp)
if command -v curl > /dev/null; then
    curl -fSL -o "$TMP" "$URL"
elif command -v wget > /dev/null; then
    wget -qO "$TMP" "$URL"
else
    echo "Error: curl or wget required" >&2; exit 1
fi

chmod +x "$TMP"

if [ "$OS" = "windows" ]; then
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
    mv "$TMP" "${INSTALL_DIR}/yamp.exe"
    echo ""
    echo "  Installed to ${INSTALL_DIR}/yamp.exe"
    echo "  Add to PATH: export PATH=\"\$HOME/bin:\$PATH\""
elif [ -w "$INSTALL_DIR" ]; then
    mv "$TMP" "${INSTALL_DIR}/yamp"
    echo "  Installed to ${INSTALL_DIR}/yamp"
else
    sudo mv "$TMP" "${INSTALL_DIR}/yamp"
    echo "  Installed to ${INSTALL_DIR}/yamp"
fi

echo ""
echo "  Done! Run: yamp"
echo ""
