package domain

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterAll(v *validator.Validate, trans ut.Translator) error {
	if err := registerValidations(v); err != nil {
		return err
	}

	if err := v.RegisterTranslation("valid_exercise_lang", trans,
		func(ut ut.Translator) error {
			return ut.Add("valid_exercise_lang", "{0} должен быть валидным языком", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("valid_exercise_lang", fe.Field())
			return t
		}); err != nil {
		return err
	}

	return nil
}

func registerValidations(v *validator.Validate) error {
	if err := registerExerciseLangValidation(v); err != nil {
		return err
	}
	return nil
}

func registerExerciseLangValidation(v *validator.Validate) error {
	return v.RegisterValidation("valid_exercise_lang", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(ExerciseLang)
		if !ok {
			return false
		}
		return value.IsValid()
	})
}
