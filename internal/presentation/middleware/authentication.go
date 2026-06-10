package middleware

import (
	"fmt"
	"strings"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	logger     logging.Logger
	cfg        *bootstrap.Config
	jwtService usecase.JwtService
}

func NewAuthMiddleWare(cfg *bootstrap.Config, jwtService usecase.JwtService) *AuthMiddleware {
	logger := logging.NewLogger(cfg)
	return &AuthMiddleware{
		logger:     logger,
		cfg:        cfg,
		jwtService: jwtService,
	}
}

func (authMiddleware *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", authMiddleware.cfg.Constant.ErrorField.Auth, authMiddleware.cfg.Constant.ErrorTag.EmptyHeader)}
			panic(authError)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", authMiddleware.cfg.Constant.ErrorField.Auth, authMiddleware.cfg.Constant.ErrorTag.InvalidHeader)}
			panic(authError)
		}

		token := parts[1]
		if token == "" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", authMiddleware.cfg.Constant.ErrorField.Jwt, authMiddleware.cfg.Constant.ErrorTag.EmptyToken)}
			panic(authError)
		}

		claims, err := authMiddleware.jwtService.ValidateToken(token)
		if err != nil {
			panic(err)
		}

		ctx.Set(authMiddleware.cfg.Constant.Context.UserID, uint(claims["userID"].(float64)))
		ctx.Next()
	}
}
