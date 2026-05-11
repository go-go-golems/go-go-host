.PHONY: build test run-dev run-host-status run-agent-status tidy lint

build:
	go build ./...

test:
	go test ./...

tidy:
	go mod tidy

lint:
	golangci-lint run -v

run-dev:
	go run ./cmd/go-go-hostd --config configs/dev.yaml

run-host-status:
	go run ./cmd/go-go-host status --api-url http://127.0.0.1:8080

run-agent-status:
	go run ./cmd/go-go-host-agent status --api-url http://127.0.0.1:8080
