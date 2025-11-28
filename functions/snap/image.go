package snap

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
)

func HexToColor(hex string) (color.RGBA, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}, nil
}

func PixelsToPng(ctx context.Context, pixels []Pixel, width, height int) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "PixelsToPng")
	defer span.End()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for _, pixel := range pixels {
		col, err := HexToColor(pixel.Color)
		if err != nil {
			slog.ErrorContext(ctx, "HexToColor", "error", err, "color", pixel.Color)
			span.RecordError(err)
			return nil, err
		}
		img.Set(pixel.X, pixel.Y, col)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		slog.ErrorContext(ctx, "png.Encode", "error", err)
		span.RecordError(err)
		return nil, err
	}

	pngData := buf.Bytes()
	slog.InfoContext(ctx, "PNG image created", "size_bytes", len(pngData))
	span.SetAttributes(attribute.Int("png.size_bytes", len(pngData)))

	return pngData, nil
}
