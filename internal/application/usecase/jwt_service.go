package usecase

import "github.com/golang-jwt/jwt/v5"

type JwtService interface {
	GenerateToken(userID uint) (string, string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
}
