---
id: run-program
title: run-program
---

This rules runs program for validation (All scripts are executed in parallel at the moment the hook is triggered).

Run program rule can be configured using the following structure:

``` yaml
- type: run-program
  when: 1 == 1
  name: Tests
  program: go
  args: ["test", "./..."]
```

- **name** - List of lines that should not be included in the commit.
- **program** - Program name or path to program binary
- **args** - List of arguments for start
- **when** - An expression on C like language. It allows you to define a condition for executing a program. See more in section [Condition expressions](../expressions.md).
- **env** - Sets additional environment variables (system environment variables also will be included) for the program.
- **output** - Indicates whether to print the command output. By default false. To display parallel output, use a prefix with script name before each output line.
- **dir** - Sets current working directory for program.
