// cmd/scripts/requests/generate_requests.go
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/scripts/requests/generate_requests.go <request_name> [version]")
	}

	requestName := os.Args[1]
	version := "v1"
	if len(os.Args) >= 3 {
		version = os.Args[2]
	}

	fileName := requestName
	if !strings.Contains(requestName, "request") {
		fileName = requestName + "_request"
	}
	outputPath := filepath.Join("app", "http", "requests", version, fileName+".go")

	// siapkan data untuk template
	data := struct {
		Name      string
		CamelCase string
		Package   string
	}{
		Name:      requestName,
		CamelCase: toCamelCase(requestName),
		Package:   strings.ToLower(requestName),
	}

	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	tmpl, err := template.New("request").Parse(requestTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	log.Println("✅ Generated:", outputPath)
}

var requestTemplate = `package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"response-std/app/pkg/response"
	"response-std/libs/external/services"
)

// {{.CamelCase}}Request represents the request structure for {{.Name}}.
type {{.CamelCase}}Request struct {
	// TODO: Add your request fields here
	// Example:
	// Name  string ` + "`json:\"name\"`" + `
	// Email string ` + "`json:\"email\"`" + `
}

// Validate validates the {{.Name}} request using govalidator.
func (r *{{.CamelCase}}Request) Validate(c *gin.Context) bool {
	// Define validation rules
	rules := govalidator.MapData{
		// TODO: Add your validation rules here
		// Example:
		// "name":  []string{"required", "min:2", "max:255"},
		// "email": []string{"required", "email"},
	}

	// Custom messages (optional, mirip Laravel)
	messages := govalidator.MapData{
		// TODO: Add your custom messages here
		// Example:
		// "name": []string{
		//     "required:Nama wajib diisi",
		//     "min:Nama minimal 2 karakter",
		//     "max:Nama maksimal 255 karakter",
		// },
		// "email": []string{
		//     "required:Email wajib diisi",
		//     "email:Format email tidak valid",
		// },
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
		response.BadRequest(c, "Invalid JSON format", err, "[{{.CamelCase}}Request]")
		return false
	}

	return true
}

// GetValidatedData returns the validated data as map.
func (r *{{.CamelCase}}Request) GetValidatedData() map[string]interface{} {
	return map[string]interface{}{
		// TODO: Map your validated fields here
		// Example:
		// "name":  r.Name,
		// "email": r.Email,
	}
}

// GetRules returns the validation rules (useful for documentation or testing).
func (r *{{.CamelCase}}Request) GetRules() govalidator.MapData {
	return govalidator.MapData{
		// TODO: Return your validation rules here
		// Example:
		// "name":  []string{"required", "min:2", "max:255"},
		// "email": []string{"required", "email"},
	}
}
`

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
