package domain

type Sentence struct {
	ID     uint64 `json:"-"`
	TextRU string `json:"text_ru"`
	TextEN string `json:"text_en"`
}
