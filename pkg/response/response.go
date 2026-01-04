// Package response provides HTTP response helpers.
package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/telemetry/logs"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Meta represents response metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalCount int64 `json:"total_count,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// Success sends a success response
func Success(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c echo.Context, data interface{}, message string) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a 204 no content response
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// Paginated sends a paginated response
func Paginated(c echo.Context, data interface{}, totalCount int64, page, pageSize int) error {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		},
	})
}

// Error sends an error response
func Error(c echo.Context, status int, code, message string) error {
	logAttrs := map[string]interface{}{
		"code":       code,
		"message":    message,
		"status":     status,
		"method":     c.Request().Method,
		"path":       c.Request().URL.Path,
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
		"remote_ip":  c.RealIP(),
	}

	if status >= 500 {
		logs.Error("API error", logAttrs)
	} else if status >= 400 {
		logs.Warn("API client error", logAttrs)
	}

	return c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// ErrorWithDetails sends an error response with details
func ErrorWithDetails(c echo.Context, status int, code, message string, details map[string]string) error {
	logAttrs := map[string]interface{}{
		"code":       code,
		"message":    message,
		"status":     status,
		"method":     c.Request().Method,
		"path":       c.Request().URL.Path,
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
		"remote_ip":  c.RealIP(),
		"details":    details,
	}

	if status >= 500 {
		logs.Error("API error", logAttrs)
	} else if status >= 400 {
		logs.Warn("API client error", logAttrs)
	}

	return c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest sends a 400 bad request response
func BadRequest(c echo.Context, message string) error {
	return Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c echo.Context, message string) error {
	return Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 forbidden response
func Forbidden(c echo.Context, message string) error {
	return Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a 404 not found response
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 conflict response
func Conflict(c echo.Context, message string) error {
	return Error(c, http.StatusConflict, "CONFLICT", message)
}

// InternalError sends a 500 internal server error response
func InternalError(c echo.Context, message string) error {
	return Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

// ValidationError sends a validation error response
func ValidationError(c echo.Context, details map[string]string) error {
	return ErrorWithDetails(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}
