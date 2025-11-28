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
	tracer := otel.Tracer("CreateCanvas")
	ctx, span := tracer.Start(parentCtx, "create-canvas")
	defer span.End()

	slog.Info("Message received", "data", string(msg.Data), "attributes", msg.Attributes)

	var input CanvasInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, "dev-rplace-database")
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
		StartDate: time.Now(),
		EndDate:   time.Time{},
	}

	docRef := fs.Collection("canvases").NewDoc()
	if _, err := docRef.Set(ctx, canvas); err != nil {
		slog.Error("Failed to create canvas", "error", err)
		return err
	}
	canvasID := docRef.ID

	slog.Info("Canvas created", "canvasID", canvasID)

	_, err = fs.Collection("canvases").Doc(canvasID).Set(ctx, canvas)
	if err != nil {
		slog.Error("Failed to create pixel document", "error", err)
		return err
	}
	return nil
}
