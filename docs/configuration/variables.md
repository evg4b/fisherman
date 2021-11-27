---
id: variables
title: Variables
---

## Common `variables` and `extract-variables` sections

Each hook configuration section contains `variables` and `extract-variables` subsection. These sections used for define user variables.

```yaml
variables:
  VariableName1: 'variable-value'
extract-variables:
  - Extract(BranchName, '^(?P<VariableName2>demo1)$')
  - Extract(BranchName, '^(?P<VariableName3>demo2)$')
```

### variables

lorem

### extract-variables

`extract-variables` section extracts variables from the name of the current branch using named groups in regular expressions.

:::info
Remember fisherman is written in the [go programming language](https://golang.org) and you should use golang notation for named grouped: `(?P<name>regex)`. You can get more information on the [official website of the package regexp](https://golang.org/pkg/regexp/).
You can also test your regular expressions on service [Rego](https://regoio.herokuapp.com/) or [regex101](https://regex101.com/) (note: select flavor golang).
:::

For example: for branch `refs/heads/N1-do-something` and regexp `^refs\/heads\/(?P<Number>N\d+)-(?P<Description>.*)$` for all templates in section `commit-msg` will be defined variables `Number` with value `N1` and `Description` with value `do-something`.

## Global predefined variables

The following variables are defined by default:

| Variable name    | Description                                                         |
|------------------|---------------------------------------------------------------------|
| FishermanVersion | The version of fisherman on which the hook is launched              |
| CWD              | Full path to current working directory                              |
| UserName         | User name from git configuration                                    |
| Email            | User email from git configuration                                   |
| OS               | Name of current operation system (linux, windows, or darwin(macos)) |
