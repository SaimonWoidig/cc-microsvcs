package metrics

/*
Contains an Echo middleware exposing Echo webserver metrics in OTEL format.
Rewritten from https://github.com/labstack/echo-contrib/blob/master/prometheus/prometheus.go.
*/

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	// DefaultMeterName represents the default name of the meter used to create metrics. It is set to "go.opentelemetry.io/otel/metric/echo-metrics".
	DefaultMeterName string = "go.opentelemetry.io/otel/metric/echo-metrics"
	// DefaultInstrumentNamePrefix represents the default prefix added to the names of the created metrics. It is set to "echo".
	DefaultInstrumentNamePrefix string = "echo"
)

const (
	// bucketKiB is representing the size of a kilobyte. Used as bucket boundaries in the creation of histograms for request and response sizes in the metrics package.
	bucketKiB float64 = 1024
	// bucketMiB is representing the size of a megabyte. Used as bucket boundaries in the creation of histograms for request and response sizes in the metrics package.
	bucketMiB float64 = 1024 * bucketKiB

	// unitByte represents the unit of measurement for byte values.
	unitByte string = "By"
	// unitSecond represents the unit of measurement for second values.
	unitSecond string = "s"
	// unitDimensionless represents the unit of measurement for dimensionless value
	unitDimensionless string = "1"
)

// secondsBucket is a slice of float64 values representing the boundaries for histogram buckets. These buckets are used in the creation of histograms for request duration in the metrics package.
var secondsBucket = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

// bytesBucket is a slice of float64 values representing the boundaries for histogram buckets. These buckets are used in the creation of histograms for request and response sizes in the metrics package.
var bytesBucket = []float64{1.0 * bucketKiB, 2.0 * bucketKiB, 5.0 * bucketKiB, 10.0 * bucketKiB, 100 * bucketKiB, 500 * bucketKiB, 1.0 * bucketMiB, 2.5 * bucketMiB, 5.0 * bucketMiB, 10.0 * bucketMiB}

// Config is a struct that represents the configuration options for a metrics middleware.
type Config struct {
	// MeterProvider is metric.MeterProvider used to create metrics.
	MeterProvider metric.MeterProvider
	// Skipper is the middleware.Skipper function used to determine if the middleware should be skipped for a request.
	Skipper middleware.Skipper
	// MeterName is the name of the meter used to create metrics.
	MeterName string
	// InstrumentNamePrefix is the prefix added to the names of the created metrics.
	InstrumentNamePrefix string
}

/*
New returns a new instance of the middleware function with default configuration.

The returned middleware function measures the latency, request count, request size, and response size of HTTP requests handled by an Echo server.
It uses OpenTelemetry to create metrics and records the measurements for each request.

The default configuration uses the following settings:
  - MeterProvider: The default OpenTelemetry meter provider obtained from otel.GetMeterProvider().
  - Skipper: The default Echo middleware skipper obtained from middleware.DefaultSkipper.
  - MeterName: The default meter name "echo_http_metrics".
  - InstrumentNamePrefix: The default instrument name prefix "echo".

Returns:
  - echo.MiddlewareFunc: The middleware function that can be used with Echo's Use() method.
  - error: An error if there was a problem creating the metrics or configuring the middleware.

Example usage:

	e := echo.New()
	e.Use(New())

Note: This function is a convenience wrapper around NewWithConfig() and provides default configuration values. If you need to customize the configuration, use NewWithConfig() instead.

See Also:
  - NewWithConfig: Creates a new instance of the middleware function with custom configuration.
  - Config: The configuration struct used to customize the middleware behavior.
*/
func New() (echo.MiddlewareFunc, error) {
	config := Config{
		MeterProvider:        otel.GetMeterProvider(),
		Skipper:              middleware.DefaultSkipper,
		MeterName:            DefaultMeterName,
		InstrumentNamePrefix: DefaultInstrumentNamePrefix,
	}
	return NewWithConfig(config)
}

