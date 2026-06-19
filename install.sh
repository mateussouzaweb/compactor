#!/bin/bash

set -euo pipefail

# Ensure the system provides apt for package installation
if ! command -v apt >/dev/null 2>&1; then
  echo "[ERROR] This installer requires apt package manager."
  exit 1
fi

# Determine the target architecture for the download
case "$(uname -m)" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  aarch64|arm64)
    ARCH="arm64"
    ;;
  *)
    echo "[ERROR] Unsupported architecture: $(uname -m)"
    exit 1
    ;;
esac

# Start process with an informational message
echo "[INFO] Installing compactor and dependencies..."

# Build URL for the architecture-specific binary
REPOSITORY="https://github.com/mateussouzaweb/compactor"
BINARY="$REPOSITORY/releases/latest/download/compactor-$ARCH"

# Prefer user-local binaries destination when available
BINARIES="/usr/local/bin"
DESTINATION="$BINARIES/compactor"

if [ -d "$HOME/.local/bin" ]; then
  BINARIES="$HOME/.local/bin"
  DESTINATION="$BINARIES/compactor"
fi

# Ensure the binary directory is on PATH for the current session
export PATH="$PATH:$BINARIES"

# Install Node.js and npm if missing
if ! command -v node >/dev/null 2>&1; then
  echo "[INFO] Installing NodeJS and NPM..."
  apt install -y nodejs npm
fi

# Install libjpeg if missing
if ! /sbin/ldconfig -p | grep -q libjpeg; then
  echo "[INFO] Installing libjpeg..."
  apt install -y libjpeg-progs
fi

# Install required npm packages globally
PKGS=(gifsicle jpegoptim-bin cwebp-bin optipng-bin sass-embedded terser typescript svgo html-minifier rollup)
INSTALLED=$(npm list -g)

echo "[INFO] Checking NPM packages..."
for PKG in "${PKGS[@]}"; do
  if ! echo "$INSTALLED" | grep -q "$PKG"; then
    npm install -s -g "$PKG"
  fi
  echo "[INFO] $PKG - OK"
done

# Download the Compactor binary
echo "[INFO] Downloading compactor $ARCH..."
curl -fsSL "$BINARY" -o "$DESTINATION"
chmod +x "$DESTINATION"

# Print final information about the installation
echo "[INFO] Compactor installed at $DESTINATION"
echo "[INFO] $(compactor --version)"
