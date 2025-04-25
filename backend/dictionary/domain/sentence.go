package domain

type Sentence struct {
	id     uint64
	TextRU string `json:"text_ru"`
	TextEN string `json:"text_en"`
}
