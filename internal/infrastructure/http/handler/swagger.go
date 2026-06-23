// Package handler provides Swagger API documentation handler.
//
// TelemetryFlow Order Service - AI-Powered Observability & Incident Response Management (IRM) Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
