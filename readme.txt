











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

Measure-Command { docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries } | % { Write-Host "Execution time: $($_.TotalMilliseconds) ms" }


docker run --rm -v .\backend\dictionary\internal\misc:/keystore openjdk:17-alpine keytool -list -rfc -keystore /keystore/keystore.jks -storepass 241186

docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go run cmd/jwk_from_pem.go











