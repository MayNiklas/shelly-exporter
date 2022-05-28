{
  description = "prometheus exporter for shelly plug s";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:

    {
      checks."x86_64-linux".example =
        self.packages."x86_64-linux".shelly_exporter;
    } //

    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};

      in
      rec {

        formatter = pkgs.nixpkgs-fmt;

        packages = flake-utils.lib.flattenTree rec {

          shelly_exporter = with pkgs.python39Packages;
            pkgs.python39Packages.buildPythonPackage rec {
              pname = "shelly_exporter";
              version = "1.0.0";

              propagatedBuildInputs = [ requests ];
              src = self;

              installCheckPhase = ''
                runHook preCheck
                $out/bin/shelly_exporter
                runHook postCheck
              '';

              meta = with pkgs.lib; {
                description = "prometheus exporter for shelly plug s";
                homepage =
                  "https://github.com/MayNiklas/shelly-plug-s-prometheus-exporter";
                platforms = platforms.unix;
                maintainers = with maintainers; [ mayniklas ];
              };
            };

        };
        defaultPackage = packages.shelly_exporter;
      });
}
