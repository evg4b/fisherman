---
id: commit-msg-hook
title: commit-msg
---

The `pre-commit` hook is run first, before you even type in a commit message. It can be used to override the commit message or validate its contents.

It has next configuration structure:

```yaml
commit-msg:
  variables:
    from-branch: 'regexp'
  not-empty: boolean
  commit-regexp: 'regexp'
  commit-prefix: 'template'
  commit-suffix: 'template'
  static-message: 'template'
```

## Configurations rules

### variables:
This section is common variables section without additional params. See more information [here](./../variables.md).

### static-message:
This is template. Users commit message will be replaced by a compilation result of this field.

:::caution
When this parament is set other validation rules will be skipped.
:::

:::caution Note
When you commit with the `--no-verify` flag this action will not work.
:::

## Validation rules:

### commit-regexp:
The regular expression to validation commit messuage.

### commit-prefix:
The template with which the message should start.

### commit-suffix:
The template with which the message should end.

### not-empty:
This boolean value indicates whether a commit with an empty message is not allowed.
Can be used to prevent commit with `--allow-empty-message` flag.