/*
NewWithConfig is a function that creates and returns an Echo middleware function with the given configuration.

Parameters:
  - config: A Config struct that contains the configuration options for the middleware.

Returns:
  - echo.MiddlewareFunc: The Echo middleware function.
  - error: An error if any error occurs during the creation of the middleware.

The NewWithConfig function initializes the required metrics for monitoring HTTP requests and responses.
It creates histograms for request duration, request count, request size, and response size.
The function then returns a middleware function that records the metrics for each incoming request and response.

The middleware function measures the duration of each request, the count of requests, the size of the request, and the size of the response.
It records these metrics using the configured meter provider and attribute set.
The function also handles any errors that occur during the execution of the next handler and sets the appropriate HTTP status code.

Example usage:

	config := Config{
		Skipper: middleware.DefaultSkipper,
		MeterProvider: otel.GetMeterProvider(),
		InstrumentNamePrefix: "myapp",
	}

	middlewareFunc, err := NewWithConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middlewareFunc)

Note: The NewWithConfig function assumes that the required dependencies, such as the Echo framework and OpenTelemetry meter provider, are already imported and configured.
*/
func NewWithConfig(config Config) (echo.MiddlewareFunc, error) {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}
	if config.MeterProvider == nil {
		config.MeterProvider = otel.GetMeterProvider()
	}
	if config.InstrumentNamePrefix != "" {
		config.InstrumentNamePrefix += "_"
	}

	meter := config.MeterProvider.Meter(config.MeterName)

	reqDuration, err := meter.Float64Histogram(
		fmt.Sprintf("%vrequest_duration_seconds", config.InstrumentNamePrefix),
		metric.WithDescription("The HTTP request latencies in seconds."),
		metric.WithUnit(unitSecond),
		metric.WithExplicitBucketBoundaries(secondsBucket...),
	)
	if err != nil {
		return nil, err
	}

	reqCount, err := meter.Int64Counter(
		fmt.Sprintf("%vrequests_total", config.InstrumentNamePrefix),
		metric.WithDescription("Number of HTTP requests processed."),
		metric.WithUnit(unitDimensionless),
	)
	if err != nil {
		return nil, err
	}

	reqSize, err := meter.Int64Histogram(
		fmt.Sprintf("%vrequest_size_bytes", config.InstrumentNamePrefix),
		metric.WithDescription("The HTTP request sizes in bytes."),
		metric.WithUnit(unitByte),
		metric.WithExplicitBucketBoundaries(bytesBucket...),
	)
	if err != nil {
		return nil, err
	}

	resSize, err := meter.Int64Histogram(
		fmt.Sprintf("%vresponse_size_bytes", config.InstrumentNamePrefix),
		metric.WithDescription("The HTTP response sizes in bytes."),
		metric.WithUnit(unitByte),
		metric.WithExplicitBucketBoundaries(bytesBucket...),
	)
	if err != nil {
		return nil, err
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			elapsed := float64(time.Since(start)) / float64(time.Second)

			status := c.Response().Status
			if err != nil {
				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					status = httpError.Code
				}
				if status == 0 || status == http.StatusOK {
					status = http.StatusInternalServerError
				}
			}

			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			attrs := attribute.NewSet(
				attribute.Int("status", status),
				attribute.String("method", c.Request().Method),
				attribute.String("path", path),
				attribute.String("host", c.Request().Host),
			)

			reqDuration.Record(c.Request().Context(), elapsed, metric.WithAttributeSet(attrs))
			reqCount.Add(c.Request().Context(), 1, metric.WithAttributeSet(attrs))
			reqSize.Record(c.Request().Context(), int64(computeApproxReqSize(c.Request())), metric.WithAttributeSet(attrs))
			resSize.Record(c.Request().Context(), c.Response().Size, metric.WithAttributeSet(attrs))

			return err
		}
	}, nil
}

/*
computeApproxReqSize calculates the approximate size of an HTTP request in bytes.

Parameters:
  - req: The *http.Request object representing the HTTP request.

Returns:
  - int: The approximate size of the HTTP request in bytes.

The computeApproxReqSize function calculates the size of an HTTP request by summing the lengths of various components, including the URL path, request method, protocol, headers, and host. It also takes into account the content length, if available.

Example usage:

	req := &http.Request{
		URL: &url.URL{
			Path: "/api/users",
		},
		Method:        "GET",
		Proto:         "HTTP/1.1",
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Host:          "example.com",
		ContentLength: 100,
	}

	size := computeApproxReqSize(req)
	fmt.Println(size) // Output: 152

Note: The computeApproxReqSize function provides an approximation of the request size and may not be exact in all cases. It is intended to be used for monitoring and analysis purposes.
*/
func computeApproxReqSize(req *http.Request) int {
	var sizeBytes int
	if req.URL != nil {
		sizeBytes = len(req.URL.Path)
	}

	sizeBytes += len(req.Method)
	sizeBytes += len(req.Proto)
	for name, values := range req.Header {
		sizeBytes += len(name)
		for _, value := range values {
			sizeBytes += len(value)
		}
	}
	sizeBytes += len(req.Host)

	if req.ContentLength != -1 {
		sizeBytes += int(req.ContentLength)
	}
	return sizeBytes
}
