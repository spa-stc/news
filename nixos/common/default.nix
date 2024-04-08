{ config, lib, pkgs, ... }: {
  imports = [
    ./tailscale.nix
  ];


  config = {
    users.groups.stc = { };

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
        # Foehammer
        "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAP0Hp1XEQ0Fx4rpuBuH7YusGU1f6s3Q4Sb1O1U0oKWA foehammer@disroot.org"

        # Github Actions Runner
        "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBwc6ZjZBYzW2bxiwIvJ1W99RQU08rcap7Fj1q8c3pVc runner@github.io"
      ];
    };
  };
}
