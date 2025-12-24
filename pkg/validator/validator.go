// Package validator provides request validation.
package validator

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator wraps go-playground validator
type Validator struct {
	validator *validator.Validate
}

// New creates a new validator
func New() *Validator {
	v := validator.New()

	// Use JSON tag names in error messages
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators here
	// v.RegisterValidation("custom", customValidation)

	return &Validator{validator: v}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return v.formatErrors(err)
	}
	return nil
}

// formatErrors converts validation errors to a user-friendly format
func (v *Validator) formatErrors(err error) *ValidationError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationError{
			Message: err.Error(),
		}
	}

	errors := make(map[string]string)
	for _, e := range validationErrors {
		field := e.Field()
		errors[field] = v.formatMessage(e)
	}

	return &ValidationError{
		Message: "Validation failed",
		Errors:  errors,
	}
}

func (v *Validator) formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + e.Param() + " characters"
	case "max":
		return "Must be at most " + e.Param() + " characters"
	case "gte":
		return "Must be greater than or equal to " + e.Param()
	case "lte":
		return "Must be less than or equal to " + e.Param()
	case "uuid":
		return "Must be a valid UUID"
	case "url":
		return "Must be a valid URL"
	case "oneof":
		return "Must be one of: " + e.Param()
	default:
		return "Invalid value"
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// HTTPError returns an echo HTTP error
func (e *ValidationError) HTTPError() *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, e)
}

// EchoValidator implements echo.Validator
type EchoValidator struct {
	validator *Validator
}

// NewEchoValidator creates a validator for Echo
func NewEchoValidator() *EchoValidator {
	return &EchoValidator{
		validator: New(),
	}
}

// Validate implements echo.Validator
func (v *EchoValidator) Validate(i interface{}) error {
	return v.validator.Validate(i)
}
