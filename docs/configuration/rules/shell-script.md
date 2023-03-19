---
id: shell-script
title: shell-script
---

This rules sets shell scripts for validation (All scripts are executed in parallel
at the moment the hook is triggered).

## Syntax

``` yaml
- type: shell-script
  when: condition
  name: check name
  shell: shell
  commands:
    - command1
    - command2
  env:
    var1: value1
    var2: value2
  encoding: cp866
  output: true
  Dir: ./src
```

You can create as many scripts to validation. Scripts will be executed on the
[shell for your system](#shell-for-system).

- **when** - An expression on C like language. It allows you to define a condition
  for executing a script. See more in section [Condition expressions](../expressions.md).
- **name** - Name of check. It will be used as prefix in output.
- **shell** - Name of shell to run script.
- **commands** - Array of strings with validation script commands. Is also supports
  tempesting based on hook variables.
- **env** - Sets additional environment variables (system environment variables
  also will be included)  for the command.
- **encoding** - IANA name of shell output encoding.
- **output** - Indicates whether to print the command output. By default false.
  To display parallel output, use a prefix with script name before each output line.
  Example:

  ``` text
  script1 | validation data set 1...
  script2 | Process started
  script1 | validation data set 2...
  script1 | validation data set 3...
  script2 | Warning something went wrong
  script1 | validation data set 3...
  ```

## Simple configuration

Also, when script does not require additional configuration (output and env variables),
it can be set with the following code:

```yaml
# single line
- type: shell-script
  commands: 'command2 arg1'

# single line array
- type: shell-script
  commands: [ 'command1', 'command2' ]

# multilane array
- type: shell-script
  commands:
    - 'command1 arg1 arg2'
    - 'command2 arg1'
```

## System related script

In the case when it is not possible to specify a universal script for all systems,
you can specify separated scripts for each system. This is possible using
the rule condition (More details see [here](./../expressions.md)).

```yaml
- type: shell-script
  when: IsLinux()
  commands: linux-bin arg1

- type: shell-script
  when: IsWindows()
  commands: windows-bin.exe arg1

- type: shell-script
  when: IsMocOs()
  commands: macos-bin arg1
```

## Shell for system

Currently, only the following system shells are supported:

- **cmd** - cmd.exe as shell host. Can be used only on Windows.
  It is also the default shell for Windows.
- **bash** - bash as shell host. Available on Windows, Linux and
  MacOs. It is the default shell for Linux and MacOs. Note: for Windows bash available
  only in package [MSYS2](https://www.msys2.org/).
- **powershell** - powershell as shell host. Can be used on Windows,
  Linux and MacOs.
