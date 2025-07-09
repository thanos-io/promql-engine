FILES_TO_FMT ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)
MDOX_VALIDATE_CONFIG ?= .mdox.validate.yaml

# if macos, use gsed
SED ?= $(shell which gsed 2>/dev/null || which sed)
BENCHSTAT = go tool -modfile go.tools.mod benchstat
MDOX = go tool -modfile go.tools.mod mdox
GCI = go tool -modfile go.tools.mod gci
FAILLINT = go tool -modfile go.tools.mod faillint
GOLANGCI_LINT = go tool -modfile go.tools.mod golangci-lint
MODERNIZE = go tool -modfile go.tools.mod modernize
COPYRIGHT = go run github.com/efficientgo/tools/copyright@v0.0.0-20220225185207-fe763185946b

GOMODULES = $(shell go list ./...)

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
	@echo ">> running unit tests (without cache)"
	@rm -rf /tmp/engine-cache
	@GORACE=atexit_sleep_ms=0 GOCACHE=/tmp/engine-cache go test -race -timeout=10m $(GOMODULES);

.PHONY: test-fast
test-fast: ## Runs all Go unit tests without race detector.
	@echo ">> running unit tests without race detection (with cache)"
	@go test -count=1 -timeout=10m $(GOMODULES);

.PHONY: test-slicelabels
test-slicelabels: ## Runs all Go unit tests with slicelabels flag.
	@export GOCACHE=/tmp/cache
	@echo ">> running unit tests with slicelabels flag (without cache)"
	@rm -rf $(GOCACHE)
	@go test -race --tags=slicelabels -timeout=10m $(GOMODULES);

.PHONY: fuzz
fuzz: ## Runs selected fuzzing tests
	@export GOCACHE=/tmp/cache
	@echo ">> running fuzz tests (without cache)"
	@rm -rf $(GOCACHE)
	@go test github.com/thanos-io/promql-engine/engine -run None -fuzz FuzzEnginePromQLSmithInstantQuery -fuzztime=90s -fuzzminimizetime 0x;
	@go test github.com/thanos-io/promql-engine/engine -run None -fuzz FuzzNativeHistogramQuery -fuzztime=90s -fuzzminimizetime 0x;
	@go test github.com/thanos-io/promql-engine/logicalplan -run None -fuzz FuzzNodesMarshalJSON -fuzztime=30s -fuzzminimizetime 0x;

.PHONY: deps
deps: ## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

.PHONY: docs
docs: ## Generates docs for all thanos commands, localise links, ensure GitHub format.
docs:
	@echo ">> generating docs"
	$(MDOX) fmt README.md
	$(MAKE) white-noise-cleanup

.PHONY: check-docs
check-docs: ## Checks docs against discrepancy with flags, links, white noise.
check-docs:
	@echo ">> checking docs"
	$(MDOX) fmt -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) README.md
	$(MAKE) white-noise-cleanup
	$(call require_clean_work_tree,'run make docs and commit changes')

.PHONY: format
format:
	@echo ">> formatting promql tests"
	@go run scripts/testvet/main.go -json -fix ./...
	@echo ">> formatting imports"
	@$(GCI) write $(shell find . -name "*.go") -s "standard" -s "prefix(github.com/thanos-io)" -s "default" -s "blank" -s "dot" --custom-order

.PHONY:lint
lint: format deps docs
	$(call require_clean_work_tree,'detected not clean work tree before running lint, previous job changed something?')
	@echo ">> verifying modules being imported"
	@# TODO(bwplotka): Add, Printf, DefaultRegisterer, NewGaugeFunc and MustRegister once exception are accepted.
	@$(FAILLINT) -paths "errors=github.com/efficientgo/core/errors,\
fmt.{Errorf}=github.com/efficientgo/core/errors.{Wrap,Wrapf},\
github.com/prometheus/prometheus/pkg/testutils=github.com/efficientgo/core/testutil,\
github.com/stretchr/testify=github.com/efficientgo/core/testutil" $(GOMODULES)
	@$(FAILLINT) -paths "fmt.{Print,Println,Errorf}" -ignore-tests $(GOMODULES)
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	@$(GOLANGCI_LINT) run
	@echo ">> ensuring Copyright headers"
	@$(COPYRIGHT) $(shell find . -name "*.go")
	@echo ">> ensuring modern go style"
	@$(MODERNIZE) -test ./...t
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
	@$(BENCHSTAT) benchmarks/old.out benchmarks/new.out
