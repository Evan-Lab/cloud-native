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
	"golang.org/x/image/draw"
)

func HexToColor(hex string) (color.RGBA, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}, nil
}

func ScaleImage(src image.Image, size int) (image.Image, error) {

	largestSide := src.Bounds().Dx()
	if src.Bounds().Dy() > largestSide {
		largestSide = src.Bounds().Dy()
	}

	scale := float64(size) / float64(largestSide)
	width := int(float64(src.Bounds().Dx()) * scale)
	height := int(float64(src.Bounds().Dy()) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst, nil
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

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	scaledImg, err := ScaleImage(img, 1024)
	if err != nil {
		slog.ErrorContext(ctx, "ScaleImage", "error", err)
		span.RecordError(err)
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaledImg); err != nil {
		slog.ErrorContext(ctx, "png.Encode", "error", err)
		span.RecordError(err)
		return nil, err
	}

	pngData := buf.Bytes()
	slog.InfoContext(ctx, "PNG image created", "size_bytes", len(pngData))
	span.SetAttributes(attribute.Int("png.size_bytes", len(pngData)))

	return pngData, nil
}
