package validator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/delivery/http"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

func NewValidator(v *validator.Validate, trans ut.Translator) (*validator.Validate, error) {
	if err := domain.RegisterAll(v, trans); err != nil {
		return nil, fmt.Errorf("failed to register domain validations: %w", err)
	}

	if err := http.RegisterAll(v, trans); err != nil {
		return nil, fmt.Errorf("failed to register http validations: %w", err)
	}

	return v, nil
}
