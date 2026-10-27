.PHONY: all install tidy build lint fmt test docker-build docker-run

GIT_COMMIT ?= $(shell git rev-parse --verify HEAD)
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
BUILDER ?= Makefile
IMAGE ?= revive:dev
PLATFORM ?=
ARGS ?= ./...
VERSION_FLAGS := -X "github.com/mgechev/revive/cli.version=$(GIT_VERSION)" -X "github.com/mgechev/revive/cli.date=$(DATE)" -X "github.com/mgechev/revive/cli.commit=$(GIT_COMMIT)" -X "github.com/mgechev/revive/cli.builtBy=$(BUILDER)"

all: test lint build

install:
	@go mod vendor

tidy:
	@go mod tidy -diff

build: tidy
	@go build -ldflags='$(VERSION_FLAGS)'

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

docker-run: docker-build
	@docker run --rm -v '$(CURDIR)':/src -w /src '$(IMAGE)' $(ARGS)
