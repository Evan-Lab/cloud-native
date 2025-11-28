package start_session

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

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

type CanvasInput struct {
	AdminID   string    `json:"adminId"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type Canvas struct {
	AdminID   string    `json:"adminId"`
	CanvasID  string    `json:"canvasId"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func init() {
	functions.CloudEvent("StartSession", StartSession)
}

func StartSession(ctx context.Context, e cloudevents.Event) error {

	var payload MessagePublishedData
	if err := e.DataAs(&payload); err != nil {
		slog.Error("Invalid CloudEvent payload", "error", err)
		return nil
	}

	msg := payload.Message
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Attributes))
	tracer := otel.Tracer("StartSession")
	ctx, span := tracer.Start(parentCtx, "start-session")
	defer span.End()

	slog.Info("Message received", "data", string(msg.Data), "attributes", msg.Attributes)

	var input Canvas
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, databaseName)
	if err != nil {
		slog.Error("Firestore init failed", "error", err)
		return err
	}
	defer fs.Close()

	status := "START"

	if input.CanvasID != "" {

		doc := fs.Collection("canvases").Doc(input.CanvasID)

		_, err := doc.Set(ctx, map[string]interface{}{
			"Status": status,
		}, firestore.MergeAll)

		if err != nil {
			slog.Error("Failed to update canvas", "canvasId", input.CanvasID, "error", err)
			return err
		}

		slog.Info("Canvas updated", "canvasId", input.CanvasID)
		return nil
	}

	if input.StartDate.IsZero() {
		slog.Error("startDate must be provided when creating a canvas")
		return nil
	}

	newCanvas := Canvas{
		AdminID:   input.AdminID,
		Name:      input.Name,
		Width:     input.Width,
		Height:    input.Height,
		Status:    status,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	docRef := fs.Collection("canvases").NewDoc()
	newCanvas.CanvasID = docRef.ID

	if _, err := docRef.Set(ctx, newCanvas); err != nil {
		slog.Error("Failed to create canvas", "error", err)
		return err
	}

	slog.Info("Canvas created", "canvasId", newCanvas.CanvasID)
	return nil
}
