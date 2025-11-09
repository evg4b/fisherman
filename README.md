<!--suppress HtmlDeprecatedAttribute -->
<p align="center">
  <a href="https://github.com/evg4b/fisherman" title="fisherman">
    <img alt="fisherman logo" width="80%" src="./.github/logo.svg">
  </a>
</p>
<p align="center">
  Small git hook management tool for developer.
</p>
<hr>
<div align="center">
    <a href="https://github.com/evg4b/fisherman/actions/workflows/rust.yml">
        <img alt="GitHub Actions Workflow Status" src="https://img.shields.io/github/actions/workflow/status/evg4b/fisherman/rust.yml?branch=master&label=Build">
    </a>
    <a href="https://github.com/evg4b/fisherman/blob/master/LICENSE">
        <img alt="GitHub License" src="https://img.shields.io/github/license/evg4b/fisherman?label=License">
    </a>
    <a href="https://codecov.io/gh/evg4b/fisherman">
        <img alt="Codecov" src="https://img.shields.io/codecov/c/github/evg4b/fisherman?label=Coverage"> 
    </a>
</div>

> **Note:** Fisherman is still in development, so you use it at your own risk.

## Install using Cargo

You can install Fisherman via [Cargo](https://doc.rust-lang.org/cargo/#the-cargo-book):

```bash
cargo install --git https://github.com/evg4b/fisherman.git
```

## Install Precompiled Binaries

1.	Go to the [Build](https://github.com/evg4b/fisherman/actions/workflows/build.yml?query=branch%3Amaster+is%3Asuccess) workflow.
2.	Select the latest successful run.
3.	Scroll down to the Artifacts section.
4.	Download the binaries matching your system and architecture.