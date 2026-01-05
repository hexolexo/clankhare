package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gorcon/rcon"
)

// Variables used for command line parameters
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

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	session.AddHandler(whitelist)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "whitelist",
			Description: "Add player to whitelist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "player",
					Description: "Minecraft username",
					Required:    true,
				},
			},
		},
		// Add more commands here
	}

	for _, cmd := range commands {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", cmd)
		if err != nil {
			fmt.Printf("Cannot create command %s: %v\n", cmd.Name, err)
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()
}
func whitelist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "whitelist" {
		player := i.ApplicationCommandData().Options[0].StringValue()

		err := whitelistPlayer(player)

		var content string
		if err != nil {
			content = fmt.Sprintf("Failed to whitelist %s: %v", player, err)
			log.Printf("RCON error: %v", err)
		} else {
			content = fmt.Sprintf("Whitelisted %s", player)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	}
}

func whitelistPlayer(username string) error {
	conn, err := rcon.Dial(Minecraft_IP+":16260", RCON_Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Execute(fmt.Sprintf("whitelist add %s", username))
	return err
}
