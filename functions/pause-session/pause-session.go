package pause_session

import (
	"context"
	"encoding/json"
	"log/slog"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type MessagePublishedData struct {
	Message PubSubMessage `json:"message"`
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

type PauseInput struct {
	CanvasID string `json:"canvasId"`
}

func init() {
	functions.CloudEvent("PauseSession", PauseSession)
}

func PauseSession(ctx context.Context, e cloudevents.Event) error {
	var payload MessagePublishedData
	if err := e.DataAs(&payload); err != nil {
		slog.Error("Invalid CloudEvent payload", "error", err)
		return nil
	}

	msg := payload.Message
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Attributes))
	tracer := otel.Tracer("PauseSession")
	ctx, span := tracer.Start(parentCtx, "pause-session")
	defer span.End()

	slog.Info("PauseSession triggered", "data", string(msg.Data))

	var input PauseInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	if input.CanvasID == "" {
		slog.Error("canvasId is required for pause")
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, databaseName)
	if err != nil {
		slog.Error("Firestore init failed", "error", err)
		return err
	}
	defer fs.Close()

	_, err = fs.Collection("canvases").Doc(input.CanvasID).Set(
		ctx,
		map[string]interface{}{
			"Status": "PAUSE",
		},
		firestore.MergeAll,
	)
	if err != nil {
		slog.Error("Failed to update canvas", "canvasId", input.CanvasID, "error", err)
		return err
	}

	slog.Info("Canvas paused", "canvasId", input.CanvasID)
	return nil
}
