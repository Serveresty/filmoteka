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
		return nil, fmt.Errorf("error wrong token signature")
	}

	tokenString := tokenArr[1]

	currentToken, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if !currentToken.Valid {
		return nil, fmt.Errorf("error token is not valid")
	}

	claims, ok := currentToken.Claims.(*models.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("error token's claims")
	}

	return claims, nil
}
