# Makefile for Touch Typing Master - Go Learning Game
# Author: Stanley Chege Thuita
# Email: altradits@gmail.com
# LinkedIn: https://www.linkedin.com/in/stanmobitech/
# Phone: +254707172370
# Description: Build, run, and package the touch typing game

# ============================================
# Variables
# ============================================

# Project information
PROJECT_NAME := touchtype-master
BINARY_NAME := $(PROJECT_NAME)
AUTHOR_NAME := "Stanley Chege Thuita"
AUTHOR_EMAIL := "altradits@gmail.com"
AUTHOR_LINKEDIN := "https://www.linkedin.com/in/stanmobitech/"
AUTHOR_PHONE := "+254707172370"

# Build variables
BUILD_DIR := ./build
DIST_DIR := ./dist
GO_VERSION := $(shell go version | cut -d' ' -f3 | sed 's/go//')
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")

# Version variables
VERSION_MAJOR := 1
VERSION_MINOR := 0
VERSION_PATCH := 0
VERSION := $(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_PATCH)

# Build metadata. Author values are stable defaults in main.go; keeping values
# without spaces here makes the flags portable across shells and platforms.
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

# Platforms to build for
PLATFORMS := windows linux darwin
ARCHITECTURES := amd64 arm64

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOTIDY := $(GOMOD) tidy
GOVET := $(GOCMD) vet

