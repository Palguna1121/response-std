package config

import (
	"fmt"
	"log"

	"github.com/Palguna1121/response-std/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ENV.DB_USER, ENV.DB_PASSWORD, ENV.DB_HOST, ENV.DB_PORT, ENV.DB_NAME)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// MigrateDB()
}

func MigrateDB() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Roles{},
		&models.Permission{},
		&models.PersonalAccessTokens{},
		&models.ModelHasPermissions{},
		&models.ModelHasRoles{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}
