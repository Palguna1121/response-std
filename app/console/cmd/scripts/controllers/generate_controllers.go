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
		log.Fatal("Usage: go run app/console/cmd/scripts/controllers/generate_controllers.go [controller_name] [version]")
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
	outputPath := filepath.Join("app", "http", "controllers", version, fileName+".go")

	// siapkan data untuk template
	data := struct {
		Name      string
		CamelCase string
		LowerCase string
	}{
		Name:      controllerName,
		CamelCase: toCamelCase(controllerName),
		LowerCase: strings.ToLower(controllerName),
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
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"response-std/app/pkg/response"
	// TODO: Import your request structs here
	// "response-std/app/http/requests"
)

//delete this section if unused===============================
// {{.CamelCase}}Request represents the request structure for {{.Name}}
type {{.CamelCase}}Request struct {
	// TODO: Define your request fields here
	// Example:
	// Name        string ` + "`" + `json:"name" form:"name"` + "`" + `
	// Email       string ` + "`" + `json:"email" form:"email"` + "`" + `
	// Phone       string ` + "`" + `json:"phone" form:"phone"` + "`" + `
	// Description string ` + "`" + `json:"description" form:"description"` + "`" + `
}

// {{.CamelCase}}UpdateRequest represents the request structure for updating {{.Name}}
type {{.CamelCase}}UpdateRequest struct {
	// TODO: Define your update request fields here
	// Example:
	// Name        *string ` + "`" + `json:"name,omitempty" form:"name"` + "`" + `
	// Email       *string ` + "`" + `json:"email,omitempty" form:"email"` + "`" + `
	// Phone       *string ` + "`" + `json:"phone,omitempty" form:"phone"` + "`" + `
	// Description *string ` + "`" + `json:"description,omitempty" form:"description"` + "`" + `
}
//delete this section if unused================================


// {{.CamelCase}}Controller implements the {{.CamelCase}} controller.
type {{.CamelCase}}Controller struct{}

// New{{.CamelCase}}Controller returns a new instance of {{.CamelCase}}Controller.
func New{{.CamelCase}}Controller() *{{.CamelCase}}Controller {
	return &{{.CamelCase}}Controller{}
}

// List{{.CamelCase}} returns a list of {{.Name}}s.
func (ctl *{{.CamelCase}}Controller) List{{.CamelCase}}(c *gin.Context) {
	data, exists := c.Get("{{.LowerCase}}")
	if !exists {
		response.NotFound(c, "{{.Name}}s not found", nil, "[List{{.CamelCase}}]")
		return
	}

	response.Success(c, "List of {{.Name}}s retrieved successfully", data)
}

// Get{{.CamelCase}}ByID returns a single {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Get{{.CamelCase}}ByID(c *gin.Context) {
	data, exists := c.Get("{{.LowerCase}}")
	if !exists {
		response.NotFound(c, "{{.Name}} not found", nil, "[Get{{.CamelCase}}ByID]")
		return
	}

	response.Success(c, "{{.Name}} retrieved successfully", data)
}

// Create{{.CamelCase}} creates a new {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Create{{.CamelCase}}(c *gin.Context) {
	// TODO: Use your request struct and validate
	// var req requests.{{.CamelCase}}Request
	// if !req.Validate(c) {
	//     return
	// }

	// TODO: Process the validated data
	// - Save to database
	// - Call service layer

	response.Created(c, "{{.Name}} created successfully", nil)
}

// Update{{.CamelCase}} updates an existing {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Update{{.CamelCase}}(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID format", err, "[Update{{.CamelCase}}]")
		return
	}

	// TODO: Use your request struct and validate
	// var req requests.{{.CamelCase}}UpdateRequest
	// if !req.Validate(c) {
	//     return
	// }

	// TODO: Process the validated data
	// - Check if record exists
	// - Update in database
	_ = id // Use the ID for database operations

	response.Success(c, "{{.Name}} updated successfully", nil)
}

// Delete{{.CamelCase}} deletes a {{.Name}}.
func (ctl *{{.CamelCase}}Controller) Delete{{.CamelCase}}(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID format", err, "[Delete{{.CamelCase}}]")
		return
	}

	// TODO: Process the deletion
	// Example:
	// - Check if record exists
	// - Check if safe to delete (no foreign key constraints)
	// - Delete from database
	// - Handle business logic
	_ = id // Use the ID for database operations

	response.Success(c, "{{.Name}} deleted successfully", nil)
}
`

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}
