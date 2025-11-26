Fisherman supports all standard Git hooks. This reference explains when each hook is triggered, what arguments it
receives, and common use cases.

## Overview

Git hooks are scripts that Git executes before or after events such as commits, pushes, and merges. Fisherman intercepts
these hooks and executes configured rules based on your `.fisherman.toml` configuration.

---

## Client-Side Hooks

Client-side hooks are triggered by local operations like committing and merging.

### `pre-commit`

**Triggered:** Before a commit is created (after staging files but before entering the commit message)

**Can abort:** Yes (exit code 1 prevents commit)

**Arguments:** None

**Common use cases:**

- Run tests
- Check code formatting
- Run linters
- Validate file contents
- Check for debugging statements or TODO comments
- Ensure no large files are being committed

**Example configuration:**

```toml
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["fmt", "--check"]

[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["clippy", "--", "-D", "warnings"]
```

**When to use:** For validations that should happen before every commit.

---

### `prepare-commit-msg`

**Triggered:** After the default commit message is generated but before the editor is opened

**Can abort:** Yes

**Arguments:**

1. Path to commit message file
2. Source of commit message (`message`, `template`, `merge`, `squash`, or `commit`)
3. Commit SHA (only for `commit` source)

**Common use cases:**

- Modify commit message template
- Add issue numbers or tags automatically
- Include branch information in commit message
- Insert commit message guidelines

**Example configuration:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.prepare-commit-msg]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = ".git/COMMIT_EDITMSG"
content = "{{IssueNumber}}: "
append = true
```

**When to use:** To prepare or modify the commit message before the user edits it.

---

### `commit-msg`

**Triggered:** After the user enters a commit message but before the commit is created

**Can abort:** Yes (exit code 1 prevents commit)

**Arguments:**

1. Path to the file containing the commit message

**Common use cases:**

- Validate commit message format (e.g., Conventional Commits)
- Require issue tracking numbers
- Enforce commit message length
- Check for forbidden words or patterns
- Validate commit message structure

**Example configuration:**

```toml
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore)(\\(.+\\))?:\\s.+"

[[hooks.commit-msg]]
type = "message-prefix"
prefix = "PROJ-"
```

**When to use:** To validate or enforce commit message standards.

---

### `post-commit`

**Triggered:** After a commit is created

**Can abort:** No (commit already created)

**Arguments:** None

**Common use cases:**

- Send notifications
- Update logs
- Trigger CI/CD pipelines
- Generate documentation
- Update statistics

**Example configuration:**

```toml
[[hooks.post-commit]]
type = "exec"
command = "notify-send"
args = ["Commit successful", "Changes have been committed"]

[[hooks.post-commit]]
type = "write-file"
path = ".git/commit-log.txt"
content = "Committed at $(date)\n"
append = true
```

**When to use:** For actions that should happen after a successful commit.

---

### `pre-rebase`

**Triggered:** Before a rebase operation starts

**Can abort:** Yes (exit code 1 prevents rebase)

**Arguments:**

1. Upstream branch being rebased onto
2. Branch being rebased (optional, only if rebasing a branch other than current)

**Common use cases:**

- Prevent rebasing protected branches
- Warn about rebasing published commits
- Check for uncommitted changes

**Example configuration:**

```toml
[[hooks.pre-rebase]]
type = "shell"
script = """
if [ "$1" = "master" ]; then
  echo "WARNING: Rebasing onto master"
fi
"""
```

**When to use:** To prevent dangerous rebase operations.

---

### `post-checkout`

**Triggered:** After checking out a branch or file

**Can abort:** No

**Arguments:**

1. Previous HEAD ref
2. New HEAD ref
3. Flag indicating whether it's a branch checkout (1) or file checkout (0)

**Common use cases:**

- Update dependencies after branch switch
- Clear caches
- Update commit message templates based on branch
- Display branch-specific information
- Trigger builds

**Example configuration:**

```toml
extract = ["branch?:^(?<IssueNumber>PROJ-\\d+)-.*$"]

[[hooks.post-checkout]]
type = "write-file"
when = "is_def_var(\"IssueNumber\")"
path = ".git/commit_msg_template.txt"
content = "{{IssueNumber}}: "

[[hooks.post-checkout]]
type = "exec"
command = "npm"
args = ["install"]
```

**When to use:** For actions that should happen after switching branches.

---

### `post-merge`

**Triggered:** After a merge completes successfully

**Can abort:** No

**Arguments:**

1. Flag indicating whether the merge was a squash merge (1) or not (0)

**Common use cases:**

- Update dependencies after merge
- Rebuild project
- Send notifications
- Update submodules

**Example configuration:**

```toml
[[hooks.post-merge]]
type = "exec"
command = "npm"
args = ["install"]

