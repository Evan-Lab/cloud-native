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
	RegisterCommand("snap", snapCmd)
}

type SnapData struct {
	CanvasID string `json:"canvas_id"`
	AuthorID string `json:"author_id"`
}

func snapCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "command.snap")
	defer span.End()

	client, err := PubSub()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Pub/Sub client", "error", err)
		return nil, err
	}

	publisher := client.Publisher("command.snap")
	defer publisher.Stop()

	payload := SnapData{
		CanvasID: interaction.GuildID + interaction.ChannelID,
		AuthorID: interaction.Member.User.ID,
	}

	slog.DebugContext(ctx, "Snap payload", "payload", payload)
	span.SetAttributes(
		attribute.String("snap.canvas_id", payload.CanvasID),
		attribute.String("snap.author_id", payload.AuthorID),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal snap payload", "error", err)
		return nil, err
	}
	slog.DebugContext(ctx, "Snap payload", "body", string(body))

	msg := &pubsub.Message{
		Data:       body,
		Attributes: make(map[string]string),
	}

	msg.Attributes["discord_interaction_token"] = interaction.Token

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))
	result := publisher.Publish(ctx, msg)

	_, err = result.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish snap message", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Published snap message", "canvas_id", payload.CanvasID, "author_id", payload.AuthorID)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Taking snapshot... :camera_with_flash:",
		},
	}, nil
}
