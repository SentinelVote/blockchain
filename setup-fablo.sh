#!/bin/sh
# shellcheck shell=dash

#
# This script downloads the latest stable release of fablo.sh
# https://github.com/hyperledger-labs/fablo/
#
# TODO: Comments on switch case and dockerized build.
#
# Fablo generates the config files for the Hyperledger Fabric network.
#

default () {
url=https://github.com/hyperledger-labs/fablo/releases/download/1.2.0/fablo.sh
script=$(basename $url)
curl -fsSL "$url" -o "$script"
chmod +x "$script"
}

fablo_rest () {
_dockerfile=$(mktemp)
cat <<EOF > "$_dockerfile"
FROM softwaremill/fablo-rest:0.1.0
RUN sed -i.bak 's/app.use(express_1.default.json({ type: () => "json" }))/app.use(express_1.default.json({ type: () => "json", limit: "500mb" }))/g' index.js
EOF
printf '%s' "$_dockerfile"
}

entrypoint () {
default
sleep 60
./fablo.sh recreate
tail -f /dev/null # Leave it running.
}

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
  entrypoint)
    entrypoint
  ;;
  fablo_rest)
    fablo_rest
  ;;
  *)
    default
  ;;
esac
}
main "$@"
