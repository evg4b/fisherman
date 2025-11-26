# Rules Reference

Fisherman provides a comprehensive set of rules for validating and automating Git hooks. Rules can be combined, use
variables, and execute conditionally.

## Common Rule Fields

All rules support these common fields:

- **`type`** (required) - The rule type identifier
- **`when`** (optional) - Conditional expression for executing the rule (
  see [Variables and Templates](Variables-and-Templates.md))
- **`extract`** (optional) - Override global variable extraction for this rule

## Execution Modes

Rules are executed in two modes:

- **Synchronous** - Rules that validate data (e.g., commit message validation, branch name checks). These execute
  quickly and block the Git operation. Synchronous rules run sequentially in the order they are defined.
- **Asynchronous** - Rules that execute external commands (e.g., `exec`, `shell`, `write-file`). These may take longer and run
  external processes. **Asynchronous rules execute in parallel** to improve performance when multiple async rules are configured.

**Parallel Execution Benefits:**

When multiple asynchronous rules are defined (e.g., running tests, linting, and building), Fisherman executes them
concurrently using parallel threads. This significantly reduces total execution time compared to running each rule
sequentially.

**Example:**
```toml
# These three async rules will run in parallel
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["clippy"]

[[hooks.pre-commit]]
type = "shell"
script = "cargo fmt --check"
```

If each rule takes 5 seconds, parallel execution completes in ~5 seconds instead of 15 seconds sequentially.

---

## Commit Message Rules

Rules for validating and enforcing commit message formats. These rules are typically used with the `commit-msg` hook.

### `message-regex`

Validates that the commit message matches a regular expression pattern.

**Type:** Synchronous

**Parameters:**

- `regex` (required) - Regular expression pattern to match against the commit message

**Example - Conventional Commits:**

```toml
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|perf|test|chore)(\\(.+\\))?:\\s.+"
```

**Example - Require Issue Number:**

```toml
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(PROJ-\\d+|NO-ISSUE):\\s.+"
```

**Use cases:**

- Enforce conventional commit format
- Require issue tracking numbers
- Prevent certain words or patterns
- Validate commit message structure

---

### `message-prefix`

Validates that the commit message starts with a specific prefix.

**Type:** Synchronous

**Parameters:**

- `prefix` (required) - String that the commit message must start with

**Example - Static Prefix:**

```toml
[[hooks.commit-msg]]
type = "message-prefix"
prefix = "feat: "
```

**Example - Dynamic Prefix with Variables:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

**Use cases:**

- Enforce consistent commit message prefixes
- Automatically add issue numbers from branch names
- Require specific commit types

---

### `message-suffix`

Validates that the commit message ends with a specific suffix.

**Type:** Synchronous

**Parameters:**

- `suffix` (required) - String that the commit message must end with

**Example:**

```toml
[[hooks.commit-msg]]
type = "message-suffix"
suffix = " [skip ci]"
```

**Example - Dynamic Suffix:**

```toml
extract = ["branch?:^(?<Username>[^/]+)/.*$"]

[[hooks.commit-msg]]
type = "message-suffix"
when = "is_def_var(\"Username\")"
suffix = " (by {{Username}})"
```

**Use cases:**

- Enforce required trailers or tags
- Add automated signatures
- Include metadata in commits

---

## Branch Name Rules

Rules for validating branch naming conventions. These rules are typically used with the `pre-push` or `pre-commit`
hooks.

### `branch-name-regex`

Validates that the current branch name matches a regular expression pattern.

**Type:** Synchronous

**Parameters:**

- `regex` (required) - Regular expression pattern to match against the branch name

**Example - Git Flow:**

```toml
[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix|release)/[a-z0-9-]+$"
```

<details open>
  <summary>Click to collapse</summary>
  This content will be visible by default.
</details>

**Example - Require Issue Prefix:**

```toml
[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(PROJ-\\d+|main|develop)-.*$"
```

**Use cases:**

- Enforce branch naming conventions (Git Flow, GitHub Flow)
- Require issue tracking numbers in branch names
- Prevent certain branch name patterns
- Ensure lowercase or specific character sets

---

### `branch-name-prefix`

Validates that the current branch name starts with a specific prefix.

**Type:** Synchronous

**Parameters:**

- `prefix` (required) - String that the branch name must start with

**Example:**

```toml
[[hooks.pre-push]]
type = "branch-name-prefix"
prefix = "feature/"
```

**Use cases:**

- Enforce branch type prefixes (feature/, bugfix/, etc.)
- Require team or project prefixes
- Implement simple branch naming policies

---

### `branch-name-suffix`

Validates that the current branch name ends with a specific suffix.

**Type:** Synchronous

**Parameters:**

- `suffix` (required) - String that the branch name must end with

**Example:**

```toml
[[hooks.pre-push]]
type = "branch-name-suffix"
suffix = "-dev"
```

**Use cases:**

- Enforce environment suffixes
- Require specific branch markers
- Distinguish branch types by suffix

---

## Execution Rules

Rules for executing external commands and scripts.

### `exec`

Executes an external command with arguments and environment variables.

**Type:** Asynchronous

**Parameters:**

- `command` (required) - Command to execute
- `args` (optional) - Array of command arguments
- `env` (optional) - Map of environment variables

**Example - Run Tests:**

```toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]
```

**Example - With Environment Variables:**

```toml
[[hooks.pre-commit]]
type = "exec"
command = "npm"
args = ["run", "test"]
env = { NODE_ENV = "test", CI = "true" }
```

**Example - With Template Variables:**

