{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    alejandra = {
      url = "github:kamadorueda/alejandra/3.0.0";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ {self, ...}:
    inputs.flake-parts.lib.mkFlake {inherit inputs;} (toplevel @ {withSystem, ...}: {
      systems = ["aarch64-darwin" "aarch64-linux" "x86_64-linux"];

      perSystem = {
        config,
        self',
        inputs',
        pkgs,
        system,
        ...
      }: {
        _module.args.pkgs = import inputs.nixpkgs {
          localSystem = system;
          config = {
            allowUnfree = true;
            allowAliases = true;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            just
            air

            # Frontend
            tailwindcss

            # GO
            go
            go-tools
            gotools
            gopls
            delve
          ];
        };

        packages.default = pkgs.buildGoModule {
          name = "newsletter";

          src = ./app;

          # Inject tailwind files.
          preBuild = ''
            ${pkgs.tailwindcss}/bin/tailwindcss -i public/in.css -o public/static/tailwind.min.css --minify
          '';

          vendorHash = "sha256-JiXYo11ciagANmmeaNLTyoIg93OENeSf7Dfhaxzyveo=";
        };

        packages.docker = pkgs.dockerTools.buildLayeredImage {
          name = "ghcr.io/spa-stc/newsletter";

          config = {
            Cmd = ["${self'.packages.default}/bin/newsletter"];
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
          };
        };

        formatter = inputs'.alejandra.packages.default;
      };
    });
}
