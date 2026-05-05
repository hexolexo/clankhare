package bridge

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/nats-io/nats.go"
)

func Start(s *discordgo.Session) (func(), error) {
	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	nc.Subscribe("mc.chat.outgoing", func(msg *nats.Msg) {
		_, err := s.ChannelMessageSend(channelID, string(msg.Data))
		if err != nil {
			log.Printf("failed to send to Discord: %v", err)
		}
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID || m.ChannelID != channelID {
			return
		}
		nc.Publish("mc.chat.incoming", []byte(m.Author.Username+": "+m.Content))
	})

	return func() { nc.Drain() }, nil
}
