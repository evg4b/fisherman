Fisherman provides a powerful variable extraction and template system that allows you to extract information from your
Git context (branch name, repository path) and use it throughout your hook configurations.

# Overview

Variables and templates enable:

- **Dynamic configurations** - Adapt rule behavior based on branch names or paths
- **Conditional execution** - Execute rules only when certain conditions are met
- **Template substitution** - Inject variables into commands, messages, and file content

---

# Variable Extraction

Variables are extracted from Git context using regex patterns with named capture groups.

## Extraction Syntax

```toml
extract = ["<source>:<regex-pattern>"]
```

**Format:** `"<source>:<regex>"`

where:

- `<source>` - The data source to extract from
- `<regex>` - A regular expression with named capture groups

## Supported Sources

## `branch` - Current Branch Name (Required Match)

Extracts variables from the current Git branch name. If the pattern doesn't match, Fisherman will throw an error.

**Example:**

```toml
extract = ["branch:^(?<IssueNumber>PROJ-\\d+)-.*$"]
```

For branch `PROJ-1234-fix-bug`, this extracts:

- `IssueNumber = "PROJ-1234"`

---

## `branch?` - Current Branch Name (Optional Match)

Extracts variables from the branch name, but doesn't fail if the pattern doesn't match.

**Example:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]
```

For branch `main`, this extracts nothing (no error).
For branch `PROJ-1234-fix-bug`, this extracts:

- `IssueNumber = "PROJ-1234"`

---

## `repo_path` - Repository Path (Required Match)

Extracts variables from the absolute path to the repository. If the pattern doesn't match, Fisherman will throw an
error.

**Example:**

```toml
extract = ["repo_path:^/home/(?<Username>[^/]+)/.*$"]
```

For path `/home/john/projects/myrepo`, this extracts:

- `Username = "john"`

---

## `repo_path?` - Repository Path (Optional Match)

Extracts variables from the repository path, but doesn't fail if the pattern doesn't match.

**Example:**

```toml
extract = ["repo_path?:^/home/(?<Username>[^/]+)/.*$"]
```

For path `/opt/projects/myrepo`, this extracts nothing (no error).
For path `/home/john/projects/myrepo`, this extracts:

- `Username = "john"`

---

## Multiple Extractions

You can extract multiple variables from different sources or using multiple patterns:

```toml
extract = [
    "branch?:^(?<Type>feature|bugfix|hotfix)/(?<IssueNumber>\\d+)-.*$",
    "repo_path?:^/home/(?<Username>[^/]+)/.*$"
]
```

For branch `feature/123-new-widget` and path `/home/alice/projects/myapp`:

- `Type = "feature"`
- `IssueNumber = "123"`
- `Username = "alice"`

---

## Extraction Examples

## Extract Issue Number from Branch

```toml
extract = ["branch?:^(?<IssueNumber>JIRA-\\d+)-.*$"]
```

| Branch              | Extracted Variables         |
|---------------------|-----------------------------|
| `JIRA-1234-fix-bug` | `IssueNumber = "JIRA-1234"` |
| `main`              | (none)                      |
| `feature/new-thing` | (none)                      |

---

## Extract Multiple Parts from Branch

```toml
extract = ["branch?:^(?<Type>[^/]+)/(?<Project>[^/]+)/(?<Description>.*)$"]
```

| Branch                         | Extracted Variables                                                       |
|--------------------------------|---------------------------------------------------------------------------|
| `feature/backend/user-auth`    | `Type = "feature"`, `Project = "backend"`, `Description = "user-auth"`    |
| `bugfix/frontend/button-style` | `Type = "bugfix"`, `Project = "frontend"`, `Description = "button-style"` |
| `main`                         | (none)                                                                    |

---

## Extract Username from Path

```toml
extract = ["repo_path?:^/home/(?<Username>[^/]+)/.*$"]
```

| Repository Path              | Extracted Variables  |
|------------------------------|----------------------|
| `/home/alice/projects/myapp` | `Username = "alice"` |
| `/home/bob/code/fisherman`   | `Username = "bob"`   |
| `/opt/projects/app`          | (none)               |

---

# Template Substitution

Once variables are extracted, you can use them in your rules using the `{{variable}}` syntax.

## Template Syntax

```toml
"{{VariableName}}"
```

Templates can be used in:

- Command arguments (`args`)
- File paths (`path`)
- File content (`content`)
- Commit message prefixes/suffixes
- Shell scripts

## Template Examples

## Simple Substitution

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

For branch `PROJ-1234-fix-bug`, commit messages must start with `PROJ-1234: `.

---

## Multiple Variables in Templates

```toml
extract = ["branch?:^(?<Type>[^/]+)/(?<Issue>\\d+)-.*$"]

