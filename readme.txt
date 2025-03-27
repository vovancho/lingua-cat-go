











https://wkrzywiec.medium.com/create-and-configure-keycloak-oauth-2-0-authorization-server-f75e2f6f6046






docker compose exec -it lcg-keycloak -Dkeycloak.migration.action=export -Dkeycloak.migration.provider=singleFile -Dkeycloak.migration.realmName=lingua-cat-go -Dkeycloak.migration.usersExportStrategy=REALM_FILE -Dkeycloak.migration.file=/temp/realm-lingua-cat-go.json



docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --help
docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --realm lingua-cat-go --users realm_file --dir /tmp/keycloak-export
docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh import --file /tmp/keycloak-export/lingua-cat-go-realm.json
docker compose exec -it lcg-keycloak ls -l /tmp/keycloak-export

docker compose exec -it lcg-keycloak /opt/keycloak/bin/kc.sh export --realm lingua-cat-go --users realm_file --dir /tmp/keycloak-export


