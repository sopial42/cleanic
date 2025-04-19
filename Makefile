APPNAME := server
GO_PATH := $(shell go env GOPATH)
COVERDIR := ./build/coverdata
test_suite_dir := ./tests/venom
integration_test_suite := "**/*.venom.yml"

REFLEX := $(GO_PATH)/bin/reflex
$(REFLEX):
	go install github.com/cespare/reflex@latest

# Detect OS (darwin for macOS, linux for Linux)
OS := $(shell uname | tr '[:upper:]' '[:lower:]')
# Detect architecture amd or arm
ARCH := $(shell uname -m)

ifeq ($(ARCH), x86_64)
	ARCH = amd64
else ifeq ($(ARCH), arm64)
	ARCH = arm64
else ifeq ($(ARCH), aarch64)
	ARCH = arm64
else
    $(error Unsupported architecture: $(ARCH))
endif

VENOM_VERSION := v1.2.0
VENOM := $(GO_PATH)/bin/venom-$(VENOM_VERSION)
$(VENOM):
	curl -sSfLo $(VENOM) https://github.com/ovh/venom/releases/download/$(VENOM_VERSION)/venom.$(OS)-$(ARCH)
	chmod +x $(VENOM)


golangci_lint_version := 2.1.2
GOLANGCILINT := $(GO_PATH)/bin/golangci-lint-$(golangci_lint_version)
$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_PATH)/bin v$(golangci_lint_version)
	cp $(GO_PATH)/bin/golangci-lint $(GOLANGCILINT)

.PHONY: build
build:
	@echo "🔸 Build test binary...";
# Use go test and '-c' arg to build the binary using the main_test.go
# the main_test.go is used by integration tests to include venom tests in coverage
	@go test -c -tags=integration -buildvcs=false -coverpkg="./..." -race -o build/$(APPNAME).test ./cmd
	@echo "🔸 Done";

lint: $(GOLANGCILINT)
	@echo "🔸 Run golangci-lint...";
	@$(GOLANGCILINT) run --timeout 1m ./...
	@echo "🔸 Done";

run: $(REFLEX)
	$(REFLEX) -r '\.go$$' --start-service -- \
  	go run -race cmd/main.go ${args}

dependencies:
	docker compose up --remove-orphans -d

integration: env=integration
integration: build $(VENOM)
	@echo "🔸 Start server...";
	@./build/$(APPNAME).test -test.coverprofile=./build/server.venom.cover.out > ./build/$(APPNAME).log 2>&1 &
	@sleep 5;
	@echo "🔸 Run integration tests...";
	@$(MAKE) env=$(env) test_suite="$(test_suite_dir)/**/$(integration_test_suite)" venom;
	@echo "🔸 Kill server";
	@pkill $(APPNAME).test 2> /dev/null || true;
	@echo "🔸 Generate coverage report";
	@go tool cover -html=./build/server.venom.cover.out -o ./build/server.venom.cover.html;
	@echo "🔸 Done";

venom: $(VENOM)
	venom_var_file="./tests/venom/vars/$(env).yml"; \
	$(VENOM) run $(test_suite) \
	--output-dir=./build \
	--var-from-file "$$venom_var_file"

dbtty:
	@echo "[INFO] Login to psql inside db container"
	@echo "[INFO] exemple command: \dt;"
	@docker exec -it db psql -U kotai
