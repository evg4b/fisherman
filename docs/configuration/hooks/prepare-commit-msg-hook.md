---
id: prepare-commit-msg-hook
title: prepare-commit-msg
---
This hook is invoked by git-commit right after preparing the default log message, and before the editor is started.

## Structure of configuration section

```yaml
commit-msg:
  variables:
    from-branch: 'regexp'
  message: 'template'
```

## Configurations rules

### variables

This section is common variables section without additional params. See more information [here](./../variables.md).

### message

This is [template](/). Users commit message will be replaced by a compilation result of this field.

:::caution Note
When you commit with the `--no-verify` flag this action will not work.
:::
