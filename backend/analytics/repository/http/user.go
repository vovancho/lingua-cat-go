package http

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"github.com/vovancho/lingua-cat-go/pkg/keycloak"
)

type Config struct {
	AdminRealmEndpoint string
	AdminToken         string
}

type userRepository struct {
	client *keycloak.AdminClient
}

func NewUserRepository(client *keycloak.AdminClient) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r userRepository) GetByID(ctx context.Context, userId auth.UserID) (*domain.User, error) {
	keycloakUser, err := r.client.GetUser(ctx, uuid.UUID(userId).String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := &domain.User{
		ID:       auth.UserID(keycloakUser.ID),
		Username: keycloakUser.Username,
	}

	return user, nil
}
