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
		fmt.Println("Usage: go run cmd/migrate.go [up|down|drop|...]")
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

	cmd := exec.Command(
		"migrate",
		"-database", url,
		"-path", fmt.Sprintf("%s/database/migrations", config.ENV.API_VERSION),
		action,
	)

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
