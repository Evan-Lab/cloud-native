package main

import (
	"github.com/EpitechPGE45-2025/G-CLO-910-PAR-9-1-cloudnative-13/lib/go/discord"
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
)

func createCommands(s *discordgo.Session, guildID string) error {
	appID := s.State.User.ID

	cmd := &discordgo.ApplicationCommand{
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
	}

	created, err := s.ApplicationCommandCreate(appID, guildID, cmd)
	if err != nil {
		slog.Error("Cannot create command", "error", err)
		return err
	}
	slog.Info("Command created", "id", created.ID)
	return nil
}

func run() error {
	guildID := os.Getenv("DISCORD_GUILD_ID")

	s, err := discord.Session()
	if err != nil {
		slog.Error("Failed to create Discord session", "error", err)
		return err
	}
	defer s.Close()
	if err := createCommands(s, guildID); err != nil {
		slog.Error("Failed to create commands", "error", err)
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}
