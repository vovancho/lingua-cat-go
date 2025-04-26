package domain

import "time"

type Sentence struct {
	ID        uint64     `json:"-" db:"id"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	TextRU    string     `json:"text_ru" db:"text_ru"`
	TextEN    string     `json:"text_en" db:"text_en"`
}
