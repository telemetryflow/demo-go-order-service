// Package api provides embedded API documentation.
package api

import (
	_ "embed"
)

//go:embed swagger.json
var SwaggerJSON []byte
