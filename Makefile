SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec
ROOT_DIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.DEFAULT_GOAL = help

VERSION := $(shell grep 'version:' plugin.yaml | cut -d '"' -f 2)
LDFLAGS := -X main.version=$(VERSION)

# Go configuration
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Test configuration
GINKGO        := $(GOBIN)/ginkgo
GINKGO_PROCS  ?= 3
GINKGO_FLAGS  ?= --silence-skips --procs=$(GINKGO_PROCS)

# Cyberpunk theme
CYBER_CACHE := .cyber.sh
CYBER_URL := https://raw.githubusercontent.com/Noksa/install-scripts/main/cyberpunk.sh

$(CYBER_CACHE):
	@curl -s $(CYBER_URL) > $(CYBER_CACHE)

.PHONY: cyber-update
cyber-update: ## Refresh cyberpunk theme cache
	@rm -f $(CYBER_CACHE)
	@curl -s $(CYBER_URL) > $(CYBER_CACHE)
	@source $(CYBER_CACHE) && cyber_ok "Cyberpunk theme updated"

# Test runner macro
define run_tests
	@if [ ! -f $(GINKGO) ]; then \
		echo "-> installing ginkgo CLI..."; \
		go install github.com/onsi/ginkgo/v2/ginkgo@latest; \
	fi
	@$(GINKGO) $(GINKGO_FLAGS) $(if $(2),--focus "$(2)",) $(1)
endef

##@ General

.PHONY: help
help: $(CYBER_CACHE) ## Show help
	@source $(CYBER_CACHE) && { \
		echo ""; \
		echo -e "$${CYBER_D}╔═══════════════════════════════════════╗$${CYBER_X}"; \
		echo -e "$${CYBER_D}║$${CYBER_X}  $${CYBER_M}🦋$${CYBER_X} $${CYBER_B}$${CYBER_C}Helm Resolve Deps$${CYBER_X}"; \
		echo -e "$${CYBER_D}╚═══════════════════════════════════════╝$${CYBER_X}"; \
	}
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[36mUsage:\033[0m make \033[35m<target>\033[0m\n\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m \033[37m%s\033[0m\n", $$1, $$2 } /^##@/ { printf "\n\033[35m⚡ %s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: build
build: ## Build the binary with version injection
	@go build -ldflags "$(LDFLAGS)" -o helm-resolve-deps ./cmd/resolve_deps.go

.PHONY: lint
lint: ## Run linters and formatters
	@./scripts/check.sh

##@ Testing

.PHONY: test-unit
test-unit: ## Run unit tests
	$(call run_tests,./...)

.PHONY: test
test: test-unit ## Run all tests (alias)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@go test -count=1 -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"
