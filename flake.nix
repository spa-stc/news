{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    # Framework For Defining "Flake Modules" with imports
    flake-parts.url = "github:hercules-ci/flake-parts";

    # Deployments
    deploy-rs.url = "github:serokell/deploy-rs";

    # Secrets 
    agenix.url = "github:ryantm/agenix";
  };

  outputs = inputs@{ self, ... }: inputs.flake-parts.lib.mkFlake { inherit inputs; } {
    systems = [ "x86_64-linux" "aarch64-linux" ];

    imports = [ ./nixos ];

    flake.nixosModules.pb = import ./pocketbase.nix self;

    perSystem = { config, self', inputs', pkgs, system, ... }: {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ dive just nodePackages_latest.wrangler ];
      };

      devShells.backend = pkgs.mkShell {
        buildInputs = with pkgs; [ just inputs.agenix.packages.${system}.default inputs.deploy-rs.packages.${system}.default ];
      };

      packages = rec {
        bin = pkgs.buildGo122Module {
          name = "newsletter";
          vendorHash = "sha256-RGJP71gBUF4W4j9QywA6ye6sP+fS9kYnNzcWf9vLqgM=";
          src = ./.;
          subPackages = [ "cmd/newsletter" ];
        };

        bin-docker = pkgs.dockerTools.buildLayeredImage {
          name = "newsletter";
          tag = "latest";
          config = {
            Cmd = [ "${self'.packages.default}/bin/newsletter" "serve" "--http" "0.0.0.0:8090" "--dir" "/pb_data" ];
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              "PRODUCTION=true"
            ];
          };
        };

        frontend = pkgs.buildNpmPackage {
          name = "newsletter-frontend";
          src = ./frontend;
          npmDepsHash = "sha256-fNgoQB2Gp21s2McjuQhWfaxoMJxathil0Fvrz9XBMtA=";

          installPhase = ''
            cp -r dist $out
          '';

          VITE_PB_URL = "https://newspb.stpaulacademy.tech";
        };

        default = bin;
      };
    };
  };
}
