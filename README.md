# shelly_exporter
A Shelly Plug S Prometeus exporter written in golang.

## How to execute

Metrics will be exposed on: http://localhost:8080/probe?target=<shelly_ip>

### Nix / NixOS
This repository contains a `flake.nix` file.
```sh
nix run .#shelly_exporter
```

### Libaries used:
- https://github.com/prometheus/client_golang

### API documentation:
- https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status
