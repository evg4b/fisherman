---
id: prepare-message
title: prepare-message
---

<!-- TODO: Add correct description -->

``` yaml
- type: 'prepare-message'
  when: '1 == 1'
  message: 'Message draft'
```

- **message** - This is [template](/). Users commit message will be replaced by
  a compilation result of this field.

:::caution Note
When you commit with the `--no-verify` flag this action will not work.
:::
