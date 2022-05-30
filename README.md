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
3. Scrape exporter with Prometheus:
```nix
{ lib, pkgs, config, ... }:
let
  shellyTargets = [
    "192.168.0.2"
    "192.168.0.3"
    "192.168.0.4"
    "192.168.0.5"
    "192.168.0.6"
    "192.168.0.7"
  ];
in {
  services.prometheus = {
    enable = true;
    scrapeConfigs = [{
      job_name = "shelly";
      scrape_interval = "10s";
      metrics_path = "/probe";
      static_configs = [{ targets = shellyTargets; }];
      relabel_configs = [
        {
          source_labels = [ "__address__" ];
          target_label = "__param_target";
        }
        {
          source_labels = [ "__param_target" ];
          target_label = "instance";
        }
        {
          target_label = "__address__";
          replacement = "127.0.0.1:8080";
        }
      ];
    }];
  };
}
```

### Libaries used:
- https://github.com/prometheus/client_golang

### Libary documentation:
- https://pkg.go.dev/github.com/prometheus/client_golang/prometheus

### API documentation:
- https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status
