











https://wkrzywiec.medium.com/create-and-configure-keycloak-oauth-2-0-authorization-server-f75e2f6f6046






docker compose exec -it lcg-keycloak -Dkeycloak.migration.action=export -Dkeycloak.migration.provider=singleFile -Dkeycloak.migration.realmName=lingua-cat-go -Dkeycloak.migration.usersExportStrategy=REALM_FILE -Dkeycloak.migration.file=/temp/realm-lingua-cat-go.json



docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --help
docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --realm lingua-cat-go --users realm_file --dir /tmp/keycloak-export
docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh import --file /tmp/keycloak-export/lingua-cat-go-realm.json
docker compose exec -it lcg-keycloak ls -l /tmp/keycloak-export

docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --realm lingua-cat-go --users realm_file --dir /tmp/keycloak-export



docker run --rm -v .\backend\example:/go golang:1.24-alpine3.21 go -C sync-wait-group mod init github.com/vovancho/lingua-cat-go/example
docker run --rm -v .\backend\example:/go golang:1.24-alpine3.21 go -C sync-wait-group mod tidy

docker run --rm -v .\backend\example:/go golang:1.24-alpine3.21 go -C sync-wait-group run main.go
docker run --rm -v .\backend\example:/go golang:1.24-alpine3.21 go -C sync-once run main.go

---
docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go mod init github.com/vovancho/lingua-cat-go/dictionary
docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go run app/main.go
docker compose run lcg-dictionary-backend go run app/main.go

docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod init github.com/vovancho/lingua-cat-go/exercise
docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod tidy

dictionary migrations:
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable create -ext sql -dir /migrations init_schema
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable up
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable down 1
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable drop -f

exercise migrations:
docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable create -ext sql -dir /migrations init_schema
docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable up
docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable down 1
docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable drop -f



GROK:
напиши пример с sync/atomic на golang с объяснением, выводом, рекомендациями


docker compose restart lcg-dictionary-backend
docker compose restart lcg-exercise-backend
docker compose logs -f lcg-dictionary-backend lcg-exercise-backend

docker run --rm -v .\backend\example\grpc:/defs namely/protoc-all:1.51_2 -f dictionary.proto -l go -o /defs/gen
docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs


grpcurl -plaintext -d '{"limit": 4, "lang": "en"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries

docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries
docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDYwMjU3MzQsImlhdCI6MTc0NTk4OTczNCwianRpIjoiMzQzZGU1NjItMmViMS00MDgwLTg0ZWQtNmQ0ZDhiMzZiZTlkIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6ImY2MDU0ZmY4LWI4ZDAtNGZhNS1hMjJkLWFjZTQxNjMyNGZjMiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.mmb0jwcjnd-Fl2gb47rKtniLbhGBT9C4CO4jO8L4Zw3gVmYF75P_Oz9h7s011bebZEn3sh9ldyaSmbZP8kMjlUY7w8bkUoC4M0K4boMoYSb10-pbHBej3yvfEZqbUZWxANoUKNfMco7XkIRK3Mb3DNJAg05lehC0UE9S7L2I7c9DX-tcSvIziw10asMDEF9UBS8mUMGyK4D57H2FjwzxeU8Py1_BDD1V4xFvl6f2-lpuQPv1qEiGAMaMG6A2JmvL9Z3FWaAlGUN99GhX7CTEPRtDd_2AsoEybTgMnP8yO5BvR6rPTCXWYfxUSmpKzlIhP3u2rSu81sYPyh_qYWHX8g' -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries

Measure-Command { docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries } | % { Write-Host "Execution time: $($_.TotalMilliseconds) ms" }


docker run --rm -v .\backend\dictionary\internal\misc:/keystore openjdk:17-alpine keytool -list -rfc -keystore /keystore/keystore.jks -storepass 241186

docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go run cmd/jwk_from_pem.go
docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go get github.com/google/wire@v0.6.0



docker run --rm -v .\backend\dictionary\example:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go run jwt.go





Конфиг slog:
```
func init() {
	// Настройка структурированного логирования
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{Key: slog.TimeKey, Value: slog.StringValue(a.Value.Time().UTC().Format(time.RFC3339Nano))}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
```

docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"
docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"


docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase -run TestDictionaryUseCase_GetByID/Timeout
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase -run TestDictionaryUseCase_GetByID/InternalTimeout

















