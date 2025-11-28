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
	CanvasID  string    `json:"canvasId"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type Canvas struct {
	AdminID   string    `json:"adminId"`
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

	var input CanvasInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	if input.CanvasID == "" {
		slog.Error("canvasId must be provided when creating a canvas")
		return nil
	}

	if input.AdminID == "" || input.Name == "" || input.Width <= 0 || input.Height <= 0 {
		slog.Error("Missing required fields", "input", input)
		return nil
	}

	if input.StartDate.IsZero() {
		slog.Error("startDate must be provided when creating a canvas")
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, databaseName)
	if err != nil {
		slog.Error("Firestore init failed", "error", err)
		return err
	}
	defer fs.Close()

	canvas := Canvas{
		AdminID:   input.AdminID,
		Name:      input.Name,
		Width:     input.Width,
		Height:    input.Height,
		Status:    "START",
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	_, err = fs.Collection("canvases").Doc(input.CanvasID).Set(ctx, canvas)
	if err != nil {
		slog.Error("Failed to create canvas", "canvasId", input.CanvasID, "error", err)
		return err
	}

	slog.Info("Canvas created", "canvasId", input.CanvasID)
	return nil
}
