---
id: pre-commit-hook
title: pre-commit
---
The `pre-commit` hook is run first, before you even type in a commit message. It’s used to inspect the snapshot that’s about to be committed, to see if you’ve forgotten something, to make sure tests run, or to examine whatever you need to inspect in the code.

## Structure of configuration section
```yaml
pre-commit:
  variables:
    from-branch: 'regexp'
  shell:
    script1:
      - command1
      - command2
    script2:
      - command3
      - command4
```

## Configurations rules:

### variables:
This section is common variables section without additional params. See more information [here](./../variables.md).

### shell:
