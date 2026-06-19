.PHONY: test build run integration-test

test:
	go test ./...

build:
	go build -o bin/swe-workshop-server ./cmd/server

run:
	docker compose up --build

integration-test:
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
