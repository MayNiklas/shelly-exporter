{
  description = "prometheus exporter for shelly plug s";

  inputs = { flake-utils.url = "github:numtide/flake-utils"; };

  outputs = { self, nixpkgs, flake-utils, ... }:

    {
      nixosModules.default = self.nixosModules.shelly_exporter;
      nixosModules.shelly_exporter = { lib, pkgs, config, ... }:
        with lib;

        let cfg = config.services.shelly-exporter;
        in
        {

          options.services.shelly-exporter = {

            enable = mkEnableOption "shelly-exporter";

            user = mkOption {
              type = types.str;
              default = "shelly-exporter";
              description = "User account under which s3photoalbum runs.";
            };

            group = mkOption {
              type = types.str;
              default = "shelly-exporter";
              description = "Group under which s3photoalbum runs.";
            };

          };

          config = mkIf cfg.enable {

            systemd.services.shelly-exporter = {
              description = "A shelly metrics exporter";
              wantedBy = [ "multi-user.target" ];
              serviceConfig = mkMerge [{
                User = cfg.user;
                Group = cfg.group;
                ExecStart = "${self.packages."${pkgs.system}".shelly_exporter}/bin/shelly-plug-s-prometheus-exporter";
                Restart = "on-failure";
              }];
            };

            users.users = mkIf (cfg.user == "shelly-exporter") {
              shelly-exporter = {
                isSystemUser = true;
                group = cfg.group;
                description = "shelly-exporter system user";
              };
            };

            users.groups =
              mkIf (cfg.group == "shelly-exporter") { shelly-exporter = { }; };

          };
          meta = { maintainers = with lib.maintainers; [ mayniklas ]; };
        };
    }

    //

    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};

      in
      rec {

        formatter = pkgs.nixpkgs-fmt;
        packages = flake-utils.lib.flattenTree rec {

          default = shelly_exporter;

          shelly_exporter = pkgs.buildGoModule rec {
            pname = "shelly-plug-s-prometheus-exporter";
            version = "1.0.0";
            src = ./.;
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
