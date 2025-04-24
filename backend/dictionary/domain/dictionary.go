package domain

import (
	"context"
	"time"
)

type Dictionary struct {
	ID           uint64        `json:"id"`
	Name         string        `json:"name"`
	Type         uint16        `json:"type"`
	Lang         string        `json:"lang"`
	CreatedAt    time.Time     `json:"created_at"`
	DeletedAt    *time.Time    `json:"deleted_at"`
	Sentences    []Sentence    `json:"sentences"`
	Translations []Translation `json:"translations"`
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
