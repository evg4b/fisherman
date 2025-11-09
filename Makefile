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
	@cargo clippy --all-targets --all-features --fix --allow-dirty

.PHONY: test
test:
	@echo "Running tests..."
	@cargo test

.PHONY: build
build:
	@echo "Building release..."
	@cargo build --release

.PHONY: install
install:
	@echo "Installing crate..."
	@cargo install --path .

.PHONY: coverage
coverage:
	@echo "Generating code coverage..."
	@cargo llvm-cov --open
