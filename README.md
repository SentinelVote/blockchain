# Hyperledger Fabric Blockchain

## Running the network

You should be in the root repository directory, i.e. `~/path/to/blockchain`.

```sh
./setup-fablo.sh
./fablo.sh recreate
```

---

## Prerequisites

- bash
- curl
- docker
- docker-compose
- fablo
- jq

## Troubleshooting

### Potential problems with docker-compose

Recent installations of docker have deprecated `docker-compose` in favor of `docker compose`,
but shell scripts that are used in Hyperledger Fabric and Fablo (the generation tool) still use `docker-compose`.

In your terminal, check if you have `docker-compose`:

```
command -V docker-compose
```

If you get `docker-compose not found`, see the next section for a patch.

#### Patching docker-compose

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

### Root permissions for docker socket

Run the following your terminal:

```sh
sudo usermod -aG docker "$USER"
```

Then, log out and log back in.

### Docker daemon

Run the following your terminal:

```sh
sudo systemctl enable --now docker.service docker.socket
```

This starts docker now and also starts docker by default on boot.
