package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"response-std/app/helpers/helper"
	"response-std/app/models/entities"
	"response-std/app/pkg/permissions"
	"response-std/app/pkg/response"
	"response-std/libs/external/services"
)

type UserController struct {
	DB         *gorm.DB
	Permission *permissions.Spatie
}

func NewUserController(db *gorm.DB, perm *permissions.Spatie) *UserController {
	return &UserController{
		DB:         db,
		Permission: perm,
	}
}

func (ctl *UserController) ListUser(c *gin.Context) {
	var users []entities.User
	if err := ctl.DB.Preload("Roles").Find(&users).Error; err != nil {
		response.InternalServerError(c, "Failed to fetch users", err, "[ListUser]")
		return
	}

	data := make([]gin.H, len(users))
	for i, user := range users {
		data[i] = gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		}
	}

	response.Success(c, "List of users retrieved successfully", data)
}

func (ctl *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user entities.User
	if err := ctl.DB.Preload("Roles").First(&user, id).Error; err != nil {
		response.NotFound(c, "User not found", err, "[GetUserByID]")
		return
	}
	response.Success(c, "User retrieved successfully", user)
}

func (ctl *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input", err, "[CreateUser]")
		return
	}

	// Check if email already exists
	var count int64
	ctl.DB.Model(&entities.User{}).Where("email = ?", input.Email).Count(&count)
	if count > 0 {
		response.UnprocessableEntity(c, "Email already in use", nil, "[CreateUser]")
		return
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		response.InternalServerError(c, "Failed to hash password", err, "[CreateUser]")
		return
	}

	user := entities.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ctl.DB.Create(&user).Error; err != nil {
		response.InternalServerError(c, "Failed to create user", err, "[CreateUser]")
		return
	}

	// Check and assign "pelanggan" role
	roleName := "pelanggan"
	_, err = ctl.Permission.FindRoleByName(roleName)
	if err != nil {
		// Role not found, create it
		_, err = ctl.Permission.CreateRole(roleName, "web")
		if err != nil {
			response.InternalServerError(c, "Failed to create role pelanggan", err, "[CreateUser]")
			return
		}
	}

	if err := ctl.Permission.AssignRole(user.ID, roleName); err != nil {
		response.InternalServerError(c, "Failed to assign role to user", err, "[CreateUser]")
		return
	}

	response.Success(c, "User created successfully", nil)
}

func (ctl *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user entities.User
	if err := ctl.DB.First(&user, id).Error; err != nil {
		response.NotFound(c, "User not found", err, "[UpdateUser]")
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Invalid input", err, "[UpdateUser]")
		return
	}

	// Check if email already exists and different from current email
	if input.Email != nil && *input.Email != user.Email {
		var count int64
		ctl.DB.Model(&entities.User{}).Where("email = ?", *input.Email).Count(&count)
		if count > 0 {
			response.UnprocessableEntity(c, "Email already in use", nil, "[UpdateUser]")
			return
		}
	}

	// Update user
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		hashed, _ := helper.HashPassword(*input.Password)
		user.Password = hashed
	}
	user.UpdatedAt = time.Now()

	if err := ctl.DB.Save(&user).Error; err != nil {
		response.InternalServerError(c, "Failed to update user", err, "[UpdateUser]")
		return
	}

	response.Success(c, "User updated successfully", user)
}

func (ctl *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user entities.User
	if err := ctl.DB.First(&user, id).Error; err != nil {
		response.NotFound(c, "User not found", err, "[UpdateUser]")
		return
	}
	if err := ctl.DB.Delete(&user, id).Error; err != nil {
		response.InternalServerError(c, "Failed to delete user", err, "[DeleteUser]")
		return
	}
	response.Success(c, "User deleted successfully", nil)
}

// sengaja error
func (ctl *UserController) ErrorDebug(c *gin.Context) {
	// Simulate an error
	services.AppLogger.Debug("Debug test message", nil)
	response.Success(c, "Debug test message", nil)
}

func (ctl *UserController) ErrorInfo(c *gin.Context) {
	// Simulate an error
	services.AppLogger.Info("Info test message", nil)
	response.Success(c, "Info test message", nil)
}

func (ctl *UserController) ErrorWarn(c *gin.Context) {
	// Simulate an error
	services.AppLogger.Warn("Warning test message", nil)
	response.Success(c, "Warning test message", nil)
}

func (ctl *UserController) ErrorError(c *gin.Context) {
	// Simulate an error
	services.AppLogger.Error("Error test message", nil, nil)
	response.Success(c, "Error test message", nil)
}

func (ctl *UserController) ErrorCritical(c *gin.Context) {
	// Simulate an error
	services.AppLogger.Critical("Critical test message", nil, nil)
	response.Success(c, "Critical test message", nil)
}

// func TestLoggingChannels(logger *services.Logger) {
// 	// logger.Debug("Debug test message", nil)
// 	// logger.Info("Info test message", nil)
// 	// logger.Warn("Warning test message", nil)
// 	// logger.Error("Error test message", nil, nil)
// 	// logger.Critical("Critical test message", nil, nil)
// }
