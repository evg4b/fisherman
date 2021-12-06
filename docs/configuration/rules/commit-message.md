---
id: commit-message
title: commit-message
---

<!-- TODO: Add correct description -->

``` yaml
- type: commit-message
  when: 1 == 1
  prefix: message-prefix
  suffix: message-suffix
  regexp: '[a-zA-Z]'
```

- **commit-regexp** - The regular expression to validation commit messuage.
- **commit-prefix** - The template with which the message should start.
- **commit-suffix** - The template with which the message should end.
- **not-empty** - This boolean value indicates whether a commit with an empty message is not allowed.

Can be used to prevent commit with `--allow-empty-message` flag.
