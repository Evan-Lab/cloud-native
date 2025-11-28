package snap

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

var (
	tracer = otel.Tracer("github.com/Evan-Lab/cloud-native/functions/snap")
)

func init() {
	functions.CloudEvent("SnapCmd", SnapCmd)
}

type SnapData struct {
	CanvasID string `json:"canvas_id"`
	AuthorID string `json:"author_id"`
}

type MessagePublishedData struct {
	Message PubSubMessage `json:"message"`
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
	MessageID  string            `json:"messageId"`
}

func SnapCmd(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		slog.ErrorContext(ctx, "event.DataAs", "error", err)
		return fmt.Errorf("event.DataAs: %w", err)
	}
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Message.Attributes))
	ctx, span := tracer.Start(ctx, "SnapCmd")
	defer span.End()

	var payload SnapData
	if err := json.Unmarshal(msg.Message.Data, &payload); err != nil {
		slog.ErrorContext(ctx, "json.Unmarshal", "error", err, "data", string(msg.Message.Data))
		span.RecordError(err)
		return fmt.Errorf("failed to unmarshal SnapData: %w", err)
	}
	span.SetAttributes(
		attribute.String("snap.canvas_id", payload.CanvasID),
		attribute.String("snap.author_id", payload.AuthorID),
	)

	slog.InfoContext(ctx, "Received SnapCmd event", "canvas_id", payload.CanvasID, "author_id", payload.AuthorID)

	canvas, pixels, err := Snapshot(ctx, &payload)
	if err != nil {
		slog.ErrorContext(ctx, "Snapshot", "error", err)
		span.RecordError(err)
		return fmt.Errorf("Snapshot failed: %w", err)
	}

	data, err := PixelsToPng(ctx, pixels, canvas.Width, canvas.Height)
	if err != nil {
		slog.ErrorContext(ctx, "PixelsToPng", "error", err)
		span.RecordError(err)
		return fmt.Errorf("PixelsToPng failed: %w", err)
	}

	pngUrl, pixelsUrl, err := UploadSnapshot(ctx, canvas, pixels, data)
	if err != nil {
		slog.ErrorContext(ctx, "UploadSnapshot", "error", err)
		span.RecordError(err)
		return fmt.Errorf("UploadSnapshot failed: %w", err)
	}

	slog.InfoContext(ctx, "Snapshot process completed successfully", "canvas_id", canvas.ID, "png_url", pngUrl, "pixels_url", pixelsUrl)

	if interaction, ok := msg.Message.Attributes["discord_interaction_token"]; ok {
		if err := RespondToInteraction(ctx, interaction, pngUrl); err != nil {
			slog.ErrorContext(ctx, "RespondToInteraction", "error", err)
			span.RecordError(err)
			return fmt.Errorf("RespondToInteraction failed: %w", err)
		}
		slog.InfoContext(ctx, "Responded to Discord interaction", "canvas_id", canvas.ID)
	} else {
		slog.WarnContext(ctx, "No discord_interaction attribute found in Pub/Sub message")
	}

	return nil
}
