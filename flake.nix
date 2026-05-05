{
  description = "Clankhare flake";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
    in {
      packages.default = pkgs.buildGoModule {
        pname = "discord-bot";
        version = "0.0.1";
        src = ./.;
        vendorHash = null;
        env.CGO_ENABLED = 1;
        nativeBuildInputs = [pkgs.pkg-config];
        meta = {
          description = "The discord bot for the forgotten wonderland";
          homepage = "https://github.com/hexolexo/clankhare";
        };
      };
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [go gopls gotools gcc];
      };
    })
    // {
      # Lives outside eachDefaultSystem — modules don't carry a system
      nixosModules.default = {
        config,
        lib,
        pkgs,
        ...
      }: let
        cfg = config.services.clankhare;
      in {
        options.services.clankhare = {
          enable = lib.mkEnableOption "Clankhare discord bot";

          package = lib.mkOption {
            type = lib.types.package;
            # HACK: self.packages pulls the package for the current system at
            # module eval time — works fine but requires the flake to be in your
            # inputs so `self` is defined
            default = self.packages.${pkgs.system}.default;
            description = "The clankhare package to use";
          };

          configFile = lib.mkOption {
            type = lib.types.path;
            description = "Path to a file containing the bot token (keeps it out of the nix store)";
          };

          user = lib.mkOption {
            type = lib.types.str;
            default = "clankhare";
          };

          group = lib.mkOption {
            type = lib.types.str;
            default = "clankhare";
          };
        };

        config = lib.mkIf cfg.enable {
          users.users.${cfg.user} = {
            isSystemUser = true;
            group = cfg.group;
          };
          users.groups.${cfg.group} = {};

          systemd.services.clankhare = {
            description = "Clankhare discord bot";
            wantedBy = ["multi-user.target"];
            after = ["network.target"];

            serviceConfig = {
              ExecStart = "${cfg.package}/bin/discord-bot";
              User = cfg.user;
              Group = cfg.group;

              # Reads the token from a file and injects it as an env var —
              # adjust the var name to whatever your bot actually reads
              EnvironmentFile = cfg.configFile;

              # Basic hardening — tighten as needed
              NoNewPrivileges = true;
              PrivateTmp = true;
              Restart = "on-failure";
            };
          };
        };
      };
    };
}
