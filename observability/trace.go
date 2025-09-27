package observability

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var OtlpEndpoint = os.Getenv("OTLP_ENDPOINT")
var Tracer = otel.Tracer("patty-wagon")

func NewOTLPTraceExporter(ctx context.Context, otlpEndpoint string) *otlptrace.Exporter {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint))
	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatalf("Failed to create the collector trace exporter: %v", err)
	}
	return traceExp
}

func SetupTracer(ctx context.Context, otlpEndpoint string) {
	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceNameKey.String("patty-wagon"),
	))
	if err != nil {
		log.Fatalf("failed to initialize resource: %v", err)
	}

	traceExp := NewOTLPTraceExporter(ctx, otlpEndpoint)
	bsp := trace.NewBatchSpanProcessor(traceExp)
	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down tracer...")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error shutting down tracer: %w", err)
		}
	}()
}
