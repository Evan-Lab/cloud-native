package proxy

import (
	"context"
	"encoding/json"
	"log/slog"

	"cloud.google.com/go/pubsub/v2"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	RegisterCommand("stop", stopCmd)
}

type StopData struct {
	CanvasID string `json:"canvasId"`
}

func stopCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "command.stop")
	defer span.End()

	client, err := PubSub()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Pub/Sub client", "error", err)
		return nil, err
	}

	publisher := client.Publisher("session-events")
	defer publisher.Stop()

	payload := PauseData{
		CanvasID: interaction.GuildID + interaction.ChannelID,
	}

	slog.DebugContext(ctx, "Stop payload", "payload", payload)
	span.SetAttributes(
		attribute.String("stop.canvas_id", payload.CanvasID),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal stop payload", "error", err)
		return nil, err
	}
	slog.DebugContext(ctx, "Stop payload", "body", string(body))

	msg := &pubsub.Message{
		Data:       body,
		Attributes: make(map[string]string),
	}

	msg.Attributes["action"] = "stop"
	msg.Attributes["discord_interaction_token"] = interaction.Token

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))
	result := publisher.Publish(ctx, msg)

	_, err = result.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish stop message", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Published stop message", "canvas_id", payload.CanvasID)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Canvas stop command received! Initializing...",
		},
	}, nil
}
