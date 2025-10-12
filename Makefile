default: lint test build

lint:
	cargo clippy --all-features --fix --allow-dirty

test:
	cargo test

test-unit:
	cargo test --lib

test-integration:
	cargo build --release
	cargo test --test '*' --release -- --test-threads=1

test-all: test-unit test-integration

build:
	cargo build --release

install:
	cargo install --path .

coverage:
	cargo llvm-cov --open

coverage-integration:
	cargo llvm-cov --release --test '*' -- --test-threads=1