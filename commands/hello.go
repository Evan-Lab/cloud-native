package commands

import "github.com/bwmarrin/discordgo"

func Hello(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Say hello",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Your name (optional)",
				Required:    false,
				// For autocomplete, set Autocomplete: true (see section 2)
			},
		},
	}, nil
}
