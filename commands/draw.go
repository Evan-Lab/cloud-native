package commands

import "github.com/bwmarrin/discordgo"

func Draw(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:        "draw",
		Description: "Draw on the current canvas",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "x",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "X coordinate",
			},
			{
				Name:        "y",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Y coordinate",
			},
			{
				Name:        "color",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionString,
				Description: "Color in hex format (e.g., #FF5733)",
			},
		},
	}, nil
}
