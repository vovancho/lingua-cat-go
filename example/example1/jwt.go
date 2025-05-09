package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func main() {
	const publicPEM = `
-----BEGIN CERTIFICATE-----
MIICqTCCAZECBgGUPBwS8zANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDDA1saW5n
dWEtY2F0LWdvMB4XDTI1MDEwNjE0NTI0MFoXDTM1MDEwNjE0NTQyMFowGDEWMBQG
A1UEAwwNbGluZ3VhLWNhdC1nbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBAMcBFrWq8IAHkPqT5zIqq4FmWWmZ+c+eQlzNRA33roZh+jMucNitiFBF92Ao
loD2WGJpj52rFJhAVwuspT48M8XRTASBflSgGof0ipDA2jEmQOTsZXn4wd+d/pY9
Or+DUvJKGi8rP0wR0w3n5JsbTWb3jHDJmbNorQf4/WfQ/TN/h6WhFy+gIj3oc1RO
8Iiaa0iio36cxxxcXPZxhBbnS7v4p52Y2UMsKWQqnTOPTFZjJCKRDtPWBYIsrxrM
Oy1uiS+k0MhFeNJTMZdvMEFwhRhN62ipaoHalTSBHH67op40dn+i16eDPkXDPpms
9PTyVdJcLuU1xcP9zHP8Vq6oDS0CAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAOAcV
7QMyxMuqyacZunjBvE/Ssq+HGEDA3WKoFdiOaKG0D1GYm+gjBB1hO5YjGl1O/yX9
iICSTA6oH+GfvvKY2xaYf7Wb76vLw1kRI62KiVE/gOD7ezhYtsZ6GFMGdFhtjJHC
p8bsi3z9GYipsby9eUWet9qW5nhflqpgF+z/ykHz77Hd3DdDYbszkFv3ujDJLyuZ
wJqruX0KxSFn9Vs2SPetuEZp5x1KwAaEb+e/w5WM24ya7hVZFmu3LeAP/6ffeHDr
/J2rLMgd0Dh8NJ9oNoUsXyaduQxSqUK+Uq1f8V+pAHDJ87h88rw0NohzFzjbnowp
GVag4ONioDHKW4d4Ew==
-----END CERTIFICATE-----
`

	data := []byte(publicPEM)

	block, _ := pem.Decode(data)
	if block == nil {
		fmt.Println("failed to decode PEM")
		return
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("public key is not RSA, got %T", cert.PublicKey)
		return
	}

	// Создание JWK с явным указанием kid
	key := jwk.NewRSAPublicKey()
	if err := key.FromRaw(rsaPubKey); err != nil {
		fmt.Println("Ошибка создания JWK:", err)
		return
	}

	cachedJWK := jwk.NewSet()
	cachedJWK.Add(key)

	fmt.Println("JWK загружен")

	const tokenStr = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDU5NzI2MzksImlhdCI6MTc0NTkzNjYzOSwianRpIjoiOTI3YzRiMjUtOTY0NC00MDk1LWI4ZTMtMjkzNzdlZDkwMTRmIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6IjU4NTc1NDhkLWM2ODYtNDNhNy1iZjY3LTUzMDBhMDdkNWZhMCIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.g_Y5yWgsLqFyFIjrhz6iKYa2wifPfiJtFN0_OjmogRvWUkJrvieT64vDy_phmh_psaocXj6eB5VLm29vELBZzP8bsVlDHMb1dSyRQEezVoukuENUVMgNH6jWnICnITZVge1kzyif7mtRTM6nLrMxpqZoEWY69wuQIAbiwPGjZi8_rMy2_I0X_H0rn9oGRe1MwGM1_FVmyWdgTElujkTUe5Wk6pWFBLGepoFyoKcxo1QRa_aha1vrNPSoLem88ULLOb7vipPxqZJN0JuTPHoU4r-uGDw-t4KqGPTE1wGKu-Xkzz0x4bbL9eeN9V03tkRbfAZs5rO5M20xZ3DyRRJf_w"

	token, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithVerify(jwa.RS256, rsaPubKey),
		jwt.WithValidate(true),
	)
	if err != nil {
		fmt.Printf("Неверный токен: %v\n", err) // Вывод подробной ошибки
		return
	}

	sub, ok := token.Get("sub")
	if !ok {
		fmt.Println("Отсутствует sub в токене")
		return
	}

	fmt.Println("UUID пользователя", sub)
}
