// v1/requests/auth/login_request.go
package auth

import (
	"fmt"
	"response-std/app/pkg/response"
	"response-std/libs/external/services"

	"github.com/davecgh/go-spew/spew"
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
		Data:     r, // pakai data yang sudah di-bind
		Rules:    rules,
		Messages: messages,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()

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

		err := fmt.Errorf("%v", errorInterface)
		response.UnprocessableValidation(c, "Validation failed", err, errorInterface, "[LoginRequest.Validate]")
		spew.Dump(errors, "Validation errors", "\n errors from validation", errorInterface)
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
