package commands

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gorcon/rcon"
)

var (
	whitelistCounts = make(map[string]int)
	countMutex      sync.Mutex
)

func init() {
	Register(Command{
		Definition: &discordgo.ApplicationCommand{
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
		Handler: handleWhitelist,
	})
}

func handleWhitelist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	countMutex.Lock()
	if whitelistCounts[userID] >= 3 {
		countMutex.Unlock()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You've used all 3 whitelist slots",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	whitelistCounts[userID]++
	countMutex.Unlock()

	player := i.ApplicationCommandData().Options[0].StringValue()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		err := rconWhitelist(player)
		content := fmt.Sprintf("Whitelisted %s", player)
		if err != nil {
			content = fmt.Sprintf("Failed to whitelist %s: %v", player, err)
			log.Printf("RCON error: %v", err)
		}
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &content})
	}()
}

func rconWhitelist(username string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`).MatchString(username) {
		return fmt.Errorf("invalid minecraft username format")
	}
	conn, err := rcon.Dial(os.Getenv("MINECRAFT_IP")+":16260", os.Getenv("RCON_Password"))
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Execute(fmt.Sprintf("whitelist add %s", username))
	return err
}
