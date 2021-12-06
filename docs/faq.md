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

For branch with name `PROJ-175-new_very_very_important_feature` commit message should be started with `PROJ-175: `

This rule will be skipped when the branch name does not match to expression. Remove `when: IssueNumber != nil` to change it.

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
      - type: run-program
        name: Linting
        program: golangci-lint
        args: [ run,  ./... ]
      - type: run-program
        name: Tests
        program: go
        args: [ test,  ./... ]
```

Related links:

<!-- TODO: Add correct links -->
- [run-program rule](./configuration/rules/run-program.md)


## I want create difference rules for each operation system.

Can be used for creation shell scripts specified for os or run program with different params.

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
- [run-program rule](./configuration/rules/run-program.md)
