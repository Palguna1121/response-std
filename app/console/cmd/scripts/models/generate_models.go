package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run app/console/cmd/scripts/controllers/generate_models.go <table_name>")
	}
	table := os.Args[1]

	filename := findMigrationFor(table)
	if filename == "" {
		log.Fatalf("Migration for table '%s' not found", table)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	structDef := parseCreateTable(string(content), table)
	writeModelFile(table, structDef)
}

func findMigrationFor(table string) string {
	pattern := fmt.Sprintf("database/migrations/*_create_%s_table.up.sql", table)
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func parseCreateTable(sql, table string) string {
	lines := strings.Split(sql, "\n")
	var fields []string
	insideTable := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		re := regexp.MustCompile(`(?i)^create table\s+` + "`?" + table + "`?" + `\s*\(`)
		if re.MatchString(line) {
			insideTable = true
			continue
		}

		if insideTable {
			if strings.HasPrefix(line, ");") || line == ")" || line == ");" {
				break
			}
			line = strings.TrimSuffix(line, ",")
			field := parseField(line)
			if field != "" {
				fields = append(fields, field)
			}
		}
	}

	structLines := []string{
		"package entities",
		"",
		"import \"time\"",
		"",
		"type " + toCamelCase(table) + " struct {",
	}
	structLines = append(structLines, fields...)
	structLines = append(structLines, "}")
	return strings.Join(structLines, "\n")
}

func parseField(line string) string {
	line = strings.TrimSuffix(line, ",")
	re := regexp.MustCompile("`?(\\w+)`?\\s+(\\w+)")
	m := re.FindStringSubmatch(line)
	if len(m) < 3 {
		return ""
	}
	column := m[1]
	sqlType := m[2]
	name := toCamelCase(column)
	goType := sqlTypeToGoType(sqlType)

	// Tambahan parsing atribut
	gormTags := []string{fmt.Sprintf("column:%s", column)}

	if strings.Contains(strings.ToLower(line), "primary key") {
		gormTags = append(gormTags, "primaryKey")
	}
	if strings.Contains(strings.ToLower(line), "auto_increment") || strings.Contains(strings.ToLower(line), "auto increment") {
		gormTags = append(gormTags, "autoIncrement")
	}
	if !strings.Contains(strings.ToLower(line), "not null") {
		gormTags = append(gormTags, "null") // bisa kamu ubah jadi 'gorm:"null"' jika ingin pakai pointer
	}

	tag := fmt.Sprintf("`gorm:\"%s\" json:\"%s\"`", strings.Join(gormTags, ";"), column)
	return fmt.Sprintf("\t%s %s %s", name, goType, tag)
}

func sqlTypeToGoType(sqlType string) string {
	sqlType = strings.ToLower(sqlType)
	switch {
	case strings.Contains(sqlType, "int"):
		return "int"
	case strings.Contains(sqlType, "char"), strings.Contains(sqlType, "text"):
		return "string"
	case strings.Contains(sqlType, "float"), strings.Contains(sqlType, "double"), strings.Contains(sqlType, "decimal"):
		return "float64"
	case strings.Contains(sqlType, "bool"):
		return "bool"
	case strings.Contains(sqlType, "datetime"), strings.Contains(sqlType, "timestamp"):
		return "time.Time"
	default:
		return "string"
	}
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func writeModelFile(table, structContent string) {
	outputDir := filepath.Join("app", "models", "entities")
	os.MkdirAll(outputDir, os.ModePerm)

	outputPath := filepath.Join(outputDir, table+".go")
	err := os.WriteFile(outputPath, []byte(structContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write model file: %v", err)
	}
	fmt.Println("✅ Generated:", outputPath)
}
