package commands

import (
	"github.com/Evan-Lab/cloud-native/lib/go/utils"
	"github.com/bwmarrin/discordgo"
)

func Start(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:                     "start",
		Description:              "Start on the current canvas",
		DefaultMemberPermissions: utils.Ptr(int64(discordgo.PermissionAdministrator)),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "width",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Width of the canvas",
			},
			{
				Name:        "height",
				Required:    true,
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "Height of the canvas",
			},
		},
	}, nil
}
