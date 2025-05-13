package validator

import (
	"fmt"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
)

func NewValidator(v *validator.Validate, trans ut.Translator) (*validator.Validate, error) {
	if err := domain.RegisterAll(v, trans); err != nil {
		return nil, fmt.Errorf("failed to register domain validations: %w", err)
	}

	return v, nil
}
