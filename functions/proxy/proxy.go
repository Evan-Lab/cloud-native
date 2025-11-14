package proxy

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func ping(ctx context.Context) (*discordgo.InteractionResponse, error) {
	slog.InfoContext(ctx, "Received Ping interaction")
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}, nil
}

type CommandHandler func(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error)

var cmds = make(map[string]CommandHandler)

func RegisterCommand(cmd string, handler CommandHandler) {
	cmds[cmd] = handler
}

func cmdProxy(ctx context.Context, interaction discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	data := interaction.ApplicationCommandData()
	handler, ok := cmds[data.Name]
	if !ok {
		slog.WarnContext(ctx, "No handler for command", "name", data.Name)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown command",
			},
		}, nil
	}

	return handler(ctx, interaction, data)
}
