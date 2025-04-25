package domain

import "time"

type Sentence struct {
	ID        uint64     `json:"-"`
	DeletedAt *time.Time `json:"-"`
	TextRU    string     `json:"text_ru"`
	TextEN    string     `json:"text_en"`
}
