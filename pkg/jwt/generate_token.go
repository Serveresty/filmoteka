package jwt

import (
	"filmoteka/internal/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(id string, roles []string) (string, error) {
	var claims = &models.JWTClaims{
		Role: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   id,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return "Bearer " + tokenString, nil
}
