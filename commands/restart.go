package commands

import (
	"github.com/Evan-Lab/cloud-native/lib/go/utils"
	"github.com/bwmarrin/discordgo"
)

func Restart(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error) {
	return &discordgo.ApplicationCommand{
		Name:                     "restart",
		Description:              "Restart on the current canvas",
		DefaultMemberPermissions: utils.Ptr(int64(discordgo.PermissionAdministrator)),
		Options:                  []*discordgo.ApplicationCommandOption{},
	}, nil
}
