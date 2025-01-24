package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	valid := true
	msg := ""
	if err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword)); err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}
