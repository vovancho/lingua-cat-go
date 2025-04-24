package domain

type Translation struct {
	id         uint64
	Dictionary Dictionary `json:"dictionary"`
}
