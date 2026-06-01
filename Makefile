BINARY := peektea
INSTALL_DIR := $(shell go env GOPATH)/bin

.PHONY: build install start config

build:
	go build -o $(BINARY) .

install:
	go build -o $(INSTALL_DIR)/$(BINARY) .

start:
	air

config:
	@if [ -f $(HOME)/.peektea.toml ]; then \
		echo "$(HOME)/.peektea.toml already exists — not overwriting"; \
	else \
		cp default.peektea.toml $(HOME)/.peektea.toml && echo "Created $(HOME)/.peektea.toml"; \
	fi
