package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"log/slog"
	"os"
)

type PublicKeyPath string

func NewAuthService(publicKeyPath PublicKeyPath) (*AuthService, error) {
	RSAPublicKey, err := loadRSAPublicKeyFromPEM(string(publicKeyPath))
	if err != nil {
		return nil, fmt.Errorf("failed to load JWK: %w", err)
	}
	slog.Info("JWK loaded successfully")
	return &AuthService{key: RSAPublicKey}, nil
}

type AuthService struct {
	key jwk.RSAPublicKey
}

func (s *AuthService) VerifyToken(tokenStr string) (string, error) {
	var rsaPubKey rsa.PublicKey
	if err := s.key.Raw(&rsaPubKey); err != nil {
		return "", fmt.Errorf("failed to extract RSA public key: %w", err)
	}

	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithVerify(jwa.RS256, rsaPubKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	sub, ok := token.Get("sub")
	if !ok {
		return "", fmt.Errorf("missing sub claim in token")
	}

	return fmt.Sprintf("%v", sub), nil
}

func loadRSAPublicKeyFromPEM(pemFile string) (jwk.RSAPublicKey, error) {
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA, got %T", cert.PublicKey)
	}

	key := jwk.NewRSAPublicKey()
	if err := key.FromRaw(rsaPubKey); err != nil {
		return nil, fmt.Errorf("failed to create JWK key: %w", err)
	}

	return key, nil
}
