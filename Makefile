BINARY := peektea
INSTALL_DIR := $(shell go env GOPATH)/bin

.PHONY: build install start

build:
	go build -o $(BINARY) .

install:
	go build -o $(INSTALL_DIR)/$(BINARY) .

start:
	air
