---
id: exec
title: exec
---

This rule runs program with passe arguments for validation. When program returns exit code different
from zero then action will be canceled.

:::caution Note
`exec` rule can runs only executable files in your system. In this one you can not use commands from
shell e.p. `cp`, `rm`, `mkdir` and others. For more details, see the information on
[the link](../../faq#difference-between-exec-and-shell-script).
:::caution

## Syntax

``` yaml
- type: run-program
  name: Check name
  env:
    VAR1: value
    VAR2: value
  dir: other-directory
  output: true
  commands:
    - program: executable1
      args: [arg1, arg2]
      env:
        VAR1: new value
      dir: program-directory
      output: true
    - program: executable2
      args: [arg1, arg2, arg3]
```

- **when** - An expression on C like language. It allows you to define a condition for executing a program.
  See more in section [Condition expressions](../expressions.md).
- **name** - List of lines that should not be included in the commit.
- **env** - Sets additional environment variables (system environment variables also will be included) for the program.
- **dir** - Sets current working directory for program.
- **commands** -
  - **program** - Program name or path to program binary
  - **args** - List of arguments for start
  - **env** - Sets additional environment variables (system environment variables also will be included) for the program.
  - **dir** - Sets current working directory for program.
  - **output** - Indicates whether to print the command output. By default false. To display parallel output,
    use a prefix with script name before each output line.

``` yaml
- type: run-program
  name: Check name
  commands:
    - program1 command arg1
    - program2 command arg1 arg2
    - program3 command arg1 arg2 arg3
```

``` yaml
- type: run-program
  name: Check name
  commands:
    - program: executable1 arg1, arg2
      env:
        VAR1: value1
    - program: executable2 arg1 arg2 arg3
      env:
        VAR2: value2
```
