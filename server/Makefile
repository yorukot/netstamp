.PHONY: dev test fmt tidy migrate-status

dev:
	go run ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

migrate-status:
	go run ./cmd/migrate -command status

