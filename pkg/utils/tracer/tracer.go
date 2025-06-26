package tracer

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	opentrace "go.opentelemetry.io/otel/trace"
)

// Config holds the configuration for Jaeger tracing
type Config struct {
	Endpoint    string
	ServiceName string
	Environment string
	Enabled     bool
}

// TracerProvider holds the tracer provider
type TracerProvider struct {
	provider *tracesdk.TracerProvider
}

// InitJaeger initializes Jaeger tracer
func InitJaeger(cfg Config) (*TracerProvider, error) {
	if !cfg.Enabled {
		// Return a no-op tracer provider
		provider := tracesdk.NewTracerProvider()
		return &TracerProvider{
			provider: provider,
		}, nil
	}

	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: tp,
	}, nil
}

// Shutdown gracefully shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp.provider != nil {
		if err := tp.provider.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
	}
	return nil
}
func GetTracer(name string) opentrace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string, opts ...opentrace.SpanStartOption) (context.Context, opentrace.Span) {
	tracer := GetTracer("go-fhir-demo")
	return tracer.Start(ctx, name, opts...)
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(span opentrace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// ...

func SetSpanError(span opentrace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
