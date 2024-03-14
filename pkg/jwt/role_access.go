package jwt

import "filmoteka/internal/models"

func HasUserAccess(claims models.JWTClaims) bool {
	for _, role := range claims.Role {
		if role == "admin" {
			return true
		}
	}
	return false
}
