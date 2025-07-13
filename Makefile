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

.PHONY: unit_test 
unit_test: ## run unit test
	${PKG_SCRIPTS}/unittest.sh

.PHONY: integration_test 
integration_test: ## run integration test
	${PKG_SCRIPTS}/integration.sh

.PHONY: api_test
api_test: ## run api tests
	${PKG_SCRIPTS}/api-test.sh

.PHONY: tidy
tidy: ## Tidy up the go.mod
	go mod tidy

.PHONY: lint
lint: ## Run linters
	golangci-lint run --timeout 10m --config .golangci.yml

testenv: ## Setup testenv and migrate db
	${PKG_SCRIPTS}/testenv.sh

localenv: ## Setup localenv and migrate db
	${PKG_SCRIPTS}/localenv.sh

.PHONY: generate
generate: ## generate api and mocks
	./scripts/gen-api.sh && rm -rf test/mocks/* && mockery && go mod tidy

.PHONY: run
run: ## start the app locally
	go run cmd/auth-service/main.go

.PHONY: deps
deps: ## install library
	go install github.com/vektra/mockery/v3@v3.4.0
	go install github.com/wadey/gocovmerge@latest

.PHONY: keys
keys: ## generate rsa keys
	${PKG_SCRIPTS}/gen-key.sh