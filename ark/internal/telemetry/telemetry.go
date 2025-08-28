package telemetry

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("telemetry")

func Initialize() func() {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		log.Info("OTEL_EXPORTER_OTLP_ENDPOINT not set, telemetry disabled")
		return func() {}
	}

	headers := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ark-controller"
	}

	log.Info("initializing telemetry", "endpoint", endpoint, "service", serviceName, "headers", headers)

	// Auto-configure OTLP exporter from environment variables:
	// OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_HEADERS, OTEL_SERVICE_NAME
	exporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		log.Error(err, "failed to create OTLP exporter")
		return func() {}
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)

	// Send a basic controller startup event to ensure telemetry events are sent properly
	tracer := otel.Tracer("ark/controller-startup")
	_, span := tracer.Start(context.Background(), "controller.startup")

	version := os.Getenv("VERSION")
	if version == "" {
		version = "dev"
	}

	span.SetAttributes(
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(version),
	)
	log.Info("sending controller startup telemetry event", "span.id", span.SpanContext().SpanID().String())
	span.End()

	log.Info("telemetry initialized successfully")

	return func() {
		log.Info("shutting down telemetry")
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Error(err, "failed to shutdown tracer provider")
		}
	}
}
