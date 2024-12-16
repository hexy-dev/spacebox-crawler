VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

LD_FLAGS = -X github.com/hexy-dev/spacebox/spacebox-crawler/version.Version=$(VERSION) \
	-X github.com/hexy-dev/spacebox/soacebox-crawler/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

.PHONY: fix dep build test race lint stats

fix: ## Fix fieldalignment
	fieldalignment -fix ./...

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v ./cmd/main.go

test: ## Run unittests
	@go test ./... -count=1

race: dep ## Run data race detector
	@go test -race ./... -count=1

install-linter: ## Install golangci-lint
	@mkdir -p bin
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin v1.62.2
	@bin/golangci-lint --version

lint: install-linter ## Lint the files
	./scripts/golint.sh

stats: ## Code analytics
	scc .
