#!/bin/sh
# shellcheck shell=dash

#
# This script downloads the latest stable release of fablo.sh
# fablo.sh is used to generate the configuration files for the
# Hyperledger Fabric network.
#

url=https://github.com/hyperledger-labs/fablo/releases/download/1.2.0/fablo.sh
script=$(basename $url)
curl -fsSL "$url" -o "$script"
chmod +x "$script"
