APPNAME := server
GO_PATH := $(shell go env GOPATH)
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

build:
	@echo "🔸 Build test binary...";
	@go build -buildvcs=false -race -o build/$(APPNAME).test ./cmd
	@echo "🔸 Test binary built";

run: $(REFLEX)
	$(REFLEX) -r '\.go$$' --start-service -- \
  	go run -race cmd/main.go ${args}

dependencies:
	docker compose up --remove-orphans -d

# Integration environment run all test independently using consistent test data
integration: env=integration
integration:
	@echo "🔸 Run integration tests...";
	@$(MAKE) dev_reset_db=true env=$(env) test_suite="$(test_suite_dir)/**/$(integration_test_suite)" venom
	@echo "🔸 Integration tests done";

# VENOM_PRESERVE_CASE=ON should be ON by default on venom 1.2.0
venom: $(VENOM)
	venom_var_file="./tests/venom/vars/$(env).yml"; \
	$(VENOM) run $(test_suite) \
	--output-dir=./build \
	--var-from-file "$$venom_var_file"
