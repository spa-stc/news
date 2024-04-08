{ config, lib, pkgs, ... }: {
  imports = [
    # Digital Ocean
    ../../hardware/do.nix
  ];

  stc.services = {
    pb.enable = true;
  };

  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = false;
    };
  };

  networking = {
    hostName = "backend";
    firewall = {
      enable = true;
      allowedTCPPorts = [ 22 ];
    };
  };

  swapDevices = [{
    device = "/var/lib/swap";
    size = 1024 * 1; # 1 GB Swap.
  }];

  system.stateVersion = "23.05";

  environment.systemPackages = with pkgs; [ vim ];

  virtualisation.digitalOcean = {
    setRootPassword = true;
  };

  stc.tailscale.enable = true;
} 
