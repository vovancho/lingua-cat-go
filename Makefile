up:
	docker compose up -d --remove-orphans

stop:
	docker compose stop

test-dictionary:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/dictionary golang:1.24-alpine3.21 go test -v ./usecase

test-exercise:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/exercise golang:1.24-alpine3.21 go test -v ./usecase

test-analytics:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/analytics golang:1.24-alpine3.21 go test -v ./usecase

test-all: test-dictionary test-exercise test-analytics

wire-dictionary:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/dictionary/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"

wire-exercise:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/exercise/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"

wire-analytics:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app/analytics/internal/wire golang:1.24-alpine3.21 sh -c "go install github.com/google/wire/cmd/wire@latest && wire"

wire-all: wire-dictionary wire-exercise wire-analytics

logs:
	docker compose logs -f lcg-dictionary-backend lcg-exercise-backend lcg-analytics-backend

restart-dictionary:
	docker compose restart lcg-dictionary-backend

restart-exercise:
	docker compose restart lcg-exercise-backend

restart-analytics:
	docker compose restart lcg-analytics-backend

restart-all:
	docker compose restart lcg-dictionary-backend lcg-exercise-backend lcg-analytics-backend

tidy-dictionary:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 sh -c "apk add --no-cache git && go -C dictionary mod tidy"

tidy-exercise:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 sh -c "apk add --no-cache git && go -C exercise mod tidy"

tidy-analytics:
	docker run --rm -v .\backend:/app -v pkgmod:/go/pkg/mod -w /app golang:1.24-alpine3.21 sh -c "apk add --no-cache git && go -C analytics mod tidy"

proto-dictionary:
	docker run --rm --entrypoint sh -v .\backend:/defs namely/protoc-all:1.51_2 -c "entrypoint.sh -f proto/dictionary.proto -l go -o dictionary/delivery/grpc --with-gateway --with-openapi-json-names --generate-unbound-methods && mv dictionary/delivery/grpc/proto/dictionary.swagger.json dictionary/docs/grpc-gw-swagger.json && rm -r dictionary/delivery/grpc/proto"

proto-exercise:
	docker run --rm --entrypoint sh -v .\backend:/defs namely/protoc-all:1.51_2 -c "entrypoint.sh -f proto/dictionary.proto -l go -o exercise/repository/grpc"

migrations-reset-dictionary:
	docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable drop -f
	docker run --rm -v .\backend\dictionary\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://dictionary:secret@localhost:54321/dictionary?sslmode=disable up

migrations-reset-exercise:
	docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable drop -f
	docker run --rm -v .\backend\exercise\migrations:/migrations --network host migrate/migrate -path /migrations -database postgres://exercise:secret@localhost:54322/exercise?sslmode=disable up

migrations-reset-analytics:
	docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database "clickhouse://localhost:59000?username=analytics&database=analytics&password=secret&x-multi-statement=true" drop -f
	docker run --rm -v .\backend\analytics\migrations:/migrations --network host migrate/migrate -path /migrations -database "clickhouse://localhost:59000?username=analytics&database=analytics&password=secret&x-multi-statement=true" up

migrations-reset-all: migrations-reset-dictionary migrations-reset-exercise migrations-reset-analytics

swagger-dictionary:
	docker run --rm -v .\backend:/code -v pkgmod:/go/pkg/mod -w /code/dictionary ghcr.io/swaggo/swag:v1.16.4 init -d /code/dictionary,/code/pkg --ot json -g cmd/main.go -o ./docs

swagger-exercise:
	docker run --rm -v .\backend:/code -v pkgmod:/go/pkg/mod -w /code/exercise ghcr.io/swaggo/swag:v1.16.4 init -d /code/exercise,/code/pkg --ot json -g cmd/main.go -o ./docs

swagger-analytics:
	docker run --rm -v .\backend:/code -v pkgmod:/go/pkg/mod -w /code/analytics ghcr.io/swaggo/swag:v1.16.4 init -d /code/analytics,/code/pkg --ot json -g cmd/main.go -o ./docs

keycloak-export:
	docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --realm lingua-cat-go --users realm_file --dir /tmp/keycloak-export

keycloak-import:
	docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh import --file /tmp/keycloak-export/lingua-cat-go-realm.json

tree:
	docker run --rm -v .\backend:/src -w /src johnfmorton/tree-cli tree -o backend.txt -l 10

ab-start:
	docker run --rm -v .\project\ab:/src -w /src --network host ricsanfre/docker-curl-jq /src/ab_test_chain.sh
