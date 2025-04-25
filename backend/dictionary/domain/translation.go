package domain

import "time"

type Translation struct {
	ID         uint64     `json:"-"`
	DeletedAt  *time.Time `json:"-"`
	Dictionary Dictionary `json:"dictionary"`
}
