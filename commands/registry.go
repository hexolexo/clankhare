package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Registry []Command

func Register(c Command) {
	Registry = append(Registry, c)
}
