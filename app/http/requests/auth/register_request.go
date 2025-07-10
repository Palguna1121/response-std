// v1/requests/auth/register_request.go
package auth

import (
	"response-std/app/pkg/response"
	"response-std/libs/external/services"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate(c *gin.Context) bool {
	// Define validation rules
	rules := govalidator.MapData{
		"name":                  []string{"required", "min:2", "max:255"},
		"email":                 []string{"required", "email"},
		"password":              []string{"required", "min:8"},
		"password_confirmation": []string{"required"},
	}

	// Custom messages
	messages := govalidator.MapData{
		"name": []string{
			"required:Nama wajib diisi",
			"min:Nama minimal 2 karakter",
			"max:Nama maksimal 255 karakter",
		},
		"email": []string{
			"required:Email wajib diisi",
			"email:Format email tidak valid",
		},
		"password": []string{
			"required:Password wajib diisi",
			"min:Password minimal 8 karakter",
		},
		"password_confirmation": []string{
			"required:Konfirmasi password wajib diisi",
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

		response.UnprocessableEntity(c, "Validation failed", nil, "[RegisterRequest]")
		return false
	}

	// Bind validated data to struct
	if err := c.ShouldBindJSON(r); err != nil {
		response.BadRequest(c, "Invalid JSON format", err, "[RegisterRequest]")
		return false
	}

	// Custom validation: password confirmation
	if r.Password != r.PasswordConfirmation {
		errors := map[string][]string{
			"password_confirmation": {"Password dan konfirmasi password tidak cocok"},
		}
		// Convert errors to map[string]interface{} for logging
		errorInterface := make(map[string]interface{}, len(errors))
		for k, v := range errors {
			errorInterface[k] = v
		}
		services.AppLogger.Debug("Validation failed", errorInterface)

		response.UnprocessableEntity(c, "Validation failed", nil, "[RegisterRequest]")
		return false
	}

	return true
}

// GetValidatedData returns the validated data
func (r *RegisterRequest) GetValidatedData() map[string]interface{} {
	return map[string]interface{}{
		"name":                  r.Name,
		"email":                 r.Email,
		"password":              r.Password,
		"password_confirmation": r.PasswordConfirmation,
	}
}
