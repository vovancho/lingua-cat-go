package domain

type Sentence struct {
	id          uint64
	Text        string `json:"text"`
	Translation string `json:"translation"`
}
