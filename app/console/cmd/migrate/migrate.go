package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"response-std/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run app/console/cmd/migrate/migrate.go [up|down|drop|force VERSION] [ver]")
		return
	}

	action := os.Args[1]

	config.InitConfig()

	url := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s",
		config.ENV.DB_USER,
		config.ENV.DB_PASSWORD,
		config.ENV.DB_HOST,
		config.ENV.DB_PORT,
		config.ENV.DB_NAME,
	)

	migrationPath := "database/migrations"

	// Handle `force` needing version number after 'force'
	args := []string{"-database", url, "-path", migrationPath, action}
	if action == "force" && len(os.Args) >= 4 {
		forceVersion := os.Args[3]
		args = append(args, forceVersion)
	}

	cmd := exec.Command("migrate", args...)

	// Handle destructive commands (down, drop)
	if action == "down" || action == "drop" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println("Failed to pipe input:", err)
			return
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, "y\n")
		}()
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Migration failed:", err)
	}
}
