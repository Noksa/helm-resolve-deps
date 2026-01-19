SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -euc

.DEFAULT_GOAL = help

##@ Help & Information
.PHONY: help
help: ## Show this help message with available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\n🚀 \033[1;34mHelm In Pod - Helm plugin to run commands inside pods\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1;33m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: lint
lint: ## Run linters and formatters
	@./scripts/check.sh

.PHONY: test
test: ## Run tests
	@go test -v ./...
