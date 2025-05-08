











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

docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod init github.com/vovancho/lingua-cat-go/dictionary
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod tidy

docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod init github.com/vovancho/lingua-cat-go/exercise
docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod tidy

docker run --rm -v .\backend\analytics:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod init github.com/vovancho/lingua-cat-go/analytics
docker run --rm -v .\backend\analytics:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go mod tidy

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

analytics migrations:
docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database clickhouse://analytics:secret@localhost:59000/analytics create -ext sql -dir /migrations init_schema
docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database "clickhouse://localhost:59000?username=analytics&database=analytics&password=secret&x-multi-statement=true" up
docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database "clickhouse://localhost:59000?username=analytics&database=analytics&password=secret&x-multi-statement=true" down 1
docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database "clickhouse://localhost:59000?username=analytics&database=analytics&password=secret&x-multi-statement=true" drop -f



GROK:
напиши пример с sync/atomic на golang с объяснением, выводом, рекомендациями


docker compose restart lcg-dictionary-backend
docker compose restart lcg-exercise-backend
docker compose restart lcg-analytics-backend
docker compose logs -f lcg-dictionary-backend lcg-exercise-backend lcg-analytics-backend

docker run --rm -v .\backend\example\grpc:/defs namely/protoc-all:1.51_2 -f dictionary.proto -l go -o /defs/gen
docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs
docker run --rm -v .\backend\exercise\exercise\repository\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs


grpcurl -plaintext -d '{"limit": 4, "lang": "en"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries

docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries
docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDYxOTMzMTQsImlhdCI6MTc0NjE1NzMxNCwianRpIjoiMGI2NTE1MjctZWE3Yy00N2JlLTk4YmEtNTU1MzgxMWNhMTBlIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6IjkzZWQzY2VkLTZhODQtNDQ0Zi1hNzM2LWFkNmIxYTM3NDlkNiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.rAr2sC8Xq7QctL71_MdP7xlPtUrcrPzMID6dZtpuSTxe37QJBkFKb8KLMGsWyrIozeg3zmEoszHr2pQ7_hKpQhuAvgG1HtrrdgAoJWbIIT1GQKfTs83XYc9XYpkTAA5m-cpx-AUGPuSuXDDUlOUdODHrlU-jJK1Pe_LnGMxBqHnQGk6ZlVL_Zx-rUwM-q53PRueKUVAXu7nD8r3m_clCLu8nE5HOiWQTve3kaaR1SNwd9NTmJkx609Qer6uFNMrvE_8o0zViplISNE8Ut_Na-2h5A3I4YC6cQXAPfCvjxk5y98yNhBw0M1VUSv6tFx4A15DPRX2_9AfxYE_JZG8Tag' -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries
docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ3c1A2RW9SZUFYYlRmWTZBMTU3NEt4SFdPZlZXUTJwNTN3eEtIUjR2N0VFIn0.eyJleHAiOjE3NDYxOTMzMTQsImlhdCI6MTc0NjE1NzMxNCwianRpIjoiMGI2NTE1MjctZWE3Yy00N2JlLTk4YmEtNTU1MzgxMWNhMTBlIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLmxvY2FsaG9zdC9yZWFsbXMvbGluZ3VhLWNhdC1nbyIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiIxMWM2NWU0MS0yNDk2LTQzYWYtYWM0Yy1kYWE4OThjMjQ2NjQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJsaW5ndWEtY2F0LWdvLWRldiIsInNpZCI6IjkzZWQzY2VkLTZhODQtNDQ0Zi1hNzM2LWFkNmIxYTM3NDlkNiIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiaHR0cDovL2xpbmd1YS1jYXQtZ28ubG9jYWxob3N0Il0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLWxpbmd1YS1jYXQtZ28iLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibGluZ3VhLWNhdC1nby1kZXYiOnsicm9sZXMiOlsiVklTSVRPUiJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJkdW1teS11c2VyIiwiZW1haWwiOiJkZXYtdXNlckBtYWlsLmRldiJ9.rAr2sC8Xq7QctL71_MdP7xlPtUrcrPzMID6dZtpuSTxe37QJBkFKb8KLMGsWyrIozeg3zmEoszHr2pQ7_hKpQhuAvgG1HtrrdgAoJWbIIT1GQKfTs83XYc9XYpkTAA5m-cpx-AUGPuSuXDDUlOUdODHrlU-jJK1Pe_LnGMxBqHnQGk6ZlVL_Zx-rUwM-q53PRueKUVAXu7nD8r3m_clCLu8nE5HOiWQTve3kaaR1SNwd9NTmJkx609Qer6uFNMrvE_8o0zViplISNE8Ut_Na-2h5A3I4YC6cQXAPfCvjxk5y98yNhBw0M1VUSv6tFx4A15DPRX2_9AfxYE_JZG8Tag' -d '{\"ids\": [55,56,57,58]}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetDictionaries

Measure-Command { docker run --rm --network host -v .\backend\dictionary\dictionary\delivery\grpc\proto:/proto fullstorydev/grpcurl:latest -plaintext -import-path /proto -proto /proto/dictionary.proto -d '{\"limit\": 4, \"lang\": \"en\"}' api.lingua-cat-go.localhost:50051 dictionary.DictionaryService/GetRandomDictionaries } | % { Write-Host "Execution time: $($_.TotalMilliseconds) ms" }


docker run --rm -v .\backend\dictionary\internal\misc:/keystore openjdk:17-alpine keytool -list -rfc -keystore /keystore/keystore.jks -storepass 241186

docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go run cmd/jwk_from_pem.go
docker run --rm -v .\backend\dictionary:/app -w /app golang:1.24-alpine3.21 go get github.com/google/wire@v0.6.0



docker run --rm -v .\backend\dictionary\example:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go run jwt.go
docker run --rm -v .\backend\exercise:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go run app/main_outbox.go





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
docker run --rm -v .\backend\analytics:/app -v pkgmod:/go/pkg/mod -w /app/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"


docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase -run TestDictionaryUseCase_GetByID/Timeout
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go test -v ./dictionary/usecase -run TestDictionaryUseCase_GetByID/InternalTimeout



'{"event":"exercise_completed","user_id":123}' | docker run --rm -i --network host apache/kafka:3.7.1 /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9093 --topic lcg_exercise_completed
'{"event":"exercise_completed","user_id":123}' | docker run --rm -i --network lingua-cat-go_default apache/kafka:3.7.1 /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server lcg-kafka:9092 --topic lcg_exercise_completed
docker run --rm --network lingua-cat-go_default apache/kafka:3.7.1 /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server lcg-kafka:9092 --topic lcg_exercise_completed --from-beginning





echo '{"event":"exercise_completed","user_id":123}' |  /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server lcg-kafka:9092 --topic lcg_exercise_completed
/opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server lcg-kafka:9092 --topic lcg_exercise_completed --from-beginning



docker ps --format "{{.Names}}"


------------------------------------- Получить KEYCLOAK_ADMIN_TOKEN ----------------------------------------------------
docker compose exec lcg-exercise-backend sh

curl -X POST --location "http://lcg-keycloak/realms/lingua-cat-go/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -H "Accept: application/json" \
    -d 'grant_type=client_credentials&client_id=lingua-cat-go-admin&client_secret=OZhDEZkDUVcDCrkhgAERGrUITRQ1LhiR'

"access_token" => KEYCLOAK_ADMIN_TOKEN
------------------------------------------------------------------------------------------------------------------------



docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
docker run --rm -v .\backend\dictionary:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs



docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/gen-grpc-gateway:1.51_2 -f proto/dictionary.proto -s DictionaryService
docker run --rm -v .\backend\dictionary\dictionary\delivery\grpc:/defs -w /defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -I /defs --go_out=/defs/gen --go_opt=paths=source_relative --go-grpc_out=/defs/gen --go-grpc_opt=paths=source_relative --grpc-gateway_out=/defs/gen --grpc-gateway_opt=paths=source_relative


docker run --rm -v $(pwd):/defs -w /defs namely/protoc-all:1.51_2 \
  -f dictionary.proto \
  -l go \
  -I /defs \
  --go_out=./gen \
  --go_opt=paths=source_relative \
  --go-grpc_out=./gen \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./gen \
  --grpc-gateway_opt=paths=source_relative


powershell:
docker run --rm -v ${PWD}/backend/dictionary:/defs namely/protoc-all:1.51_2 -f dictionary/delivery/grpc/proto/dictionary.proto -l openapi -o . --with-openapi-json-names

Вывод:
Language openapi is not supported. Please specify one of: go ruby csharp java python objc gogo php node typescript web cpp descriptor_set scala


docker run --rm -v ${PWD}/backend/dictionary:/defs namely/gen-grpc-gateway:1.51_2 -f dictionary/delivery/grpc/proto/dictionary.proto -s DictionaryService


docker run --rm quay.io/skopeo/stable list-tags docker://ghcr.io/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2


powershell:
docker run --rm -v ${PWD}:/defs ghcr.io/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2:latest -I /defs -I /usr/include --openapiv2_out /defs/backend/dictionary/openapi --openapiv2_opt logtostderr=true /defs/backend/dictionary/delivery/grpc/proto/dictionary.proto

Вывод:
Unable to find image 'ghcr.io/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2:latest' locally
docker: Error response from daemon: Head "https://ghcr.io/v2/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2/manifests/latest": denied

docker run --rm -v ${PWD}:/defs ghcr.io/grpc-ecosystem/grpc-gateway/dev protoc -I /defs -I /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v2.25.1/third_party/googleapis --openapiv2_out /defs/backend/dictionary/openapi --openapiv2_opt logtostderr=true /defs/backend/dictionary/delivery/grpc/proto/dictionary.proto


docker run --rm -v ${PWD}/backend/dictionary/dictionary/delivery/grpc:/defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs

Сгенерировать grpc-gateway и swagger.json:
docker run --rm -v ${PWD}/backend/dictionary/dictionary/delivery/grpc:/defs namely/protoc-all:1.51_2 -f proto/dictionary.proto -l go -o /defs --with-gateway --with-openapi-json-names --generate-unbound-methods










