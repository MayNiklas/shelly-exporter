{
  description = "prometheus exporter for shelly plug s";

  inputs = { flake-utils.url = "github:numtide/flake-utils"; };

  outputs = { self, nixpkgs, flake-utils, ... }:

    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};

      in
      rec {

        formatter = pkgs.nixpkgs-fmt;
        defaultPackage = packages.shelly_exporter;
        packages = flake-utils.lib.flattenTree rec {

          shelly_exporter = pkgs.buildGoModule rec {
            pname = "shelly_exporter";
            version = "1.0.0";
            src = ./.;
            subPackages = [ "cmd/shelly_exporter" ];
            vendorSha256 =
              "sha256-IBgntTSqgjvi6ddOLenB1rS+Pfs3MKZfn8OnAWUYgkk=";
            meta = with pkgs.lib; {
              description = "prometheus exporter for shelly plug s";
              homepage =
                "https://github.com/MayNiklas/shelly-plug-s-prometheus-exporter";
              platforms = platforms.unix;
              maintainers = with maintainers; [ mayniklas ];
            };
          };

        };
      });
}
