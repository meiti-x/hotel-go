# Variables
BINARY_PATH=bin/
BINARY_NAME=go-hotel
SRC_DIR=./src
PKG_DIR=$(SRC_DIR)/...
CMD_DIR=$(SRC_DIR)/cmd/...

all: build

build:
	go build -o $(BINARY_PATH)$(BINARY_NAME) $(CMD_DIR)

dev:
	go run $(CMD_DIR)

test:
	go test -v ./...

clean:
	rm -rf $(BINARY_PATH)

fmt:
	go fmt $(PKG_DIR)

run: build
	./$(BINARY_PATH)$(BINARY_NAME)
