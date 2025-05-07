package domain

import (
	"context"
	"github.com/vovancho/lingua-cat-go/analytics/internal/auth"
)

type User struct {
	ID       auth.UserID `json:"id"`
	Username string      `json:"username"`
}

type UserUseCase interface {
	GetByID(ctx context.Context, userId auth.UserID) (*User, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, userId auth.UserID) (*User, error)
}
