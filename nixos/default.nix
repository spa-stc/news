{ self, inputs, lib, ... }:
let
  mkSystem = extraModules: inputs.nixpkgs.lib.nixosSystem {
    system = "x86_64-linux";
    modules = with inputs; [
      # Agenix Secrets
      agenix.nixosModules.default


      # Other Common Modules
      ./common
    ] ++ extraModules;
  };

  mkSystemDeployNode = configPath: hostname: sshUser: {
    inherit hostname;
    fastConnection = true;
    profiles = {
      system = {
        inherit sshUser;
        path = inputs.deploy-rs.lib.x86_64-linux.activate.nixos configPath;
        user = "root";
      };
    };
  };
in
{

  flake.nixosConfigurations = {
    backend = mkSystem [ ./hosts/backend ];
  };

  flake.deploy.nodes = {
    backend = mkSystemDeployNode self.nixosConfigurations.backend "backend" "root";
  };
}
