#!/bin/sh
# CTXbank 1-Click Installer for macOS & Linux
# Usage: curl -fsSL https://raw.githubusercontent.com/hamziCodes/CTXbank/main/scripts/install.sh | sh

set -e

REPO="hamziCodes/CTXbank"
INSTALL_DIR="/usr/local/bin"

echo "========================================"
echo "  CTXbank Fast Installer (macOS/Linux)"
echo "========================================"

# 1. Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *) echo "Unsupported operating system: $OS" && exit 1 ;;
esac

# 2. Detect Architecture
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac

echo "[1/4] Detected system: $OS / $ARCH"

# 3. Query Latest Release
echo "[2/4] Querying latest release from GitHub ($REPO)..."
LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "v0.1.0")
if [ -z "$LATEST_TAG" ]; then
  LATEST_TAG="v0.1.0"
fi

ASSET="ctx-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${LATEST_TAG}/${ASSET}"

# 4. Download and Extract
TEMP_DIR=$(mktemp -d)
echo "[3/4] Downloading $ASSET ($LATEST_TAG)..."
curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$ASSET"

echo "      Extracting binaries..."
tar -xzf "$TEMP_DIR/$ASSET" -C "$TEMP_DIR"

# Check write permissions for /usr/local/bin
if [ -w "$INSTALL_DIR" ]; then
  cp "$TEMP_DIR/ctx" "$INSTALL_DIR/ctx"
  cp "$TEMP_DIR/ctx-mcp" "$INSTALL_DIR/ctx-mcp"
  chmod +x "$INSTALL_DIR/ctx" "$INSTALL_DIR/ctx-mcp"
else
  USER_BIN="$HOME/.ctxbank/bin"
  mkdir -p "$USER_BIN"
  cp "$TEMP_DIR/ctx" "$USER_BIN/ctx"
  cp "$TEMP_DIR/ctx-mcp" "$USER_BIN/ctx-mcp"
  chmod +x "$USER_BIN/ctx" "$USER_BIN/ctx-mcp"
  INSTALL_DIR="$USER_BIN"

  # Add to PATH in rc files
  if [ -f "$HOME/.zshrc" ] && ! grep -q ".ctxbank/bin" "$HOME/.zshrc"; then
    echo 'export PATH="$HOME/.ctxbank/bin:$PATH"' >> "$HOME/.zshrc"
  fi
  if [ -f "$HOME/.bashrc" ] && ! grep -q ".ctxbank/bin" "$HOME/.bashrc"; then
    echo 'export PATH="$HOME/.ctxbank/bin:$PATH"' >> "$HOME/.bashrc"
  fi
fi

rm -rf "$TEMP_DIR"

echo "[4/4] Binaries installed to $INSTALL_DIR"
echo ""
echo "+-------------------------------------------------------------+"
echo "|  CTXbank installed successfully!                            |"
echo "+-------------------------------------------------------------+"
echo ""
echo "Run 'ctx --help' in your terminal to get started."
