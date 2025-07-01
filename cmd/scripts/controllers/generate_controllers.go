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
		log.Fatal("Usage: go run cmd/scripts/controllers/generate_controllers.go <controller_name> [version]")
	}

	controllerName := os.Args[1]
	version := "v1"
	if len(os.Args) >= 3 {
		version = os.Args[2]
	}

	fileName := controllerName
	if !strings.Contains(controllerName, "controller") {
		fileName = controllerName + "_controller"
	}
	outputPath := filepath.Join(version, "controllers", fileName+".go")

	// siapkan data untuk template
	data := struct {
		Name      string
		CamelCase string
	}{
		Name:      controllerName,
		CamelCase: toCamelCase(controllerName),
	}

	err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	tmpl, err := template.New("controller").Parse(controllerTemplate)
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

var controllerTemplate = `package controllers

import (
	"github.com/gin-gonic/gin"

	"response-std/core/response"
)

// {{.CamelCase}}Controller implements the {{.CamelCase}} controller.
type {{.CamelCase}}Controller struct{}

// New{{.CamelCase}}Controller returns a new instance of {{.CamelCase}}Controller.
func New{{.CamelCase}}Controller() *{{.CamelCase}}Controller {
	return &{{.CamelCase}}Controller{}
}

// List{{.CamelCase}} returns a list of {{.Name}}s.
func (ctl *{{.CamelCase}}Controller) List{{.Name}}(c *gin.Context) {
	data, exists := c.Get("{{.Name}}s")
	if !exists {
		response.NotFound(c, "{{.Name}}s not found", nil, "[List{{.CamelCase}}]")
		return
	}

	response.Success(c, "List of {{.Name}}s retrieved successfully", data)
}

// Get{{.CamelCase}}ByID returns a single {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Get{{.CamelCase}}ByID(c *gin.Context) {
	data, exists := c.Get("{{.Name}}")
	if !exists {
		response.NotFound(c, "Not found", nil, "[Get{{.CamelCase}}ByID]")
		return
	}

	response.Success(c, "Get {{.Name}} by ID", data)
}

// Create{{.CamelCase}} creates a new {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Create{{.CamelCase}}(c *gin.Context) {
	// TODO: validate input
	// TODO: sanitize input
	// TODO: save to database

	response.Created(c, "Create {{.Name}}", nil)
}

// Update{{.CamelCase}} updates an existing {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Update{{.CamelCase}}(c *gin.Context) {
	// TODO: validate input
	// TODO: sanitize input
	// TODO: update in database

	response.Success(c, "Update {{.Name}}", nil)
}

// Delete{{.CamelCase}} deletes a {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Delete{{.CamelCase}}(c *gin.Context) {
	// TODO: validate input
	// TODO: delete from database

	response.Success(c, "Delete {{.Name}}", nil)
}
`

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
