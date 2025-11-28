package reset_session

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

type ResetInput struct {
	CanvasID string `json:"canvasId"`
}

func init() {
	functions.CloudEvent("ResetSession", ResetSession)
}

func ResetSession(ctx context.Context, e cloudevents.Event) error {
	var payload MessagePublishedData
	if err := e.DataAs(&payload); err != nil {
		slog.Error("Invalid CloudEvent payload", "error", err)
		return nil
	}

	msg := payload.Message

	parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Attributes))
	tracer := otel.Tracer("ResetSession")
	ctx, span := tracer.Start(parentCtx, "reset-session")
	defer span.End()

	slog.Info("ResetSession triggered", "data", string(msg.Data))

	var input ResetInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	if input.CanvasID == "" {
		slog.Error("canvasId is required for reset")
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, "dev-rplace-database")
	if err != nil {
		slog.Error("Firestore init failed", "error", err)
		return err
	}
	defer fs.Close()

	bw := fs.BulkWriter(ctx)

	pixels := fs.Collection("canvases").Doc(input.CanvasID).Collection("pixels")
	iter := pixels.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		bw.Delete(doc.Ref)
	}

	bw.End()

	slog.Info("Pixels reset", "canvasId", input.CanvasID)
	return nil
}
