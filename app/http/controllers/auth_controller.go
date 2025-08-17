// v1/controllers/auth_controller.go
package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"response-std/app/helpers/helper"
	"response-std/app/http/requests/auth"
	"response-std/app/models/entities"
	"response-std/app/pkg/permissions"
	"response-std/app/pkg/response"

	"github.com/davecgh/go-spew/spew"
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
		if db == nil {
			response.InternalServerError(c, "Database connection is not initialized", nil, "[Login]")
			return
		}
		// Validate request using LoginRequest
		var loginReq auth.LoginRequest
		// BIND JSON DULU
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			response.BadRequest(c, "Invalid request format", err, "[Login]")
			return
		}

		// VALIDATE
		if !loginReq.Validate(c) {
			return
		}

		// Determine if username is email or name (same logic as Laravel)
		var loginField string
		if strings.Contains(loginReq.Username, "@") {
			loginField = "email"
		} else {
			loginField = "name"
		}

		var user entities.User
		err := db.Preload("Roles.Permissions").Preload("Permissions").
			Where(loginField+" = ?", loginReq.Username).First(&user).Error

		if err != nil {
			response.UnprocessableEntity(c, "Invalid credentials", err, "[Login]")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
			response.UnprocessableEntity(c, "Invalid credentials", err, "[Login]")
			return
		}

		// Generate token
		plainToken := generateSanctumToken()
		hashedToken := sha256.Sum256([]byte(plainToken))
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		// Set expires (12 jam dari sekarang)
		expiresAt := time.Now().Add(12 * time.Hour)

		token := entities.PersonalAccessTokens{
			TokenableID:   user.ID,
			TokenableType: "App\\entities\\User",
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
			response.InternalServerError(c, "Failed to create token", err, "[Login]")
			return
		}

		if token.ID == 0 {
			response.InternalServerError(c, "Failed to get token", nil, "[Login]")
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
			response.Unauthorized(c, "Token tidak valid", nil, "[Logout]")
			return
		}

		parts := strings.SplitN(strings.TrimPrefix(authHeader, "Bearer "), "|", 2)
		if len(parts) != 2 {
			response.Unauthorized(c, "Token tidak lengkap", nil, "[Logout]")
			return
		}

		tokenID, rawToken := parts[0], parts[1]
		hashed := sha256.Sum256([]byte(rawToken))
		hashedHex := hex.EncodeToString(hashed[:])

		var token entities.PersonalAccessTokens
		err := db.Where("id = ? AND token = ?", tokenID, hashedHex).First(&token).Error
		if err != nil {
			response.Unauthorized(c, "Token tidak dikenali", err, "[Logout]")
			return
		}

		// Hapus token dari database
		if err := db.Delete(&token).Error; err != nil {
			response.InternalServerError(c, "Failed to logout", err, "[Logout]")
			return
		}

		response.Success(c, "Logout berhasil", nil)
	}
}

// ---------------------------
// REGISTER
// ---------------------------
func (a *AuthController) Register(db *gorm.DB, spatie *permissions.Spatie) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate request using RegisterRequest
		var registerReq auth.RegisterRequest
		// BIND JSON DULU
		if err := c.ShouldBindJSON(&registerReq); err != nil {
			response.BadRequest(c, "Invalid request format", err, "[Register]")
			return
		}
		// VALIDATE
		// Gunakan method Validate dari RegisterRequest
		if !registerReq.Validate(c) {
			return // Response sudah di-handle di dalam Validate()
		}

		// Cek apakah email sudah digunakan
		var count int64
		db.Model(&entities.User{}).Where("email = ?", registerReq.Email).Count(&count)
		if count > 0 {
			response.UnprocessableEntity(c, "Email sudah digunakan", nil, "[Register]")
			return
		}

		// Hash password menggunakan bcrypt (sama seperti Laravel)
		// Laravel default menggunakan cost 12, tapi 10-12 sudah cukup aman
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), 12)
		if err != nil {
			response.InternalServerError(c, "Gagal memproses password", err, "[Register]")
			return
		}

		user := entities.User{
			Name:     registerReq.Name,
			Email:    registerReq.Email,
			Password: string(hashedPassword), // Simpan password yang sudah di-hash
		}

		if err := db.Create(&user).Error; err != nil {
			response.UnprocessableEntity(c, "Gagal mendaftar", err, "[Register]")
			return
		}

		// assign default role (misal: user)
		defaultRole := "user"
		if err := spatie.AssignRole(user.ID, defaultRole); err != nil {
			log.Printf("Failed to assign role %s to %s: %v", defaultRole, user.Name, err)
		}

		response.Created(c, "Akun berhasil didaftarkan, silahkan login!", nil)
	}
}

// ---------------------------
// PROFILE / ME
// ---------------------------
func (a *AuthController) Me(c *gin.Context, spatie *permissions.Spatie) {
	user, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthenticated", nil, "[Me]")
		return
	}

	u, ok := user.(entities.User)
	if !ok {
		response.NotFound(c, "User not found", nil, "[Me]")
		return
	}

	userRoles, err := spatie.GetUserRoles(u.ID)
	if err != nil {
		spew.Dump("apa nih", err)
		data := gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"roles": u.Roles, // Tambahkan roles jika dibutuhkan
		}

		response.Success(c, "User fetched!", data)
		return
	}

	data := gin.H{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"roles": userRoles, // Tambahkan roles jika dibutuhkan
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
			response.Unauthorized(c, "Unauthenticated", nil, "[RefreshToken]")
			return
		}

		u, ok := user.(entities.User)
		if !ok {
			response.NotFound(c, "User not found", nil, "[RefreshToken]")
			return
		}

		// Hapus token lama
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			parts := strings.SplitN(strings.TrimPrefix(authHeader, "Bearer "), "|", 2)
			if len(parts) == 2 {
				tokenID := parts[0]
				db.Where("id = ?", tokenID).Delete(&entities.PersonalAccessTokens{})
			}
		}

		// Generate token baru
		plainToken := generateSanctumToken()
		hashedToken := sha256.Sum256([]byte(plainToken))
		hashedTokenHex := hex.EncodeToString(hashedToken[:])

		expiresAt := time.Now().Add(24 * time.Hour)

		token := entities.PersonalAccessTokens{
			TokenableID:   u.ID,
			TokenableType: "App\\entities\\User",
			Name:          "go-client",
			Token:         hashedTokenHex,
			Abilities:     helper.StringPtr("['*']"),
			ExpiresAt:     &expiresAt,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if err := db.Create(&token).Error; err != nil {
			response.InternalServerError(c, "Failed to create token", err, "[RefreshToken]")
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
func getPrimaryRole(roles []entities.Roles) string {
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
