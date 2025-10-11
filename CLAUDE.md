# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Fisherman is a Git hook management tool written in Rust that allows developers to configure and manage Git hooks through declarative configuration files (.fisherman.toml/.yaml/.json). It supports hierarchical configuration (global, repository, and local scopes), template variables extracted from Git context, conditional rule execution, and parallel execution of asynchronous rules.

## Development Commands

### Building and Testing
```bash
# Default: lint, test, and build
make

# Run tests
cargo test
make test

# Run a single test (by name pattern)
cargo test test_name

# Run tests with coverage report
cargo llvm-cov --open
make coverage

# Lint with Clippy (auto-fix)
cargo clippy --all-targets --all-features --fix --allow-dirty
make lint

# Build release binary
cargo build --release
make build

# Install from local source
cargo install --path .
make install
```

### Running the CLI
```bash
# Install hooks in current repository
fisherman install [--force]

# Handle a specific hook (called by Git hooks)
fisherman handle <hook-name>

# Explain configured rules for a hook
fisherman explain <hook-name>
```

## Architecture

### Core Components

**Context System** (`src/context/`)
- `Context` trait provides abstraction for Git repository operations and configuration access
- `GitRepoContext` implements the context using git2 library
- Provides access to: repository path, hooks directory, configuration, binary location, and extracted variables
- Uses `MockContext` for testing (mockall)

**Configuration Loading** (`src/configuration/`)
- Hierarchical loading: global (~/.fisherman.toml) → repository (.fisherman.toml) → local (.git/.fisherman.toml)
- Supports TOML, YAML, and JSON formats using figment
- Configurations are merged by concatenating rules per hook (not overriding)
- Structure: top-level `extract` array + `hooks.<hook-name>` arrays containing rule definitions

**Rules System** (`src/rules/`)
- Rules are defined in config and compiled into `CompiledRule` trait objects at runtime
- `Rule::compile()` processes:
  1. Variable extraction from Git context (branch name, repo path)
  2. Conditional evaluation (`when` expression using Rhai scripting)
  3. Template rendering for rule parameters (using simple {{var}} syntax)
- Two execution modes:
  - **Synchronous**: Validation rules (message-regex, branch-name-regex, etc.) run sequentially
  - **Asynchronous**: External execution rules (exec, shell, write-file) run in parallel using rayon
- Available rule types:
  - Commit message: `message-regex`, `message-prefix`, `message-suffix`
  - Branch name: `branch-name-regex`, `branch-name-prefix`, `branch-name-suffix`
  - Execution: `exec` (run commands), `shell` (run shell scripts)
  - File operations: `write-file`

**Templates & Variables** (`src/templates/`)
- Simple template engine supporting `{{variable}}` syntax
- Variables extracted from Git context using regex with named capture groups
- Extraction patterns: `branch:regex`, `branch?:regex` (optional), `repo_path:regex`, `repo_path?:regex`
- Variables are resolved during rule compilation and injected into rule parameters

**Scripting** (`src/scripting/`)
- Conditional execution using Rhai scripting language
- `Expression` wrapper evaluates `when` conditions with variable context
- Built-in functions: `is_def_var("name")` checks if variable is defined

**Hooks** (`src/hooks/`)
- `GitHook` enum represents all Git hook types (pre-commit, commit-msg, pre-push, etc.)
- Hook installation: creates shell script in `.git/hooks/` that calls `fisherman handle <hook>`
- For commit-msg hook, passes arguments ($@) to fisherman
- Supports --force flag to backup and overwrite existing hooks

**Commands** (`src/commands/`)
- Three CLI commands:
  - `install`: Install hooks into repository
  - `handle`: Execute configured rules for a hook (called by installed Git hooks)
  - `explain`: Display configured rules and their order for debugging

### Data Flow

1. **Installation**: `fisherman install` → loads config → creates hook scripts in `.git/hooks/`
2. **Hook Execution**: Git triggers hook → hook script calls `fisherman handle <hook>` → loads config → extracts variables → compiles rules (evaluates `when`, renders templates) → executes synchronous rules sequentially → executes asynchronous rules in parallel → reports results
3. **Variable Extraction**: Git context (branch/repo path) → regex matching → named groups → variable map → available to templates and conditions

### Key Design Patterns

- **Trait-based abstraction**: `Context` trait enables testing and modularity
- **Visitor pattern**: `CompiledRule` trait allows polymorphic rule execution
- **Builder/Factory pattern**: Rule compilation transforms configuration into executable rule objects
- **Strategy pattern**: Different rule types implement the same `CompiledRule` interface

## Important Implementation Details

### Parallel Execution
- Async rules (exec, shell, write-file) execute concurrently using rayon's parallel iterators
- Designed to reduce total execution time when multiple commands/scripts are configured
- Ensure async rules are independent and don't conflict (e.g., writing to the same file)

### Template Rendering
- Uses `t!()` macro (defined in templates module) for template rendering
- Templates are rendered during rule compilation, not execution
- Failed template rendering causes rule compilation to fail

### Error Handling
- Uses anyhow for flexible error propagation
- Custom error types in `src/hooks/errors.rs` for specific failure cases
- Rule execution failures stop hook execution and abort the Git operation

### Testing Strategy
- Unit tests use MockContext to isolate components
- Integration tests use tempdir for filesystem operations
- rstest for parameterized testing (e.g., testing all hook types)
- mockall for mocking the Context trait

## Common Patterns

### Adding a New Rule Type
1. Create rule implementation in `src/rules/<rule_name>.rs` implementing `CompiledRule`
2. Add enum variant to `RuleParams` in `src/rules/rule_def.rs` with serde rename
3. Add compilation case in `Rule::compile()` match statement
4. Add display case in `RuleParams::name()` for debugging
5. Write tests in the rule module and in `rule_def.rs`

### Configuration Merging
- Configurations from different scopes are concatenated, not merged/overridden
- Each scope adds its rules to the end of the hook's rule list
- Execution order: global rules → repository rules → local rules

### Variable Extraction
- Global `extract` in config applies to all rules unless overridden
- Per-rule `extract` overrides global extraction for that specific rule
- Use `branch?:` or `repo_path?:` for optional matching (won't fail if regex doesn't match)
- Use `is_def_var("VarName")` in `when` conditions before using variables in templates
