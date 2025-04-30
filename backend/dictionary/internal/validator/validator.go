package validator

import (
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

func NewValidator() (*validator.Validate, ut.Translator, error) {
	validate := validator.New()
	uni := ut.New(ru.New(), ru.New())
	trans, _ := uni.GetTranslator("ru")

	if err := rutranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	if err := domain.RegisterAll(validate, trans); err != nil {
		return nil, nil, fmt.Errorf("failed to register domain validations: %w", err)
	}

	return validate, trans, nil
}
