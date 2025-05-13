package domain

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterAll(v *validator.Validate, trans ut.Translator) error {
	if err := registerValidations(v); err != nil {
		return err
	}

	if trans == nil {
		return nil
	}

	if err := v.RegisterTranslation("valid_dictionary_type", trans,
		func(ut ut.Translator) error {
			return ut.Add("valid_dictionary_type", "{0} должен быть валидным типом", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_dictionary_type", fe.Field())
			return t
		}); err != nil {
		return err
	}

	if err := v.RegisterTranslation("valid_dictionary_lang", trans,
		func(ut ut.Translator) error {
			return ut.Add("valid_dictionary_lang", "{0} должен быть валидным языком", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_dictionary_lang", fe.Field())
			return t
		}); err != nil {
		return err
	}

	return nil
}

func registerValidations(v *validator.Validate) error {
	if err := registerDictionaryTypeValidation(v); err != nil {
		return err
	}
	if err := registerDictionaryLangValidation(v); err != nil {
		return err
	}
	return nil
}

func registerDictionaryTypeValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_dictionary_type", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(DictionaryType)
		if !ok {
			return false
		}
		return value.IsValid()
	})
}

func registerDictionaryLangValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_dictionary_lang", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(DictionaryLang)
		if !ok {
			return false
		}
		return value.IsValid()
	})
}
