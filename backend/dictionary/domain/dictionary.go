package domain

import (
	"context"
	"time"
)

type DictionaryID uint64
type DictionaryType uint16

const (
	SimpleDictionary        DictionaryType = 1
	PhrasalVerbDictionary   DictionaryType = 2
	IrregularVerbDictionary DictionaryType = 3
	PhraseDictionary        DictionaryType = 4
)

func (t DictionaryType) IsValid() bool {
	return t >= SimpleDictionary && t <= PhraseDictionary
}

type DictionaryLang string

const (
	RuDictionary DictionaryLang = "ru"
	EnDictionary DictionaryLang = "en"
)

func (l DictionaryLang) IsValid() bool {
	return l == RuDictionary || l == EnDictionary
}

type Dictionary struct {
	ID           DictionaryID   `json:"id" db:"id"`
	DeletedAt    *time.Time     `json:"-" db:"deleted_at"`
	Lang         DictionaryLang `json:"lang" db:"lang" validate:"required,valid_dictionary_lang"`
	Name         string         `json:"name" db:"name" validate:"min=2"`
	Type         DictionaryType `json:"type" db:"type" validate:"required,valid_dictionary_type"`
	Translations []Translation  `json:"translations,omitempty" db:"-" validate:"dive"`
	Sentences    []Sentence     `json:"sentences" db:"-" validate:"dive"`
}

type DictionaryUseCase interface {
	GetByIDs(ctx context.Context, ids []DictionaryID) ([]Dictionary, error)
	GetRandomDictionaries(ctx context.Context, lang DictionaryLang, limit uint8) ([]Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}

type DictionaryRepository interface {
	GetByIDs(ctx context.Context, ids []DictionaryID) ([]Dictionary, error)
	GetRandomDictionaries(ctx context.Context, lang DictionaryLang, limit uint8) ([]Dictionary, error)
	IsExistsByNameAndLang(ctx context.Context, name string, lang DictionaryLang) (bool, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}
