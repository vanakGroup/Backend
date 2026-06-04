package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/presentation/controller/base"

	"github.com/gin-gonic/gin"
)

const (
	genericError = "errors.generic"
)

type RecoveryMiddleware struct {
	constants *bootstrap.Constants
}

func NewRecoveryMiddleware(constants *bootstrap.Constants) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		constants: constants,
	}
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
	if validationErrors, ok := err.(*exception.ValidationError); ok {
		handleValidationError(ctx, *validationErrors, rm.constants.Context.Translator)
	} else if BindingErrors, ok := err.(*exception.BindingError); ok {
		handleBindingError(ctx, *BindingErrors, rm.constants.Context.Translator)
	} else if AppErrors, ok := err.(*exception.UnauthorizedError); ok {
		handleUnauthorizedError(ctx, *AppErrors, rm.constants.Context.Translator)
	} else if NotFoundError, ok := err.(*exception.NotFoundError); ok {
		handleNotFoundError(ctx, *NotFoundError, rm.constants.Context.Translator)
	} else if DuplicateError, ok := err.(*exception.DuplicateError); ok {
		handleDuplicateError(ctx, *DuplicateError, rm.constants.Context.Translator)
	} else {
		handleDefaultError(ctx, err)
	}
}

func handleDuplicateError(ctx *gin.Context, duplicateError exception.DuplicateError, transKey string) {
	trans := base.GetTranslator(ctx, transKey)
	itemName, _ := trans.Translate(duplicateError.Item)
	message, _ := trans.Translate("errors.duplicate", itemName)

	base.Response(ctx, 409, message, nil)
}

func handleNotFoundError(ctx *gin.Context, notFoundError exception.NotFoundError, transKey string) {
	trans := base.GetTranslator(ctx, transKey)
	itemName, _ := trans.Translate(notFoundError.Item)
	message, _ := trans.Translate("errors.notFound", itemName)

	base.Response(ctx, 404, message, nil)
}

func handleDefaultError(ctx *gin.Context, err error) {
	base.Response(
		ctx,
		409,
		err.Error(),
		nil,
	)
}

func handleValidationError(ctx *gin.Context, validationError exception.ValidationError, transKey string) {
	trans := base.GetTranslator(ctx, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, validationError := range validationError.FieldErros {
		if _, ok := errorMessages[validationError.Field]; !ok {
			errorMessages[validationError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(validationError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", validationError.Tag), fieldName)
		errorMessages[validationError.Field][validationError.Tag] = message
	}

	base.Response(ctx, 422, errorMessages, nil)
}

func handleBindingError(ctx *gin.Context, bindingError exception.BindingError, transKey string) {
	trans := base.GetTranslator(ctx, transKey)
	message, _ := trans.Translate(genericError)

	if numError, ok := bindingError.Err.(*strconv.NumError); ok {
		message, _ = trans.Translate("errors.numeric", numError.Num)
	} else if bindingError == http.ErrMissingFile {
		message, _ = trans.Translate("errors.fileRequired")
	}

	base.Response(ctx, 400, message, nil)
}

func handleUnauthorizedError(ctx *gin.Context, AppError exception.UnauthorizedError, transKey string) {
	trans := base.GetTranslator(ctx, transKey)
	message, _ := trans.Translate(AppError.Message)

	base.Response(ctx, 401, message, nil)
}
