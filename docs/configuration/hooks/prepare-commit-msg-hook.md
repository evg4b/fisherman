---
id: prepare-commit-msg-hook
title: prepare-commit-msg
---
:::danger
This page under construction...
:::

## Structure of configuration section

```yaml
commit-msg:
    variables:
      from-branch: regexp
    message: template
```

## Configurations rules:

### variables

### message
This is [template](/). Users commit message will be replaced by a compilation result of this field.

:::caution
When you commit with the `--no-verify` flag this action will not work.
:::