[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Type\") && is_def_var(\"Issue\")"
command = "echo"
args = ["Working on {{Type}} issue #{{Issue}}"]
```

For branch `feature/123-new-widget`, outputs: `Working on feature issue #123`.

---

## Templates in Shell Scripts

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.pre-push]]
type = "shell"
when = "is_def_var(\"IssueNumber\")"
script = """
echo "Validating issue {{IssueNumber}}..."
curl -f "https://api.example.com/issues/{{IssueNumber}}"
"""
```

---

## Templates in File Content

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.post-checkout]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = "{{IssueNumber}}: "
```

---

# Conditional Expressions

The `when` parameter allows you to execute rules conditionally based on extracted variables and expressions.

## Expression Syntax

Fisherman uses the [Rhai scripting language](https://rhai.rs/) for expressions. Common operations:

## Check if Variable is Defined

```toml
when = "is_def_var(\"VariableName\")"
```

## Check if Variable is NOT Defined

```toml
when = "!is_def_var(\"VariableName\")"
```

## Compare Variable Values

```toml
when = "is_def_var(\"Type\") && Type == \"feature\""
```

## Numeric Comparisons

```toml
when = "is_def_var(\"IssueNumber\") && parse_int(IssueNumber) > 1000"
```

## Logical Operators

```toml
when = "is_def_var(\"Type\") && (Type == \"feature\" || Type == \"bugfix\")"
```

---

## Conditional Expression Examples

## Execute Only When Variable Exists

```toml
[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

---

## Execute When Variable Doesn't Exist

```toml
[[hooks.post-checkout]]
type = "write-file"
when = "!is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = ""
```

---

## Execute Based on Variable Value

```toml
extract = ["branch?:^(?<Type>[^/]+)/.*$"]

[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Type\") && Type == \"feature\""
command = "cargo"
args = ["test", "--features", "integration"]
```

---

## Complex Conditions

```toml
extract = [
    "branch?:^(?<Type>[^/]+)/(?<Issue>\\d+)-.*$"
]

[[hooks.pre-push]]
type = "shell"
when = "is_def_var(\"Type\") && is_def_var(\"Issue\") && Type == \"hotfix\" && parse_int(Issue) < 100"
script = "echo 'Warning: Old hotfix branch!'"
```

---

# Available Expression Functions

Fisherman provides these functions for use in `when` expressions:

## `is_def_var(variable_name)`

Check if a variable is defined.

**Parameters:**

- `variable_name` (string) - Name of the variable to check

**Returns:** `true` if the variable exists, `false` otherwise

**Example:**

```toml
when = "is_def_var(\"IssueNumber\")"
```

---

## `parse_int(string_value)`

Parse a string to an integer.

**Parameters:**

- `string_value` (string) - String to parse

**Returns:** Integer value

**Example:**

```toml
when = "parse_int(IssueNumber) > 1000"
```

---

## Standard Rhai Functions

Fisherman also supports standard Rhai functions and operators:

- **Comparison:** `==`, `!=`, `>`, `<`, `>=`, `<=`
- **Logical:** `&&` (and), `||` (or), `!` (not)
- **String:** `contains()`, `starts_with()`, `ends_with()`, `len()`
- **Math:** `+`, `-`, `*`, `/`, `%`

**Example:**

```toml
when = "is_def_var(\"Type\") && Type.starts_with(\"feature\")"
```

---

# Complete Examples

## Example 1: Dynamic Commit Message Prefix

```toml
# Extract issue number from branch name
extract = ["branch?:^(?<IssueNumber>JIRA-\\d+)-.*$"]

# Require issue number in commit message if branch has one
[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

**Behavior:**

- Branch `JIRA-1234-fix-bug`: Commit messages must start with `JIRA-1234: `
- Branch `main`: No prefix required

---

## Example 2: Conditional Test Execution

```toml
# Extract branch type
extract = ["branch?:^(?<Type>[^/]+)/.*$"]

# Run unit tests on all branches
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test", "--lib"]

# Run integration tests only on feature branches
[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Type\") && Type == \"feature\""
command = "cargo"
args = ["test", "--test", "integration"]
```

---

## Example 3: Dynamic Commit Template for Git GUI

```toml
# Extract issue number from branch
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

# Set commit template when branch has issue number
[[hooks.post-checkout]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = "{{IssueNumber}}: \n\n"

# Clear commit template when branch has no issue number
[[hooks.post-checkout]]
type = "write-file"
when = "!is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = ""
```

---

## Example 4: Multi-Team Repository

```toml
# Extract team and issue from branch
extract = ["branch?:^(?<Team>frontend|backend)/(?<Issue>\\d+)-.*$"]

# Frontend team: Run JS tests
[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Team\") && Team == \"frontend\""
command = "npm"
args = ["test"]

# Backend team: Run Rust tests
[[hooks.pre-commit]]
type = "exec"
when = "is_def_var(\"Team\") && Team == \"backend\""
command = "cargo"
args = ["test"]

# Both teams: Validate commit message format
[[hooks.commit-msg]]
type = "message-regex"
when = "is_def_var(\"Issue\")"
regex = "^\\[#{{Issue}}\\].*$"
```

---

## Example 5: User-Specific Configuration

```toml
# Extract username from repository path
extract = ["repo_path?:^/home/(?<Username>[^/]+)/.*$"]

# Send desktop notification for specific user
[[hooks.post-commit]]
type = "exec"
when = "is_def_var(\"Username\") && Username == \"alice\""
command = "notify-send"
args = ["Commit successful!", "Changes committed by {{Username}}"]
```

---

# Best Practices

1. **Use optional extraction (`?`) for flexibility** - Prevents errors when patterns don't match
2. **Always check variables before using them** - Use `when = "is_def_var(\"Variable\")"` to avoid errors
3. **Use descriptive variable names** - `IssueNumber` is better than `ID`
4. **Keep regex patterns simple** - Complex patterns are hard to maintain
5. **Test your patterns** - Use tools like [regex101.com](https://regex101.com/) to test patterns
6. **Document complex patterns** - Add comments explaining what each pattern extracts
7. **Use named groups** - `(?<Name>...)` makes intent clear
8. **Combine multiple extractions** - Extract from both branch and path when needed

---

# Troubleshooting

## Variable Not Extracted

**Problem:** Variable is not available in rules

**Solutions:**

- Check the extraction pattern matches your branch/path
- Verify the named group syntax: `(?<Name>...)`
- Use optional extraction (`branch?:` or `repo_path?:`) to avoid errors
- Test the regex pattern with your actual branch names

---

## Template Not Substituted

**Problem:** Template `{{Variable}}` appears literally in output

**Solutions:**

- Ensure variable is extracted (check `extract` configuration)
- Verify variable name matches exactly (case-sensitive)
- Check the rule has access to extracted variables
- Use `when = "is_def_var(\"Variable\")"` to verify extraction

---

## Conditional Not Working

**Problem:** Rule executes when it shouldn't (or vice versa)

**Solutions:**

- Verify expression syntax (use double quotes for strings)
- Check variable names are correct
- Test with simpler conditions first
- Use logical operators correctly (`&&`, `||`, `!`)
- Remember: `is_def_var()` requires double-escaped quotes in TOML: `\"Variable\"`

---

## Pattern Matching Issues

**Problem:** Regex pattern doesn't match expected branches

**Solutions:**

- Escape special regex characters (`\d`, `\.`, `\^`, etc.)
- In TOML, backslashes need double escaping: `\\d` for `\d`
- Test patterns on [regex101.com](https://regex101.com/) with Rust flavor
- Use `branch?:` during testing to avoid errors

---

# See Also

- [Configuration](./Configuration) - Configuration file structure and scopes
- [Rules Reference](./Rules-reference) - Complete list of available rules
- [Examples](./Examples-of-usage) - Real-world usage examples
