# Clankhare

Discord bot for the forgotten wonderland discord

## Setup

1. add to flake.nix

```nix
{
    inputs = {
        nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
        agenix = {
            url = "github:ryantm/agenix";
            inputs.nixpkgs.follows = "nixpkgs";
        };

        clankhare.url = "github:hexolexo/clankhare"; 
    };

    outputs = { self, nixpkgs, clankhare, ... }: {
        nixosConfigurations.your-hostname = nixpkgs.lib.nixosSystem {
            system = "x86_64-linux";
            modules = [
                ./configuration.nix
                {
                    environment.systemPackages = [ 
                        clankhare.packages.x86_64-linux.default 
                    ];
                }
            ];
        };
    };
}
```

2. setup secrets.nix

```nix
let
    user = "ssh-ed25519 AAAAC3Nza..."; # Local SSH pubkey
    host = "ssh-ed25519 AAAAC3Nza..."; # Remote SSH pubkey
in {
    "clankhare-env.age".publicKeys = [ user host ];
}
```

```bash
agenix -e clankhare-env.age
TOKEN=DISCORD_TOKEN
MINECRAFT_IP=10.0.0.255
RCON_Password=PASSWORD
DISCORD_CHANNEL_ID=1111111111111111111
NATS_URL=127.0.0.1
```

3. setup clankhare.nix module

```nix
{ pkgs, inputs, config, ... }: {
  imports = [ inputs.agenix.nixosModules.default ];

  # Decrypt the secret
  age.secrets.clankhare-env = {
    file = ./clankhare-env.age;
    mode = "0444";
  };

  systemd.services.clankhare = {
    description = "Clankhare Discord Bot";
    after = [ "network-online.target" ];
    wantedBy = [ "multi-user.target" ];

    serviceConfig = {
      ExecStart = "${inputs.clankhare.packages.${pkgs.system}.default}/bin/discord-bot";
      
      EnvironmentFile = config.age.secrets.clankhare-env.path;

      Restart = "always";
      RestartSec = "5s";
      DynamicUser = true;
      
      ProtectSystem = "strict";
      ProtectHome = true;
      PrivateTmp = true;
    };
  };
}
```

4. import clankhare.nix in you configuration.nix


## Commands

- `/whitelist <player>` - Add player to Minecraft whitelist via RCON
- More to come
