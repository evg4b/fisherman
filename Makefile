default: lint test build

lint:
	cargo clippy --fix --allow-dirty --allow-staged

test:
	cargo test

build:
	cargo build --release