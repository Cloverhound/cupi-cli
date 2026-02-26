#!/usr/bin/env bash
# install-local.sh — build from source and install cupi + Claude Code skill
set -euo pipefail

BINARY="cupi-cli"
SKILL_NAME="cupi-cli"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
SKILL_DIR="$HOME/.claude/skills/$SKILL_NAME"

# Must be run from the project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "==> Building $BINARY..."
go build -o "$BINARY" .

echo "==> Installing binary to $INSTALL_DIR/$BINARY"
mkdir -p "$INSTALL_DIR"
cp "$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"
rm "$BINARY"

echo "==> Installing Claude Code skill to $SKILL_DIR/"
mkdir -p "$SKILL_DIR"
cp skill/SKILL.md "$SKILL_DIR/SKILL.md"

echo ""
echo "Done."
echo ""
echo "  Binary : $INSTALL_DIR/$BINARY"
echo "  Skill  : $SKILL_DIR/SKILL.md"
echo ""

# PATH check
if ! command -v "$BINARY" &>/dev/null; then
  echo "NOTE: $INSTALL_DIR is not in your PATH."
  echo "Add this to your shell profile (~/.zshrc or ~/.bashrc):"
  echo ""
  echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
  echo ""
fi

echo "Run 'cupi-cli --help' to get started."
