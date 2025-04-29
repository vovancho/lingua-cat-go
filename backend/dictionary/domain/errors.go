package domain

import "errors"

var (
	DictNotFoundError    = errors.New("DICTIONARY_NOT_FOUND")
	DictsNotFoundError   = errors.New("DICTIONARIES_NOT_FOUND")
	DictLangInvalidError = errors.New("DICTIONARY_LANG_INVALID")
)
