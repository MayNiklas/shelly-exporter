# shelly-exporter

A Shelly power metrics exporter written in golang.
Currently only tested for Shelly Plug S.

[![Go](https://github.com/MayNiklas/shelly-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/MayNiklas/shelly-exporter/actions/workflows/go.yml)

## Available metrics

Name     | Description
---------|------------
shelly_power_current | Current real AC power being drawn, in Watts
shelly_power_total | Total energy consumed by the attached electrical appliance in Watt-minute
shelly_temperature | internal device temperature in °C
shelly_update_available | Info whether newer firmware version is available
shelly_uptime | Seconds elapsed since boot

All metrics include the following labels:

* device IP
* device name
* device hostname

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
  inputs.shelly-exporter = {
    url = "github:MayNiklas/shelly-exporter";
    inputs = { nixpkgs.follows = "nixpkgs"; };
  };
}
```

2. Enable the service & the prometheus scraper in your configuration:

```nix
{ shelly-exporter, ... }: {

  imports = [ shelly-exporter.nixosModules.default ];

  services.shelly-exporter = {
    enable = true;
    port = "8080";
    listen = "localhost";
    user = "shelly-exporter";
    group = "shelly-exporter";

    configure-prometheus = true;
    targets = [
      "http://192.168.0.2"
      "http://192.168.0.3"
      "http://192.168.0.4"
      "http://192.168.0.5"
      "http://192.168.0.6"
      "http://192.168.0.7"
    ];
  };
}
```

### Libaries used

- https://github.com/prometheus/client_golang

### Libary documentation

- https://pkg.go.dev/github.com/prometheus/client_golang/prometheus

### API documentation

- https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-status
