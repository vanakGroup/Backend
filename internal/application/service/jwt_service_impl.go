package service

import (
	"fmt"
	"time"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"
	domainJWT "github.com/Vanak/Backend/internal/domain/jwt"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	logger     logging.Logger
	cfg        *bootstrap.Config
	keyManager domainJWT.KeyManager
}

func NewJwtService(cfg *bootstrap.Config, keyManager domainJWT.KeyManager) *JwtService {
	logger := logging.NewLogger(cfg)
	return &JwtService{
		keyManager: keyManager,
		cfg:        cfg,
		logger:     logger,
	}
}

func (js *JwtService) GenerateToken(userID uint) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)

	accessTokenString, err := accessToken.SignedString(js.keyManager.GetPrivateKey())

	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
		"iat":    time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)

	refreshTokenString, err := refreshToken.SignedString(js.keyManager.GetPrivateKey())

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (js *JwtService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", js.cfg.Constant.ErrorField.Jwt, js.cfg.Constant.ErrorTag.InvalidSignedToken)}
		}
		return js.keyManager.GetPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", js.cfg.Constant.ErrorField.Jwt, js.cfg.Constant.ErrorTag.InvalidAuthToken)}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", js.cfg.Constant.ErrorField.Jwt, js.cfg.Constant.ErrorTag.ClaimReject)}
	}

	return claims, nil
}
