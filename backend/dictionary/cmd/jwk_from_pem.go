package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
)

func loadJWKFromPEM(pemFile string) (jwk.Set, error) {
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA, got %T", cert.PublicKey)
	}

	key := jwk.NewRSAPublicKey()
	if err := key.FromRaw(rsaPubKey); err != nil {
		return nil, err
	}

	set := jwk.NewSet()
	set.Add(key)
	return set, nil
}

func main() {
	jwkSet, err := loadJWKFromPEM("./internal/misc/public.pem")
	if err != nil {
		log.Fatalf("Ошибка загрузки JWK: %v", err)
	}
	log.Println("JWK успешно загружен:", jwkSet)
}
