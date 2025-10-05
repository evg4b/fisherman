# Configuration

Fisherman supports a variety of configuration options to configure its behavior and customize the hooks it manages.

## Configuration File Formats

Fisherman supports three configuration file formats:

- **TOML** (`.fisherman.toml`) - The default and recommended format. Human-readable and easy to edit.
- **YAML** (`.fisherman.yaml` or `.fisherman.yml`) - A popular format for configuration files.
- **JSON** (`.fisherman.json`) - Widely-used format for data interchange.

All formats are functionally equivalent - choose the one you prefer.

## Configuration Scopes

Fisherman uses a hierarchical configuration system where multiple configuration files can be loaded in sequence. This
allows you to define global rules that apply to all repositories, while still allowing repository-specific
customization.

### 1. Global Scope

**Location:** `~/.fisherman.toml` (in your home directory)

**Purpose:** Rules that apply to all repositories on your system

**Use case:** Common rules you want across all projects (e.g., preventing commits to `main` branch, enforcing commit
message format)

**Example:**

```toml
# ~/.fisherman.toml
[[hooks.pre-commit]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"
```

### 2. Repository Scope

**Location:** `/path/to/repo/.fisherman.toml` (in the root of your Git repository)

**Purpose:** Repository-specific rules shared with all developers via Git

**Use case:** Project-specific requirements (e.g., running tests, building documentation)

**Example:**

```toml
# /path/to/repo/.fisherman.toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]
```

### 3. Local Scope

**Location:** `/path/to/repo/.git/.fisherman.toml` (in the `.git` directory)

**Purpose:** Personal rules for a specific repository that are **not** shared with others

**Use case:** Personal development preferences that you don't want to commit to the repository

**Example:**

```toml
# /path/to/repo/.git/.fisherman.toml
[[hooks.pre-commit]]
type = "exec"
command = "notify-send"
args = ["Running pre-commit hook"]
```

### Configuration Merging

When multiple configuration files exist, Fisherman loads them in this order:

1. Global scope (`~/.fisherman.toml`)
2. Repository scope (`/path/to/repo/.fisherman.toml`)
3. Local scope (`/path/to/repo/.git/.fisherman.toml`)

**Important:** Fisherman merges configurations by **concatenating** rules for each hook. Rules are executed in the order
they appear across all configuration files.

## Configuration File Structure

### Basic Structure

```toml
# Optional: Extract variables from context (branch name, repo path)
extract = ["<source>:<regex>", "<source>:<regex>"]

# Define hooks and their rules
[[hooks.<hook-name>]]
type = "<rule-type>"
# ... rule-specific parameters

[[hooks.<hook-name>]]
type = "<rule-type>"
when = "<condition>"
# ... rule-specific parameters
```

### Top-Level Fields

#### `extract` (optional)

An array of variable extraction patterns. Each pattern extracts named groups from Git context (branch name or repository
path) and makes them available to rules as template variables.

**Format:** `"<source>:<regex>"`

**Supported sources:**

- `branch` - Current Git branch name (required match)
- `branch?` - Current Git branch name (optional match)
- `repo_path` - Repository path (required match)
- `repo_path?` - Repository path (optional match)

**Example:**

```toml
extract = [
    "branch:^(?<IssueNumber>PROJ-\\d+)-.*$",
    "repo_path?:^/home/(?<username>[^/]+)/.*$"
]
```

See [Variables and Templates](Variables-and-Templates.md) for more details.

### Hook Configuration

#### `[[hooks.<hook-name>]]`

Defines a rule for a specific Git hook. You can define multiple rules for the same hook - they will execute in sequence.

**Supported hooks:** See [Git Hooks](Git-Hooks.md) for the complete list.

**Example:**

```toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["clippy"]
```

### Rule Fields

All rules support these common fields:

#### `type` (required)

The type of rule to execute. See [Rules Reference](Rules.md) for available rule types.

#### `when` (optional)

A conditional expression that determines if the rule should execute. The rule only runs if the expression evaluates to
`true`.

**Example:**

```toml
[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

See [Variables and Templates](Variables-and-Templates.md) for expression syntax.

#### `extract` (optional)

Override the global `extract` configuration for this specific rule. If specified, only these extraction patterns will be
available to the rule.

**Example:**

```toml
[[hooks.pre-commit]]
type = "exec"
extract = ["branch:^(?<Feature>[^/]+)/.*$"]
command = "echo"
args = ["Working on feature: {{Feature}}"]
```

## Configuration Examples

### Minimal Configuration

```toml
# .fisherman.toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]
```

### Multi-Hook Configuration

```toml
# .fisherman.toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["fmt", "--check"]

[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"

[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+$"
```

### Configuration with Variables

```toml
# .fisherman.toml
extract = ["branch?:^(?<IssueNumber>JIRA-\\d+)-.*$"]

[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "

[[hooks.pre-push]]
type = "exec"
when = "is_def_var(\"IssueNumber\")"
command = "echo"
args = ["Pushing changes for {{IssueNumber}}"]
```

### YAML Configuration Example

```yaml
# .fisherman.yaml
extract:
  - "branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"

hooks:
  pre-commit:
    - type: exec
      command: cargo
      args:
        - test
    - type: exec
      command: cargo
      args:
        - fmt
        - --check

  commit-msg:
    - type: message-prefix
      when: is_def_var("IssueNumber")
      prefix: "{{IssueNumber}}: "
```

### JSON Configuration Example

```json
{
  "extract": [
    "branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"
  ],
  "hooks": {
    "pre-commit": [
      {
        "type": "exec",
        "command": "cargo",
        "args": [
          "test"
        ]
      }
    ],
    "commit-msg": [
      {
        "type": "message-prefix",
        "when": "is_def_var(\"IssueNumber\")",
        "prefix": "{{IssueNumber}}: "
      }
    ]
  }
}
```

## Configuration Best Practices

1. **Start simple** - Begin with basic rules and add complexity as needed
2. **Use global scope for common rules** - Avoid repeating the same rules in every repository
3. **Use repository scope for project-specific rules** - Share rules with your team
4. **Use local scope for personal preferences** - Keep personal rules out of version control
5. **Test your configuration** - Use `fisherman explain <hook>` to verify your configuration
6. **Use conditional execution** - Make rules flexible with `when` expressions
7. **Document complex rules** - Add comments to explain non-obvious configurations

## Troubleshooting

### Multiple Configuration Files

If you have multiple configuration files in the same scope (e.g., both `.fisherman.toml` and `.fisherman.yaml` in the
same directory), Fisherman will return an error. Keep only one configuration file per scope.

### Invalid Configuration

If your configuration file has syntax errors or invalid rule definitions, Fisherman will display an error message when
you try to install hooks. Use a TOML/YAML/JSON validator to check your syntax.

### Variables Not Available

If variables extracted in the global configuration are not available in rules, make sure:

- The extraction pattern matches the current branch or repo path
- The variable name in the template matches the named group in the regex
- The `when` condition checks if the variable is defined before using it

See [Variables and Templates](Variables-and-Templates.md) for more details.
