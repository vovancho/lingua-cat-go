package domain

type Translation struct {
	ID         uint64     `json:"-"`
	Dictionary Dictionary `json:"dictionary"`
}
