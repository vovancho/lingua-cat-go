package domain

import "errors"

var (
	DictNotFoundError               = errors.New("DICTIONARY_NOT_FOUND")
	DictExistsError                 = errors.New("DICTIONARY_EXISTS")
	DictsNotFoundError              = errors.New("DICTIONARIES_NOT_FOUND")
	DictsRandomCountError           = errors.New("DICTIONARIES_RANDOM_COUNT_INVALID")
	DictLangInvalidError            = errors.New("DICTIONARY_LANG_INVALID")
	DictTranslationLangInvalidError = errors.New("DICTIONARY_TRANSLATION_LANG_INVALID")
	DictTranslationRequiredError    = errors.New("DICTIONARY_TRANSLATION_REQUIRED")
)
