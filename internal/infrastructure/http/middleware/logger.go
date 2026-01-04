// Package middleware provides HTTP middleware.
package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/telemetry/logs"
	"github.com/telemetryflow/order-service/telemetry/metrics"
	"github.com/telemetryflow/order-service/telemetry/traces"
)

// Logger returns a logging middleware with tracing support
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			ctx := req.Context()

			// Start trace span for this request
			spanID, _ := traces.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", req.Method, req.URL.Path), map[string]interface{}{
				"http.method":      req.Method,
				"http.url":         req.URL.String(),
				"http.path":        req.URL.Path,
				"http.user_agent":  req.UserAgent(),
				"http.remote_addr": c.RealIP(),
				"http.request_id":  c.Response().Header().Get(echo.HeaderXRequestID),
			})

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)
			res := c.Response()

			// Record metrics
			metrics.RecordHTTPRequest(req.Method, req.URL.Path, res.Status, duration.Seconds())

			// Log request at appropriate level based on status
			logAttrs := map[string]interface{}{
				"method":     req.Method,
				"path":       req.URL.Path,
				"status":     res.Status,
				"duration":   duration.String(),
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				"user_agent": req.UserAgent(),
				"remote_ip":  c.RealIP(),
			}

			if res.Status >= 500 {
				logs.Error("HTTP request failed", logAttrs)
			} else if res.Status >= 400 {
				logs.Warn("HTTP request client error", logAttrs)
			} else {
				logs.Info("HTTP request", logAttrs)
			}

			// End trace span
			var spanErr error
			if res.Status >= 500 {
				spanErr = fmt.Errorf("HTTP %d", res.Status)
			}
			_ = traces.EndSpan(ctx, spanID, spanErr)

			return err
		}
	}
}
