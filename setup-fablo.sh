#!/bin/sh
# shellcheck shell=dash

fail () { printf %s\\n "Error: $1" >&2 ; exit 1 ; }

# Fetch the latest stable release of https://github.com/hyperledger-labs/fablo/
default () {
command -v curl || fail "Please install curl."

url=https://github.com/hyperledger-labs/fablo/releases/download/1.2.0/fablo.sh
script=$(basename $url)
curl -fsSL "$url" -o "$script"
chmod +x "$script"
}

# Deploy the blockchain with 9 peer nodes instead of 2 (default).
# This option is used in the production server deployment script.
patch_production() {
sed -i 's/instances: [2-9]/instances: 9/' fablo-config.yaml
}

main () {
test -z "$1" && default && exit 0
case "$1" in
  prod)
    patch_production
  ;;
  production)
    patch_production
  ;;
  *)
    default
  ;;
esac
}
main "$@"
