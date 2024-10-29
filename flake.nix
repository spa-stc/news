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

        packages = let
          inherit (pkgs) lib;

          isIncluded = includePaths: path: type: let
            pathStr = toString path;
            matchesPath = includePath: let
              includePath' = toString includePath;
            in
              lib.hasPrefix includePath' pathStr || lib.hasPrefix pathStr includePath';
          in
            builtins.any matchesPath includePaths;

          filterSourceFiles = root: includePaths: let
            absolutePaths =
              map (
                path:
                  if lib.hasPrefix "/" (toString path)
                  then path
                  else root + "/${toString path}"
              )
              includePaths;
          in
            builtins.filterSource (path: type: isIncluded absolutePaths path type) root;
        in {
          default = pkgs.buildGoModule {
            name = "spa-newsletter";

            src = filterSourceFiles ./. [
              ./app
              ./cmd
              ./config
              ./cron
              ./db
              ./go.mod
              ./go.sum
              ./main.go
              ./resource
              ./util
              ./web
            ];

            vendorHash = "sha256-M53e6aFThhScGoDp4RqM5JQFhLHfOm5nV4RssfQ1Mdc=";
            doCheck = false;
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/spa-stc/news";
            tag = "latest";

            config = {
              Entrypoint = pkgs.writeShellScript "newsletter-entrypoint" ''
                set -eo pipefail

                 if [[ -z $NEWSLETTER_DATABASE_URL ]]; then
                 	echo "NEWSLETTER_DATABASE_URL must be set in the environment."
                   	exit
                 fi

                ${pkgs.go-migrate}/bin/migrate -path="${self'.packages.migrations}" -database="$NEWSLETTER_DATABASE_URL" up

                ${self'.packages.default}/bin/newsletter
              '';
              Env = [
                "NEWSLETTER_PUBLIC_DIR=${self'.packages.public}"
                "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              ];
            };
          };

          public = pkgs.stdenv.mkDerivation {
            name = "spa-newsletter-assets";

            buildInputs = with pkgs; [tailwindcss];

            src = ./public;

            buildPhase = ''
              tailwindcss -i tailwind_in.css -o assets/main.css -c tailwind.config.js --minify
              rm -rf tailwind_in.css tailwind.config.js
            '';

            installPhase = ''
              mkdir $out
              cp -r ./* $out
            '';
          };

          migrations = pkgs.stdenv.mkDerivation {
            name = "spa-newsletter-migrations";
            src = ./migrations;

            installPhase = ''
              mkdir $out
              cp -r ./*.up.sql $out
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
