package reset_session

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	gtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/api/option"
)

func init() {
	setupLogging()
	shutdown, err := setupOpenTelemetry(context.Background())
	if err != nil {
		slog.Error("Failed to set up OpenTelemetry", "error", err)
		return
	}
	// Ensure OpenTelemetry is shut down properly on application exit.
	go func() {
		<-context.Background().Done()
		if err := shutdown(context.Background()); err != nil {
			slog.Error("Failed to shut down OpenTelemetry", "error", err)
		}
	}()
}

func setupOpenTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown combines shutdown functions from multiple OpenTelemetry
	// components into a single function.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	exporter, err := gtrace.New(
		gtrace.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")), // or leave empty to autodetect
		gtrace.WithTimeout(2*time.Second),                       // keep exports bounded
		gtrace.WithTraceClientOptions([]option.ClientOption{
			option.WithTelemetryDisabled(),
		}),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("discord-proxy"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)),
	)
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
	otel.SetTracerProvider(tp)

	return shutdown, nil
}

// func setupOpenTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
// 	var shutdownFuncs []func(context.Context) error

// 	// shutdown combines shutdown functions from multiple OpenTelemetry
// 	// components into a single function.
// 	shutdown = func(ctx context.Context) error {
// 		var err error
// 		for _, fn := range shutdownFuncs {
// 			err = errors.Join(err, fn(ctx))
// 		}
// 		shutdownFuncs = nil
// 		return err
// 	}

// 	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

// 	exporter, err := gtrace.New(
// 		gtrace.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")), // or leave empty to autodetect
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := resource.New(ctx,
// 		resource.WithAttributes(
// 			semconv.ServiceName("discord-proxy"),
// 		),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	tp := trace.NewTracerProvider(
// 		trace.WithSampler(trace.AlwaysSample()),
// 		trace.WithResource(res),
// 		trace.WithSpanProcessor(trace.NewSimpleSpanProcessor(exporter)),
// 	)
// 	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
// 	otel.SetTracerProvider(tp)

// 	return shutdown, nil

// 	// // Configure Metric Export to send metrics as OTLP
// 	// mreader, err := autoexport.NewMetricReader(ctx)
// 	// if err != nil {
// 	// 	err = errors.Join(err, shutdown(ctx))
// 	// 	return
// 	// }
// 	// mp := metric.NewMeterProvider(
// 	// 	metric.WithReader(mreader),
// 	// )
// 	// shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
// 	// otel.SetMeterProvider(mp)
// }

func setupLogging() {
	// Use json as our base logging format.
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replacer})
	// Add span context attributes when Context is passed to logging calls.
	instrumentedHandler := handlerWithSpanContext(jsonHandler)
	// Set this handler as the global slog handler.
	slog.SetDefault(slog.New(instrumentedHandler))
}
