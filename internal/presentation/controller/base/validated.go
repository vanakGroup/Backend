package base

import (
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/infrastructure/validation"

	"github.com/gin-gonic/gin"
)

func Validated[T any](ctx *gin.Context) T {
	var params T
	if err := ctx.ShouldBind(&params); err != nil {
		NewBindingError := exception.NewBinding(err)
		BindingError := NewBindingError.Error()
		panic(BindingError)
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		NewBindingError := exception.NewBinding(err)
		BindingError := NewBindingError.Error()
		panic(BindingError)
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		NewBindingError := exception.NewBinding(err)
		BindingError := NewBindingError.Error()
		panic(BindingError)
	}

	validation.ValidatedField(params)

	return params
}
