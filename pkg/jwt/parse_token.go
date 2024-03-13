package jwt

import (
	"filmoteka/internal/models"
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(token string) (*models.JWTClaims, error) {
	tokenArr := strings.Split(token, " ")
	if len(tokenArr) < 2 {
		return nil, fmt.Errorf("error wrong token")
	}

	tokenString := tokenArr[1]
	secretKey := os.Getenv("SECRET_KEY")

	currentToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := currentToken.Claims.(*models.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("error token's claims")
	}

	return claims, nil
}
