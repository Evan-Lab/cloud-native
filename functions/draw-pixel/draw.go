package draw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type PixelInput struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
	AuthorID string `json:"authorId"`
	CanvasID string `json:"canvasId"`
}

type Pixel struct {
	AuthorID  string    `json:"authorId"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updatedAt"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
}

type MessagePublishedData struct {
	Message PubSubMessage `json:"message"`
}

type PubSubMessage struct {
	Data       []byte            `json:"data"`
	Attributes map[string]string `json:"attributes"`
}

func init() {
	functions.CloudEvent("DrawPixel", DrawPixel)
}

type CanvasSize struct {
	Width  int `firestore:"width"`
	Height int `firestore:"height"`
}

func GetCanvasAdminID(ctx context.Context, canvasID string, database string, projectID string) (string, error) {
	if projectID == "" {
		return "", errors.New("GOOGLE_CLOUD_PROJECT missing")
	}
	if database == "" {
		return "", errors.New("database missing")
	}
	if canvasID == "" {
		return "", errors.New("canvasID missing")
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err != nil {
		return "", err
	}
	defer fs.Close()

	doc, err := fs.Collection("canvases").Doc(canvasID).Get(ctx)
	if err != nil {
		return "", err
	}

	data := doc.Data()
	if data == nil {
		return "", errors.New("canvas document is empty")
	}

	slog.Info("Debug canvas doc", "data", doc.Data())

	admin, ok := data["AdminID"]
	if !ok {
		return "", errors.New("AdminId missing or invalid")
	}

	return admin.(string), nil
}

func GetCanvasSize(ctx context.Context, collection string, canvasID string, database string, projectId string) (CanvasSize, error) {
	var size CanvasSize
	if canvasID == "" {
		return size, errors.New("canvasID is required")
	}
	if projectId == "" {
		return size, errors.New("GOOGLE_CLOUD_PROJECT env var is required")
	}
	if database == "" {
		return size, errors.New("projectId is required")
	}
	if collection == "" {
		return size, errors.New("collection is required")
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectId, database)
	if err != nil {
		return size, err
	}
	defer fs.Close()

	doc, err := fs.Collection(collection).Doc(canvasID).Get(ctx)
	if err != nil {
		return size, err
	}

	if err := doc.DataTo(&size); err != nil {
		return size, err
	}

	slog.Info("Canvas size fetched",
		"canvasID", canvasID,
		"width", size.Width,
		"height", size.Height,
	)

	return size, nil
}

func GetTimeLastPixel(ctx context.Context, collection string, authorId string, database string, projectId string) (time.Time, error) {
	if projectId == "" {
		return time.Time{}, errors.New("GOOGLE_CLOUD_PROJECT env var is required")
	}
	if database == "" {
		return time.Time{}, errors.New("database is required")
	}
	if collection == "" {
		return time.Time{}, errors.New("collection is required")
	}
	if authorId == "" {
		return time.Time{}, errors.New("authorId is required")
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectId, database)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to init Firestore: %w", err)
	}
	defer fs.Close()

	doc, err := fs.Collection(collection).Doc(authorId).Get(ctx)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read document: %w", err)
	}

	val, ok := doc.Data()["updatedAt"]
	if !ok {
		return time.Time{}, errors.New("field updatedAt not found in document")
	}

	timestamp, ok := val.(time.Time)
	if !ok {
		return time.Time{}, errors.New("updatedAt is not a valid timestamp")
	}

	return timestamp, nil
}

func SaveLastPixelTime(ctx context.Context, collection string, authorId string, database string, projectId string, t time.Time) error {
	if projectId == "" {
		return errors.New("projectId is required")
	}
	if database == "" {
		return errors.New("database is required")
	}
	if collection == "" {
		return errors.New("collection is required")
	}
	if authorId == "" {
		return errors.New("authorId is required")
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectId, database)
	if err != nil {
		return fmt.Errorf("failed to init Firestore: %w", err)
	}
	defer fs.Close()

	_, err = fs.Collection(collection).
		Doc(authorId).
		Set(ctx, map[string]interface{}{
			"updatedAt": t,
		})

	if err != nil {
		return fmt.Errorf("failed to write rate limit: %w", err)
	}

	return nil
}

func DrawPixel(ctx context.Context, e cloudevents.Event) error {
	var payload MessagePublishedData
	if err := e.DataAs(&payload); err != nil {
		slog.Error("Invalid CloudEvent payload", "error", err)
		return nil
	}

	msg := payload.Message

	parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Attributes))
	tracer := otel.Tracer("DrawPixel")
	ctx, span := tracer.Start(parentCtx, "event")
	defer span.End()

	slog.Info("Message received", "data", string(msg.Data), "attributes", msg.Attributes)

	var input PixelInput
	if err := json.Unmarshal(msg.Data, &input); err != nil {
		slog.Error("Invalid JSON", "raw", string(msg.Data))
		return nil
	}

	if input.CanvasID == "" || input.Color == "" || input.AuthorID == "" {
		slog.Error("Missing required fields", "input", input)
		return nil
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		slog.Error("Missing GOOGLE_CLOUD_PROJECT")
		return nil
	}

	adminID, err := GetCanvasAdminID(ctx, input.CanvasID, "dev-rplace-database", projectID)
	if err != nil {
		slog.Error("Failed to fetch adminId", "error", err)
		return nil
	}

	if input.AuthorID != adminID {
		t, err := GetTimeLastPixel(ctx, "rate_limits", input.AuthorID, "dev-rplace-database", projectID)
		if err == nil {
			if elapsed := time.Since(t); elapsed < 35*time.Second {
				slog.Warn("Cooldown not finished", "remaining", 35*time.Second-elapsed)
				return nil
			}
		}
	} else {
		slog.Info("Admin bypass: cooldown ignored", "adminId", adminID)
	}

	size, err := GetCanvasSize(ctx, "canvases", input.CanvasID, "dev-rplace-database", projectID)
	if err != nil {
		slog.Error("Canvas not found", "canvas_id", input.CanvasID)
		return nil
	}

	if input.X < 0 || input.Y < 0 || input.X >= size.Width || input.Y >= size.Height {
		slog.Error("Pixel out of bounds", "input", input)
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, "dev-rplace-database")
	if err != nil {
		slog.Error("Firestore init failed", "error", err)
		return err
	}
	defer fs.Close()

	pixel := Pixel{
		AuthorID:  input.AuthorID,
		Color:     input.Color,
		UpdatedAt: time.Now(),
		X:         input.X,
		Y:         input.Y,
	}

	docID := fmt.Sprintf("%d_%d", pixel.X, pixel.Y)
	_, err = fs.Collection("canvases").
		Doc(input.CanvasID).
		Collection("pixels").
		Doc(docID).
		Set(ctx, pixel)

	if err != nil {
		slog.Error("Failed to write pixel", "error", err)
		return err
	}

	if err := SaveLastPixelTime(ctx, "rate_limits", input.AuthorID, "dev-rplace-database", projectID, time.Now()); err != nil {
		slog.Warn("Failed to update rate limit", "error", err)
	}

	slog.Info("Pixel written", "canvas_id", input.CanvasID, "x", input.X, "y", input.Y)
	return nil
}
