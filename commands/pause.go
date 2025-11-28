package commands

import (
	"github.com/Evan-Lab/cloud-native/lib/go/utils"
	"github.com/bwmarrin/discordgo"
)

func Pause(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:                     "pause",
		Description:              "Pause on the current canvas",
		DefaultMemberPermissions: utils.Ptr(int64(discordgo.PermissionAdministrator)),
		Options:                  []*discordgo.ApplicationCommandOption{},
	}, nil
}
