{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-parts.url = "github:hercules-ci/flake-parts";

    alejandra = {
      url = "github:kamadorueda/alejandra/3.0.0";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    naersk = {
      url = "github:nix-community/naersk";
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
      }: let
        rustToolchain = pkgs.rust-bin.fromRustupToolchainFile ./rust-toolchain.toml;

        naersk = inputs.naersk.lib.${system}.override {
          cargo = rustToolchain;
          rustc = rustToolchain;
        };

        buildInputs = with pkgs; [
          sqlite.dev
          openssl
        ];

        naitiveBuildInputs = with pkgs; [
          pkg-config
        ];
      in {
        _module.args.pkgs = import inputs.nixpkgs {
          localSystem = system;
          config = {
            allowUnfree = true;
            allowAliases = true;
          };
          overlays = [
            inputs.rust-overlay.overlays.default
          ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs;
            [
              just
              rustToolchain
              sqlx-cli
              cargo-expand
            ]
            ++ buildInputs
            ++ naitiveBuildInputs;
        };

        packages.default = naersk.buildPackage {
          name = "newsletter";
          src = ./.;

          inherit buildInputs naitiveBuildInputs;
        };

        packages.public-dir = pkgs.runCommand "public-dir" {src = ./public;} "mkdir -p $out; cp -r $src/* $out";

        packages.docker = pkgs.dockerTools.buildLayeredImage {
          name = "newsletter-docker";
          tag = "latest";

          config = {
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              "NEWSLETTER_PUBLIC=${self'.packages.public-dir}"
            ];
            Cmd = ["${self'.packages.default}"];
          };
        };

        formatter = inputs'.alejandra.packages.default;
      };
    });
}
