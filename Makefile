.PHONY: all gifs docker-lint lint lintmax gosec govulncheck test build goreleaser tag-major tag-minor tag-patch release bump-glazed install run-dev run-host-status run-agent-status web-install web-dev web-build web-embed storybook storybook-build oidc-e2e

all: build

VERSION=v0.1.14
GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target

TAPES=$(wildcard doc/vhs/*tape)
gifs: $(TAPES)
	for i in $(TAPES); do vhs < $$i; done

docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

lint:
	GOWORK=off golangci-lint run -v

lintmax:
	GOWORK=off golangci-lint run -v --max-same-issues=100

gosec:
	GOWORK=off go install github.com/securego/gosec/v2/cmd/gosec@latest
	GOWORK=off gosec -exclude-generated -exclude=G101,G304,G301,G306 -exclude-dir=.history ./...

govulncheck:
	GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

test:
	GOWORK=off go test ./...

build:
	GOWORK=off go generate ./...
	GOWORK=off go build ./...

goreleaser:
	GOWORK=off goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $(shell svu major)

tag-minor:
	git tag $(shell svu minor)

tag-patch:
	git tag $(shell svu patch)

release:
	git push origin --tags
	GOWORK=off GOPROXY=proxy.golang.org go list -m github.com/go-go-golems/go-go-host@$(shell svu current)

bump-glazed:
	GOWORK=off go get github.com/go-go-golems/glazed@latest
	GOWORK=off go get github.com/go-go-golems/clay@latest
	GOWORK=off go mod tidy

GO_GO_HOST_BINARY=$(shell which go-go-host)
install:
	GOWORK=off go build -o ./dist/go-go-host ./cmd/go-go-host && \
		cp ./dist/go-go-host $(GO_GO_HOST_BINARY)

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

web-embed:
	go run ./cmd/build-web

storybook:
	cd web/admin && pnpm storybook

storybook-build:
	cd web/admin && pnpm storybook:build

oidc-e2e:
	GO_GO_HOST_OIDC_E2E=1 node scripts/oidc-login-playwright.mjs

.PHONY: logcopter-generate
logcopter-generate:
	GOWORK=off go generate ./...

.PHONY: logcopter-check
logcopter-check:
	GOWORK=off go tool logcopter-gen -area-prefix go-go-golems.go-go-host -strip-prefix github.com/go-go-golems/go-go-host -check ./cmd/... ./internal/... ./pkg/...
