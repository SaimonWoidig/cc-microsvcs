package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/agoda-com/opentelemetry-logs-go/logs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	commonconfig "github.com/SaimonWoidig/cc-microsvcs/common/config"
	"github.com/SaimonWoidig/cc-microsvcs/common/logging"
	commonotel "github.com/SaimonWoidig/cc-microsvcs/common/otel"
	otellogging "github.com/SaimonWoidig/cc-microsvcs/common/otel/logging"
	otelmetrics "github.com/SaimonWoidig/cc-microsvcs/common/otel/metrics"
	oteltracing "github.com/SaimonWoidig/cc-microsvcs/common/otel/tracing"
	"github.com/SaimonWoidig/cc-microsvcs/service.auth/pkg/config"
)

const AppName = "auth-service"

type Container struct {
	Config         *config.Config
	Logger         *slog.Logger
	Resource       *resource.Resource
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	LoggerProvider logs.LoggerProvider
}

func NewContainer() *Container {
	c := new(Container)

	cfg, err := initConfig()
	if err != nil {
		panic(err.Error())
	}
	c.Config = cfg

	c.Logger = initLogger(c.Config.Logging.LogLevel, c.Config.Logging.Pretty)
	c.Logger.Info("base logger initialized")

	otel.SetErrorHandler(commonotel.NewOtelSlogErrorHandler(c.Logger))
	iid := uuid.New().String()
	res, err := initResource(iid, c.Config.Server.Addr, c.Config.Server.Port)
	if err != nil {
		panic(err.Error())
	}
	c.Resource = res
	tp, err := initOtelTracing(
		c.Resource,
		c.Config.Telemetry.Tracing.OTLPEndpoint,
		c.Config.Telemetry.Tracing.ExportTimeoutSeconds,
		c.Config.Telemetry.Tracing.SamplingRatio,
	)
	if err != nil {
		panic(err.Error())
	}
	c.TracerProvider = tp
	mp, err := initOtelMetrics(c.Resource, c.Config.Telemetry.Metrics.OTLPEndpoint, c.Config.Telemetry.Metrics.ExportTimeoutSeconds, c.Config.Telemetry.Metrics.ExportIntervalSeconds, c.Config.Telemetry.Metrics.MemStatsIntervalSeconds)
	if err != nil {
		panic(err.Error())
	}
	c.MeterProvider = mp
	lp, err := initOtelLogging(c.Resource, c.Config.Telemetry.Logging.OTLPEndpoint, c.Config.Telemetry.Logging.ExportTimeoutSeconds, c.Config.Telemetry.Logging.BatchTimeoutSeconds)
	if err != nil {
		panic(err.Error())
	}
	c.LoggerProvider = lp
	cl := otellogging.NewSlogOtelCompositeLogger(c.Logger.Handler(), otellogging.NewSlogOtelHandler(c.LoggerProvider, c.Config.Telemetry.Logging.LogLevel))
	c.Logger = cl
	c.Logger.Info("composite OTLP logger initialized")

	return c
}

func (c *Container) Shutdown() error {
	return nil
}

func initConfig() (*config.Config, error) {
	c := new(config.Config)
	v := commonconfig.NewViper()
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}

func initLogger(logLevel string, pretty bool) *slog.Logger {
	l := logging.NewSlogLogger(logLevel, pretty)
	return l
}

func initResource(instanceID string, serverAddr string, serverPort int) (*resource.Resource, error) {
	return commonotel.NewResource(context.Background(), commonotel.ResourceConfig{
		Name:       AppName,
		Version:    "1.0.0",
		InstanceID: instanceID,
		Addr:       serverAddr,
		Port:       serverPort,
	})
}

func initOtelTracing(resource *resource.Resource, otlpEndpoint string, exportTimeoutSeconds int, samplingRatio float64) (trace.TracerProvider, error) {
	traceExporter, err := oteltracing.NewOTLPTraceExporter(
		context.Background(),
		otlpEndpoint,
		time.Duration(exportTimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	traceSampler := oteltracing.NewRatioSampler(samplingRatio)
	tp, err := oteltracing.NewTraceProvider(
		resource,
		traceExporter,
		traceSampler,
	)
	if err != nil {
		return nil, err
	}
	return tp, nil
}

func initOtelMetrics(resource *resource.Resource, otlpEndpoint string, exportTimeoutSeconds, exportIntervalSeconds, memStatsIntervalSeconds int) (metric.MeterProvider, error) {
	metricExporter, err := otelmetrics.NewOTLPMetricExporter(
		context.Background(),
		otlpEndpoint,
		time.Duration(exportTimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	mp, err := otelmetrics.NewMeterProvider(
		resource,
		metricExporter,
		time.Duration(exportIntervalSeconds)*time.Second,
		time.Duration(memStatsIntervalSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func initOtelLogging(resource *resource.Resource, otlpEndpoint string, exportTimeoutSeconds, batchTimeoutSeconds int) (logs.LoggerProvider, error) {
	logExporter, err := otellogging.NewOTLPLogsExporter(
		context.Background(),
		otlpEndpoint,
		time.Duration(exportTimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	lp, err := otellogging.NewLogProvider(
		resource,
		logExporter,
		time.Duration(batchTimeoutSeconds)*time.Second,
	)
	if err != nil {
		return nil, err
	}
	return lp, nil
}
