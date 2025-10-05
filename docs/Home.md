<!--suppress HtmlDeprecatedAttribute -->
<p align="center">
  <a href="https://github.com/evg4b/fisherman" title="fisherman">
    <img alt="fisherman logo" width="80%" src="https://raw.githubusercontent.com/evg4b/fisherman/refs/heads/master/.github/logo.svg">
  </a>
</p>
<p align="center">
  Small git hook management tool for developer.
</p>

## Overview

Fisherman is a lightweight, declarative Git hook manager that simplifies automation of your local Git workflows.
While many hook managers focus on basic tasks like running tests or linting code,
Fisherman enables more sophisticated hook configurations with minimal effort.

## Key Features

- ✅ **Declarative Configuration** - Define hooks using TOML, YAML, or JSON files
- ✅ **Context-Aware Rules** - Create hooks that respond to branch names, file paths, and more
- ✅ **Variable Extraction** - Pull information from your environment using regex patterns
- ✅ **Multiple Hook Types** - Support for all standard Git hooks (pre-commit, commit-msg, pre-push, etc.)
- ✅ **Conditional Execution** - Execute rules based on expressions using the `when` parameter
- ✅ **Template Support** - Use variables in your configurations with `{{variable}}` syntax
- ✅ **Multiple Scopes** - Global, repository, and local configuration support
- ✅ **Parallel Execution** - Async rules (exec, shell, write-file) run in parallel for faster hook execution
- ✅ **Easy Installation** - Single binary, no runtime dependencies

## Quick Start

### Installation

Install Fisherman using Cargo:

```bash
cargo install --git https://github.com/evg4b/fisherman.git
```

### Basic Usage

1. Create a configuration file in your repository (`.fisherman.toml`, `.fisherman.yaml`, or `.fisherman.json`)
2. Define your hooks and rules
3. Install the hooks with `fisherman install`

### Example Configuration

```toml
# .fisherman.toml

# Run tests before committing
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

# Validate commit message format
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"
```

### Install Hooks

```bash
# Install configured hooks
fisherman install

# Install specific hooks
fisherman install pre-commit pre-push

# Force install (override existing hooks)
fisherman install --force
```

## Documentation

- [Installation](Installation.md) - How to install Fisherman
- [CLI Commands](CLI-Commands.md) - Command-line interface reference
- [Configuration](Configuration.md) - Configuration file format and scopes
- [Rules Reference](Rules.md) - Complete list of available rules
- [Variables and Templates](Variables-and-Templates.md) - How to use variables and templates
- [Examples](Examples-of-usage.md) - Common use cases and examples
- [Git Hooks](Git-Hooks.md) - Supported Git hooks reference

## Why Fisherman?

Fisherman was designed to be:

- **Simple** - Easy-to-read configuration without complex scripting
- **Flexible** - Supports multiple rule types and conditional execution
- **Consistent** - Same configuration format across all hooks
- **Portable** - Single binary written in Rust, no runtime dependencies
- **Global** - Define rules once in your home directory for all repositories

## Contributing

Contributions are welcome! Please check the [GitHub repository](https://github.com/evg4b/fisherman) for more
information.

## License

MIT License - see the [LICENSE](https://github.com/evg4b/fisherman/blob/master/LICENSE) file for details.