package main

import (
	"log"
	"response-std/config"
	"response-std/core/helper"
	"response-std/core/models/entities"
	"time"

	"gorm.io/gorm"
)

func main() {
	config.InitConfig()
	config.LoadDB()

	// Run seeder
	seedUsers(config.DB)
	seedRoles(config.DB)
	assignRoles(config.DB)

	log.Println("Seeding completed successfully!")
}

func seedUsers(db *gorm.DB) {
	hashedPassword1, err := helper.HashPassword("kuroneko")
	if err != nil {
		log.Fatalf("Failed to hash password for kuroneko: %v", err)
	}
	hashedPassword2, err := helper.HashPassword("shiroinu")
	if err != nil {
		log.Fatalf("Failed to hash password for shiroinu: %v", err)
	}

	users := []entities.User{
		{
			ID:        1,
			Name:      "kuroneko",
			Email:     "kuroneko@mail.com",
			Password:  hashedPassword1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "shiroinu",
			Email:     "shiroinu@mail.com",
			Password:  hashedPassword2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		db.Create(&user)
	}
	log.Println("Users seeded")
}

func seedRoles(db *gorm.DB) {
	roles := []entities.Roles{
		{
			ID:        1,
			Name:      "admin",
			GuardName: "web",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "staff",
			GuardName: "web",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, role := range roles {
		db.Create(&role)
	}
	log.Println("Roles seeded")
}

func assignRoles(db *gorm.DB) {
	// Assign admin role to user 1
	adminAssignment := entities.ModelHasRoles{
		RoleID:    1,
		ModelType: "entities.User", // atau "App\\entities\\User" sesuai kebutuhan
		ModelID:   1,
	}
	db.Create(&adminAssignment)

	// Assign staff role to user 2
	staffAssignment := entities.ModelHasRoles{
		RoleID:    2,
		ModelType: "entities.User",
		ModelID:   2,
	}
	db.Create(&staffAssignment)

	log.Println("Roles assigned to users")
}
