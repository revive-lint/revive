.PHONY: all tidy tidy-tools build goreleaser-build lint fmt test docker-build docker-run docker-smoke

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
IMAGE ?= ghcr.io/revive-lint/revive:dev
PLATFORM ?=
ARGS ?= ./...
VERSION_FLAGS := -X "github.com/mgechev/revive/cli.version=$(GIT_VERSION)" -X "github.com/mgechev/revive/cli.date=$(DATE)" -X "github.com/mgechev/revive/cli.commit=$(GIT_COMMIT)" -X "github.com/mgechev/revive/cli.builtBy=$(BUILDER)"

all: test lint build

tidy:
	@go mod tidy -diff

tidy-tools:
	@set -e; for dir in $(dir $(wildcard tools/*/go.mod)); do go -C $$dir mod tidy -diff; done

build: tidy
	@go build -ldflags='$(VERSION_FLAGS)'

goreleaser-build:
	@go tool -modfile=tools/goreleaser/go.mod goreleaser build --snapshot --clean

lint:
	revive --config revive.toml ./...
	golangci-lint run

fmt:
	golangci-lint fmt

test:
	@go test -v -race ./...

docker-build:
	@docker buildx build \
		--build-arg VERSION='$(GIT_VERSION)' \
		--build-arg REVISION='$(GIT_COMMIT)' \
		--build-arg BUILDTIME='$(DATE)' \
		--build-arg BUILDER='$(BUILDER)' \
		$(if $(PLATFORM),--platform '$(PLATFORM)') \
		--load -t '$(IMAGE)' .

docker-smoke: docker-build
	@./test/docker/smoke.sh '$(IMAGE)'

docker-run: docker-build
	@docker run --rm -v '$(CURDIR)':/app -w /app '$(IMAGE)' $(ARGS)
