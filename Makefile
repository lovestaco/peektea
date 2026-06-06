BINARY      := peektea
INSTALL_DIR := $(shell go env GOPATH)/bin
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -s -w -X main.version=$(VERSION)

.PHONY: build install start rm bump release-notes-init release snapshot

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

# Create a git tag and push it.
# Usage: make bump VERSION=v0.2.0
bump:
	@if [ -z "$(VERSION_ARG)" ]; then \
		read -p "version (e.g. v0.2.0): " v; \
	else \
		v=$(VERSION_ARG); \
	fi; \
	git tag -a $$v -m "release $$v" && \
	git push origin $$v && \
	echo "tagged and pushed $$v"

# Scaffold a new release notes file from the template.
# Usage: make release-notes-init VERSION=v0.2.0
release-notes-init:
	@if [ -z "$(VERSION)" ]; then echo "usage: make release-notes-init VERSION=v0.x.x"; exit 1; fi
	@if [ -f "docs/releases/$(VERSION).md" ]; then echo "docs/releases/$(VERSION).md already exists"; exit 1; fi
	@DATE=$$(date -u +%Y-%m-%d); \
	sed \
		-e "s/__VERSION__/$(VERSION)/g" \
		-e "s/__DATE__/$$DATE/g" \
		docs/releases/_template.md > docs/releases/$(VERSION).md
	@echo "created docs/releases/$(VERSION).md — fill it in and run 'make release VERSION=$(VERSION)'"

# Preview rendered release notes without publishing.
# Usage: make release-notes-preview VERSION=v0.1.0
release-notes-preview:
	@if [ -z "$(VERSION)" ]; then echo "usage: make release-notes-preview VERSION=v0.x.x"; exit 1; fi
	python3 scripts/release_gh.py $(VERSION) --render

# Build release archives locally without publishing (no tag needed).
snapshot:
	goreleaser release --snapshot --clean

# Build and publish a GitHub release.
# Usage: make release VERSION=v0.1.0
release:
	@if [ -z "$(VERSION)" ]; then echo "usage: make release VERSION=v0.x.x"; exit 1; fi
	python3 scripts/release_gh.py $(VERSION) --render > /tmp/peektea_release_notes.md
	goreleaser release --clean --release-notes /tmp/peektea_release_notes.md

rm:
	rm -f ~/.peektea.toml && echo "Removed ~/.peektea.toml"
