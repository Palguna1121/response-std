package auth

import (
	"response-std/app/models/entities"

	"github.com/gin-gonic/gin"
)

func GetAuthenticatedUser(c *gin.Context) (entities.User, bool) {
	userInterface, exists := c.Get("user")
	if !exists {
		return entities.User{}, false
	}

	user, ok := userInterface.(entities.User)
	return user, ok
}
