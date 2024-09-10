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
          overlays = [];
        };

        packages = {
          default = pkgs.buildGoModule {
            name = "spa-newsletter";

            src = ./.;

            vendorHash = "sha256-1ofVvgiH9zf5x8CLuFjEpJFwadelZ8GzbgZcCJFjnCk=";
            doCheck = false;
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "spa-newsletter";

            config = {
              Cmd = ["${self'.packages.default}/bin/newsletter"];
              Env = [
                "NEWSLETTER_PUBLIC_DIR=${self'.packages.public}"
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              ];
            };
          };

          public = pkgs.stdenv.mkDerivation {
            name = "public-dir";

            buildInputs = with pkgs; [tailwindcss];

            src = ./.;

            buildPhase = ''
              tailwindcss -i public/tailwind_in.css -o public/assets/main.css -c tailwind.config.js --minify
              rm -rf public/tailwind_in.css
            '';

            installPhase = ''
              mkdir $out
              cp -r ./public/* $out
            '';
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            just
            dive

            # Frontend.
            tailwindcss

            # Golang
            go
            go-tools
            gotools
            gopls
            golangci-lint
            go-migrate
            delve

            #DB
            sqlc
            postgresql
          ];
        };

        formatter = inputs'.alejandra.packages.default;
      };
    });
}
