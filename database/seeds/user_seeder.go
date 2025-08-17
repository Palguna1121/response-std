package seeds

import (
	"log"

	"response-std/app/helpers/helper"
	"response-std/app/models/entities"
	"response-std/app/pkg/permissions"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, spatie *permissions.Spatie) {
	users := []entities.User{
		{Name: "kuroneko", Email: "kuroneko@gmail.com"},
		{Name: "shiroinu", Email: "shiroinu@gmail.com"},
	}

	var saved []entities.User

	for _, u := range users {
		hashed, _ := helper.HashPassword(u.Name)
		u.Password = hashed

		var existing entities.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			log.Printf("User %s already exists", existing.Name)
			saved = append(saved, existing)
			continue
		}

		if err := db.Create(&u).Error; err != nil {
			log.Fatalf("Failed to create user %s: %v", u.Name, err)
		}

		saved = append(saved, u)
	}

	log.Println("Users seeded.")

	for _, u := range saved {
		var role string
		if u.Name == "kuroneko" {
			role = "admin"
		} else {
			role = "user"
		}
		if err := spatie.AssignRole(u.ID, role); err != nil {
			log.Printf("Failed to assign role %s to %s: %v", role, u.Name, err)
		}
	}
}
