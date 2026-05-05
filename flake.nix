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

        env = {
          CGO_ENABLED = 1;
        };

        nativeBuildInputs = [pkgs.pkg-config];

        meta = {
          description = "The discord bot for the forgotten wonderland";
          homepage = "https://github.com/hexolexo/clankhare";
        };
      };

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [go gopls gotools gcc];
      };
    });
}
