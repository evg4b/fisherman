# ===========================================
# Default
# ===========================================
.PHONY: default
default: lint unit-test build

# ===========================================
# Rust Tasks
# ===========================================
.PHONY: lint
lint:
	@echo "Running Clippy..."
	@cargo clippy --all-targets --all-features --fix --allow-dirty -- -D warnings

.PHONY: unit-test
unit-test:
	@echo "Running unit tests..."
	@cargo test -p fisherman_core

.PHONY: unit-test-coverage
unit-test-coverage:
	@echo "Generating code coverage..."
	@cargo llvm-cov -p fisherman_core --no-fail-fast

.PHONY: e2e-test
e2e-test:
	@echo "Running e2e tests..."
	@cargo test -p fisherman --no-fail-fast

.PHONY: build
build:
	@echo "Building release..."
	@cargo build --release

.PHONY: install
install:
	@echo "Installing crate..."
	@cargo install --path .