```toml
extract = ["branch?:^(?<Feature>[^/]+)/.*$"]

[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Feature\")"
command = "echo"
args = ["Running tests for {{Feature}}"]
```

**Success/Failure:**

- **Success:** Command exits with code 0
- **Failure:** Command exits with non-zero code

**Use cases:**

- Run tests before committing
- Execute linters and formatters
- Build projects before pushing
- Run custom validation scripts
- Integrate with external tools

---

### `shell`

Executes a shell script with full shell features (pipes, redirections, variables, etc.).

**Type:** Asynchronous

**Parameters:**

- `script` (required) - Shell script to execute
- `env` (optional) - Map of environment variables

**Example - Simple Script:**

```toml
[[hooks.pre-commit]]
type = "shell"
script = "cargo test && cargo clippy"
```

**Example - Multi-line Script:**

```toml
[[hooks.pre-commit]]
type = "shell"
script = """
#!/bin/bash
set -e
echo "Running pre-commit checks..."
cargo fmt --check
cargo test
cargo clippy -- -D warnings
echo "All checks passed!"
"""
```

**Example - With Variables:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.pre-push]]
type = "shell"
when = "is_def_var(\"IssueNumber\")"
script = """
echo "Validating issue {{IssueNumber}}..."
curl -f "https://api.example.com/issues/{{IssueNumber}}" || exit 1
"""
```

**Example - With Environment Variables:**

```toml
[[hooks.pre-commit]]
type = "shell"
env = { DATABASE_URL = "sqlite::memory:", LOG_LEVEL = "debug" }
script = "cargo test --features integration"
```

**Success/Failure:**

- **Success:** Script exits with code 0
- **Failure:** Script exits with non-zero code

**Use cases:**

- Complex validation logic with multiple commands
- Scripts requiring shell features (pipes, conditionals)
- Integration with external services
- Multi-step validation processes
- Custom business logic

---

## File System Rules

Rules for file system operations.

### `write-file`

Writes content to a file, optionally appending to existing content.

**Type:** Asynchronous

**Parameters:**

- `path` (required) - File path to write (relative to repository root)
- `content` (required) - Content to write to the file
- `append` (optional) - If `true`, append to existing file; if `false` or omitted, overwrite the file

**Example - Create Commit Template:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.post-checkout]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = "{{IssueNumber}}: "
```

**Example - Clear Commit Template:**

```toml
[[hooks.post-checkout]]
type = "write-file"
when = "!is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = ""
```

**Example - Append to Log File:**

```toml
[[hooks.post-commit]]
type = "write-file"
path = ".git/commit-log.txt"
content = "Committed at {{timestamp}}\n"
append = true
```

**Success/Failure:**

- **Success:** File is written successfully
- **Failure:** Cannot write file (permissions, invalid path, etc.)

**Use cases:**

- Update commit message templates dynamically
- Log hook executions
- Generate configuration files
- Update build metadata
- Create or modify Git configuration

---

## Rule Combinations

Rules can be combined to create sophisticated workflows:

### Conditional Execution Based on Branch

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

# Require issue number in commit message if branch has one
[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "

# Run integration tests only for feature branches
[[hooks.pre-push]]
type = "exec"
when = "is_def_var(\"IssueNumber\")"
command = "cargo"
args = ["test", "--features", "integration"]
```

### Multiple Validation Steps

```toml
# Check branch name format
[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(feature|bugfix)/[a-z0-9-]+$"

# Check commit message format
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix):\\s.+"

# Run tests
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

# Check formatting
[[hooks.pre-commit]]
type = "shell"
script = "cargo fmt --check"
```

### Dynamic Configuration

```toml
extract = [
    "branch?:^(?<Team>[^/]+)/(?<Issue>[^/]+)/.*$",
    "repo_path?:^/home/(?<Username>[^/]+)/.*$"
]

# Different rules based on team
[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Team\") && Team == \"frontend\""
command = "npm"
args = ["run", "test:ui"]

[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Team\") && Team == \"backend\""
command = "cargo"
args = ["test"]
```

---

## Best Practices

1. **Start with simple rules** - Begin with basic validation before adding complex logic
2. **Use meaningful rule names** - Comments help explain complex configurations
3. **Test rules individually** - Use `fisherman handle <hook>` to test rules
4. **Handle failures gracefully** - Ensure error messages are clear and actionable
5. **Use conditional execution** - Avoid running unnecessary rules with `when` conditions
6. **Leverage variables** - Extract information once and reuse it across rules
7. **Keep scripts focused** - Each rule should have a single, clear responsibility
8. **Consider performance** - Minimize expensive operations in frequently-used hooks; leverage parallel execution for async rules
9. **Document complex regex** - Add comments explaining non-obvious patterns
10. **Use appropriate hooks** - Match rules to the right Git hook for best results
11. **Design for parallelism** - Ensure async rules are independent and can safely run concurrently
12. **Avoid shared state** - When using multiple async rules, avoid operations that conflict (e.g., writing to the same file)

---

## Troubleshooting

### Rule Not Executing

- Check that the hook is configured for the correct Git hook
- Verify the `when` condition evaluates to `true`
- Ensure variables are extracted correctly
- Use `fisherman explain <hook>` to see configured rules

### Command Failures

- Check command exists and is in PATH
- Verify arguments are correct
- Review environment variables
- Test command manually outside Fisherman

### Template Variables Not Working

- Ensure variables are extracted with `extract`
- Check variable names match regex named groups
- Verify variable is defined before using (use `when` condition)
- Use `{{variable}}` syntax, not `${variable}`

See [Examples](Examples-of-usage.md) for more real-world use cases.
