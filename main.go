package main

import (
	"Clankhare/commands"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token         string
	Minecraft_IP  string
	RCON_Password string
)

func init() {
	Token = os.Getenv("TOKEN")
	Minecraft_IP = os.Getenv("MINECRAFT_IP")
	RCON_Password = os.Getenv("RCON_Password")
}

func main() {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session:", err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		for _, cmd := range commands.Registry {
			if cmd.Definition.Name == i.ApplicationCommandData().Name {
				cmd.Handler(s, i)
				return
			}
		}
	})

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	if err = session.Open(); err != nil {
		log.Fatal("error opening connection:", err)
	}

	for _, cmd := range commands.Registry {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", cmd.Definition)
		if err != nil {
			log.Fatalf("cannot create command %s: %v", cmd.Definition.Name, err)
		}
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	session.Close()
}
