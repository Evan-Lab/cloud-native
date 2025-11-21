package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel/attribute"
)

var (
	projectID string
	topicID   string
)

func init() {
	topicID = os.Getenv("DRAW_TOPIC_ID")
	RegisterCommand("draw", drawCmd)
}

type DrawData struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
	CanvasID string `json:"canvas_id"`
	AuthorID string `json:"author_id"`
}

func drawCmd(ctx context.Context, interaction discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	ctx, span := tracer.Start(ctx, "command.draw")
	defer span.End()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Pub/Sub client", "error", err)
		return nil, err
	}

	publisher := client.Publisher(topicID)
	defer publisher.Stop()

	payload := DrawData{
		CanvasID: interaction.GuildID,
		AuthorID: interaction.Member.User.ID,
	}

	colorOpt := data.GetOption("color")
	xOpt := data.GetOption("x")
	yOpt := data.GetOption("y")

	if colorOpt == nil || xOpt == nil || yOpt == nil {
		slog.WarnContext(ctx, "Missing required options", "color", colorOpt, "x", xOpt, "y", yOpt)
		return nil, fmt.Errorf("missing required options")
	}

	payload.Color = colorOpt.StringValue()
	payload.X = int(xOpt.IntValue())
	payload.Y = int(yOpt.IntValue())

	slog.DebugContext(ctx, "Draw command options", "x", payload.X, "y", payload.Y, "color", payload.Color)
	span.SetAttributes(
		attribute.Int("draw.x", payload.X),
		attribute.Int("draw.y", payload.Y),
		attribute.String("draw.color", payload.Color),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal draw payload", "error", err)
		return nil, err
	}
	slog.DebugContext(ctx, "Draw payload", "body", string(body))

	result := publisher.Publish(ctx, &pubsub.Message{
		Data: body,
	})

	_, err = result.Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to publish draw message", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Published draw message", "canvas_id", payload.CanvasID, "x", payload.X, "y", payload.Y, "color", payload.Color)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: ":thumbsup:",
		},
	}, nil
}
