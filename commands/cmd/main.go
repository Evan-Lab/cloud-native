package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Evan-Lab/cloud-native/commands"
	"github.com/Evan-Lab/cloud-native/lib/go/discord"

	"github.com/bwmarrin/discordgo"
)

type commandFunc func(s *discordgo.Session, guildID string) (*discordgo.ApplicationCommand, error)

var cmds = []commandFunc{
	commands.Hello,
	commands.Draw,
}

func createCommands(s *discordgo.Session, guildID string) error {
	appID := s.State.User.ID

	discordCmds := make([]*discordgo.ApplicationCommand, 0, len(cmds))
	var errs []error
	for _, cmdFunc := range cmds {
		cmd, err := cmdFunc(s, guildID)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to create command %s: %w", cmd.Name, err))
			slog.Error("Failed to create command", "name", cmd.Name, "error", err)
			continue
		}
		discordCmds = append(discordCmds, cmd)
		slog.Info("Command ready to be created", "name", cmd.Name)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered errors creating commands: %v", errs)
	}

	createdCmds, err := s.ApplicationCommandBulkOverwrite(appID, guildID, discordCmds)
	if err != nil {
		return fmt.Errorf("failed to bulk overwrite commands: %w", err)
	}

	for _, cmd := range createdCmds {
		slog.Info("Created command", "name", cmd.Name, "id", cmd.ID, "version", cmd.Version)
	}

	if len(createdCmds) == 0 {
		slog.Warn("No commands were created")
		return nil
	}

	slog.Info("All commands created successfully", "count", len(createdCmds))

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
