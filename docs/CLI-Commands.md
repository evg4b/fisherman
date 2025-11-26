Fisherman provides several commands to help you manage Git hooks. This document outlines the available commands and
their usage.

# Commands Overview

| Command   | Purpose                                    | Can Abort Operation           |
|-----------|--------------------------------------------|-------------------------------|
| `install` | Install configured hooks into `.git/hooks` | N/A                           |
| `explain` | Display hook configuration and rules       | N/A                           |
| `handle`  | Execute a hook (used internally by Git)    | Yes (for hooksthat can abort) |

---

# `install` - Install Hooks

Installs Git hooks into your repository's `.git/hooks` directory based on your configuration files.

## Syntax

```bash
fisherman install [OPTIONS] [HOOKS...]
```

## Options

- `-f`, `--force` - Overwrite existing hook scripts. Original scripts are backed up with `.bkp` extension
- `-h`, `--help` - Display help information

## Arguments

- `HOOKS` (optional) - Space-separated list of specific hooks to install (e.g., `pre-commit`, `commit-msg`)

If omitted, installs all hooks that have rules configured in your `.fisherman.toml` files.

## Behavior

1. **Reads configuration** from all scopes (global, repository, local)
2. **Creates hook scripts** in `.git/hooks/` for each configured hook
3. **Sets execute permissions** on hook scripts automatically
4. **Backs up existing hooks** when using `--force` (saved as `.bkp` files)

## Examples

**Install all configured hooks:**

```bash
fisherman install
```

**Install specific hooks:**

```bash
fisherman install pre-commit commit-msg
```

**Force install (override existing hooks):**

```bash
fisherman install --force
```

**Force install specific hooks:**

```bash
fisherman install --force pre-commit pre-push
```

## Exit Codes

- `0` - Success
- `1` - Error (e.g., not in a Git repository, configuration error)

## Notes

- Must be run from within a Git repository
- Hooks are installed per-repository (not globally)
- Original hook scripts are preserved with `.bkp` extension when using `--force`
- Hook scripts call `fisherman handle <hook>` to execute configured rules

---

# `explain` - Explain Hook Configuration

Displays detailed information about a specific hook's configuration, including all rules that will execute when the hook
is triggered.

## Syntax

```bash
fisherman explain <HOOK>
```

## Arguments

- `HOOK` (required) - Name of the Git hook to explain (e.g., `pre-commit`, `commit-msg`)

## Behavior

1. **Reads configuration** from all scopes (global, repository, local)
2. **Displays hook information**:
    - Hook name
    - Configuration file sources
    - All rules that will execute
    - Rule types and parameters
3. **Shows merged configuration** from all scopes

## Examples

**Explain pre-commit hook:**

```bash
fisherman explain pre-commit
```

**Explain commit-msg hook:**

```bash
fisherman explain commit-msg
```

## Output Format

```
Hook: pre-commit
Configured in:
  - /home/user/.fisherman.toml
  - /home/user/project/.fisherman.toml

Rules:
  1. exec cargo test
  2. exec cargo fmt --check
  3. exec cargo clippy -- -D warnings
```

## Exit Codes

- `0` - Success
- `1` - Error (e.g., invalid hook name, configuration error)

## Notes

- Does not execute any rules, only displays configuration
- Useful for debugging configuration issues
- Shows the merged configuration from all scopes

---

# `handle` - Execute Hook (Internal Command)

Executes the specified hook and its associated rules. This command is primarily used internally when Git triggers a
hook, but can be useful for testing or debugging.

## Syntax

```bash
fisherman handle <HOOK> [ARGS...]
```

## Arguments

- `HOOK` (required) - Name of the Git hook to execute (e.g., `pre-commit`, `commit-msg`)
- `ARGS` (optional) - Arguments passed by Git to the hook (varies by hook type)

## Behavior

1. **Reads configuration** from all scopes
2. **Extracts variables** based on `extract` configuration
3. **Compiles rules** that apply to the hook
4. **Executes rules sequentially**:
    - Synchronous rules first (validation)
    - Asynchronous rules second (commands, scripts)
