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
  add-to-index:
    globs:
      - glob: dist/*
        required: true
      - glob: styles/**/*.css
        required: false
  suppress-commit-files:
    globs: [glob1, glob2, glob3]
    remove-from-index: false
```

## Configurations rules

### variables

This section is common variables section without additional params. See more information [here](./../variables.md).

### shell

Section with [shell scripts](../shell-script.md) for validation or new file generation. When script finished with non zero code then push will be rejected and other scrips canceled.

### add-to-index

**globs** - List of files to be added to index before commit but after when all validations and shell scripts finished. This files always will be added to index.

**required** - This flag marks this action as an required or not. When glob masked as and where there are no files matched to it, commit will be rejected.

You can also use a short note (optional will be false).:

``` yaml
pre-commit:
  shell: # ...
  add-to-index: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
```

It can be used for adding auto generated files (mocks, styles and other). This example before commit generate new mocks add add updated files to commit index:

``` yaml
pre-commit:
  shell:
    generate: go generate ./mocks/...
  add-to-index:
    - mocks/**/*
```

### add-to-index

**globs** - List of files to be checked in index before commit. Commit fill be rejected when file will be founded.

**remove-from-index** - When this flag is `true` then files founded in index will be removed from it and commit
well be continued.

Also you can use short syntax for this section (strings array, `remove-from-index` is this case will be `false`):

``` yaml
suppress-commit-files: [glob1, glob2, glob3]
```
