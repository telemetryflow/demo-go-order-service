// Package handler provides HTTP handlers for Swagger documentation.
package handler

import (
	_ "embed"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/docs/api"
)

//go:embed swagger_ui.html
var swaggerUIHTML string

// SwaggerHandler handles swagger documentation requests
type SwaggerHandler struct {
	title string
}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler(title string) *SwaggerHandler {
	return &SwaggerHandler{title: title}
}

// RegisterRoutes registers swagger routes
func (h *SwaggerHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/docs", h.SwaggerUI)
	e.GET("/docs/", h.SwaggerUI)
	e.GET("/docs/spec/swagger.json", h.SwaggerSpec)
}

// SwaggerUI serves the Swagger UI page
func (h *SwaggerHandler) SwaggerUI(c echo.Context) error {
	tmpl, err := template.New("swagger").Parse(swaggerUIHTML)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to load Swagger UI")
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return tmpl.Execute(c.Response().Writer, map[string]string{
		"SpecURL": "/docs/spec/swagger.json",
		"Title":   h.title,
	})
}

// SwaggerSpec serves the embedded OpenAPI specification
func (h *SwaggerHandler) SwaggerSpec(c echo.Context) error {
	return c.JSONBlob(http.StatusOK, api.SwaggerJSON)
}
