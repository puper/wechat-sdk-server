package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

func CheckPassword(old, new string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(old), []byte(new))
	return err == nil
}
