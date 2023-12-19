# Hyperledger Fabric Blockchain

## Prerequisites

- bash
- curl
- docker
- docker-compose
- fablo
- go 1.20 or earlier
- jq

## Potential problems with docker-compose

Recent distributions/installations of docker have deprecated `docker-compose` in favor of `docker compose`,
but shell scripts that are used in the Hyperledger Fabric project and Fablo (the generation tool) still use `docker-compose`.

In your terminal, check if you have `docker-compose`:

```
command -V docker-compose
```

If you get `docker-compose not found`, see the next section for a patch.

### Patching docker-compose

You can patch this docker-compose by doing the following in your terminal:

```sh
if [ -f /usr/bin/docker-compose ]; then
  printf %s\\n 'docker-compose already exists.'
else
  sudo printf '%s\n%s' '#!/bin/sh' 'docker compose "$@"' > /usr/bin/docker-compose
  chmod +x /usr/bin/docker-compose
  USER_CURRENT=$(whoami)
  sudo chown "$USER_CURRENT":"$USER_CURRENT" /usr/bin/docker-compose
  unset USER_CURRENT
fi
```

## Root permissions problems for docker socket

Run the following your terminal:

```sh
sudo usermod -aG docker "$USER"
```

Then, log out and log back in.

## Docker daemon

```sh
sudo systemctl enable --now docker.service docker.socket
```

This starts docker now, and docker starts by default on boot.

## Fablo

Fablo is the tool used to generate the Hyperledger Fabric network topology. \
Get the latest version of Fablo: run `./setup-fablo.sh` in your terminal.
