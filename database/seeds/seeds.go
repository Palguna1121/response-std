package seeds

import (
	"log"

	"response-std/app/pkg/permissions"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB, spatie *permissions.Spatie) {
	log.Println("Seeding: Roles")
	SeedRoles(spatie)

	log.Println("Seeding: Users")
	SeedUsers(db, spatie)
}
