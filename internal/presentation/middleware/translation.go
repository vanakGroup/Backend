package middleware

import (
	"strings"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/localization"

	"github.com/gin-gonic/gin"
)

func TranslatorMiddleware(translatorService localization.Translator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := "fa_IR"
		if lang := ctx.GetHeader("Accept-Language"); lang != "" {
			if strings.Contains(lang, "en") {
				locale = "en_US"
			}
		}

		translatorInstance := translatorService.GetTranslator(locale)

		ctx.Set(bootstrap.Run().Constant.Context.Translator, translatorInstance)

		ctx.Next()
	}
}
