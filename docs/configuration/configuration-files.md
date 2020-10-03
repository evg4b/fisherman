---
id: configuration-files
title: Configuration files
---

:::danger
This page under construction...
:::

```yaml
variables:
  name: value

hooks:
  commit-msg:
    variables:
      from-branch: 'regexp'
    not-empty: false
    commit-regexp: 'regexp'
    commit-prefix: 'template'
    commit-suffix: 'template'
    static-message: 'template'
  prepare-commit-msg:
    variables:
      from-branch: 'regexp'
    message: 'template'

output:
  level: None
```
