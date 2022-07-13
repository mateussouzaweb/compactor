#!/bin/bash

set -e
set -u

VERSION="v0.1.2"
REPOSITORY="https://github.com/mateussouzaweb/compactor"
BINARY="${REPOSITORY}/releases/download/${VERSION}/compactor"
BINARIES="${HOME}/.local/bin"

echo "[INFO] Installing compactor and dependencies..."

# Install dependencies
if [[ $(node -v) == "" ]]; then
  echo "[INFO] Installing NodeJS and NPM..."
  apt install -y nodejs npm
fi

if [[ $(/sbin/ldconfig -p | grep libjpeg) == "" ]]; then
  echo "[INFO] Installing libjpeg..."
  apt install -y libjpeg-progs
fi

# Install NPM packages
PKGS=( gifsicle jpegoptim-bin cwebp-bin optipng-bin sass terser typescript svgo html-minifier )
INSTALLED=$(npm list -g)

echo "[INFO] Checking NPM packages..."
for PKG in "${PKGS[@]}"; do
  echo $INSTALLED | grep "${PKG}" || npm install -s -g "${PKG}"
  echo "[INFO] ${PKG} - OK"
done

# Make sure binaries path works
export PATH="$PATH:$BINARIES"

# Install compactor
mkdir -p $BINARIES

echo "[INFO] Downloading compactor..."
wget -q $BINARY -O $BINARIES/compactor
chmod +x $BINARIES/compactor

echo "[INFO] Compactor ${VERSION} installed!"