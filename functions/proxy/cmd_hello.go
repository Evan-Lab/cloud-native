package proxy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel/attribute"
)

func init() {
	RegisterCommand("hello", helloCmd)
}

func helloCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "discord.command.hello")
	defer span.End()

	name := "World"
	if opt := data.GetOption("name"); opt != nil {
		if s, ok := opt.Value.(string); ok && s != "" {
			name = s
		}
	}
	span.SetAttributes(attribute.String("options.name", name))
	slog.InfoContext(ctx, "Received hello command", "name", name)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello, %s!", name),
		},
	}, nil
}
