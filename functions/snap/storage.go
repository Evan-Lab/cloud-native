package snap

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"go.opentelemetry.io/otel/attribute"
)

var bucketName string

func init() {
	bucketName = os.Getenv("SNAPSHOT_BUCKET")
	if bucketName == "" {
		panic("SNAPSHOT_BUCKET environment variable is not set")
	}
}

// func uploadObject(ctx context.Context, obj *storage.ObjectHandle, data []byte, contentType string) (err error) {
// 	ctx, span := tracer.Start(ctx, "uploadObject")
// 	defer span.End()

// 	writer := obj.NewWriter(ctx)
// 	defer func() {
// 		err = writer.Close()
// 		if err != nil {
// 			slog.ErrorContext(ctx, "writer.Close", "error", err)
// 			span.RecordError(err)
// 			err = fmt.Errorf("writer.Close: %w", err)
// 		}
// 	}()

// 	writer.ContentType = contentType
// 	writer.

// 	_, err = writer.Write(data)
// 	if err != nil {
// 		slog.ErrorContext(ctx, "writer.Write", "error", err)
// 		span.RecordError(err)
// 		return fmt.Errorf("writer.Write: %w", err)
// 	}

// 	return nil
// }

func WritePng(w io.WriteCloser, pngData []byte) error {
	err := errors.Join(
		func() error { _, err := w.Write(pngData); return err }(),
		w.Close(),
	)
	if err != nil {
		return fmt.Errorf("failed to write PNG data: %w", err)
	}
	return nil
}

func uploadPng(ctx context.Context, bucket *storage.BucketHandle, canvasID string, pngData []byte) (string, error) {
	ctx, span := tracer.Start(ctx, "uploadPng")
	defer span.End()

	path := fmt.Sprintf("canvas_%s.png", canvasID)
	slog.DebugContext(ctx, "Uploading PNG to GCS", "path", path)
	span.SetAttributes(attribute.String("snapshot.png_path", path))

	obj := bucket.Object(path)
	writer := obj.NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = "image/png"

	if err := WritePng(writer, pngData); err != nil {
		slog.ErrorContext(ctx, "writePng", "error", err)
		span.RecordError(err)
		return "", fmt.Errorf("failed to write PNG to GCS: %w", err)
	}

	url, err := bucket.SignedURL(obj.ObjectName(), &storage.SignedURLOptions{
		Method:      "GET",
		Expires:     time.Now().Add(24 * time.Hour),
		ContentType: "image/png",
	})
	if err != nil {
		slog.ErrorContext(ctx, "bucket.SignedURL", "error", err)
		span.RecordError(err)
		return "", fmt.Errorf("failed to generate signed URL for PNG: %w", err)
	}

	return url, nil
}

func WritePixels(w io.WriteCloser, pixels []Pixel) error {
	gz := gzip.NewWriter(w)
	err := errors.Join(
		json.NewEncoder(gz).Encode(pixels),
		gz.Close(),
		w.Close(),
	)

	if err != nil {
		return fmt.Errorf("failed to write pixels: %w", err)
	}
	return nil
}

func uploadPixels(ctx context.Context, bucket *storage.BucketHandle, canvasID string, pixels []Pixel) (string, error) {
	ctx, span := tracer.Start(ctx, "uploadPixels")
	defer span.End()

	path := fmt.Sprintf("canvas_%s.json.gz", canvasID)
	slog.DebugContext(ctx, "Uploading pixels to GCS", "path", path)
	span.SetAttributes(attribute.String("snapshot.pixels_path", path))

	obj := bucket.Object(path)

	writer := obj.NewWriter(ctx)

	writer.ContentType = "application/gzip"

	if err := WritePixels(writer, pixels); err != nil {
		slog.ErrorContext(ctx, "writePixels", "error", err)
		span.RecordError(err)
		return "", fmt.Errorf("failed to write pixels to GCS: %w", err)
	}

	url, err := bucket.SignedURL(obj.ObjectName(), &storage.SignedURLOptions{
		Method:      "GET",
		Expires:     time.Now().Add(24 * time.Hour),
		ContentType: "application/gzip",
	})
	if err != nil {
		slog.ErrorContext(ctx, "bucket.SignedURL", "error", err)
		span.RecordError(err)
		return "", fmt.Errorf("failed to generate signed URL for PNG: %w", err)
	}

	return url, nil
}

func UploadSnapshot(ctx context.Context, canvas *Canvas, pixels []Pixel, pngData []byte) (string, string, error) {
	ctx, span := tracer.Start(ctx, "UploadSnapshot")
	defer span.End()

	client, err := storage.NewClient(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "storage.NewClient", "error", err)
		span.RecordError(err)
		return "", "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	pngUrl, err1 := uploadPng(ctx, bucket, canvas.ID, pngData)

	pixelsUrl, err2 := uploadPixels(ctx, bucket, canvas.ID, pixels)

	return pngUrl, pixelsUrl, errors.Join(err1, err2)
}
