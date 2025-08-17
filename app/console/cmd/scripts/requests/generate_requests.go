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
	"github.com/davecgh/go-spew/spew"
	"response-std/app/pkg/response"
	"response-std/libs/external/services"
)

// {{.CamelCase}}Request represents the request structure for {{.Name}}.
type {{.CamelCase}}Request struct {
	// TODO: Tambahkan field-field request-mu di sini
	// Contoh:
	// Name  string  ` + "`json:\"name\"`" + `
	// Email *string ` + "`json:\"email\"`" + `
}

// Validate validates the {{.Name}} request using govalidator.
func (r *{{.CamelCase}}Request) Validate(c *gin.Context) bool {
	rules := govalidator.MapData{
		// TODO: Tambahkan aturan validasi
		// "name": []string{"required", "min:2"},
	}

	messages := govalidator.MapData{
		// TODO: Tambahkan pesan error kustom
		// "name": {"required:Nama wajib diisi", "min:Minimal 2 karakter"},
	}

	opts := govalidator.Options{
		Rules:    rules,
		Messages: messages,
		Data:     r, // struct sudah dibind dari controller
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()

	// ✅ Validasi tambahan manual jika perlu
	// contoh: cek urutan tanggal, cek kondisi khusus, dll

	if len(e) > 0 {
		errors := make(map[string][]string)
		for field, msgs := range e {
			errors[field] = msgs
		}

		errorInterface := make(map[string]interface{}, len(errors))
		for k, v := range errors {
			errorInterface[k] = v
		}

		services.AppLogger.Debug("Validation failed", errorInterface)
		spew.Dump(errors, "Validation errors", "\n errors from validation", errorInterface)

		response.UnprocessableValidation(c, "Validation failed", nil, errorInterface, "[{{.CamelCase}}Request.Validate]")
		return false
	}

	return true
}

// GetValidatedData returns the validated data as map.
func (r *{{.CamelCase}}Request) GetValidatedData() map[string]interface{} {
	return map[string]interface{}{
		// TODO: Petakan field validated-mu
		// "name":  r.Name,
	}
}

// GetRules returns the validation rules (useful for documentation or testing).
func (r *{{.CamelCase}}Request) GetRules() govalidator.MapData {
	return govalidator.MapData{
		// TODO: Kembalikan aturan validasi di sini
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
