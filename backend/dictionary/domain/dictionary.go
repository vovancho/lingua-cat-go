package domain

import (
	"context"
	"time"
)

type Dictionary struct {
	ID           uint64        `json:"id" db:"id"`
	DeletedAt    *time.Time    `json:"-" db:"deleted_at"`
	Lang         string        `json:"lang" db:"lang"`
	Name         string        `json:"name" db:"name"`
	Type         uint16        `json:"type" db:"type"`
	Sentences    []Sentence    `json:"sentences" db:"-"`
	Translations []Translation `json:"translations" db:"-"`
}

type DictionaryUseCase interface {
	GetByID(ctx context.Context, id uint64) (Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id uint64, name string) error
	Delete(ctx context.Context, id uint64) error
}

type DictionaryRepository interface {
	GetByID(ctx context.Context, id uint64) (Dictionary, error)
	Store(ctx context.Context, d *Dictionary) error
	ChangeName(ctx context.Context, id uint64, name string) error
	Delete(ctx context.Context, id uint64) error
}
