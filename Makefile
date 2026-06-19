.PHONY: test build run integration-test keycloak-run keycloak-token

test:
	go test ./...

build:
	go build -o bin/swe-workshop-server ./cmd/server

run:
	docker compose up --build

integration-test:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

keycloak-run:
	docker compose -f docker-compose.keycloak.yml up

keycloak-token:
	./scripts/keycloak-token.sh
