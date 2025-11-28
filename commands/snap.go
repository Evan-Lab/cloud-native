package commands

import "github.com/bwmarrin/discordgo"

func Snap(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:        "snap",
		Description: "Take a snapshot of the current canvas",
		Options:     []*discordgo.ApplicationCommandOption{},
	}, nil
}
