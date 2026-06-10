package validation

import (
	"reflect"

	"github.com/Vanak/Backend/bootstrap"
	"github.com/Vanak/Backend/internal/domain/exception"

	"github.com/go-playground/validator/v10"
)

var Val = validator.New()

func init() {
	Val.RegisterValidation(bootstrap.NewConstant().ValidationTag.IranianPhone, IranianMobileValidator, true)
	Val.RegisterValidation(bootstrap.NewConstant().ValidationTag.IraninaPostal, IranianPostalCodeValidator, true)
}

func ValidatedField[T any](params T) {

	if err := Val.Struct(params); err != nil {
		validationError, _ := err.(validator.ValidationErrors)
		customValidationError := exception.NewValidationErrors()

		for _, err := range validationError {
			field := formatValidationError[T](err)

			customValidationError.Add(field, err.Tag())
		}
		panic(customValidationError)
	}
}

func formatValidationError[T any](err validator.FieldError) string {
	var params T

	tagType := []string{"json", "uri", "form"}

	reflectType := reflect.TypeOf(params)

	field, _ := reflectType.FieldByName(err.StructField())

	tagValue := getAnyTag(field, tagType...)

	return tagValue
}

func getAnyTag(field reflect.StructField, tagNames ...string) string {

	for _, tagName := range tagNames {
		if tag := field.Tag.Get(tagName); tag != "" {
			return tag
		}
	}

	return field.Name
}
