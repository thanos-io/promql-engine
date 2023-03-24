include .bingo/Variables.mk

FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)
MDOX_VALIDATE_CONFIG ?= .mdox.validate.yaml

# if macos, use gsed
SED ?= $(shell which gsed 2>/dev/null || which sed)

LINT_DIRS = $(shell go list ./... | grep -v "internal/prometheus")

define require_clean_work_tree
	@git update-index -q --ignore-submodules --refresh

	@if ! git diff-files --quiet --ignore-submodules --; then \
		echo >&2 "cannot $1: you have unstaged changes."; \
		git diff-files --name-status -r --ignore-submodules -- >&2; \
		echo >&2 "Please commit or stash them."; \
		exit 1; \
	fi

	@if ! git diff-index --cached --quiet HEAD --ignore-submodules --; then \
		echo >&2 "cannot $1: your index contains uncommitted changes."; \
		git diff-index --cached --name-status -r --ignore-submodules HEAD -- >&2; \
		echo >&2 "Please commit or stash them."; \
		exit 1; \
	fi
endef

help: ## Displays help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-z0-9A-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: test
test: ## Runs all Go unit tests.
export GOCACHE=/tmp/cache
test:
	@echo ">> running unit tests (without cache)"
	@rm -rf $(GOCACHE)
	@go test -v -race -timeout=5m $(shell go list ./...);

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
	@# TODO(bwplotka): Add, Printf, DefaultRegisterer, NewGaugeFunc and MustRegister once exception are accepted.
	@$(FAILLINT) -paths "errors=github.com/efficientgo/core/errors,\
fmt.{Errorf}=github.com/efficientgo/core/errors.{Wrap,Wrapf},\
github.com/prometheus/prometheus/pkg/testutils=github.com/efficientgo/core/testutil,\
github.com/stretchr/testify=github.com/efficientgo/core/testutil" $(LINT_DIRS)
	@$(FAILLINT) -paths "fmt.{Print,Println,Sprint,Errorf}" -ignore-tests $(LINT_DIRS)
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	@$(GOLANGCI_LINT) run
	@echo ">> ensuring Copyright headers"
	@$(COPYRIGHT) $(shell echo $LINT_DIRS | xargs -i find "{}" -name "*.go")
	$(call require_clean_work_tree,'detected files without copyright, run make lint and commit changes')

.PHONY: white-noise-cleanup
white-noise-cleanup: ## Cleans up white noise in docs.
white-noise-cleanup:
	@echo ">> cleaning up white noise"
	@find . -type f \( -name "*.md" \) | SED_BIN="$(SED)" xargs scripts/cleanup-white-noise.sh

benchmarks:
	@mkdir -p benchmarks

.PHONY: bench-old
bench-old: benchmarks
	@echo "Benchmarking old engine"
	@go test ./... -bench 'BenchmarkRangeQuery/.*/old_engine'  -run none -count 5 | sed -u 's/\/old_engine//' > benchmarks/old.out
	@go test ./... -bench 'BenchmarkNativeHistograms/.*/old_engine'  -run none -count 5 | sed -u 's/\/old_engine//' >> benchmarks/old.out

.PHONY: bench-new
bench-new: benchmarks
	@echo "Benchmarking new engine"
	@go test ./... -bench 'BenchmarkRangeQuery/.*/new_engine'  -run none -count 5 | sed -u 's/\/new_engine//' > benchmarks/new.out
	@go test ./... -bench 'BenchmarkNativeHistograms/.*/new_engine'  -run none -count 5 | sed -u 's/\/new_engine//' >> benchmarks/new.out

.PHONY: benchmark
benchmark: bench-old bench-new
	@benchstat benchmarks/old.out benchmarks/new.out

.PHONY: sync-parser
sync-parser:
	@echo "Cleaning existing directories"
	@rm -rf internal/prometheus/parser
	@mkdir -p tmp
	@rm -rf tmp/prometheus
	@echo "Cloning prometheus"
	@git clone git@github.com:prometheus/prometheus.git tmp/prometheus
	@echo "Copying parser"
	cp -r tmp/prometheus/promql/parser internal/prometheus
	@echo "Cleaning up"
	@rm -rf tmp
