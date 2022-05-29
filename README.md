# shelly_exporter
A Shelly Plug S Prometeus exporter written in golang.

WORK IN PROGRESS!

## How to execute

### Nix / NixOS
This repository contains a `flake.nix` file.
```sh
nix run .#shelly_exporter
```

### Libaries used:
- https://github.com/prometheus/client_golang

### API documentation:
- https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status
