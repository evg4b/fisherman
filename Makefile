# ===========================================
# Variables
# ===========================================
CARGO = cargo
DOCKER = docker
PWD = $(shell pwd)
DOCKER_IMAGE_PREFIX = evg4b/rust
DOCKERFILE_DIR = docker

# ===========================================
# Default
# ===========================================
.PHONY: default
default: lint test build

# ===========================================
# Rust Tasks
# ===========================================
.PHONY: lint
lint:
	@echo "Running Clippy..."
	$(CARGO) clippy --all-targets --all-features --fix --allow-dirty

.PHONY: test
test:
	@echo "Running tests..."
	$(CARGO) test

.PHONY: build
build:
	@echo "Building release..."
	$(CARGO) build --release

.PHONY: install
install:
	@echo "Installing crate..."
	$(CARGO) install --path .

.PHONY: coverage
coverage:
	@echo "Generating code coverage..."
	$(CARGO) llvm-cov --open

# ===========================================
# Docker contauner helpres
# ===========================================
define docker_build_run
	@IMAGE=$(DOCKER_IMAGE_PREFIX)-$(1) ; \
	echo "Building Docker image $$IMAGE..." ; \
	$(DOCKER) build --platform linux/amd64 -f $(DOCKERFILE_DIR)/Dockerfile.$(1) -t $$IMAGE . ; \
	echo "Running build for target $(1)..." ; \
	$(DOCKER) run --platform linux/amd64 -it --rm -v "$(PWD)":/app -w /app $$IMAGE $(2)
endef

# ===========================================
# Clen all targets
# ===========================================

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf target

# ===========================================
# aarch64-unknown-linux-gnu
# ===========================================
.PHONY: debug-aarch64-unknown-linux-gnu
debug-aarch64-unknown-linux-gnu:
	$(call docker_build_run,aarch64-unknown-linux-gnu,$(CARGO) build --target aarch64-unknown-linux-gnu)

.PHONY: debug-aarch64-unknown-linux-gnu
release-aarch64-unknown-linux-gnu:
	$(call docker_build_run,aarch64-unknown-linux-gnu,$(CARGO) build --target aarch64-unknown-linux-gnu --release)

.PHONY: container-aarch64-unknown-linux-gnu
container-aarch64-unknown-linux-gnu:
	@echo "Starting interactive container..."
	$(call docker_build_run,aarch64-unknown-linux-gnu,bash)

# ===========================================
# x86_64-pc-windows-gnu
# ===========================================
.PHONY: debug-x86_64-pc-windows-gnu
debug-x86_64-pc-windows-gnu:
	$(call docker_build_run,x86_64-pc-windows-gnu,$(CARGO) build --target x86_64-pc-windows-gnu)

.PHONY: debug-x86_64-pc-windows-gnu
release-x86_64-pc-windows-gnu:
	$(call docker_build_run,x86_64-pc-windows-gnu,$(CARGO) build --target x86_64-pc-windows-gnu --release)

.PHONY: container-x86_64-pc-windows-gnu
container-x86_64-pc-windows-gnu:
	@echo "Starting interactive container..."
	$(call docker_build_run,x86_64-pc-windows-gnu,bash)


