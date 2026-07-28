.PHONY: all build run test clean install

BINARY_NAME=matt
BUILD_DIR=bin
SRC_DIR=./cmd/matt
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v2.0.0")
LDFLAGS=-ldflags "-X github.com/Chintanpatel24/Matt/internal/version.Version=$(VERSION)"

all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME) 2>/dev/null || cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
