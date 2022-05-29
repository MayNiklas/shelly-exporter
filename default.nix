{ lib, pkgs, config, ... }:
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
        ExecStart = "${pkgs.shelly_exporter}/bin/shelly_exporter";
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
}
