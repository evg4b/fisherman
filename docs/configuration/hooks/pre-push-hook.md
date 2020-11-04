---
id: pre-push-hook
title: pre-push
---
The `pre-push` hook runs during git push, after the remote refs have been updated but before any objects have been transferred. It receives the name and location of the remote as parameters, and a list of to-be-updated refs through stdin. You can use it to validate a set of ref updates before a push occurs.

## Structure of configuration section
```yaml
pre-push:
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
