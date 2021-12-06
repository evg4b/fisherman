---
id: suppress-commit-files
title: suppress-commit-files
---

<!-- TODO: Add correct description -->

``` yaml
- type: suppress-commit-files
  when: 1 == 1
  globs: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
  remove-from-index: false
```

- **globs** - List of files to be checked in index before commit. Commit fill be rejected when file will be founded.
- **remove-from-index** - When this flag is `true` then files founded in index will be removed from it and commit well be continued.
