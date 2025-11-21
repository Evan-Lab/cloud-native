package proxy

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-slog/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// var (
// 	tracer = otel.Tracer("<my-module>")
// )

var (
	tp *sdktrace.TracerProvider
)

func init() {
	ctx := context.Background()

	texporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		panic(err)
	}

	ssp := sdktrace.NewSimpleSpanProcessor(texporter)
	tp = sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(ssp),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	handler := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(otelslog.NewHandler(handler)))
}
