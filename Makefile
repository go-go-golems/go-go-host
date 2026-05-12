.PHONY: build test run-dev run-host-status run-agent-status tidy lint web-install web-dev web-build web-embed storybook storybook-build

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

web-install:
	cd web/admin && pnpm install

web-dev:
	cd web/admin && pnpm dev

web-build:
	cd web/admin && pnpm build

web-embed: web-build
	rm -rf internal/webadmin/dist
	mkdir -p internal/webadmin/dist
	cp -R web/admin/dist/. internal/webadmin/dist/

storybook:
	cd web/admin && pnpm storybook

storybook-build:
	cd web/admin && pnpm storybook:build
