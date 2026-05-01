package base

import (
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
