package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(byte), err
}

func ChechHashedPassword(Password, hashedpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(hashedpassword))
	return err == nil
}
