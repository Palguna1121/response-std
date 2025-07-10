// v1/requests/auth/login_request.go
package auth

import (
	"response-std/app/pkg/response"
	"response-std/libs/external/services"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate validates the login request
func (r *LoginRequest) Validate(c *gin.Context) bool {
	// Define validation rules
	rules := govalidator.MapData{
		"username": []string{"required", "min:3"},
		"password": []string{"required", "min:8"},
	}

	// Custom messages (optional, mirip Laravel)
	messages := govalidator.MapData{
		"username": []string{
			"required:Username wajib diisi",
			"min:Username minimal 3 karakter",
		},
		"password": []string{
			"required:Password wajib diisi",
			"min:Password minimal 8 karakter",
		},
	}

	// Create validator options
	opts := govalidator.Options{
		Request:  c.Request,
		Rules:    rules,
		Messages: messages,
	}

	// Validate request
	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		// Format error mirip Laravel
		errors := make(map[string][]string)
		for field, msgs := range e {
			errors[field] = msgs
		}

		// Convert errors to map[string]interface{} for logging
		errorInterface := make(map[string]interface{}, len(errors))
		for k, v := range errors {
			errorInterface[k] = v
		}
		services.AppLogger.Debug("Validation failed", errorInterface)

		response.UnprocessableEntity(c, "Validation failed", nil, "[LoginRequest]")
		return false
	}

	// Bind validated data to struct
	if err := c.ShouldBindJSON(r); err != nil {
		response.BadRequest(c, "Invalid JSON format", err, "[LoginRequest]")
		return false
	}

	return true
}

// GetValidatedData returns the validated data
func (r *LoginRequest) GetValidatedData() map[string]interface{} {
	return map[string]interface{}{
		"username": r.Username,
		"password": r.Password,
	}
}
