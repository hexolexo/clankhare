# Clankhare

Discord bot for the forgotten wonderland discord

## Setup

1. Create `.env` file

```env
TOKEN=DISCORD_TOKEN
MINECRAFT_IP=10.0.0.255
RCON_Password=PASSWORD
```

2. Install dependencies

```bash
go mod tidy
```

3. Enable bot intents in Discord Developer Portal

Bot → Privileged Gateway Intents → Message Content Intent (Good luck I'm not a discord dev)

4. Run the bot

```bash
go run main.go
```

## Commands

- `/whitelist <player>` - Add player to Minecraft whitelist via RCON
- More to come
