---
id: shell-script
title: Shell scripts section
---

This section sets shell scripts for validation. All scripts are executed in parallel at the moment the hook is triggered.

## Scripts section configuration

Shell script can be configured using the following structure:

```yaml
script-name:
  commands:
    - command2 arg1
    - command3 arg1 arg2
  env:
    VAR1: VALUE1
    VAR2: VALUE2
  output: true
```
You can create as many scripts to validation. Scripts will be executed on the [shell for your system](#shell-for-system).

### commands

Array of strings with validation script commands. Is also supports tempesting based on hook variables.

### env

Sets additional environment variables (system environment variables also will be included)  for the command.

### output

Indicates whether to print the command output. By default false. To display parallel output, use a prefix with script name before each output line.

``` text
script1 | validation data set 1...
script2 | Process started
script1 | validation data set 2...
script1 | validation data set 3...
script2 | Warning something went wrong
script1 | validation data set 3...
```

## Simple configuration

Also, when script does not require additional configuration (output and env variables), it can be set with the following code:

```yaml
single-line: command2 arg1
single-line-array: [command1, command2]
multilane-array:
  - command1 arg1 arg2
  - command2 arg1
```

## System related script

In the case when it is not possible to specify a universal script for all systems, you can specify separated scripts for each system.

:::caution Note
Be sure to specify scripts for each system, otherwise the hook will be skipped for the system without configuration
:::

```yaml
shell:
  windows:
    windows-script1: # ... script definition
    windows-script1: # ... script definition
  linux:
    linux-script1: # ... script definition
    linux-script2: # ... script definition
  darwin:
    darwin-script1: # ... script definition
```

## Shell for system

Currently, only the following system shells are supported (They can be globally with field `default-shell` or with field `shell` fore each script directly):

- **Linux** - `bash`
- **Mac OS** - `bash`
- **Windows** - `powershell`

:::caution Note on windows
Powershell by default do not return non zero exit code on fail. [See more](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_preference_variables?view=powershell-7#erroractionpreference).
:::
