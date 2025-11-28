package snap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/api/iterator"
)

type Canvas struct {
	ID      string `firestore:"-"`
	AdminID string `firestore:"adminId"`
	Name    string `firestore:"name"`
	Status  string `firestore:"status"`

	Height int `firestore:"height"`
	Width  int `firestore:"width"`

	EndDate   time.Time `firestore:"endDate"`
	StartDate time.Time `firestore:"startDate"`
}

type Pixel struct {
	X         int       `firestore:"X" json:"x"`
	Y         int       `firestore:"Y" json:"y"`
	Color     string    `firestore:"color" json:"color"`
	AuthorID  string    `firestore:"authorID" json:"author_id"`
	UpdatedAt time.Time `firestore:"UpdatedAt" json:"updated_at"`
}

const DEFAULT_COLOR = "#FFFFFF"

func Snapshot(ctx context.Context, data *SnapData) (*Canvas, []Pixel, error) {
	ctx, span := tracer.Start(ctx, "command.snap")
	defer span.End()
	client, err := Firestore(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer client.Close()

	// Fetch pixels ordered by Y, then by X
	doc, err := client.Collection("canvases").Doc(data.CanvasID).Get(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "firestore.Get", "error", err, "canvas_id", data.CanvasID)
		span.RecordError(err)
		return nil, nil, err
	}

	if !doc.Exists() {
		err = fmt.Errorf("firestore document does not exist")
		slog.ErrorContext(ctx, "firestore document does not exist", "canvas_id", data.CanvasID)
		span.RecordError(err)
		return nil, nil, err
	}

	var canvas Canvas
	if err := doc.DataTo(&canvas); err != nil {
		slog.ErrorContext(ctx, "doc.DataTo", "error", err, "canvas_id", data.CanvasID)
		span.RecordError(err)
		return nil, nil, err
	}
	canvas.ID = doc.Ref.ID

	slog.InfoContext(ctx, "Creating snapshot", "canvas_id", data.CanvasID, "author_id", data.AuthorID)
	span.AddEvent("Snapshot creation started", trace.WithAttributes(
		attribute.String("canvas_id", data.CanvasID),
		attribute.String("author_id", data.AuthorID),
	))

	query := client.Collection("canvases").Doc(data.CanvasID).Collection("pixels").
		Where("X", ">=", 0).
		Where("X", "<", canvas.Width).
		Where("Y", ">=", 0).
		Where("Y", "<", canvas.Height).
		OrderBy("Y", firestore.Asc).
		OrderBy("X", firestore.Asc)

	iter := query.Documents(ctx)
	defer iter.Stop()

	pixels := make([]Pixel, canvas.Width*canvas.Height)
	now := time.Now()
	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			index := y*canvas.Width + x
			pixels[index] = Pixel{
				X:         x,
				Y:         y,
				Color:     DEFAULT_COLOR, // Default color
				UpdatedAt: now,
				AuthorID:  "",
			}
		}
	}
	pcCount := 0
	for i := 0; i < len(pixels); i++ {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			slog.ErrorContext(ctx, "iter.Next", "error", err, "canvas_id", data.CanvasID)
			span.RecordError(err)
			return nil, nil, err
		}
		var pixel Pixel
		if err := doc.DataTo(&pixel); err != nil {
			slog.ErrorContext(ctx, "doc.DataTo", "error", err, "canvas_id", data.CanvasID)
			span.RecordError(err)
			return nil, nil, err
		}
		if pixel.X < 0 || pixel.X >= canvas.Width || pixel.Y < 0 || pixel.Y >= canvas.Height {
			slog.WarnContext(ctx, "Pixel out of bounds", "x", pixel.X, "y", pixel.Y, "canvas_width", canvas.Width, "canvas_height", canvas.Height)
			continue
		}
		slog.DebugContext(ctx, "Fetched pixel", "x", pixel.X, "y", pixel.Y, "color", pixel.Color)
		index := pixel.Y*canvas.Width + pixel.X
		pixels[index] = pixel
		pcCount++
	}

	slog.InfoContext(ctx, "Snapshot created successfully", "canvas_id", data.CanvasID, "author_id", data.AuthorID, "pixels_count", pcCount)

	return &canvas, pixels, nil
}
