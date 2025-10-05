default: lint test build

lint:
	cargo clippy --all-targets --all-features --fix --allow-dirty

test:
	cargo test

build:
	cargo build --release