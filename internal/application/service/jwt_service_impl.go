package service

import (
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	logger logging.Logger
	cfg    *bootstrap.Config
}

func NewJwtService(cfg *bootstrap.Config) *JwtService {
	logger := logging.NewLogger(cfg)

	return &JwtService{
		logger: logger,
		cfg:    cfg,
	}
}

func (jwtService *JwtService) GenerateToken(userID uint) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(jwtService.cfg.Env.Jwt.SecretKey))

	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(jwtService.cfg.Env.Jwt.SecretKey))

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (jwtService *JwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", jwtService.cfg.Constant.ErrorField.Jwt, jwtService.cfg.Constant.ErrorTag.InvalidSignedToken)}
		}
		return []byte(jwtService.cfg.Env.Jwt.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", jwtService.cfg.Constant.ErrorField.Jwt, jwtService.cfg.Constant.ErrorTag.InvalidAuthToken)}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", jwtService.cfg.Constant.ErrorField.Jwt, jwtService.cfg.Constant.ErrorTag.ClaimReject)}
	}

	return claims, nil
}
