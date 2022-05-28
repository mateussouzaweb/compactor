#!/bin/bash

set -e
set -u

VERSION="v0.1.2"
REPOSITORY="https://github.com/mateussouzaweb/compactor"
BINARY="${REPOSITORY}/releases/download/${VERSION}/compactor"
BINARIES="$HOME/.local/bin"

echo "[INFO] Installing compactor and dependencies..."

# Install dependencies
if [[ $(node -v) == "" ]]; then
  apt install -y nodejs npm
fi

if [[ $(/sbin/ldconfig -p | grep libjpeg) == "" ]]; then
  apt install -y libjpeg-progs
fi

# Install NPM packages
p=( gifsicle jpegoptim-bin cwebp-bin optipng-bin sass terser typescript svgo html-minifier )
for i in "${p[@]}"; do
  npm list -g | grep "$i" || npm install --silence -g "$i"
done

# Install compactor
mkdir -p $BINARIES
wget $BINARY -O $BINARIES/compactor
chmod +x $BINARIES/compactor

echo "[INFO] Compactor ${VERSION} installed!"