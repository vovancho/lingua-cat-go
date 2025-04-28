package domain

import (
	"context"
	"time"
)

type DictionaryID uint64
type DictionaryType uint16

const (
	SimpleDictionary        DictionaryType = 1
	PhrasalVerbDictionary                  = 2
	IrregularVerbDictionary                = 3
	PhraseDictionary                       = 4
)

func (t DictionaryType) IsValid() bool {
	return t >= SimpleDictionary && t <= PhraseDictionary
}

type DictionaryLang string

const (
	RuDictionary DictionaryLang = "ru"
	EnDictionary                = "en"
)

func (l DictionaryLang) IsValid() bool {
	return l == RuDictionary || l == EnDictionary
}

type Dictionary struct {
	ID           DictionaryID   `json:"id" db:"id"`
	DeletedAt    *time.Time     `json:"-" db:"deleted_at"`
	Lang         DictionaryLang `json:"lang" db:"lang"`
	Name         string         `json:"name" db:"name"`
	Type         DictionaryType `json:"type" db:"type"`
	Sentences    []Sentence     `json:"sentences" db:"-"`
	Translations []Translation  `json:"translations,omitempty" db:"-"`
}

type DictionaryUseCase interface {
	GetByID(ctx context.Context, id DictionaryID) (*Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}

type DictionaryRepository interface {
	GetByID(ctx context.Context, id DictionaryID) (*Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}
