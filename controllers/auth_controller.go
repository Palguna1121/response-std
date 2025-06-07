package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"response-std/helper"
	"response-std/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// ---------------------------
// LOGIN
// ---------------------------
func (a *AuthController) Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			helper.UnprocessableEntity(c, "Invalid credentials")
			return
		}

		var user models.User
		err := db.Preload("Roles.Permissions").Preload("Permissions").
			Where("email = ?", input.Email).First(&user).Error

		if err != nil {
			helper.NotFound(c, "User not found")
			return
		}

		if !user.CheckPassword(input.Password) {
			helper.UnprocessableEntity(c, "Invalid credentials")
			return
		}

		// Generate token
		plainToken := generateSanctumToken()
		//old
		hashedToken := sha256.Sum256([]byte(plainToken))
		// Encode hashed token ke hex
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		// Set expires at (1 hari dari sekarang)
		expiresAt := time.Now().Add(24 * time.Hour)

		token := models.PersonalAccessToken{
			TokenableID:   user.ID,
			TokenableType: "App\\Models\\User",
			Name:          "go-client",
			Token:         hashedTokenHex,
			Abilities:     helper.StringPtr("['*']"),
			ExpiresAt:     &expiresAt,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		// Sebelum create
		// Use transaction to ensure data consistency
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&token).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("Error creating token: %v", err)
			helper.InternalServerError(c, "Failed to create token")
			return
		}

		// Verify token ID was generated
		if token.ID == 0 {
			helper.InternalServerError(c, "Failed to get token")
			return
		}

		accessToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

		// Format response seperti Laravel
		response := gin.H{
			"name":  user.Name,
			"email": user.Email,
			"token": accessToken,
			"role":  getPrimaryRole(user.Roles),
			"session": gin.H{
				"expires_at": expiresAt.Format(time.RFC3339Nano),
				"expired_in": 24,
			},
		}

		helper.Success(c, "Login successful. Welcome Bro 🔥✌️", response)
	}
}

// ---------------------------
// LOGOUT
// ---------------------------
func (a *AuthController) Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			helper.Unauthorized(c, "Token tidak valid")
			return
		}

		parts := strings.SplitN(strings.TrimPrefix(authHeader, "Bearer "), "|", 2)
		if len(parts) != 2 {
			helper.Unauthorized(c, "Token tidak lengkap")
			return
		}

		tokenID, rawToken := parts[0], parts[1]
		hashed := sha256.Sum256([]byte(rawToken))
		hashedHex := hex.EncodeToString(hashed[:])

		var token models.PersonalAccessToken
		err := db.Where("id = ? AND token = ?", tokenID, hashedHex).First(&token).Error
		if err != nil {
			helper.Unauthorized(c, "Token tidak dikenali")
			return
		}

		db.Delete(&token)
		helper.Success(c, "Logout berhasil", nil)
	}
}

// ---------------------------
// REGISTER
// ---------------------------
func (a *AuthController) Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			helper.BadRequest(c, "Data tidak valid")
			return
		}

		var count int64
		db.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
		if count > 0 {
			helper.BadRequest(c, "Email sudah digunakan")
			return
		}

		user := models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
		}
		if err := db.Create(&user).Error; err != nil {
			helper.UnprocessableEntity(c, "Gagal mendaftar: "+err.Error())
			return
		}

		helper.Created(c, "Register berhasil", nil)
	}
}

// ---------------------------
// PROFILE / ME
// ---------------------------
func (a *AuthController) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		helper.Unauthorized(c, "Unauthenticated")
		return
	}

	u, exists := user.(models.User)
	if !exists {
		helper.NotFound(c, "User not found")
		return
	}

	data := gin.H{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	}

	helper.Success(c, "User fetched!", data)
}

// ---------------------------
// UTILITIES
// ---------------------------
func getPrimaryRole(roles []models.Role) string {
	if len(roles) > 0 {
		return roles[0].Name
	}
	return ""
}

func generateSanctumToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 48) // 48 karakter seperti Sanctum default
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
