package helper

import (
	"os"
	"response-std/config"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to check if method is valid
func IsValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

// function to easily save an image on the server and return the path
func FileHelper(fileName string, fileData []byte) (string, error) {
	monthYear := time.Now().Format("2006-01")
	fileName = monthYear + "-" + uuid.New().String() + "-" + fileName

	// Ambil versi API dari config
	apiVersion := config.ENV.API_VERSION

	// Path file berdasarkan versi API
	filePath := "../../public/storage/" + apiVersion + "/images/" + fileName

	// Pastikan folder tujuan ada
	err := os.MkdirAll("../../public/storage/"+apiVersion+"/images", os.ModePerm)
	if err != nil {
		return "", err
	}

	// Simpan file
	err = os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}