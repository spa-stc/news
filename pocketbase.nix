self:
{ config, lib, pkgs, ... }:
let
  cfg = config.stc.services.pb;
in
{
  options.stc.services.pb = {
    enable = lib.mkEnableOption "Enable Pocketbase";

    port = lib.mkOption {
      type = lib.types.port;
      default = 3000;
      description = lib.mdDoc "The port to listen on";
    };

    datapath = lib.mkOption {
      type = lib.types.path;
      default = "/srv/pocketbase/pb_data";
      description = lib.mdDoc "Where to place pocketbase data.";
    };
  };

  config = lib.mkIf cfg.enable {
    age.secrets = {
      pbkey = {
        file = ./nixos/secrets/pb_key.age;
        owner = "pocketbase";
        mode = "600";
        group = "stc";
      };
    };


    services.caddy = {
      enable = true;
      email = "stc@students.spa.edu";

      virtualHosts."newspb.stpaulacademy.tech" = {
        extraConfig = ''
          request_body {
                 max_size 10MB
          }
          reverse_proxy 0.0.0.0:${toString cfg.port} {
            transport http {
              read_timeout 360s
            }
          }
        '';
      };
    };

    networking.firewall.allowedTCPPorts = [ 80 443 ];

    users.users.pocketbase = {
      createHome = true;
      group = "stc";
      description = "User for the pocketbase service";
      isSystemUser = true;
      home = "/srv/pocketbase";
    };

    systemd.services.pocketbase = {
      wantedBy = [ "multi-user.target" ];

      serviceConfig = {
        User = "pocketbase";
        Group = "stc";
        Restart = "on-failure";
        WorkingDirectory = "/srv/pocketbase";
        RestartSec = "30s";
        Type = "exec";
        StandardOutput = "/var/pocketbase.log";
        StandardError = "/var/pocketbase-error.log";
      };

      script =
        let
          pb = self.packages.${pkgs.system}.bin;
          keypath = config.age.secrets.pbkey.path;
        in
        ''
          export ENCKEY=(cat ${keypath})
          export PRODUCTION=true
          exec ${pb}/bin/newsletter serve --http 0.0.0.0:${toString cfg.port} --dir ${cfg.datapath} --encryptionEnv $ENCKEY --https=""
        '';
    };
  };
}
