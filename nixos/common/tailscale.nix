{ config, lib, pkgs, ... }:
let cfg = config.stc.tailscale;
in
{
  options.stc.tailscale = {
    enable = lib.mkEnableOption "Enables Tailscale";
  };

  config = lib.mkIf cfg.enable {
    age.secrets = {
      tskey.file = ../secrets/ts_key.age;
    };


    # Setup Tailscale With Auto-Provisioning.
    services.tailscale = {
      enable = true;
      authKeyFile = config.age.secrets.tskey.path;
      openFirewall = true;
    };


    # Tell the firewall to implicitly trust packets routed over Tailscale:
    networking.firewall.trustedInterfaces = [ "tailscale0" ];
  };
}

