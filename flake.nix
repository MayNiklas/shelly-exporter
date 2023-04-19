{
  description = "prometheus exporter for shelly plug s";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:

    {
      nixosModules.default = self.nixosModules.shelly-exporter;
      nixosModules.shelly-exporter = { lib, pkgs, config, ... }:
        with lib;

        let cfg = config.services.shelly-exporter;
        in
        {

          options.services.shelly-exporter = {

            enable = mkEnableOption "shelly-exporter";

            configure-prometheus = mkEnableOption "enable shelly-exporter in prometheus";

            port = mkOption {
              type = types.str;
              default = "8080";
              description = "Port under which shelly-exporter is accessible.";
            };

            listen = mkOption {
              type = types.str;
              default = "localhost";
              example = "127.0.0.1";
              description = "Address under which shelly-exporter is accessible.";
            };

            targets = mkOption {
              type = types.listOf types.str;
              default = [ "http://192.168.15.2" ];
              example = [ "http://192.168.15.2" ];
              description = "Shelly's to monitor";
            };

            user = mkOption {
              type = types.str;
              default = "shelly-exporter";
              description = "User account under which shelly-exporter runs.";
            };

            group = mkOption {
              type = types.str;
              default = "shelly-exporter";
              description = "Group under which shelly-exporter runs.";
            };

          };

          config = mkIf cfg.enable {

            systemd.services.shelly-exporter = {
              description = "A shelly metrics exporter";
              wantedBy = [ "multi-user.target" ];
              serviceConfig = mkMerge [{
                User = cfg.user;
                Group = cfg.group;
                ExecStart = "${self.packages."${pkgs.system}".shelly-exporter}/bin/shelly-exporter -port ${cfg.port} -listen ${cfg.listen}";
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

            services.prometheus = mkIf cfg.configure-prometheus {
              scrapeConfigs = [{
                job_name = "shelly";
                scrape_interval = "15s";
                metrics_path = "/probe";
                static_configs = [{ targets = cfg.targets; }];
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
                    replacement =
                      "127.0.0.1:${cfg.port}";
                  }
                ];
              }];
            };

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

          default = shelly-exporter;

          shelly-exporter = pkgs.buildGoModule rec {
            pname = "shelly-exporter";
            version = "1.0.0";
            src = self;
            vendorSha256 =
              "sha256-Oq4Y5BJOft4fNSnoqExUeupWbKk6YOK068MaUMgdeuk=";
            installCheckPhase = ''
              runHook preCheck
              $out/bin/shelly-exporter -h
              runHook postCheck
            '';
            doCheck = true;
            meta = with pkgs.lib; {
              description = "prometheus exporter";
              homepage =
                "https://github.com/MayNiklas/shelly-exporter";
              platforms = platforms.unix;
              maintainers = with maintainers; [ mayniklas ];
            };
          };

          # Documenation for this feature: https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/docker/examples.nix
          # nix build .#docker-image
          # docker load < result
          docker-image = pkgs.dockerTools.buildLayeredImage {

            name = "mayniki/shelly-exporter";
            tag = "latest";

            # Using "now" breaks reproducibility, the resulting image
            # will not be identical, but contain a useful timestamp:
            created = "now";

            contents = [
              self.packages."${pkgs.system}".shelly-exporter
            ];

            config.Cmd = [ "${self.packages."${pkgs.system}".shelly-exporter}/bin/shelly-exporter" ];

            config.ExposedPorts = {
              "8080/tcp" = { };
            };

            config.Env = [
              "port=8080"
              "listen=0.0.0.0"
            ];

          };

        };
      });
}
