package proxy

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	RegisterCommand("start", startCmd)
	RegisterCommand("restart", startCmd)

}

type StartData struct {
	CanvasID  string    `json:"canvasId"`
	AdminID   string    `json:"adminId"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func startCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "command.start")
	defer span.End()

	client, err := PubSub()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Pub/Sub client", "error", err)
		return nil, err
	}

	publisher := client.Publisher("session-events")
	defer publisher.Stop()

	payload := StartData{
		CanvasID:  interaction.GuildID + interaction.ChannelID,
		AdminID:   interaction.Member.User.ID,
		Name:      interaction.Member.User.Username + "'s Canvas",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	payload.Width = int(data.GetOption("width").IntValue())
	payload.Height = int(data.GetOption("height").IntValue())

	slog.DebugContext(ctx, "Start payload", "payload", payload)
	span.SetAttributes(
		attribute.String("start.canvas_id", payload.CanvasID),
		attribute.String("start.author_id", payload.AdminID),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal start payload", "error", err)
		return nil, err
	}
	slog.DebugContext(ctx, "Start payload", "body", string(body))

	msg := &pubsub.Message{
		Data:       body,
		Attributes: make(map[string]string),
	}

	msg.Attributes["action"] = "start"
	msg.Attributes["discord_interaction_token"] = interaction.Token

	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(msg.Attributes))
	result := publisher.Publish(ctx, msg)

	_, err = result.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish snap message", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Published snap message", "canvas_id", payload.CanvasID, "author_id", payload.AdminID)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Canvas start command received! Initializing...",
		},
	}, nil
}
