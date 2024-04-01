{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    # Framework For Defining "Flake Modules" with imports
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs@{ self, ... }: inputs.flake-parts.lib.mkFlake { inherit inputs; } {
    systems = [ "x86_64-linux" "aarch64-linux" ];

    perSystem = { config, self', inputs', pkgs, system, ... }: {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [ dive just ];
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
            Cmd = [ "${self'.packages.default}/bin/newsletter" "serve" ];
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              "PRODUCTION=true"
            ];
          };
        };

        frontend = pkgs.buildNpmPackage {
          name = "newsletter-frontend";
          src = ./frontend;
          npmDepsHash = "sha256-DRFOrx3ONYED47NG7eB9s30ZZ0xpPJJnG5v9L53V7EU=";

          installPhase = ''
            cp -r dist $out
          '';
        };

        default = bin;
      };
    };
  };
}
