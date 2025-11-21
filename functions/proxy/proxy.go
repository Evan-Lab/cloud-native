package proxy

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ping(ctx context.Context) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "ping")
	defer span.End()

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
	ctx, span := tracer.Start(ctx, "command.dispatch",
		trace.WithAttributes(
			attribute.String("discord.command.name", data.Name),
			attribute.String("discord.command.id", data.ID),
		))
	defer span.End()

	handler, ok := cmds[data.Name]
	if !ok {
		span.SetStatus(codes.Error, "no handler for command")
		slog.WarnContext(ctx, "No handler for command", "name", data.Name)
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown command",
			},
		}, nil
	}

	resp, err := handler(ctx, interaction, data)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "command handler returned error")
	}
	return resp, err
}
