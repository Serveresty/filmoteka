package jwt

import (
	"filmoteka/internal/models"
	"regexp"
	"unicode"
)

func ValidateAuthData(user *models.User) bool {
	first := hasLetters(user.FirstName)
	last := hasLetters(user.LastName)
	email := isValidEmail(user.Email)
	if user.Password == "" || !first || !last || !email {
		return false
	}
	return true
}

func hasLetters(input string) bool {
	for _, char := range input {
		if unicode.IsLetter(char) {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
