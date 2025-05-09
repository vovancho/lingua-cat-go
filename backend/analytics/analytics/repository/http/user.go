package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/vovancho/lingua-cat-go/analytics/domain"
	"github.com/vovancho/lingua-cat-go/pkg/auth"
	"net/http"
)

type keycloakUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Config struct {
	AdminRealmEndpoint string
	AdminToken         string
}

type httpUserRepository struct {
	config Config
	client *http.Client
}

func NewHttpUserRepository(config Config, client *http.Client) domain.UserRepository {
	return &httpUserRepository{
		config: config,
		client: client,
	}
}

func (u httpUserRepository) GetByID(ctx context.Context, userId auth.UserID) (*domain.User, error) {
	url := fmt.Sprintf("%susers/%s", u.config.AdminRealmEndpoint, uuid.UUID(userId).String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+u.config.AdminToken)

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: status %d, %v", resp.StatusCode, url)
	}

	var ku keycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&ku); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	parsedID, err := uuid.Parse(ku.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	return &domain.User{
		ID:       auth.UserID(parsedID),
		Username: ku.Username,
	}, nil
}
