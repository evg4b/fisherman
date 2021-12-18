---
id: commit-message
title: commit-message
---

The rule is used to check a commit message. It sets the criteria for the content
of the message. Can be used to prevent committing with `--allow-empty-message` flag.

One field or all fields can be filled in a rule. Several such ones can also be defined,
in which case the message must satisfy all of them.

## Syntax

``` yaml
- type: commit-message
  prefix: message-prefix
  suffix: message-suffix
  regexp: '[a-zA-Z]'
```

- **commit-regexp** - The regular expression to validate a commit message.
- **commit-prefix** - The template with which the message should start.
- **commit-suffix** - The template with which the message should end.
- **not-empty** - This boolean value indicates whether a commit with an empty
  message is not allowed.
