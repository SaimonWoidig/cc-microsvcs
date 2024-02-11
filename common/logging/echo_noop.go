package logging

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// NoopLogger is an implementation of echo.Logger. It discards all logs.
type NoopLogger struct{}

// Debug implements echo.Logger.
func (*NoopLogger) Debug(i ...interface{}) {}

// Debugf implements echo.Logger.
func (*NoopLogger) Debugf(format string, args ...interface{}) {}

// Debugj implements echo.Logger.
func (*NoopLogger) Debugj(j log.JSON) {}

// Error implements echo.Logger.
func (*NoopLogger) Error(i ...interface{}) {}

// Errorf implements echo.Logger.
func (*NoopLogger) Errorf(format string, args ...interface{}) {}

// Errorj implements echo.Logger.
func (*NoopLogger) Errorj(j log.JSON) {}

// Fatal implements echo.Logger.
func (*NoopLogger) Fatal(i ...interface{}) {}

// Fatalf implements echo.Logger.
func (*NoopLogger) Fatalf(format string, args ...interface{}) {}

// Fatalj implements echo.Logger.
func (*NoopLogger) Fatalj(j log.JSON) {}

// Info implements echo.Logger.
func (*NoopLogger) Info(i ...interface{}) {}

// Infof implements echo.Logger.
func (*NoopLogger) Infof(format string, args ...interface{}) {}

// Infoj implements echo.Logger.
func (*NoopLogger) Infoj(j log.JSON) {}

// Level implements echo.Logger.
func (*NoopLogger) Level() log.Lvl { return log.OFF }

// Output implements echo.Logger.
func (*NoopLogger) Output() io.Writer { return io.Discard }

// Panic implements echo.Logger.
func (*NoopLogger) Panic(i ...interface{}) {}

// Panicf implements echo.Logger.
func (*NoopLogger) Panicf(format string, args ...interface{}) {}

// Panicj implements echo.Logger.
func (*NoopLogger) Panicj(j log.JSON) {}

// Prefix implements echo.Logger.
func (*NoopLogger) Prefix() string { return "" }

// Print implements echo.Logger.
func (*NoopLogger) Print(i ...interface{}) {}

// Printf implements echo.Logger.
func (*NoopLogger) Printf(format string, args ...interface{}) {}

// Printj implements echo.Logger.
func (*NoopLogger) Printj(j log.JSON) {}

// SetHeader implements echo.Logger.
func (*NoopLogger) SetHeader(h string) {}

// SetLevel implements echo.Logger.
func (*NoopLogger) SetLevel(v log.Lvl) {}

// SetOutput implements echo.Logger.
func (*NoopLogger) SetOutput(w io.Writer) {}

// SetPrefix implements echo.Logger.
func (*NoopLogger) SetPrefix(p string) {}

// Warn implements echo.Logger.
func (*NoopLogger) Warn(i ...interface{}) {}

// Warnf implements echo.Logger.
func (*NoopLogger) Warnf(format string, args ...interface{}) {}

// Warnj implements echo.Logger.
func (*NoopLogger) Warnj(j log.JSON) {}

// Interface guard.
var _ echo.Logger = (*NoopLogger)(nil)

// NewNoopLogger creates a new NoopLogger and returns a pointer to it.
func NewNoopLogger() *NoopLogger {
	return new(NoopLogger)
}
