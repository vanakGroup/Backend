package base

import (
	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"
	"github.com/Vanak/Backend/internal/domain/localization"

	"github.com/gin-gonic/gin"
)

func GetTranslator(ctx *gin.Context, key string) localization.TranslatorInstance {
	translator, exists := ctx.Get(key)
	if !exists {
		panic("translator not registered!")
	}

	return translator.(localization.TranslatorInstance)
}

func GetLocalizedTemplateFile(ctx *gin.Context, key, persianTemplateFile, englishTemplateFile string) string {
	trans := GetTranslator(ctx, key)
	switch trans.Locale() {
	case "fa_IR":
		return persianTemplateFile
	case "en_US":
		return englishTemplateFile
	default:
		return persianTemplateFile
	}
}

func GetID(ctx *gin.Context) uint {
	id, ok := ctx.Get(bootstrap.Run().Constant.Context.UserID)
	if !ok {
		panic(exception.NewNotFoundError(bootstrap.Run().Constant.ErrorField.User))
	}
	ID, _ := id.(uint)
	return ID
}
