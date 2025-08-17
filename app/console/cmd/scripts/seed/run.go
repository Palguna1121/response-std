package main

import (
	"log"

	"response-std/app/pkg/permissions"
	"response-std/config"
	"response-std/database/seeds"
)

func main() {
	config.InitConfig()
	config.LoadDBMysql()

	if config.DB == nil {
		log.Fatal("config.DB is nil after InitDB")
	}

	log.Println("Running seeders...")
	seeds.SeedAll(config.DB, permissions.NewSpatie(config.DB))
	log.Println("Seeding completed.")
}