[[hooks.post-merge]]
type = "exec"
command = "cargo"
args = ["build"]
```

**When to use:** To update project state after merging.

---

### `pre-push`

**Triggered:** Before pushing commits to a remote

**Can abort:** Yes (exit code 1 prevents push)

**Arguments:**

1. Name of remote
2. Location of remote

**Stdin:** Lines with format: `<local ref> <local sha> <remote ref> <remote sha>`

**Common use cases:**

- Validate branch names
- Run tests before pushing
- Prevent pushing to protected branches
- Check for secrets or sensitive data
- Verify all tests pass
- Lint code

**Example configuration:**

```toml
[[hooks.pre-push]]
type = "branch-name-regex"
regex = "^(feature|bugfix|hotfix)/[a-z0-9-]+$"

[[hooks.pre-push]]
type = "exec"
command = "cargo"
args = ["test", "--all"]

[[hooks.pre-push]]
type = "shell"
script = """
if git rev-parse --abbrev-ref HEAD | grep -q "^main$"; then
  echo "ERROR: Cannot push directly to main branch"
  exit 1
fi
"""
```

**When to use:** To validate changes before they're shared with others.

---

### `pre-applypatch`

**Triggered:** Before applying a patch

**Can abort:** Yes

**Arguments:** None

**Common use cases:**

- Validate patch format
- Check for conflicts
- Run tests

**When to use:** When working with email-based patches (uncommon in modern workflows).

---

### `applypatch-msg`

**Triggered:** After extracting the patch message but before applying the patch

**Can abort:** Yes

**Arguments:**

1. Path to file containing the proposed commit message

**Common use cases:**

- Validate commit message from patch
- Modify commit message before applying

**When to use:** When working with email-based patches (uncommon in modern workflows).

---

### `post-applypatch`

**Triggered:** After a patch is applied

**Can abort:** No

**Arguments:** None

**Common use cases:**

- Send notifications
- Update logs

**When to use:** When working with email-based patches (uncommon in modern workflows).

---

### `pre-merge-commit`

**Triggered:** After a merge succeeds but before the merge commit is created

**Can abort:** Yes (exit code 1 aborts the merge commit)

**Arguments:** None

**Common use cases:**

- Validate merge commit message
- Check for merge conflicts in specific files
- Run tests on merged code

**Example configuration:**

```toml
[[hooks.pre-merge-commit]]
type = "exec"
command = "cargo"
args = ["test"]
```

**When to use:** Similar to `pre-commit`, but specifically for merge commits.

---

### `post-rewrite`

**Triggered:** After commits are rewritten (e.g., `git commit --amend`, `git rebase`)

**Can abort:** No

**Arguments:**

1. Name of command that triggered the rewrite (`amend` or `rebase`)

**Stdin:** Lines with format: `<old sha> <new sha>`

**Common use cases:**

- Update references to rewritten commits
- Send notifications about history changes

**When to use:** To track or respond to commit history changes.

---

### `pre-auto-gc`

**Triggered:** Before automatic garbage collection

**Can abort:** Yes (exit code 1 prevents garbage collection)

**Arguments:** None

**Common use cases:**

- Prevent garbage collection during critical operations
- Log garbage collection events

**When to use:** Rarely needed in normal workflows.

---

## Server-Side Hooks

Server-side hooks are triggered by network operations and run on the Git server.

### `pre-receive`

**Triggered:** When receiving a push, before any references are updated

**Can abort:** Yes (prevents entire push)

**Arguments:** None

**Stdin:** Lines with format: `<old sha> <new sha> <ref name>`

**Common use cases:**

- Enforce project-wide policies
- Validate all commits in push
- Check permissions
- Prevent force pushes

**When to use:** On Git servers to enforce policies for all developers.

---

### `update`

**Triggered:** Once per ref being updated in a push

**Can abort:** Yes (prevents update of specific ref)

**Arguments:**

1. Name of ref being updated
2. Old object name
3. New object name

**Common use cases:**

- Per-branch permissions
- Enforce branch-specific policies
- Prevent updates to protected branches

**When to use:** On Git servers for per-branch policies.

---

### `post-receive`

**Triggered:** After all references are updated on the server

**Can abort:** No

**Arguments:** None

**Stdin:** Lines with format: `<old sha> <new sha> <ref name>`

**Common use cases:**

- Trigger CI/CD pipelines
- Send notifications
- Update documentation
- Deploy applications
- Update issue trackers

**When to use:** On Git servers to trigger post-push actions.

---

### `post-update`

**Triggered:** After all references are updated (similar to post-receive)

**Can abort:** No

**Arguments:** List of updated refs

**Common use cases:**

- Update Git repository metadata
- Send email notifications

**When to use:** On Git servers (post-receive is generally preferred).

---

### `push-to-checkout`

**Triggered:** When pushing to a repository with a checked-out working directory

**Can abort:** Yes

**Arguments:**

1. Name of ref being updated

**Common use cases:**

- Update working directory after push
- Deploy to production

**When to use:** On Git servers with checked-out working directories.

---

### `proc-receive`

**Triggered:** When receiving a push for refs matching `receive.procReceiveRefs`

**Can abort:** Yes

**Arguments:** None

**Stdin:** Lines with format: `<old sha> <new sha> <ref name>`

**Common use cases:**

- Custom reference handling
- Integration with external systems

**When to use:** Advanced server-side scenarios with custom ref processing.

---

### `reference-transaction`

**Triggered:** At various points during a reference transaction

**Can abort:** Yes (during prepared state only)

**Arguments:**

1. Transaction state (`prepared`, `committed`, `aborted`)

**Stdin:** Lines with format: `<old sha> <new sha> <ref name>`

**Common use cases:**

- Audit reference changes
- Implement custom reference policies

**When to use:** Advanced use cases requiring transaction-level control.

---

## Other Hooks

### `post-index-change`

**Triggered:** After the index (staging area) changes

**Can abort:** No

**Arguments:**

1. Flag indicating if `index.lock` file existed
2. Flag indicating if `index` file was modified

**Common use cases:**

- Update IDE state
- Regenerate file indexes
- Update build caches

**When to use:** For development tools that need to respond to staging changes.

---

### `fsmonitor-watchman`

**Triggered:** When Git needs to know which files have changed (if fsmonitor is enabled)

**Can abort:** No

**Arguments:**

1. Version (currently "2")
2. Timestamp of last update

**Common use cases:**

- Integrate with file system monitoring tools (Watchman)
- Optimize Git operations in large repositories

**When to use:** For performance optimization in very large repositories.

---

### `sendemail-validate`

**Triggered:** Before sending an email via `git send-email`

**Can abort:** Yes

**Arguments:**

1. Path to file containing email to be sent

**Common use cases:**

- Validate email format
- Check for required headers

**When to use:** When using `git send-email` for patch submission.

---

## Perforce Hooks

These hooks are specific to Git's Perforce (P4) integration.

### `p4-changelist`

**Triggered:** When creating a Perforce changelist message

**Arguments:**

1. Name of changelist template file

**When to use:** When using `git-p4` for Perforce integration.

---

### `p4-prepare-changelist`

**Triggered:** After preparing the Perforce changelist message but before editor is opened

**Arguments:**

1. Name of changelist template file

**When to use:** When using `git-p4` for Perforce integration.

---

### `p4-post-changelist`

**Triggered:** After a Perforce changelist is submitted

**Arguments:**

1. Name of changelist template file

**When to use:** When using `git-p4` for Perforce integration.

---

### `p4-pre-submit`

**Triggered:** Before submitting to Perforce

**Can abort:** Yes

**Arguments:** None

**When to use:** When using `git-p4` for Perforce integration.

---

## Commonly Used Hooks Summary

For most projects, you'll primarily use these hooks:

| Hook            | Phase                  | Can Abort | Common Use                     |
|-----------------|------------------------|-----------|--------------------------------|
| `pre-commit`    | Before commit          | Yes       | Run tests, check formatting    |
| `commit-msg`    | After entering message | Yes       | Validate commit message        |
| `post-commit`   | After commit           | No        | Notifications, logging         |
| `post-checkout` | After branch switch    | No        | Update dependencies, templates |
| `pre-push`      | Before push            | Yes       | Validate branch, run tests     |
| `post-merge`    | After merge            | No        | Update dependencies            |

---

## Best Practices

1. **Keep hooks fast** - Slow hooks interrupt developer workflow
2. **Use appropriate hooks** - Match actions to the right hook lifecycle
3. **Provide clear error messages** - Help developers understand failures
4. **Test hooks thoroughly** - Use `fisherman handle <hook>` to test
5. **Document hook behavior** - Explain what each hook does in your configuration
6. **Use client-side hooks for validation** - Server-side hooks for enforcement
7. **Make hooks skippable for emergencies** - Use `git commit --no-verify` when necessary
8. **Don't run expensive operations in every hook** - Use conditional execution

---

## See Also

- [Configuration](Configuration.md) - How to configure hooks
- [Rules Reference](Rules.md) - Available rule types
- [Examples](Examples-of-usage.md) - Real-world examples
- [Git Hooks Documentation](https://git-scm.com/docs/githooks) - Official Git hooks reference
