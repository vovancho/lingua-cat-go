











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




GROK:
напиши пример с sync/atomic на golang с объяснением, выводом, рекомендациями













