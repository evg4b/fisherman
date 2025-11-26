This document contains a quite complex example of how to use fisherman,
which may not be entirely clear from the documentation.

# Commit Message Validation Based on Branch Name

Extract an issue number from your branch name and automatically enforce it as a prefix in commit messages.

```toml .fisherman.toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.commit-msg]]
type = "message-prefix"
when = "is_def_var(\"IssueNumber\")"
prefix = "{{IssueNumber}}: "
```

**How it works:**

- For branch `PROJ-1234-fix-bug`, commits must start with `PROJ-1234: `
- If the branch doesn't match the pattern, no prefix is required
- The `?` in the extract rule makes this pattern optional

# Add a Commit Message Prefix to Fork

[Fork](https://git-fork.com/) is a UI git client, and it has a commit message template feature.
But it doesn't support dynamic templates based on the branch name or other variables.
You can use fisherman to fix this.
Fork keeps all configuration files in the `.git` directory including `commit_msg_template.txt` which contains the
template.
Just 2 rules with conditionals are enough to set the template based on the branch name.

```toml .fisherman.toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.post-checkout]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = "./.git/commit_msg_template.txt"
content = "{{IssueNumber}}: "

[[hooks.post-checkout]]
type = "write-file"
when = "!is_def_var(\"IssueNumber\")"
path = "./.git/commit_msg_template.txt"
content = ""
```

**How it works:**

- For branch `PROJ-1234-fix-bug`, the extracted variable `IssueNumber` will be `PROJ-1234`.
- The first rule has a condition that checks if the variable `IssueNumber` is defined.
- If it is defined, the fisherman writes the content `PROJ-1234: ` to the file `.git/commit_msg_template.txt`.
- The second rule has a condition that checks if the variable `IssueNumber` is not defined.
- If it is not defined, the fisherman an empty string to the file `.git/commit_msg_template.txt`.
- Fork will use this file as a commit message template.

# Branch name validation

Ensure that your branch name follows a specific pattern, such as starting with `feature/` or `bugfix/` and
does not contain non-ASCII characters.

```toml .fisherman.toml
[[hooks.pre-push]]
type = "branch-name-regex"
expression = "^(feature|bugfix)/[a-zA-Z0-9-_]+$"
```

**How it works:**

- Validates branch names like `feature/new-feature` or `bugfix/fix-issue`
- Rejects names with non-ASCII characters or incorrect prefixes
- The `pre-push` hook ensures this validation before pushing changes and creates a pull request if the branch is valid.

# Parallel Execution for Faster Pre-Commit Checks

Run multiple time-consuming checks concurrently to speed up your pre-commit workflow.

```toml .fisherman.toml
# These async rules will execute in parallel

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["clippy", "--", "-D", "warnings"]

[[hooks.pre-commit]]
type = "shell"
script = "cargo fmt --check"

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["build", "--release"]
```

**How it works:**

- All four async rules execute **in parallel** rather than sequentially
- If each command takes ~5 seconds, total execution time is ~5 seconds instead of ~20 seconds
- Synchronous rules (like `message-regex`) still run sequentially before async rules
- Failed rules are still reported even when running in parallel
- Best for independent tasks that don't depend on each other's results

**Performance comparison:**

| Execution Mode | Time for 4 rules @ 5s each |
|----------------|---------------------------|
| Sequential     | ~20 seconds               |
| Parallel       | ~5 seconds                |

This dramatically improves developer experience when multiple validation steps are required.