5. **Reports results** for each rule
6. **Exits with appropriate code**:
    - `0` if all rules pass
    - `1` if any rule fails (aborts Git operation for pre-hooks)

## Examples

**Manually trigger pre-commit hook:**

```bash
fisherman handle pre-commit
```

**Manually trigger commit-msg hook:**

```bash
fisherman handle commit-msg .git/COMMIT_EDITMSG
```

## Output Format

```
Hook: pre-commit
Configured in:
  - /home/user/.fisherman.toml

exec cargo test executed successfully
exec cargo fmt --check executed successfully
```

Or if a rule fails:

```
Hook: pre-commit
Configured in:
  - /home/user/.fisherman.toml

exec cargo test: some tests failed
```

## Exit Codes

- `0` - All rules passed
- `1` - One or more rules failed

## Notes

- **Internal use**: Typically called by Git hooks installed via `fisherman install`
- **Manual testing**: Useful for testing hook configuration without triggering Git
- **Arguments**: Some hooks receive arguments from Git (e.g., `commit-msg` receives path to commit message file)
- **Can abort operations**: For hooks that can abort Git operations (e.g., `pre-commit`, `commit-msg`, `pre-push`), a
  failure prevents the Git operation

## Common Use Cases for Manual Testing

**Test pre-commit hook:**

```bash
# Make some changes
git add .

# Test pre-commit hook manually
fisherman handle pre-commit

# If successful, commit
git commit -m "test commit"
```

**Test commit-msg hook:**

```bash
# Create a test commit message
echo "feat: test message" > /tmp/test-commit-msg

# Test commit-msg hook manually
fisherman handle commit-msg /tmp/test-commit-msg
```

---

# Global Options

These options work with all commands:

- `-h`, `--help` - Display help information for the command
- `-V`, `--version` - Display Fisherman version information

## Examples

```bash
# Get help for install command
fisherman install --help

# Get help for Fisherman
fisherman --help

# Display version
fisherman --version
```

---

# Environment Variables

Fisherman respects these environment variables:

- **`HOME`** - Used to locate global configuration (`~/.fisherman.toml`)
- **`GIT_DIR`** - Git directory location (usually `.git`)

---

# Exit Codes

All Fisherman commands use these exit codes:

- `0` - Success
- `1` - Error (configuration error, rule failure, invalid arguments, etc.)

---

# Typical Workflow

## Initial Setup

```bash
# 1. Create configuration
cat > .fisherman.toml <<EOF
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"
EOF

# 2. Install hooks
fisherman install

# 3. Verify configuration
fisherman explain pre-commit
fisherman explain commit-msg
```

## Testing Hooks

```bash
# Test hooks manually
fisherman handle pre-commit
fisherman handle commit-msg .git/COMMIT_EDITMSG
```

## Updating Configuration

```bash
# 1. Edit configuration
vim .fisherman.toml

# 2. Reinstall hooks (if needed)
fisherman install --force

# 3. Test new configuration
fisherman explain pre-commit
fisherman handle pre-commit
```

---

# Troubleshooting Commands

## Check Installed Hooks

```bash
# List installed hooks
ls -la .git/hooks/

# View hook script content
cat .git/hooks/pre-commit
```

## Verify Configuration

```bash
# Explain all configured hooks
fisherman explain pre-commit
fisherman explain commit-msg
fisherman explain pre-push
```

## Test Hooks Manually

```bash
# Test individual hooks
fisherman handle pre-commit
fisherman handle commit-msg .git/COMMIT_EDITMSG
```

## Check Hook Execution

```bash
# Add verbose output to see what's happening
git commit -v
```

---

# See Also

- [Installation](./Installation) - How to install Fisherman
- [Configuration](./Configuration) - Configuration file format
- [Rules Reference](./Rules-reference) - Available rule types
- [Git Hooks](./Git-hooks) - Git hooks reference
