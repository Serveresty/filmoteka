package models

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	Role []string
	jwt.StandardClaims
}
