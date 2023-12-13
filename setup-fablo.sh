#!/bin/sh

url=https://github.com/hyperledger-labs/fablo/releases/download/1.2.0/fablo.sh
script=$(basename $url)
curl -fsSL "$url" -o "$script"
chmod +x "$script"
