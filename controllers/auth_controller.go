package controllers

import (
	"fmt"
	"strings"
	"time"

	"response-std/config"
	"response-std/helper"
	"response-std/models"
	"response-std/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Secret untuk JWT
func getJwtSecret() []byte {
	return []byte(config.ENV.JWT_SECRET)
}

// Format token seperti Laravel Sanctum: "id|jwt"
func generateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, err := token.SignedString(getJwtSecret())
	if err != nil {
		return "", err
	}
	return formatSanctumToken(userID, jwtStr), nil
}

func formatSanctumToken(userID uint, jwt string) string {
	return fmt.Sprintf("%d|%s", userID, jwt)
}

func parseSanctumToken(rawToken string) (*jwt.Token, uint, error) {
	parts := strings.Split(rawToken, "|")
	if len(parts) != 2 {
		return nil, 0, fmt.Errorf("invalid token format")
	}
	userID := parts[0]
	jwtToken := parts[1]
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return getJwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, 0, err
	}
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)
	return token, uid, nil
}

func (a *AuthController) Login(c *gin.Context) {
	type LoginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.BadRequest(c, "Invalid input")
		return
	}

	var user models.User
	db := config.DB
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		helper.NotFound(c, "Email not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		helper.BadRequest(c, "Invalid password")
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		helper.InternalServerError(c, "Failed to generate token")
		return
	}

	data := gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}

	helper.Success(c, "Login successful", data)
}

func (a *AuthController) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		helper.Unauthorized(c, "Unauthenticated")
		return
	}

	data := gin.H{
		"id":    user.(models.User).ID,
		"name":  user.(models.User).Name,
		"email": user.(models.User).Email,
	}

	helper.Success(c, "User fetched!", data)
}

func (a *AuthController) Register(c *gin.Context) {
	type RegisterInput struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,gte=8"`
	}
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.BadRequest(c, "Invalid input")
		return
	}

	var user models.User
	db := config.DB
	if err := db.Where("email = ?", input.Email).First(&user).Error; err == nil {
		helper.BadRequest(c, "Email already taken")
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		helper.InternalServerError(c, "Failed to hash password")
		return
	}

	user = models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   3, // participant
	}

	if err := db.Create(&user).Error; err != nil {
		helper.InternalServerError(c, "Failed to create user")
		return
	}

	data := gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}

	helper.Created(c, "User created!", data)
}
