include .bingo/Variables.mk

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)
MDOX_VALIDATE_CONFIG ?= .mdox.validate.yaml

help: ## Displays help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: test
test:
	go test -race ./...


.PHONY: deps
deps: ## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

.PHONY: docs
docs: ## Generates docs for all thanos commands, localise links, ensure GitHub format.
docs: $(MDOX)
	@echo ">> generating docs"
	PATH="${PATH}:$(GOBIN)" $(MDOX) fmt README.md
	$(MAKE) white-noise-cleanup

.PHONY: check-docs
check-docs: ## Checks docs against discrepancy with flags, links, white noise.
check-docs: $(MDOX)
	@echo ">> checking docs"
	PATH="${PATH}:$(GOBIN)" $(MDOX) fmt -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) README.md
	$(MAKE) white-noise-cleanup
	$(call require_clean_work_tree,'run make docs and commit changes')

.PHONY: format
format: $(GOIMPORTS)
	@echo ">> formatting go code"
	@gofmt -s -w $(FILES_TO_FMT)
	@$(GOIMPORTS) -w $(FILES_TO_FMT)

.PHONY:lint
lint: format deps $(GOLANGCI_LINT) $(FAILLINT) $(COPYRIGHT) docs
	$(call require_clean_work_tree,'detected not clean work tree before running lint, previous job changed something?')
	@echo ">> verifying modules being imported"
	@# TODO(bwplotka): Add, Printf, DefaultRegisterer, NewGaugeFunc and MustRegister once exception are accepted. Add fmt.{Errorf}=github.com/pkg/errors.{Errorf} once https://github.com/fatih/faillint/issues/10 is addressed.
	@$(FAILLINT) -paths "errors=github.com/efficientgo/core/errors,\
github.com/prometheus/prometheus/pkg/testutils=github.com/efficientgo/core/testutil,\
github.com/stretchr/testify=github.com/efficientgo/core/testutil" ./...
	@$(FAILLINT) -paths "fmt.{Print,Println,Sprint}" -ignore-tests ./...
	@echo ">> GOLANGCI_LINT is disabled as this module is in development."
	#@echo ">> linting all of the Go files GOGC=${GOGC}"
	#@$(GOLANGCI_LINT) run
	@echo ">> ensuring Copyright headers"
	@$(COPYRIGHT) $(shell go list -f "{{.Dir}}" ./... | xargs -i find "{}" -name "*.go")
	$(call require_clean_work_tree,'detected files without copyright, run make lint and commit changes')

.PHONY: white-noise-cleanup
white-noise-cleanup: ## Cleans up white noise in docs.
white-noise-cleanup:
	@echo ">> cleaning up white noise"
	@find . -type f \( -name "*.md" \) | SED_BIN="$(SED)" xargs scripts/cleanup-white-noise.sh
