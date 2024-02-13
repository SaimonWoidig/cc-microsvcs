package config

type LoggingConfig struct {
	LogLevel string `mapstructure:"logLevel"`
	Pretty   bool   `mapstructure:"pretty"`
}

type TracingConfig struct {
	OTLPEndpoint  string  `mapstructure:"otlpEndpoint"`
	SamplingRatio float64 `mapstructure:"samplingRatio"`
}
type MetricsConfig struct {
	OTLPEndpoint            string `mapstructure:"otlpEndpoint"`
	RequestTimeoutSeconds   int    `mapstructure:"requestTimeout"`
	ExportIntervalSeconds   int    `mapstructure:"exportInterval"`
	MemStatsIntervalSeconds int    `mapstructure:"memStatsInterval"`
}
type LogsConfig struct {
	OTLPEndpoint string `mapstructure:"otlpEndpoint"`
	LogLevel     string `mapstructure:"logLevel"`
}

type TelemetryConfig struct {
	Tracing TracingConfig `mapstructure:"tracing"`
	Metrics MetricsConfig `mapstructure:"metrics"`
	Logging LogsConfig    `mapstructure:"logs"`
}

type ServerConfig struct {
	Addr string `mapstructure:"addr"`
	Port int    `mapstructure:"port"`
}

type Config struct {
	Logging   LoggingConfig   `mapstructure:"logging"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	Server    ServerConfig    `mapstructure:"server"`
}
