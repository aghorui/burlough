GO              := go
DELVE           := dlv
NPM             := npm
SRC_DIRECTORY   := src
BIN_DIRECTORY   := build
BIN_NAME        := burlough
DEBUG_BIN_NAME  := burlough_debug
PROJECT_ROOT    := $(shell pwd)


.PHONY: all clean build run run-debug test

build:
	mkdir -p $(BIN_DIRECTORY)
	cd $(SRC_DIRECTORY) && \
	$(GO) build -o $(PROJECT_ROOT)/$(BIN_DIRECTORY)/$(BIN_NAME)

build-release:
	mkdir -p $(BIN_DIRECTORY)
	cd $(SRC_DIRECTORY) && \
	$(GO) build -o $(PROJECT_ROOT)/$(BIN_DIRECTORY)/$(BIN_NAME) -ldflags '-s'  -tags RELEASE

run-debug:
	mkdir -p $(BIN_DIRECTORY)
	# You need to have delve in PATH.
	cd $(SRC_DIRECTORY) && \
	$(DELVE) debug --output $(PROJECT_ROOT)/$(BIN_DIRECTORY)/$(DEBUG_BIN_NAME) --wd ../$(BIN_DIRECTORY)

run: $(BIN_DIRECTORY)/$(BIN_NAME)
	cd $(BIN_DIRECTORY) && \
	./$(BIN_NAME)

test:
	cd $(SRC_DIRECTORY) && \
	$(GO) test ./...

clean:
	rm $(BIN_DIRECTORY)/$(BIN_NAME)
	rm $(BIN_DIRECTORY)/$(DEBUG_BIN_NAME)
