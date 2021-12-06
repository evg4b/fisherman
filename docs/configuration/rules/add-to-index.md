---
id: add-to-index
title: add-to-index
---

The rule is used to automatically add files to the index after user or previous scripts.
It can be used for adding auto generated files (mocks, styles and other).

It has next syntax:

``` yaml
- type: 'add-to-index'
  when: 1 == 1
  globs:
    - glob: ~/project/*.go
      required: true
    - glob: ~/mocks/*.go
      required: false
```

- **globs** - List of files to be added to index before commit but after when all validations and shell scripts finished. This files always will be added to index.
- **required** - This flag marks this action as an required or not. When glob masked as required and where there are no files matched to it, commit will be rejected.

You can also use a short syntax (required will be false):

``` yaml
- type: 'add-to-index'
  when: '1 == 1'
  globs: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
```
