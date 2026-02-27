#!/usr/bin/env bash
set -euo pipefail

REPO="Cloverhound/cupi-cli"
BINARY="cupi"
SKILL_NAME="cupi"
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

# Install binary
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

# ── Skill installation wizard ─────────────────────────────────────────────────
echo ""

SKILL_URL="https://raw.githubusercontent.com/${REPO}/${TAG}/skill/SKILL.md"
SKILL_TMP="$(mktemp)"

if ! curl -fsSL "$SKILL_URL" -o "$SKILL_TMP" 2>/dev/null; then
  echo "Warning: Could not download skill file. Skipping skill installation."
  rm -f "$SKILL_TMP"
  echo ""
  echo "Run 'cupi --help' to get started."
  exit 0
fi

echo "==> Install skill file for AI coding assistants?"
echo "    Press Enter to install at the shown path, type a new path to override, or 'n' to skip."
echo ""

# Format: "Display Name|default/path/SKILL.md"
TOOLS=(
  "Claude Code|${HOME}/.claude/skills/${SKILL_NAME}/SKILL.md"
  "Codex|${HOME}/.codex/skills/${SKILL_NAME}/SKILL.md"
  "Gemini|${HOME}/.gemini/skills/${SKILL_NAME}/SKILL.md"
)

for entry in "${TOOLS[@]}"; do
  TOOL_NAME="${entry%%|*}"
  DEFAULT_PATH="${entry#*|}"

  echo "  ${TOOL_NAME}"
  printf "    Path : %s\n    > " "$DEFAULT_PATH"

  # Read from /dev/tty so this works when the script is piped via curl | bash
  read -r RESPONSE </dev/tty

  case "$RESPONSE" in
    n|N|no|No|NO|skip)
      echo "    Skipped."
      ;;
    "")
      mkdir -p "$(dirname "$DEFAULT_PATH")"
      cp "$SKILL_TMP" "$DEFAULT_PATH"
      echo "    Installed → ${DEFAULT_PATH}"
      ;;
    *)
      mkdir -p "$(dirname "$RESPONSE")"
      cp "$SKILL_TMP" "$RESPONSE"
      echo "    Installed → ${RESPONSE}"
      ;;
  esac
  echo ""
done

rm -f "$SKILL_TMP"

echo "Run 'cupi --help' to get started."
