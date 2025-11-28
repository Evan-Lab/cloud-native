package snap_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/Evan-Lab/cloud-native/functions/snap"
)

func buildPixelData(canvas snap.Canvas) []snap.Pixel {
	pixels := make([]snap.Pixel, 0, canvas.Width*canvas.Height)
	for y := 0; y < canvas.Height; y++ {
		for x := 0; x < canvas.Width; x++ {
			color := fmt.Sprintf("#%06x", rand.Intn(0xFFFFFF+1))
			pixels = append(pixels, snap.Pixel{
				X:         x,
				Y:         y,
				Color:     color,
				UpdatedAt: time.Now().Add(-time.Duration(rand.Intn(1000)) * time.Second),
				AuthorID:  fmt.Sprintf("%d%d%d", rand.Uint64(), rand.Uint64(), rand.Uint64()),
			})
		}
	}
	return pixels
}

type NoopWriteCloser struct {
	io.Writer
}

func (n NoopWriteCloser) Close() error {
	return nil
}

func testSize(t *testing.T, width, height int) {
	ctx := context.Background()

	canvas := snap.Canvas{
		ID:     fmt.Sprintf("test-canvas-%dx%d", width, height),
		Width:  width,
		Height: height,
	}

	pixels := buildPixelData(canvas)

	pngData, err := snap.PixelsToPng(ctx, pixels, canvas.Width, canvas.Height)
	if err != nil {
		t.Fatalf("PixelsToPng failed: %v", err)
	}
	t.Logf("Generated PNG data size: %d bytes", len(pngData))

	var buf bytes.Buffer
	w := NoopWriteCloser{Writer: &buf}

	if err := snap.WritePng(w, pngData); err != nil {
		t.Fatalf("WritePng failed: %v", err)
	}
	t.Logf("Written PNG buffer size: %d bytes", buf.Len())
	if buf.Len() == 0 {
		t.Fatalf("PNG buffer is empty")
	}

	buf.Reset()

	if err := snap.WritePixels(w, pixels); err != nil {
		t.Fatalf("WritePixels failed: %v", err)
	}
	t.Logf("Written Pixels buffer size: %d bytes", buf.Len())
	if buf.Len() == 0 {
		t.Fatalf("Pixels buffer is empty")
	}
}

func Test10x10(t *testing.T) {
	testSize(t, 10, 10)
}

func Test50x50(t *testing.T) {
	testSize(t, 50, 50)
}

func Test50x100(t *testing.T) {
	testSize(t, 50, 100)
}

func Test100x100(t *testing.T) {
	testSize(t, 100, 100)
}

func Test250x250(t *testing.T) {
	testSize(t, 250, 250)
}

func Test500x500(t *testing.T) {
	testSize(t, 500, 500)
}

func Test1000x1000(t *testing.T) {
	testSize(t, 1000, 1000)
}
