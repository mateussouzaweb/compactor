#!/bin/bash

set -e
set -u

echo "Installing compactor..."

# Install dependencies
if [[ $(node -v) == "" ]]; then
  sudo apt install -y nodejs npm
fi

if [[ $(/sbin/ldconfig -p | grep libjpeg) == "" ]]; then
  sudo apt install -y libjpeg-progs
fi

# Install NPM packages
p=( gifsicle jpegoptim-bin cwebp-bin optipng-bin sass terser typescript svgo html-minifier )
for i in "${p[@]}"; do
  npm list -g | grep "$i" || npm install --quiet -g "$i"
done

# Install compactor
REPOSITORY="https://github.com/mateussouzaweb/compactor/releases/download/latest"
sudo wget $REPOSITORY/compactor -O /usr/local/bin/compactor
sudo chmod +x /usr/local/bin/compactor

echo "Compactor installed!"