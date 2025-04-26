package domain

import "time"

type Translation struct {
	ID         uint64     `json:"-" db:"id"`
	DeletedAt  *time.Time `json:"-" db:"deleted_at"`
	Dictionary Dictionary `json:"dictionary" db:"dictionary"`
}
