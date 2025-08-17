package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

const seedTemplate = `package seeds

import (
	"log"
	"gorm.io/gorm"
)

func {{.FuncName}}(db *gorm.DB) {
	log.Println("Seeding: {{.Name}}")
	// TODO: implement seeder logic here
}
`

type SeederData struct {
	Name     string // e.g. user
	FuncName string // e.g. UserSeeder
}

// Convert to PascalCase
func pascalCase(s string) string {
	return strings.Map(func(r rune) rune {
		return r
	}, strings.Title(s))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run generate_seeder.go <name>")
		os.Exit(1)
	}

	name := os.Args[1]                            // misal: user
	fileName := fmt.Sprintf("%s_seeder.go", name) // user_seeder.go
	funcName := toPascalCase(name) + "Seeder"     // UserSeeder

	outputPath := filepath.Join("database", "seeds", fileName)

	// Cek jika sudah ada
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("Seeder %s already exists\n", fileName)
		os.Exit(1)
	}

	// Buat file
	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	// Template
	tmpl, err := template.New("seeder").Parse(seedTemplate)
	if err != nil {
		fmt.Printf("Template parse error: %v\n", err)
		os.Exit(1)
	}

	data := SeederData{
		Name:     name,
		FuncName: funcName,
	}

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Printf("Template execution error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Seeder %s generated at %s\n", fileName, outputPath)
}

// Convert "user" to "User", "post_tag" to "PostTag"
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = string(unicode.ToUpper(rune(part[0]))) + part[1:]
	}
	return strings.Join(parts, "")
}
