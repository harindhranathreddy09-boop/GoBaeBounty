.PHONY: all build clean test lint run

BINARY_NAME=gobaebounty
BUILD_DIR=bin
CMD_PATH=./cmd/bbgolang

all: build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	which golangci-lint > /dev/null || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

run: build
	$(BUILD_DIR)/$(BINARY_NAME) --help
