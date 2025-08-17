package seeds

import (
	"log"

	"response-std/app/pkg/permissions"
)

func SeedRoles(spatie *permissions.Spatie) {
	roles := []string{"admin", "user"}

	for _, role := range roles {
		_, err := spatie.FindRoleByName(role)
		if err != nil {
			if _, err := spatie.CreateRole(role, "web"); err != nil {
				log.Fatalf("Failed to create role %s: %v", role, err)
			}
			log.Printf("Created role: %s", role)
		} else {
			log.Printf("Role %s already exists", role)
		}
	}
}
