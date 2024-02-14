package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const TracerName string = "github.com/SaimonWoidig/cc-microsvcs/common/otel/tracing"

func NewStdoutTraceExporter() (*stdouttrace.Exporter, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	return exporter, err
}

func NewOTLPTraceExporter(ctx context.Context, otlpEndpoint string, exportTimeout time.Duration) (*otlptrace.Exporter, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:        true,
			MaxElapsedTime: time.Minute,
		}),
		otlptracehttp.WithTimeout(exportTimeout),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(otlpEndpoint),
	)
	return exporter, err
}

func NewAlwaysSampleSampler() sdktrace.Sampler {
	return sdktrace.AlwaysSample()
}

func NewRatioSampler(ratio float64) sdktrace.Sampler {
	return sdktrace.TraceIDRatioBased(ratio)
}

func NewTraceProvider(res *sdkresource.Resource, exporter sdktrace.SpanExporter, sampler sdktrace.Sampler) (*sdktrace.TracerProvider, error) {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
	)
	return tp, nil
}

func NewTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewTracer(tp *sdktrace.TracerProvider) trace.Tracer {
	return tp.Tracer(TracerName)
}
