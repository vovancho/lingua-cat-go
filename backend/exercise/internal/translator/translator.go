package translator

import (
	"fmt"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
)

func NewTranslator() (ut.Translator, error) {
	uni := ut.New(ru.New(), ru.New())
	trans, found := uni.GetTranslator("ru")
	if !found {
		return nil, fmt.Errorf("RU locale not found")

	}
	return trans, nil
}
