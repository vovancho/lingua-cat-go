package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type AdminClientConfig struct {
	TokenEndpoint      string // Эндпойнт для получения и обновления токенов
	AdminRealmEndpoint string // Эндпойнт для работы с пользователями
	ClientID           string // Идентификатор клиента
	ClientSecret       string // Секрет клиента
	RefreshToken       string // Долгоживущий refresh-токен
}

type AdminClient struct {
	config      AdminClientConfig
	client      *http.Client
	accessToken string
	mu          sync.Mutex // Мьютекс для синхронизации обновления токенов
}

func NewAdminClient(config AdminClientConfig, client *http.Client) *AdminClient {
	if client == nil {
		client = http.DefaultClient
	}

	return &AdminClient{
		config: config,
		client: client,
	}
}

// GetUser возвращает пользователя по его ID
func (c *AdminClient) GetUser(ctx context.Context, userID string) (*User, error) {
	reqUrl := fmt.Sprintf("%susers/%s", c.config.AdminRealmEndpoint, userID)
	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	var user *User
	for attempts := 0; attempts < 2; attempts++ {
		// Получаем текущий access_token с блокировкой
		c.mu.Lock()
		token := c.accessToken
		c.mu.Unlock()

		req.Header.Set("Authorization", "Bearer "+token)

		// Выполняем запрос
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
		}
		defer resp.Body.Close()

		// Если токен истек, обновляем его и повторяем запрос
		if resp.StatusCode == http.StatusUnauthorized {
			if err := c.refreshToken(); err != nil {
				return nil, fmt.Errorf("ошибка обновления токена: %w", err)
			}
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("пользователь не найден")
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("ошибка получения пользователя: статус %d", resp.StatusCode)
		}

		// Декодируем ответ
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
		}
		return user, nil
	}

	return nil, fmt.Errorf("не удалось получить пользователя после обновления токена")
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// refreshToken обновляет access_token с использованием refresh_token
func (c *AdminClient) refreshToken() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Формируем данные для запроса
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("refresh_token", c.config.RefreshToken)

	// Создаем запрос на обновление токена
	req, err := http.NewRequest("POST", c.config.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Выполняем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка обновления токена: статус %d", resp.StatusCode)
	}

	// Декодируем ответ
	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	// Обновляем токен
	c.accessToken = tokenResp.AccessToken

	return nil
}
