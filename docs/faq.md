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

Can be used to check that message starts with issue number.

```yaml
hooks:
  pre=commit:
    rules:
      - type: suppressed-text
        substrings: [ 'TODO: ' ]
        exclude: [ 'README.md' ]
```

Related links:

<!-- TODO: Add correct links -->
- [suppressed-text rule](./configuration/rules.md#suppressed-text)
