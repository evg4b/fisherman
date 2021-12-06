---
id: add-to-index
title: add-to-index
---

The rule is used to automatically add files to the index after rules or previous scripts.
This rule can be used for adding auto-generated files (mocks, styles, and others) in the index.

## Syntax

``` yaml
- type: add-to-index
  globs:
    - glob: ~/project/*.go
      required: true
    - glob: ~/mocks/*.go
      required: false
    - glob: ~/assets/*.css
    - glob: ~/assets/*.js
```

- **globs** - List of files that will be added to the git index, before commit but after when all
  validations rules finished. These files always will be added to the git index.
- **required** - This flag marks this action as required or not. When glob is masked as required
  and where there are no files matched to it, commit will be rejected.

You can also use a short syntax (in this case required will be false):

``` yaml
- type: add-to-index
  globs: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
```
