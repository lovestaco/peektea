BINARY      := peektea
INSTALL_DIR := $(shell go env GOPATH)/bin
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -s -w -X main.version=$(VERSION)

.PHONY: build install start rm release snapshot

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install:
	go build -ldflags "$(LDFLAGS)" -o $(INSTALL_DIR)/$(BINARY) .
	@if command -v $(BINARY) >/dev/null 2>&1; then \
		echo "$(BINARY) installed — try '$(BINARY) -h'"; \
	elif echo ":$$PATH:" | grep -q ":$$HOME/.local/bin:"; then \
		mkdir -p "$$HOME/.local/bin"; \
		ln -sf $(INSTALL_DIR)/$(BINARY) "$$HOME/.local/bin/$(BINARY)"; \
		echo "linked into ~/.local/bin — '$(BINARY) -h' works right now"; \
	else \
		case "$$SHELL" in */zsh) rc="$$HOME/.zshrc" ;; *) rc="$$HOME/.bashrc" ;; esac; \
		if ! grep -qs 'go/bin' "$$rc"; then \
			printf '\nexport PATH="%s:$$PATH"\n' '$(INSTALL_DIR)' >> "$$rc"; \
			echo "added $(INSTALL_DIR) to PATH in $$rc"; \
		fi; \
		echo "run 'source $$rc' (or open a new terminal), then try '$(BINARY) -h'"; \
	fi

start:
	air

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean

rm:
	rm -f ~/.peektea.toml && echo "Removed ~/.peektea.toml"
