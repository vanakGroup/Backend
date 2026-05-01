package middleware

import (
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"

	"github.com/gin-gonic/gin"
)

type RecoveryMiddleware struct{}

func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

func (rm *RecoveryMiddleware) Recovery(ctx *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				rm.handlerError(ctx, err)
				ctx.Abort()
			}
		}
	}()
	ctx.Next()
}

func (rm *RecoveryMiddleware) handlerError(ctx *gin.Context, err error) {
	if validationErrors, ok := err.(exception.ValidationError); ok {
		handleValidationError(ctx, validationErrors)
	} else if BindingErrors, ok := err.(exception.BindingError); ok {
		handleBindingError(ctx, BindingErrors)
	} else {
		handleDefaultError(ctx, err)
	}
}

func handleDefaultError(ctx *gin.Context, err error) {
	base.Response(
		ctx,
		409,
		err.Error(),
		nil,
	)
}

func handleValidationError(ctx *gin.Context, validationError exception.ValidationError) {
	base.Response(
		ctx,
		400,
		validationError.Error(),
		nil,
	)
}

func handleBindingError(ctx *gin.Context, BindingError exception.BindingError) {
	base.Response(
		ctx,
		422,
		BindingError.Error(),
		nil,
	)
}
