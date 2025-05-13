package http

import (
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

func RegisterAll(v *validator.Validate, trans ut.Translator) error {
	if err := registerValidations(v); err != nil {
		return err
	}

	if trans == nil {
		return nil
	}

	registerFn := func(ut ut.Translator) error {
		return ut.Add("valid_dict_translation_lang", "{0} перевод не может быть того же языка, что и основной словарь", true)
	}
	translationFn := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_dict_translation_lang", fe.Field())

		return t
	}

	if err := v.RegisterTranslation("valid_dict_translation_lang", trans, registerFn, translationFn); err != nil {
		return err
	}

	return nil
}

func registerValidations(v *validator.Validate) error {
	if err := registerDictionaryTranslationLangValidation(v); err != nil {
		return err
	}

	return nil
}

func registerDictionaryTranslationLangValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_dict_translation_lang", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(domain.DictionaryLang)
		if !ok {
			return false
		}

		dictionaryStoreRequest, ok := fl.Top().Interface().(*DictionaryStoreRequest)
		if !ok {
			return false // Если структура не DictionaryStoreRequest, валидация не проходит
		}

		return dictionaryStoreRequest.Lang != value
	})
}
