package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"response-std/core/helper"
	"response-std/core/models"
	"response-std/core/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	// Seed random generator untuk keamanan
	rand.Seed(time.Now().UnixNano())
	return &AuthController{}
}

// ---------------------------
// LOGIN
// ---------------------------
func (a *AuthController) Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"` // Changed from email to username
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			response.UnprocessableEntity(c, "Invalid credentials")
			return
		}

		// Determine if username is email or name (same logic as Laravel)
		var loginField string
		if strings.Contains(input.Username, "@") {
			loginField = "email"
		} else {
			loginField = "name"
		}

		var user models.User
		err := db.Preload("Roles.Permissions").Preload("Permissions").
			Where(loginField+" = ?", input.Username).First(&user).Error

		if err != nil {
			response.UnprocessableEntity(c, "Invalid credentials")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			response.UnprocessableEntity(c, "Invalid credentials")
			return
		}

		// Generate token
		plainToken := generateSanctumToken()
		hashedToken := sha256.Sum256([]byte(plainToken))
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		// Set expires (12 jam dari sekarang)
		expiresAt := time.Now().Add(12 * time.Hour)

		token := models.PersonalAccessTokens{
			TokenableID:   user.ID,
			TokenableType: "App\\Models\\User",
			Name:          "go-client",
			Token:         hashedTokenHex,
			Abilities:     helper.StringPtr("['*']"),
			ExpiresAt:     &expiresAt,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Use transaction to ensure data consistency
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&token).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Printf("Error creating token: %v", err)
			response.InternalServerError(c, "Failed to create token")
			return
		}

		if token.ID == 0 {
			response.InternalServerError(c, "Failed to get token")
			return
		}

		accessToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

		res := gin.H{
			"name":  user.Name,
			"email": user.Email,
			"token": accessToken,
			"role":  getPrimaryRole(user.Roles),
			"session": gin.H{
				"expires_at": expiresAt.Format(time.RFC3339Nano),
				"expired_in": 24,
			},
		}

		response.Success(c, "Login successful. Welcome Bro 🔥✌️", res)
	}
}

// ---------------------------
// LOGOUT
// ---------------------------
func (a *AuthController) Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, "Token tidak valid")
			return
		}

		parts := strings.SplitN(strings.TrimPrefix(authHeader, "Bearer "), "|", 2)
		if len(parts) != 2 {
			response.Unauthorized(c, "Token tidak lengkap")
			return
		}

		tokenID, rawToken := parts[0], parts[1]
		hashed := sha256.Sum256([]byte(rawToken))
		hashedHex := hex.EncodeToString(hashed[:])

		var token models.PersonalAccessTokens
		err := db.Where("id = ? AND token = ?", tokenID, hashedHex).First(&token).Error
		if err != nil {
			response.Unauthorized(c, "Token tidak dikenali")
			return
		}

		// Hapus token dari database
		if err := db.Delete(&token).Error; err != nil {
			log.Printf("Error deleting token: %v", err)
			response.InternalServerError(c, "Failed to logout")
			return
		}

		response.Success(c, "Logout berhasil", nil)
	}
}

// ---------------------------
// REGISTER
// ---------------------------
func (a *AuthController) Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name                 string `json:"name" binding:"required"`
			Email                string `json:"email" binding:"required,email"`
			Password             string `json:"password" binding:"required,min=8"`
			PasswordConfirmation string `json:"password_confirmation" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			response.BadRequest(c, "Data tidak valid")
			return
		}

		if input.Password != input.PasswordConfirmation {
			response.UnprocessableEntity(c, "Password dan konfirmasi password tidak cocok")
			return
		}

		// Cek apakah email sudah digunakan
		var count int64
		db.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
		if count > 0 {
			response.UnprocessableEntity(c, "Email sudah digunakan")
			return
		}

		// Hash password menggunakan bcrypt (sama seperti Laravel)
		// Laravel default menggunakan cost 12, tapi 10-12 sudah cukup aman
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			response.InternalServerError(c, "Gagal memproses password")
			return
		}

		user := models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: string(hashedPassword), // Simpan password yang sudah di-hash
		}

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Error creating user: %v", err)
			response.UnprocessableEntity(c, "Gagal mendaftar")
			return
		}

		response.Created(c, "Register berhasil", gin.H{
			"message": "Akun berhasil dibuat, silakan login",
		})
	}
}

// ---------------------------
// PROFILE / ME
// ---------------------------
func (a *AuthController) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthenticated")
		return
	}

	u, ok := user.(models.User)
	if !ok {
		response.NotFound(c, "User not found")
		return
	}

	data := gin.H{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"roles": u.Roles, // Tambahkan roles jika dibutuhkan
	}

	response.Success(c, "User fetched!", data)
}

// ---------------------------
// REFRESH TOKEN (Optional)
// ---------------------------
func (a *AuthController) RefreshToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dapatkan user dari middleware auth
		user, exists := c.Get("user")
		if !exists {
			response.Unauthorized(c, "Unauthenticated")
			return
		}

		u, ok := user.(models.User)
		if !ok {
			response.NotFound(c, "User not found")
			return
		}

		// Hapus token lama
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			parts := strings.SplitN(strings.TrimPrefix(authHeader, "Bearer "), "|", 2)
			if len(parts) == 2 {
				tokenID := parts[0]
				db.Where("id = ?", tokenID).Delete(&models.PersonalAccessTokens{})
			}
		}

		// Generate token baru
		plainToken := generateSanctumToken()
		hashedToken := sha256.Sum256([]byte(plainToken))
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		expiresAt := time.Now().Add(24 * time.Hour)

		token := models.PersonalAccessTokens{
			TokenableID:   u.ID,
			TokenableType: "App\\Models\\User",
			Name:          "go-client",
			Token:         hashedTokenHex,
			Abilities:     helper.StringPtr("['*']"),
			ExpiresAt:     &expiresAt,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := db.Create(&token).Error; err != nil {
			response.InternalServerError(c, "Failed to create token")
			return
		}

		accessToken := fmt.Sprintf("%d|%s", token.ID, plainToken)

		res := gin.H{
			"token": accessToken,
			"session": gin.H{
				"expires_at": expiresAt.Format(time.RFC3339Nano),
				"expired_in": 24,
			},
		}

		response.Success(c, "Token berhasil di-refresh", res)
	}
}

// ---------------------------
// UTILITIES
// ---------------------------
func getPrimaryRole(roles []models.Roles) string {
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
