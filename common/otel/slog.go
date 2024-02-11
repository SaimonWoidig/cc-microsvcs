package otel

import (
	"log/slog"

	"go.opentelemetry.io/otel"
)

/*
OtelSlogErrorHandler is a struct representing an OpenTelemetry error handler that uses the slog logging package.

Fields:
  - l: A pointer to a slog.Logger instance.

Implemented Interfaces:
  - otel.ErrorHandler

Example Usage:

	o := OtelSlogErrorHandler{}
*/
type OtelSlogErrorHandler struct {
	l *slog.Logger
}

var _ otel.ErrorHandler = OtelSlogErrorHandler{}

/*
NewOtelSlogErrorHandler creates a new instance of OtelErrorSlogLogger.

Parameters:
  - logger: A pointer to a slog.Logger instance.

Returns:
  - OtelErrorSlogLogger: The newly created OtelErrorSlogLogger instance.

Example:

	logger := slog.New()
	otelLogger := otel.NewOtelSlogErrorHandler(logger)
*/
func NewOtelSlogErrorHandler(logger *slog.Logger) OtelSlogErrorHandler {
	return OtelSlogErrorHandler{
		l: logger,
	}
}

/*
Handle implements otel.ErrorHandler. Logs the error using the slog logging package.
It takes an error as a parameter and logs the error message using the slog.Logger instance.

Parameters:
  - err: The error to be logged.

Example:

	o := OtelErrorSlogLogger{}
	o.Handle(err)
*/
func (o OtelSlogErrorHandler) Handle(err error) {
	o.l.Error("opentelemetry error", "error", err.Error())
}
