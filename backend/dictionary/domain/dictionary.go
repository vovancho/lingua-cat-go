package domain

import (
	"context"
	"time"
)

type DictionaryID uint64

type Dictionary struct {
	ID           DictionaryID  `json:"id" db:"id"`
	DeletedAt    *time.Time    `json:"-" db:"deleted_at"`
	Lang         string        `json:"lang" db:"lang"`
	Name         string        `json:"name" db:"name"`
	Type         uint16        `json:"type" db:"type"`
	Sentences    []Sentence    `json:"sentences" db:"-"`
	Translations []Translation `json:"translations,omitempty" db:"-"`
}

type DictionaryUseCase interface {
	GetByID(ctx context.Context, id DictionaryID) (Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}

type DictionaryRepository interface {
	GetByID(ctx context.Context, id DictionaryID) (Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id DictionaryID, name string) error
	Delete(ctx context.Context, id DictionaryID) error
}
