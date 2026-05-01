package validation

import (
	"fmt"
	"reflect"

	"github.com/Vanak/Backend/internal/domain/exception"

	"github.com/go-playground/validator/v10"
)

var val = validator.New()

func ValidatedField[T any](params T) {

	if err := val.Struct(params); err != nil {
		fmt.Println(err)
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
