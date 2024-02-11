package otel

import (
	"context"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type ResourceConfig struct {
	Name       string
	Version    string
	Addr       string
	Port       int
	InstanceID string
}

func NewResource(ctx context.Context, config ResourceConfig) (*sdkresource.Resource, error) {
	res, err := sdkresource.New(ctx,
		sdkresource.WithFromEnv(),
		sdkresource.WithProcess(),
		sdkresource.WithHost(),
		sdkresource.WithContainer(),
		sdkresource.WithOS(),
		sdkresource.WithSchemaURL(semconv.SchemaURL),
		sdkresource.WithAttributes(
			semconv.ServiceName(config.Name),
			semconv.ServiceVersion(config.Version),
			semconv.ServiceInstanceID(config.InstanceID),
			semconv.ServerAddress(config.Addr),
			semconv.ServerPort(config.Port),
		),
	)
	if err != nil {
		return nil, err
	}

	return sdkresource.Merge(sdkresource.Default(), res)
}
