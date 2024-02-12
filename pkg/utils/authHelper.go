package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Function that hashes passwords for users
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// Function that compares both provided password and saved password to verify the user
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Email or password is incorrect"
		check = false
	}
	return check, msg
}
