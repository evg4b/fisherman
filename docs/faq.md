---
id: faq
title: FAQ
slug: /faq
---

## I want check commit prefix based on current brach name

Can be used to check that message starts with issue number.

```yaml
hooks:
  commit-msg:
    extract-variables:
      - variable: BranchName
        expression: '^refs/heads/(?P<IssueNumber>PROJ-\d+)-.*$'
    rules:
      - type: commit-message
        when: IssueNumber != nil
        prefix: '{{IssueNumber}}: '
```

For branch with name `PROJ-175-new_very_very_important_feature` commit message
should be started with `PROJ-175:`

This rule will be skipped when the branch name does not match to expression.
Remove `when: IssueNumber != nil` to change it.

Related links:

<!-- TODO: Add correct links -->

- [extract-variables](/)
- [commit-message](/)

## I want check TODOs in my code before commit

Can be used to check that message starts with the issue number.

```yaml
hooks:
  pre-commit:
    rules:
      - type: suppressed-text
        substrings: [ 'TODO: ' ]
        exclude: [ 'README.md' ]
```

Related links:

<!-- TODO: Add correct links -->

- [suppressed-text rule](./configuration/rules/suppressed-text.md)

## I want run tests or lint before push changes to remote repo

Can be used for final validation before pr create/update operation.

```yaml
hooks:
  pre-push:
    rules:
      - type: exec
        name: Linting
        program: golangci-lint
        args: [ run,  ./... ]
      - type: exec
        name: Tests
        program: go
        args: [ test,  ./... ]
```

Related links:

<!-- TODO: Add correct links -->

- [exec rule](./configuration/rules/exec)

## I want create difference rules for each operation system

Can be used for creation shell scripts specified for os or run program with
different params.

```yaml
hooks:
  commit-msg:
    rules:
      - type: commit-message
        when: IsWindows()
        suffix: ' (Committed on windows)'
      - type: commit-message
        when: IsLinux()
        suffix: ' (Committed on linux)'
      - type: commit-message
        when: IsMac()
        suffix: ' (Committed on mac)'
```

Related links:

<!-- TODO: Add correct links -->

- [exec rule](./configuration/rules/exec)

## Difference between exec and shell-script

Rules [`exec`](./configuration/rules/exec) and [`shell-script`](./configuration/rules/shell-script)
are very similar, at first glance. But they have one fundamental difference.
Exec the more lighter rule, and one has no shell overhead.

**In case wen you wan only exec program (lint or test tool)** and abort the
actions if the exit code is not 0. Then you should use `exec` command.

**In the case when you need more complex actions** for example: move files,
communicate with operation system or execute commands pipeline then you should
use `shell-scrip` command.