# Colors for output
COLOR_RESET := \033[0m
COLOR_RED := \033[31m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_CYAN := \033[36m
COLOR_BOLD := \033[1m
COLOR_MAGENTA := \033[35m

# ============================================
# Default target
# ============================================

.PHONY: all
all: help

# ============================================
# Build targets
# ============================================

.PHONY: build
build: ## Build the application for current platform
	@echo "$(COLOR_CYAN)Building $(PROJECT_NAME) v$(VERSION)...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(COLOR_GREEN)✓ Build complete!$(COLOR_RESET)"
	@echo "  Binary: $(BUILD_DIR)/$(BINARY_NAME)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build-all
build-all: ## Build for all supported platforms
	@echo "$(COLOR_CYAN)Building for all platforms...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)
	@$(MAKE) build-windows
	@$(MAKE) build-linux
	@$(MAKE) build-darwin
	@$(MAKE) build-darwin-arm64
	@echo "$(COLOR_GREEN)✓ All builds complete!$(COLOR_RESET)"
	@ls -lh $(DIST_DIR)/

.PHONY: build-windows
build-windows: ## Build for Windows (amd64)
	@echo "$(COLOR_YELLOW)Building for Windows...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/windows/$(BINARY_NAME).exe .
	@echo "$(COLOR_GREEN)✓ Windows build complete: $(DIST_DIR)/windows/$(BINARY_NAME).exe$(COLOR_RESET)"

.PHONY: build-linux
build-linux: ## Build for Linux (amd64)
	@echo "$(COLOR_YELLOW)Building for Linux...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/linux/$(BINARY_NAME) .
	@echo "$(COLOR_GREEN)✓ Linux build complete: $(DIST_DIR)/linux/$(BINARY_NAME)$(COLOR_RESET)"

.PHONY: build-darwin
build-darwin: ## Build for macOS (amd64)
	@echo "$(COLOR_YELLOW)Building for macOS...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/darwin
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/darwin/$(BINARY_NAME)-mac .
	@echo "$(COLOR_GREEN)✓ macOS build complete: $(DIST_DIR)/darwin/$(BINARY_NAME)-mac$(COLOR_RESET)"

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS ARM64 (Apple Silicon)
	@echo "$(COLOR_YELLOW)Building for macOS ARM64...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/darwin-arm64
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/darwin-arm64/$(BINARY_NAME)-mac-arm64 .
	@echo "$(COLOR_GREEN)✓ macOS ARM64 build complete: $(DIST_DIR)/darwin-arm64/$(BINARY_NAME)-mac-arm64$(COLOR_RESET)"

.PHONY: build-docker
build-docker: ## Build Docker image
	@echo "$(COLOR_CYAN)Building Docker image...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	docker build --build-arg AUTHOR_NAME="$(AUTHOR_NAME)" \
	             --build-arg AUTHOR_EMAIL="$(AUTHOR_EMAIL)" \
	             --build-arg AUTHOR_LINKEDIN="$(AUTHOR_LINKEDIN)" \
	             --build-arg AUTHOR_PHONE="$(AUTHOR_PHONE)" \
	             -t $(PROJECT_NAME):$(VERSION) -t $(PROJECT_NAME):latest .
	@echo "$(COLOR_GREEN)✓ Docker image built: $(PROJECT_NAME):$(VERSION)$(COLOR_RESET)"

# ============================================
# Development targets
# ============================================

.PHONY: run
run: ## Run the application
	@echo "$(COLOR_CYAN)Running $(PROJECT_NAME)...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	$(GOCMD) run .

.PHONY: dev
dev: ## Run with hot-reload (requires air)
	@echo "$(COLOR_CYAN)Running in development mode...$(COLOR_RESET)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(COLOR_RED)❌ 'air' not found. Install with: go install github.com/air-verse/air@latest$(COLOR_RESET)"; \
		$(GOCMD) run .; \
	fi

.PHONY: test
test: ## Run tests
	@echo "$(COLOR_CYAN)Running tests...$(COLOR_RESET)"
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(COLOR_CYAN)Running tests with coverage...$(COLOR_RESET)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report: coverage.html$(COLOR_RESET)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(COLOR_CYAN)Running go vet...$(COLOR_RESET)"
	$(GOVET) ./...

.PHONY: fmt
fmt: ## Format code
	@echo "$(COLOR_CYAN)Formatting code...$(COLOR_RESET)"
	$(GOCMD) fmt ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "$(COLOR_CYAN)Tidying go modules...$(COLOR_RESET)"
	$(GOTIDY)

.PHONY: deps
deps: tidy ## Install dependencies
	@echo "$(COLOR_CYAN)Installing dependencies...$(COLOR_RESET)"
	$(GOCMD) mod download

# ============================================
# Clean targets
# ============================================

.PHONY: clean
clean: ## Clean build files
	@echo "$(COLOR_CYAN)Cleaning build files...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	$(GOCLEAN)
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

.PHONY: clean-all
clean-all: clean ## Clean everything including dependencies
	@echo "$(COLOR_CYAN)Cleaning everything...$(COLOR_RESET)"
	$(GOCLEAN) -cache -modcache -testcache
	@echo "$(COLOR_GREEN)✓ Complete clean$(COLOR_RESET)"

# ============================================
# Package targets
# ============================================

.PHONY: package
package: build-all ## Create release packages
	@echo "$(COLOR_CYAN)Creating release packages...$(COLOR_RESET)"
	@mkdir -p $(DIST_DIR)/packages
	@$(MAKE) package-windows
	@$(MAKE) package-linux
	@$(MAKE) package-darwin
	@$(MAKE) package-darwin-arm64
	@echo "$(COLOR_GREEN)✓ Packages created in $(DIST_DIR)/packages/$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Packaged by: $(AUTHOR_NAME)$(COLOR_RESET)"

.PHONY: package-windows
package-windows:
	@echo "$(COLOR_YELLOW)Packaging Windows...$(COLOR_RESET)"
	cd $(DIST_DIR)/windows && zip -r ../packages/$(PROJECT_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME).exe
	@echo "$(COLOR_GREEN)✓ Windows package created$(COLOR_RESET)"

.PHONY: package-linux
package-linux:
	@echo "$(COLOR_YELLOW)Packaging Linux...$(COLOR_RESET)"
	cd $(DIST_DIR)/linux && tar -czf ../packages/$(PROJECT_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)
	@echo "$(COLOR_GREEN)✓ Linux package created$(COLOR_RESET)"

.PHONY: package-darwin
package-darwin:
	@echo "$(COLOR_YELLOW)Packaging macOS...$(COLOR_RESET)"
	cd $(DIST_DIR)/darwin && tar -czf ../packages/$(PROJECT_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-mac
	@echo "$(COLOR_GREEN)✓ macOS package created$(COLOR_RESET)"

.PHONY: package-darwin-arm64
package-darwin-arm64:
	@echo "$(COLOR_YELLOW)Packaging macOS ARM64...$(COLOR_RESET)"
	cd $(DIST_DIR)/darwin-arm64 && tar -czf ../packages/$(PROJECT_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-mac-arm64
	@echo "$(COLOR_GREEN)✓ macOS ARM64 package created$(COLOR_RESET)"

# ============================================
# Install targets
# ============================================

.PHONY: install
install: build ## Install the application locally
	@echo "$(COLOR_CYAN)Installing $(PROJECT_NAME)...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME) 2>/dev/null || \
	 (echo "$(COLOR_YELLOW)⚠️  Could not copy to /usr/local/bin. Installing to ~/bin...$(COLOR_RESET)" && \
	  mkdir -p ~/bin && cp $(BUILD_DIR)/$(BINARY_NAME) ~/bin/)
	@echo "$(COLOR_GREEN)✓ Installed!$(COLOR_RESET)"

# ============================================
# Docker targets
# ============================================

.PHONY: docker-build
docker-build: build-docker ## Alias for build-docker

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "$(COLOR_CYAN)Running Docker container...$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	docker run -p 8080:8080 $(PROJECT_NAME):latest

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	@echo "$(COLOR_CYAN)Stopping Docker container...$(COLOR_RESET)"
	docker stop $(shell docker ps -q --filter ancestor=$(PROJECT_NAME):latest) 2>/dev/null || true

# ============================================
# Utility targets
# ============================================

.PHONY: info
info: ## Display project information
	@echo "$(COLOR_BOLD)Project Information$(COLOR_RESET)"
	@echo "  $(COLOR_CYAN)Name:$(COLOR_RESET)          $(PROJECT_NAME)"
	@echo "  $(COLOR_CYAN)Version:$(COLOR_RESET)       v$(VERSION)"
	@echo "  $(COLOR_CYAN)Author:$(COLOR_RESET)        $(AUTHOR_NAME)"
	@echo "  $(COLOR_CYAN)Email:$(COLOR_RESET)         $(AUTHOR_EMAIL)"
	@echo "  $(COLOR_CYAN)LinkedIn:$(COLOR_RESET)      $(AUTHOR_LINKEDIN)"
	@echo "  $(COLOR_CYAN)Phone:$(COLOR_RESET)         $(AUTHOR_PHONE)"
	@echo "  $(COLOR_CYAN)Go Version:$(COLOR_RESET)    $(GO_VERSION)"
	@echo "  $(COLOR_CYAN)Build Time:$(COLOR_RESET)    $(BUILD_TIME)"
	@echo "  $(COLOR_CYAN)Git Commit:$(COLOR_RESET)    $(GIT_COMMIT)"
	@echo "  $(COLOR_CYAN)Git Branch:$(COLOR_RESET)    $(GIT_BRANCH)"

.PHONY: author
author: ## Display author information
	@echo "$(COLOR_BOLD)Author Information$(COLOR_RESET)"
	@echo "  $(COLOR_CYAN)Name:$(COLOR_RESET)      $(AUTHOR_NAME)"
	@echo "  $(COLOR_CYAN)Email:$(COLOR_RESET)     $(AUTHOR_EMAIL)"
	@echo "  $(COLOR_CYAN)LinkedIn:$(COLOR_RESET)  $(AUTHOR_LINKEDIN)"
	@echo "  $(COLOR_CYAN)Phone:$(COLOR_RESET)     $(AUTHOR_PHONE)"

.PHONY: size
size: ## Show binary sizes
	@echo "$(COLOR_BOLD)Binary Sizes$(COLOR_RESET)"
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "  $(COLOR_CYAN)Current build:$(COLOR_RESET) $$(du -h $(BUILD_DIR)/$(BINARY_NAME) | cut -f1)"; \
	fi
	@if [ -d $(DIST_DIR) ]; then \
		echo "  $(COLOR_CYAN)Distributed builds:$(COLOR_RESET)"; \
		find $(DIST_DIR) -name "$(BINARY_NAME)*" -exec du -h {} \; 2>/dev/null | sed 's/^/    /'; \
	fi

.PHONY: help
help: ## Show this help message
	@echo "$(COLOR_BOLD)$(PROJECT_NAME) - Touch Typing Master Game$(COLOR_RESET)"
	@echo "$(COLOR_CYAN)Version: v$(VERSION)$(COLOR_RESET)"
	@echo "$(COLOR_MAGENTA)Author: $(AUTHOR_NAME)$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Usage:$(COLOR_RESET)"
	@echo "  make [target]"
	@echo ""
	@echo "$(COLOR_BOLD)Targets:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(COLOR_BOLD)Examples:$(COLOR_RESET)"
	@echo "  $(COLOR_CYAN)make build$(COLOR_RESET)        - Build for current platform"
	@echo "  $(COLOR_CYAN)make run$(COLOR_RESET)          - Run the game"
	@echo "  $(COLOR_CYAN)make build-all$(COLOR_RESET)    - Build for all platforms"
	@echo "  $(COLOR_CYAN)make package$(COLOR_RESET)      - Create release packages"
	@echo "  $(COLOR_CYAN)make clean$(COLOR_RESET)        - Clean build files"
	@echo "  $(COLOR_CYAN)make author$(COLOR_RESET)       - Show author information"
	@echo "  $(COLOR_CYAN)make info$(COLOR_RESET)         - Show full project information"
	@echo "  $(COLOR_CYAN)make help$(COLOR_RESET)         - Show this help"

# ============================================
# Default target
# ============================================

.DEFAULT_GOAL := help

# ============================================
# Phony targets (always run)
# ============================================

.PHONY: help build build-all build-windows build-linux build-darwin build-darwin-arm64 \
        build-docker run dev test test-coverage vet fmt tidy deps clean clean-all \
        package package-windows package-linux package-darwin package-darwin-arm64 \
        install docker-build docker-run docker-stop info size author

# ============================================
# Include additional makefiles if they exist
# ============================================

-include .make/*.mk
