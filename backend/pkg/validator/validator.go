package validator

import (
	"fmt"

	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/ru"
)

func NewValidator(trans ut.Translator) (*validator.Validate, error) {
	validate := validator.New()

	if err := ru.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, fmt.Errorf("failed to register default translations: %w", err)
	}

	return validate, nil
}
