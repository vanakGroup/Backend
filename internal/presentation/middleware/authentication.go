package middleware

import (
	"fmt"
	"strings"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/application/usecase"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/logging"
	"github.com/Vanak/Backend/internal/domain/ports"
	"github.com/Vanak/Backend/internal/presentation/controller/base"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	logger     logging.Logger
	cfg        *bootstrap.Config
	jwtService usecase.JwtService
	unitOfWork ports.UnitOfWork
}

func NewAuthMiddleWare(cfg *bootstrap.Config, jwtService usecase.JwtService, unitOfWork ports.UnitOfWork) *AuthMiddleware {
	logger := logging.NewLogger(cfg)
	return &AuthMiddleware{
		logger:     logger,
		cfg:        cfg,
		jwtService: jwtService,
		unitOfWork: unitOfWork,
	}
}

func (am *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", am.cfg.Constant.ErrorField.Auth, am.cfg.Constant.ErrorTag.EmptyHeader)}
			panic(authError)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", am.cfg.Constant.ErrorField.Auth, am.cfg.Constant.ErrorTag.InvalidHeader)}
			panic(authError)
		}

		token := parts[1]
		if token == "" {
			authError := &exception.UnauthorizedError{Message: fmt.Sprintf("%s.%s", am.cfg.Constant.ErrorField.Jwt, am.cfg.Constant.ErrorTag.EmptyToken)}
			panic(authError)
		}

		claims, err := am.jwtService.ValidateToken(token)
		if err != nil {
			panic(err)
		}

		ctx.Set(am.cfg.Constant.Context.UserID, uint(claims["userID"].(float64)))
		ctx.Next()
	}
}

func (am *AuthMiddleware) UserExistMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := base.GetID(ctx)

		_, err := am.unitOfWork.Factory().UserRepository().FindUserById(userID)

		if err != nil {
			panic(exception.NewNotFoundError(am.cfg.Constant.Context.UserID))
		}

		ctx.Next()
	}
}
