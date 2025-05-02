package domain

import "context"

type DictionaryID uint64
type DictionaryType uint16

const (
	SimpleDictionary        DictionaryType = 1
	PhrasalVerbDictionary   DictionaryType = 2
	IrregularVerbDictionary DictionaryType = 3
	PhraseDictionary        DictionaryType = 4
)

type DictionaryLang string

const (
	RuDictionary DictionaryLang = "ru"
	EnDictionary DictionaryLang = "en"
)

type Dictionary struct {
	ID           DictionaryID   `json:"id"`
	Lang         DictionaryLang `json:"lang"`
	Name         string         `json:"name"`
	Type         DictionaryType `json:"type"`
	Translations []Translation  `json:"translations,omitempty"`
	Sentences    []Sentence     `json:"sentences"`
}

type DictionaryUseCase interface {
	GetRandomDictionaries(ctx context.Context, lang DictionaryLang, limit uint8) ([]Dictionary, error)
	GetDictionariesByIds(ctx context.Context, dictIds []DictionaryID) ([]Dictionary, error)
}

type DictionaryRepository interface {
	GetRandomDictionaries(ctx context.Context, lang DictionaryLang, limit uint8) ([]Dictionary, error)
	GetDictionariesByIds(ctx context.Context, dictIds []DictionaryID) ([]Dictionary, error)
}
