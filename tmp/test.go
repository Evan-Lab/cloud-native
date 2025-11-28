package draw-test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type PixelInput struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
	AuthorID string `json:"author_id"`
	CanvasID string `json:"canvas_id"`
}

type Pixel struct {
	AuthorID  string    `json:"author_id"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updated_at"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
}

func init() {
	functions.HTTP("DrawPixel", DrawPixel)
}

type CanvasSize struct {
	Width  int `firestore:"width"`
	Height int `firestore:"height"`
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

func DrawPixel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input PixelInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		slog.Error("Failed to decode JSON", "error", err)
		return
	}

	if input.CanvasID == "" {
		http.Error(w, "canvas_id is required", http.StatusBadRequest)
		return
	}
	if input.Color == "" {
		http.Error(w, "color is required", http.StatusBadRequest)
		return
	}
	if input.AuthorID == "" {
		http.Error(w, "authorId is required", http.StatusBadRequest)
		return
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		http.Error(w, "GOOGLE_CLOUD_PROJECT missing", http.StatusInternalServerError)
		return
	}

	t, err := GetTimeLastPixel(ctx, "rate_limits", input.AuthorID, "dev-rplace-database", "serverless-epitech-dev-476110")
	if err != nil {
		slog.Warn("No last pixel timestamp found (first pixel allowed)", "error", err)
	} else {
		slog.Info("Last pixel updated at", "time", t)

		elapsed := time.Since(t)

		if elapsed < 35*time.Second {
			remaining := 35*time.Second - elapsed

			http.Error(w,
				fmt.Sprintf("Cooldown not finished. Wait %d more seconds", int(remaining.Seconds())),
				http.StatusTooManyRequests,
			)
			return
		}
		slog.Info("Cooldown OK: user can draw", "elapsed_seconds", elapsed.Seconds())
	}

	size, err := GetCanvasSize(ctx, "canvases", input.CanvasID, "dev-rplace-database", projectID)
	if err != nil {
		http.Error(w, "Failed to get canvas size", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Canvas %s : %dx%d\n", input.CanvasID, size.Width, size.Height)

	if input.X > size.Height || input.Y > size.Width {
		http.Error(w, "The pixel position is not in the canvas", http.StatusBadRequest)
		return
	}

	fs, err := firestore.NewClientWithDatabase(ctx, projectID, "dev-rplace-database")
	if err != nil {
		http.Error(w, "Firestore init failed", http.StatusInternalServerError)
		slog.Error("Failed Firestore init", "error", err)
		return
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
		http.Error(w, "Failed to write pixel", http.StatusInternalServerError)
		slog.Error("Failed to write pixel to Firestore", "error", err)
		return
	}

	err = SaveLastPixelTime(ctx, "rate_limits", input.AuthorID, "dev-rplace-database", projectID, time.Now())
	if err != nil {
		slog.Error("Failed to store last pixel timestamp", "error", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Pixel written"))
}
