# shelly_exporter
A Shelly Plug S Prometeus exporter written in golang.

## How to execute for development purposes

Metrics will be exposed on: http://localhost:8080/probe?target=<shelly_ip>

### Nix / NixOS
This repository contains a `flake.nix` file.
```sh
nix run .#shelly_exporter
```

## How to install

### NixOS
1. Add this repository to your `flake.nix`:
```nix
{
  inputs.shelly-prometheus-exporter = {
    url = "github:MayNiklas/shelly-plug-s-prometheus-exporter";
    inputs = { nixpkgs.follows = "nixpkgs"; };
  };
}
```
2. Enable the service in your configuration:
```nix
{ shelly-prometheus-exporter, ... }: {

  imports = [ shelly-prometheus-exporter.nixosModules.default ];

  services.shelly-exporter = {
    enable = true;
    port = "8080";
    listen = "localhost";
    user = "shelly-exporter";
    group = "shelly-exporter";
  };
}
```

### Libaries used:
- https://github.com/prometheus/client_golang

### Libary documentation:
- https://pkg.go.dev/github.com/prometheus/client_golang/prometheus

### API documentation:
- https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status
