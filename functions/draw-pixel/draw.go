package draw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

type CanvasSize struct {
	Width  int `firestore:"width"`
	Height int `firestore:"height"`
}

func init() {
	functions.CloudEvent("DrawPixel", DrawPixel)
}

func GetCanvasAdminID(ctx context.Context, canvasID string, database string, projectID string) (string, error) {
	if projectID == "" {
		return "", errors.New("projectID missing")
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
	admin, ok := data["AdminID"]
	if !ok {
		return "", errors.New("AdminID missing")
	}

	return admin.(string), nil
}

func GetCanvasStatus(ctx context.Context, canvasID string, database string, projectID string) (string, error) {
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
	statusRaw, ok := data["Status"]
	if !ok {
		return "", errors.New("Status missing")
	}

	status, ok := statusRaw.(string)
	if !ok {
		return "", errors.New("Status invalid")
	}

	return status, nil
}

func GetCanvasSize(ctx context.Context, collection string, canvasID string, database string, projectID string) (CanvasSize, error) {
	var size CanvasSize

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err != nil {
		return size, err
	}
	defer fs.Close()

	doc, err := fs.Collection(collection).Doc(canvasID).Get(ctx)
	if err != nil {
		return size, err
	}

	err = doc.DataTo(&size)
	return size, err
}

func GetTimeLastPixel(ctx context.Context, collection string, authorId string, database string, projectID string) (time.Time, error) {
	fs, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err != nil {
		return time.Time{}, err
	}
	defer fs.Close()

	doc, err := fs.Collection(collection).Doc(authorId).Get(ctx)
	if err != nil {
		return time.Time{}, err
	}

	val, ok := doc.Data()["updatedAt"]
	if !ok {
		return time.Time{}, errors.New("updatedAt missing")
	}

	t, ok := val.(time.Time)
	if !ok {
		return time.Time{}, errors.New("updatedAt invalid")
	}

	return t, nil
}

func SaveLastPixelTime(ctx context.Context, collection string, authorId string, database string, projectID string, t time.Time) error {
	fs, err := firestore.NewClientWithDatabase(ctx, projectID, database)
	if err != nil {
		return err
	}
	defer fs.Close()

	_, err = fs.Collection(collection).Doc(authorId).Set(ctx, map[string]interface{}{
		"updatedAt": t,
	})

	return err
}

func DrawPixel(ctx context.Context, e cloudevents.Event) error {
	var payload MessagePublishedData
	if err := e.DataAs(&payload); err != nil {
		slog.Error("Invalid CloudEvent", "error", err)
		return nil
	}

	msg := payload.Message
	parentCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(msg.Attributes))
	tracer := otel.Tracer("DrawPixel")
	ctx, span := tracer.Start(parentCtx, "draw-pixel")
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

	status, err := GetCanvasStatus(ctx, input.CanvasID, databaseName, projectID)
	if err != nil {
		slog.Error("Failed status fetch", "error", err)
		return nil
	}

	if status != "START" {
		slog.Warn("Canvas not in START state", "status", status)
		return nil
	}

	adminID, err := GetCanvasAdminID(ctx, input.CanvasID, databaseName, projectID)
	if err != nil {
		slog.Error("Failed adminId fetch", "error", err)
		return nil
	}

	if input.AuthorID != adminID {
		last, err := GetTimeLastPixel(ctx, "rate_limits", input.AuthorID, databaseName, projectID)
		if err == nil {
			if elapsed := time.Since(last); elapsed < 35*time.Second {
				slog.Warn("Cooldown not finished", "remaining", 35*time.Second-elapsed)
				return nil
			}
		}
	} else {
		slog.Info("Admin bypass cooldown")
	}

	size, err := GetCanvasSize(ctx, "canvases", input.CanvasID, databaseName, projectID)
	if err != nil {
		slog.Error("Canvas size error", "error", err)
		return nil
	}

	if input.X < 0 || input.Y < 0 || input.X >= size.Width || input.Y >= size.Height {
		slog.Error("Pixel out of bounds", "input", input)
		return nil
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, databaseName)
	if err != nil {
		slog.Error("Firestore init fail", "error", err)
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

	docID := fmt.Sprintf("%d_%d", input.X, input.Y)

	_, err = fs.Collection("canvases").Doc(input.CanvasID).
		Collection("pixels").Doc(docID).Set(ctx, pixel)

	if err != nil {
		slog.Error("Pixel write failed", "error", err)
		return err
	}

	if err := SaveLastPixelTime(ctx, "rate_limits", input.AuthorID, databaseName, projectID, time.Now()); err != nil {
		slog.Warn("Failed to update rate limit", "error", err)
	}

	slog.Info("Pixel written", "canvas", input.CanvasID, "x", input.X, "y", input.Y)
	return nil
}
