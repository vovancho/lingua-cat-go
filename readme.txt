











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

migrations:
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable create -ext sql -dir /migrations init_schema
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable up
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable down 1
docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable drop -f



GROK:
напиши пример с sync/atomic на golang с объяснением, выводом, рекомендациями


docker compose restart lcg-dictionary-backend
docker compose logs -f lcg-dictionary-backend

docker run --rm -v .\backend\example\grpc:/defs namely/protoc-all:1.51_2 -f dictionary.proto -l go -o /defs/gen
docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs


grpcurl -plaintext -d '{"limit": 4, "lang": "en"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries

docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries
docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDU5NjkzOTQsImlhdCI6MTc0NTkzMzM5NCwianRpIjoiNWU3ZTY4ODItN2RiNy00NGIzLThjZDctNjkwOWQ5MTVlZjdiIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6ImJmNjcxNDQ5LTEwNzAtNDdjMy04ZDE1LWVkMTFiMTZjMzE2NCIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.jmvjumSq0AIZnFsyGWk5VgZ_F-EOrym6cD_6ehRnl5R6HOS2ytZu1T4-MZDl_lQVuDfbOXRHoePuZOBV7QR47SyXPV3AmuTMtzzSn1diQT7v_Mz5-GD3-AgZxETAnTBSNp-dcdzWArrCXMFyozzoH6eREjrRDWF_DzTBgNpnp92vr2tTxw936arwy6JpAuTg9JOujg3BmJmpIKoFErYlMqtAvRmzwKuIqlkj07ituVFeAYscn48JUkJjEDMHftDGNlQ9nkFk9GuLMAaV2E_z0_Jv2gYCmHXpYG9lgGO9X2Zxkq4RsZ5_l6SD1psAcC0fgvuXX53Q925p68EbK0_nMw' -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries

Measure-Command { docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries } | % { Write-Host "Execution time: $($_.TotalMilliseconds) ms" }


docker run --rm -v .\backend\dictionary\internal\misc:/keystore openjdk:17-alpine keytool -list -rfc -keystore /keystore/keystore.jks -storepass 241186

docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go run cmd/jwk_from_pem.go



docker run --rm -v .\backend\dictionary\example:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go run jwt.go







