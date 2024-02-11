package logging

import (
	"context"
	"time"

	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/agoda-com/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	sdklogs "github.com/agoda-com/opentelemetry-logs-go/sdk/logs"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

func NewOTLPLogsExporter(ctx context.Context, otlpEndpoint string) (*otlplogs.Exporter, error) {
	exporter, err := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient(
		otlplogshttp.WithRetry(otlplogshttp.RetryConfig{
			Enabled:        true,
			MaxElapsedTime: time.Minute,
		}),
		otlplogshttp.WithTimeout(time.Second),
		otlplogshttp.WithInsecure(),
		otlplogshttp.WithEndpoint(otlpEndpoint),
	)))
	return exporter, err
}

func NewLogProvider(res *sdkresource.Resource, exporter sdklogs.LogRecordExporter) (*sdklogs.LoggerProvider, error) {
	lp := sdklogs.NewLoggerProvider(
		sdklogs.WithResource(res),
		sdklogs.WithLogRecordProcessor(sdklogs.NewBatchLogRecordProcessor(
			exporter,
			sdklogs.WithBatchTimeout(30*time.Second),
		)),
	)
	return lp, nil
}
