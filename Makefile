PKG_SCRIPTS=scripts

default: help

help: ## Show help for each of the Makefile commands
	@awk 'BEGIN \
		{FS = ":.*##"; printf "Usage: make ${cyan}<command>\n${white}Commands:\n"} \
		/^[a-zA-Z_-]+:.*?##/ \
		{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
		$(MAKEFILE_LIST)

.PHONY: coverage 
coverage: ## code coverage
	${PKG_SCRIPTS}/coverage.sh

.PHONY: tidy
tidy: ## Tidy up the go.mod
	go mod tidy

.PHONY: lint
lint: ## Run linters
	golangci-lint run --timeout 10m --config .golangci.yml

testenv: ## Setup testenv and migrate db
	${PKG_SCRIPTS}/testenv.sh

.PHONY: generate
generate: ## generate api and mocks
	./scripts/gen-api.sh && rm -rf test/mocks/* && mockery && go mod tidy

.PHONY: run
run: ## start the app locally
	go run cmd/auth-service/main.go

.PHONY: deps
deps: ## install library
	go install github.com/vektra/mockery/v3@v3.4.0
