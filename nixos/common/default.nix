{ config, lib, pkgs, ... }: {
  imports = [
    ./tailscale.nix
  ];


  config = {
    nix = {
      extraOptions = ''
        experimental-features = nix-command flakes
      '';

      gc = {
        automatic = true;
        dates = "weekly";
        options = "--delete-older-than 7d";
      };
    };

    nixpkgs.config = {
      allowUnfree = true;
    };

    users.users.root = {
      openssh.authorizedKeys.keys = [
        # Foehammer.
        "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAP0Hp1XEQ0Fx4rpuBuH7YusGU1f6s3Q4Sb1O1U0oKWA foehammer@disroot.org"
      ];
    };
  };
}
