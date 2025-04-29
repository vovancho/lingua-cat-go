package domain

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterAll(v *validator.Validate, trans ut.Translator) error {
	if err := registerValidations(v); err != nil {
		return err
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

	//if err := v.RegisterTranslation("valid_dict_translation_lang", trans,
	//	func(ut ut.Translator) error {
	//		return ut.Add("valid_dict_translation_lang", "{0} перевод не может быть того же языка, что и основной словарь", true)
	//	},
	//	func(ut ut.Translator, fe validator.FieldError) string {
	//		t, _ := ut.T("valid_dict_translation_lang", fe.Field())
	//		return t
	//	}); err != nil {
	//	return err
	//}

	return nil
}

func registerValidations(v *validator.Validate) error {
	if err := registerDictionaryTypeValidation(v); err != nil {
		return err
	}
	if err := registerDictionaryLangValidation(v); err != nil {
		return err
	}
	//if err := registerDictionaryTranslationLangValidation(v); err != nil {
	//	return err
	//}
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

//func registerDictionaryTranslationLangValidation(v *validator.Validate) error {
//	return v.RegisterValidation("valid_dict_translation_lang", func(fl validator.FieldLevel) bool {
//		value, ok := fl.Field().Interface().(DictionaryLang)
//		if !ok {
//			return false
//		}
//
//		// TODO вынести регистрацию валидатора из домена (http.DictionaryStoreRequest)
//		//dictionaryStoreRequest, ok := fl.Top().Interface().(DictionaryStoreRequest)
//		//if !ok {
//		//	return false // Если структура не DictionaryStoreRequest, валидация не проходит
//		//}
//		//return dictionaryStoreRequest.Lang != value
//
//		// Получаем верхнеуровневую структуру
//		dictionaryStoreRequest := fl.Top().Interface()
//
//		// Используем reflect для динамической проверки
//		val := reflect.ValueOf(dictionaryStoreRequest)
//		if val.Kind() == reflect.Ptr {
//			val = val.Elem() // Получаем структуру, на которую указывает указатель
//		}
//
//		if val.Kind() != reflect.Struct {
//			return false // Не структура
//		}
//
//		// Проверяем наличие поля Lang
//		langField := val.FieldByName("Lang")
//		if !langField.IsValid() {
//			return false // Поле Lang не существует
//		}
//
//		// Проверяем, что тип поля Lang совместим с DictionaryLang
//		langValue, ok := langField.Interface().(DictionaryLang)
//		if !ok {
//			return false
//		}
//
//		// Проверяем, что Lang != Translations[].Lang
//		return langValue != value
//	})
//}
