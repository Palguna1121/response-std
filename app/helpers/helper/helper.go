package helper

import (
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword menggunakan bcrypt untuk meng-hash password
// mengembalikan string hash dari password yang diberikan
// jika terjadi error, akan mengembalikan error tersebut
// contoh penggunaan:
// hashedPassword, err := HashPassword("password123")
func HashPassword(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash membandingkan password yang diberikan dengan hash yang disimpan
// mengembalikan true jika password cocok dengan hash, false jika tidak cocok
// contoh penggunaan:
// isValid := CheckPasswordHash("password123", hashedPassword)
// // isValid akan true jika password cocok dengan hash
// jika terjadi error, isi isValid akan mengembalikan false
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// StringPtr adalah helper function untuk membuat pointer dari string misalnya untuk JSON response
// jika string kosong, maka akan mengembalikan nil
// jika string tidak kosong, maka akan mengembalikan pointer ke string tersebut
// ini berguna untuk menghindari field kosong di JSON response
// contoh penggunaan:
// var name string
//
//	if err := c.ShouldBindJSON(&name); err != nil {
//	    response.UnprocessableEntity(c, err.Error())
//	    return
//	}
//
// response.Success(c, "Data berhasil disimpan", gin.H{"name": StringPtr(name)})
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

// FileHelper adalah helper function untuk menyimpan file ke folder berdasarkan tahun dan bulan
// mengembalikan path file yang disimpan atau error jika terjadi kesalahan
// contoh penggunaan:
// filePath, err := FileHelper("image.jpg", fileData)
// jika terjadi error, akan mengembalikan error tersebut
// jika berhasil, akan mengembalikan path file yang disimpan
func FileHelper(fileName string, fileData []byte) (string, error) {
	// Buat folder berdasarkan tahun dan bulan sekarang
	currentTime := time.Now()
	folderName := currentTime.Format("2006_01") + "_images"
	dirPath := "../../public/storage/" + folderName

	// Buat nama file dengan tambahan UUID dan tanggal biar unik
	fileName = currentTime.Format("2006-01") + "-" + uuid.New().String() + "-" + fileName
	filePath := dirPath + "/" + fileName

	// buat folder jika belum ada
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Simpan file ke path
	err = os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
