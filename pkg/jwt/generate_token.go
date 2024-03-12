package jwt

import (
	"filmoteka/internal/models"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(id int, roles []string) (string, error) {
	strID := strconv.Itoa(id)

	var claims = &models.JWTClaims{
		Role: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   strID,
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
