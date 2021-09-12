---
id: rules
title: Rules
---

## add-to-index

The rule is used to automatically add files to the index after user or previous scripts. It can be used for adding auto generated files (mocks, styles and other). It has next syntax:

``` yaml
- type: 'add-to-index'
  when: '1 == 1'
  globs:
    - glob: '~/project/*.go'
      required: true
    - glob: '~/mocks/*.go'
      required: false
```

**globs** - List of files to be added to index before commit but after when all validations and shell scripts finished. This files always will be added to index.

**required** - This flag marks this action as an required or not. When glob masked as required and where there are no files matched to it, commit will be rejected.

You can also use a short syntax (required will be false):

``` yaml
- type: 'add-to-index'
  when: '1 == 1'
  globs: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
```

## commit-message

``` yaml
- type: 'commit-message'
  when: '1 == 1'
  prefix: 'message-prefix'
  suffix: 'message-suffix'
  regexp: '[a-zA-Z]'
```

**commit-regexp** - The regular expression to validation commit messuage.

**commit-prefix** - The template with which the message should start.

**commit-suffix** - The template with which the message should end.

**not-empty** - This boolean value indicates whether a commit with an empty message is not allowed.
Can be used to prevent commit with `--allow-empty-message` flag.

## prepare-message

``` yaml
- type: 'prepare-message'
  when: '1 == 1'
  message: 'Message draft'
```

**message** - This is [template](/). Users commit message will be replaced by a compilation result of this field.

:::caution Note
When you commit with the `--no-verify` flag this action will not work.
:::

## shell-script

This rules sets shell scripts for validation (All scripts are executed in parallel at the moment the hook is triggered).

Shell script can be configured using the following structure:

``` yaml
- type: commit-suffix
  when: 1 == 1
  suffix: string
```

You can create as many scripts to validation. Scripts will be executed on the [shell for your system](#shell-for-system).

**commands** - Array of strings with validation script commands. Is also supports tempesting based on hook variables.

**when** - An expression on C like language. It allows you to define a condition for executing a script. See more in section [Condition expressions](./expressions.md).

**env** - Sets additional environment variables (system environment variables also will be included)  for the command.

**output** - Indicates whether to print the command output. By default false. To display parallel output, use a prefix with script name before each output line.

Example:

``` text
script1 | validation data set 1...
script2 | Process started
script1 | validation data set 2...
script1 | validation data set 3...
script2 | Warning something went wrong
script1 | validation data set 3...
```

### Simple configuration

Also, when script does not require additional configuration (output and env variables), it can be set with the following code:

```yaml
single-line: 'command2 arg1'
single-line-array: ['command1', 'command2']
multilane-array:
  - 'command1 arg1 arg2'
  - 'command2 arg1'
```

### System related script

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

### Shell for system

Currently, only the following system shells are supported (They can be globally with field `default-shell` or with field `shell` fore each script directly):

- **Linux** - `bash`
- **Mac OS** - `bash`
- **Windows** - `powershell`

:::caution Note on windows
Powershell by default do not return non zero exit code on fail. [See more](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_preference_variables?view=powershell-7#erroractionpreference).
:::

## suppress-commit-files

``` yaml
- type: suppress-commit-files
  when: 1 == 1
  globs: [ 'mocks/**/*', 'go.sum', 'go.mod' ]
  remove-from-index: false
```

**globs** - List of files to be checked in index before commit. Commit fill be rejected when file will be founded.

**remove-from-index** - When this flag is `true` then files founded in index will be removed from it and commit well be continued.
## suppressed-text

The rule that prohibits committing if forbidden lines were added.

``` yaml
- type: suppressed-text
  when: 1 == 1
  substrings: [ 'suppressed', 'text' ]
  exclude: [ 'some/excluded/file.go' ]
```

**substrings** - List of lines that should not be included in the commit.

**exclude** - List of globs in which you do not need to check this rule.
