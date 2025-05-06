package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
)

func NewValidator(trans ut.Translator) (*validator.Validate, error) {
	validate := validator.New()

	if err := rutranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	if err := domain.RegisterAll(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to register domain validations: %w", err)
	}

	return validate, nil
}
