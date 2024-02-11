package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
)

const MeterName string = "github.com/SaimonWoidig/cc-microsvcs/common/otel/metrics"

func NewStdoutMetricExporter() (sdkmetric.Exporter, error) {
	exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	return exporter, err
}

func NewOTLPMetricExporter(ctx context.Context, otlpEndpoint string, reqTimeout time.Duration) (*otlpmetrichttp.Exporter, error) {
	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{
			Enabled:        true,
			MaxElapsedTime: time.Minute,
		}),
		otlpmetrichttp.WithTimeout(reqTimeout),
		otlpmetrichttp.WithEndpoint(otlpEndpoint),
		otlpmetrichttp.WithInsecure(),
	)
	return exporter, err
}

func NewMeterProvider(res *sdkresource.Resource, exporter sdkmetric.Exporter, exportInterval time.Duration, memStatsInterval time.Duration) (*sdkmetric.MeterProvider, error) {
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(exportInterval),
			),
		),
	)

	if err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(memStatsInterval), runtime.WithMeterProvider(mp)); err != nil {
		return nil, err
	}

	if err := host.Start(host.WithMeterProvider(mp)); err != nil {
		return nil, err
	}

	return mp, nil
}

func NewMeter(mp *sdkmetric.MeterProvider) metric.Meter {
	return mp.Meter(MeterName)
}
