package logging

import (
	"log/slog"

	"github.com/SaimonWoidig/cc-microsvcs/common/logging"
	"github.com/agoda-com/opentelemetry-go/otelslog"
	"github.com/agoda-com/opentelemetry-logs-go/logs"
	slogmulti "github.com/samber/slog-multi"
)

func NewSlogOtelHandler(logProvider logs.LoggerProvider, logLevel string) *otelslog.OtelHandler {
	otelHandler := otelslog.NewOtelHandler(logProvider, &otelslog.HandlerOptions{Level: logging.LogLevelStringToSlogLevel(logLevel)})
	return otelHandler
}

func NewSlogOtelCompositeLogger(slogHandler slog.Handler, otelHandler *otelslog.OtelHandler) *slog.Logger {
	h := slogmulti.Fanout(slogHandler, otelHandler)
	return slog.New(h)
}
