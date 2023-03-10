{
  description = "snaprd daemon";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    devshell.url = "github:numtide/devshell";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, flake-utils, devshell, nixpkgs }:
    # Create system-specific outputs for the standard Nix systems
    # https://github.com/numtide/flake-utils/blob/main/default.nix#L3-L9
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; };
      in rec {
        defaultPackage = packages.snaprd;
        packages.snaprd = with pkgs;
          buildGoModule rec {
            pname = "snaprd";
            version = "0.2.2";

            src = fetchFromGitHub {
              owner = "glenn-m";
              repo = pname;
              rev = "v${version}";
              sha256 = "sha256-1BB9JH3QBHH4DAWwsx3DkP0Oq1MeFeTyuz+V5PvzrZc=";
            };

            vendorSha256 =
              "sha256-cyMLyihBQF1sILnUsfBRcDm2XABw8PwY6W+nC/jqGuo=";

            meta = with lib; {
              description =
                "Daemon that runs snaprd on a schedule and surfaces Prometheus metrics";
              homepage = "https://github.com/glenn-m/snaprd";
              platforms = platforms.linux;
              license = licenses.mit;
            };
          };

        nixosModules.default = self.nixosModules.${system}.snaprd;
        nixosModules.snaprd = { config, lib, pkgs, ... }:
          with lib;
          let
            cfg = config.services.snaprd;
            mkConfigFile =
              pkgs.writeText "snaprd.yml" (builtins.toJSON cfg.configuration);
            mkSnapraidConfigFile =
              pkgs.writeText "snapraid.conf" cfg.snapraidConfig;
          in {
            options = {
              services.snaprd = {
                enable = lib.mkEnableOption "snaprd daemon";

                package = mkOption {
                  type = types.package;
                  default = self.packages.${system}.snaprd;
                  defaultText = literalExpression "pkgs.snaprd";
                  description = lib.mdDoc ''
                    Package that should be used for snaprd.
                  '';
                };

                port = mkOption {
                  type = types.int;
                  default = 9086;
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

                snapraidConfig = mkOption {
                  type = types.lines;
                  default = "";
                  description = lib.mdDoc ''
                    Contents of the snapraid configuration file
                  '';
                };
              };
            };
            config = lib.mkIf cfg.enable {
              systemd.services.snaprd = {
                description = "snaprd daemon";
                wantedBy = [ "multi-user.target" ];
                after = [ "network.target" ];
                environment = {
                  SNAPRD_CONFIG_FILE = "${mkConfigFile}";
                  SNAPRD_METRICS_PORT = "${toString cfg.port}";
                  SNAPRD_METRICS_PATH = "${cfg.metricsPath}";
                  SNAPRD_SNAPRAID_CONFIGPATH =
                    lib.mkIf (cfg.snapraidConfig != "") mkSnapraidConfigFile;
                };
                serviceConfig = {
                  ExecStart = "${cfg.package}/bin/snaprd";
                  Restart = "always";
                };
              };
            };
          };
        devShell = let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ devshell.overlays.default ];
          };
        in pkgs.devshell.mkShell {
          imports = [ (pkgs.devshell.importTOML ./devshell.toml) ];
        };
      });
}
