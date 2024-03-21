package models

import (
	"encoding/base64"
	"os"

	"golang.org/x/crypto/argon2"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func EncodePassword(password string) string {
	salt := os.Getenv("PASSWORD_SALT")
	encrypted := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func Login(username string, password string) (*User, error) {
	var user User
	result := DB.Where(
		"username = ? AND password = ?",
		username,
		EncodePassword(password),
	).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
