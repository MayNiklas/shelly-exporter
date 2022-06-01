# shelly_exporter
A Shelly Plug S Prometeus exporter written in golang.

[![Go](https://github.com/MayNiklas/shelly-plug-s-prometheus-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/MayNiklas/shelly-plug-s-prometheus-exporter/actions/workflows/go.yml)

## How to execute for development purposes

Metrics will be exposed on: http://localhost:8080/probe?target=http://<shelly_ip>

### Nix / NixOS
This repository contains a `flake.nix` file.
```sh
# run the package
nix run .#shelly_exporter

# build the package
nix build .#shelly_exporter
```

### General
Make sure [golang](https://go.dev) is installed.
```sh
# run application
go run .

# build application
go build

# execute tests
go test -v ./...

# show test coverage
go test -covermode=count -coverpkg=./... -coverprofile cover.out -v ./...
go tool cover -html cover.out -o cover.html
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
    "http://192.168.0.2"
    "http://192.168.0.3"
    "http://192.168.0.4"
    "http://192.168.0.5"
    "http://192.168.0.6"
    "http://192.168.0.7"
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
