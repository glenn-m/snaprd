{
  description = "snaprd daemon";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    devshell.url = "github:numtide/devshell";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, flake-utils, devshell, nixpkgs }:
    flake-utils.lib.eachDefaultSystem (system: {
      packages = let pkgs = import nixpkgs { inherit system; };
      in { snaprd = pkgs.callPackage ./snaprd.nix; };
      nixosModules = {
        snaprd = { config, lib, pkgs, ... }:
          with lib;
          let
            cfg = config.services.snaprd;
            mkConfigFile =
              pkgs.writeText "snaprd.yml" (builtins.toJSON cfg.configuration);
          in {
            options = {
              services.snaprd = {
                enable = lib.mkEnableOption "snaprd daemon";

                package = mkOption {
                  type = types.package;
                  default = self.packages.${pkgs.system}.snaprd;
                  defaultText = literalExpression "pkgs.snaprd";
                  description = lib.mdDoc ''
                    Package that should be used for snaprd.
                  '';
                };

                port = mkOption {
                  type = types.int;
                  default = 3000;
                  description = lib.mdDoc ''
                    Port the Prometheus metrics will be exposed on.
                  '';
                };

                metricsPath = mkOption {
                  type = types.str;
                  default = "/metrics";
                  description = lib.mdDoc ''
                    Path the Prometheus metrics will be exposed on.
                  '';
                };

                configuration = mkOption {
                  type = types.nullOr types.attrs;
                  default = null;
                  description = lib.mdDoc ''
                    Snaprd configuration as nix attribute set.
                  '';
                };
              };
            };
            config = lib.mkIf cfg.enable {
              #environment.systemPackages =
              #  [ self.packages.${pkgs.system}.snaprd ];
              systemd.services.snaprd = {
                description = "snaprd daemon";
                wantedBy = [ "multi-user.target" ];
                after = "network.target";
                serviceConfig = {
                  ExecStart = "${cfg.package}/bin/snaprd"
                    + "--configFile=${mkConfigFile}"
                    + "--metricsPort=${cfg.port}"
                    + "--metricsPath=${cfg.metricsPath}";
                  Restart = "always";
                  WorkDirectory = "/tmp";
                };
              };
            };
          };
      };
      devShell = let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ devshell.overlay ];
        };
      in pkgs.devshell.mkShell {
        imports = [ (pkgs.devshell.importTOML ./devshell.toml) ];
      };
    });
}
