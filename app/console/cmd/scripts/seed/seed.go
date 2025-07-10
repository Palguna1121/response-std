package main

import (
	"log"
	"response-std/app/helpers/helper"
	"response-std/app/models/entities"
	"response-std/app/pkg/permissions"
	"response-std/config"

	"gorm.io/gorm"
)

var spatie *permissions.Spatie

func InitPermissions(db *gorm.DB) {
	spatie = permissions.NewSpatie(db)
}

func main() {
	config.InitConfig()
	config.LoadDBMysql()
	InitPermissions(config.DB)

	// Run seeder
	seedUsers(config.DB, spatie)

	log.Println("Seeding completed successfully!")
}

func seedUsers(db *gorm.DB, spatie *permissions.Spatie) {
	// Data user awal
	rawUsers := []entities.User{
		{Name: "kuroneko", Email: "kuroneko@gmail.com"},
		{Name: "shiroinu", Email: "shiroinu@gmail.com"},
	}

	var savedUsers []entities.User // akan menyimpan user yang sudah pasti valid ID-nya

	// Hash password dan buat user
	for _, user := range rawUsers {
		hashed, _ := helper.HashPassword(user.Name) // pake nama sebagai password
		user.Password = hashed

		// check if user exists
		var existing entities.User
		err := db.Where("email = ?", user.Email).First(&existing).Error
		if err == nil {
			log.Printf("User %s already exists", existing.Name)
			savedUsers = append(savedUsers, existing)
			continue
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user %s: %v", user.Name, err)
		}

		savedUsers = append(savedUsers, user)
	}

	log.Println("Users seeded")

	// Seed roles
	seedRoles(spatie)

	// Assign roles secara aman
	for _, user := range savedUsers {
		var roleName string
		if user.Name == "kuroneko" {
			roleName = "admin"
		} else if user.Name == "shiroinu" {
			roleName = "staff"
		}
		if err := spatie.AssignRole(user.ID, roleName); err != nil {
			log.Printf("Gagal assign role %s ke user %s: %v", roleName, user.Name, err)
		}
	}

	log.Println("Roles assigned")
}

func seedRoles(spatie *permissions.Spatie) {
	roles := []string{"admin", "staff"}

	for _, roleName := range roles {
		_, err := spatie.FindRoleByName(roleName)
		if err != nil {
			if _, createErr := spatie.CreateRole(roleName, "web"); createErr != nil {
				log.Fatalf("Gagal membuat role %s: %v", roleName, createErr)
			} else {
				log.Printf("Role %s berhasil dibuat", roleName)
			}
		} else {
			log.Printf("Role %s sudah ada, tidak perlu dibuat ulang", roleName)
		}
	}
}
