#!/usr/bin/env bash
set -euo pipefail

REPO="Cloverhound/cupi-cli"
BINARY="cupi"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)
    echo "Unsupported OS: $OS (use install.ps1 for Windows)" >&2
    exit 1
    ;;
esac

# Fetch latest release tag
echo "Fetching latest release..."
TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"

if [[ -z "$TAG" ]]; then
  echo "Failed to fetch latest release tag." >&2
  exit 1
fi

echo "Installing cupi ${TAG} (${OS}/${ARCH})..."

ARCHIVE="${BINARY}_${TAG#v}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${TAG}/${ARCHIVE}"
TMPDIR="$(mktemp -d)"

# Download and extract
curl -fsSL "$URL" -o "${TMPDIR}/${ARCHIVE}"
tar -xzf "${TMPDIR}/${ARCHIVE}" -C "$TMPDIR"

# Install
mkdir -p "$INSTALL_DIR"
mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
chmod +x "${INSTALL_DIR}/${BINARY}"
rm -rf "$TMPDIR"

echo "Installed to ${INSTALL_DIR}/${BINARY}"

# Check PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo ""
  echo "NOTE: ${INSTALL_DIR} is not in your PATH."
  echo "Add this to your shell profile:"
  echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
fi

echo ""
echo "Run 'cupi --help' to get started."
