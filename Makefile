APPNAME := server
GO_PATH := $(shell go env GOPATH)


REFLEX := $(GO_PATH)/bin/reflex
$(REFLEX):
	go install github.com/cespare/reflex@latest

run: $(REFLEX)
	$(REFLEX) --start-service -- \
  go run -race cmd/main.go ${args}
