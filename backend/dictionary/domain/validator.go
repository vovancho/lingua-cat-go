package domain

import (
	"github.com/go-playground/validator/v10"
)

func RegisterDictionaryTypeValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_dictionary_type", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(DictionaryType)
		if !ok {
			return false
		}
		return value.IsValid()
	})
}

func RegisterDictionaryLangValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_dictionary_lang", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(DictionaryLang)
		if !ok {
			return false
		}
		return value.IsValid()
	})
}
