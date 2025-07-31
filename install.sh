#!/usr/bin/env bash
set -e

# Project configuration - modify these for different projects
REPO_OWNER="sourceradar"
REPO_NAME="outline"
BINARY_NAME="outline"
HELP_TEXT="Try running: outline --help\nOr for MCP mode: outline --mcp"

# Colors for pretty output
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
  echo -e "${BLUE}==>${NC} $1"
}

print_success() {
  echo -e "${GREEN}==>${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}==>${NC} $1"
}

print_error() {
  echo -e "${RED}==>${NC} $1"
}

# Check if curl is installed
if ! command -v curl &> /dev/null; then
  print_error "curl is required but not installed. Please install curl and try again."
  exit 1
fi

# Create ~/.local/bin if it doesn't exist
INSTALL_DIR="$HOME/.local/bin"
if [ ! -d "$INSTALL_DIR" ]; then
  print_step "Creating directory $INSTALL_DIR"
  mkdir -p "$INSTALL_DIR"
fi

# Get OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize architecture names
if [ "$ARCH" = "x86_64" ] || [ "$ARCH" = "amd64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
else
  print_error "Unsupported architecture: $ARCH"
  exit 1
fi

# Convert OS names to match our release naming
if [ "$OS" = "darwin" ]; then
  OS_NAME="darwin"
elif [ "$OS" = "linux" ]; then
  OS_NAME="linux"
else
  print_error "Unsupported OS: $OS (only Linux and macOS are supported)"
  exit 1
fi

# Form the expected archive name based on our release workflow
ARCHIVE_NAME="${BINARY_NAME}-${OS_NAME}-${ARCH}.tar.gz"

# Get the download URL for the appropriate archive
print_step "Fetching latest release information..."
LATEST_RELEASE_JSON=$(curl -s "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest")
LATEST_RELEASE_URL=$(echo "$LATEST_RELEASE_JSON" | grep -o "\"browser_download_url\":[[:space:]]*\"[^\"]*${ARCHIVE_NAME}\"" | cut -d '"' -f 4)

if [ -z "$LATEST_RELEASE_URL" ]; then
  print_error "Could not find a release for your platform ($OS_NAME, $ARCH)"
  print_error "Available releases:"
  echo "$LATEST_RELEASE_JSON" | grep -o "\"browser_download_url\":[[:space:]]*\"[^\"]*\"" | cut -d '"' -f 4 | sed 's/.*\///g'
  exit 1
fi

print_step "Downloading from: $LATEST_RELEASE_URL"

# Create a temporary directory for the download
TMP_DIR=$(mktemp -d)
TMP_ARCHIVE="$TMP_DIR/$ARCHIVE_NAME"

# Download the archive
curl -L -o "$TMP_ARCHIVE" "$LATEST_RELEASE_URL"

# Extract the binary
print_step "Extracting $ARCHIVE_NAME"
tar -xzf "$TMP_ARCHIVE" -C "$TMP_DIR"

# The binary should be named ${BINARY_NAME}-${OS_NAME}-${ARCH} inside the archive
EXTRACTED_BINARY="$TMP_DIR/${BINARY_NAME}-${OS_NAME}-${ARCH}"
if [ ! -f "$EXTRACTED_BINARY" ]; then
  print_error "Expected binary not found in archive: $EXTRACTED_BINARY"
  print_error "Archive contents:"
  tar -tzf "$TMP_ARCHIVE"
  exit 1
fi

mv "$EXTRACTED_BINARY" "$INSTALL_DIR/$BINARY_NAME"

# Clean up the temporary directory
rm -rf "$TMP_DIR"

# Make the binary executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

print_success "Downloaded and installed ${BINARY_NAME} to $INSTALL_DIR/$BINARY_NAME"

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
  print_warning "$INSTALL_DIR is not in your PATH"

  # Determine shell and provide appropriate command
  SHELL_NAME="$(basename "$SHELL")"
  case "$SHELL_NAME" in
    bash)
      print_step "Run this command to add to your PATH:"
      echo "echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc && source ~/.bashrc"
      ;;
    zsh)
      print_step "Run this command to add to your PATH:"
      echo "echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.zshrc && source ~/.zshrc"
      ;;
    fish)
      print_step "Run this command to add to your PATH:"
      echo "fish_add_path $INSTALL_DIR && source ~/.config/fish/config.fish"
      ;;
    *)
      print_step "Run this command to add to your PATH:"
      echo "echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.profile && source ~/.profile"
      ;;
  esac
fi

print_success "Installation complete!"
echo -e "$HELP_TEXT"
