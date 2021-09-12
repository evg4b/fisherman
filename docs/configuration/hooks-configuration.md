---
id: hooks-configuration
title: Hooks configuration
---

## Common hook configuration

### variables

This section is common variables section without additional params. See more information [here](./variables.md).

## applypatch-msg

This hook has following rules:

- [shell-script](./rules#shell-script)

## commit-msg

The `commit-msg` hook is run first, before you even type in a commit message. It can be used to override the commit message or validate its contents. This hook has following rules:

- [shell-script](./rules#shell-script)
- [commit-message](./rules#commit-message)

## fsmonitor-watchman

This hook has following rules:

- [shell-script](./rules#shell-script)

## post-update

This hook has following rules:

- [shell-script](./rules#shell-script)

## pre-applypatch

This hook has following rules:

- [shell-script](./rules#shell-script)

## pre-commit

The `pre-commit` hook is run first, before you even type in a commit message. It’s used to inspect the snapshot that’s about to be committed, to see if you’ve forgotten something, to make sure tests run, or to examine whatever you need to inspect in the code. This hook has following rules:

- [shell-script](./rules#shell-script)
- [add-to-index](./rules#add-to-index)
- [suppress-commit-files](./rules#suppress-commit-files)
- [suppress-text](./rules#suppress-text)

## pre-push

The `pre-push` hook runs during git push, after the remote refs have been updated but before any objects have been transferred. It receives the name and location of the remote as parameters, and a list of to-be-updated refs through stdin. You can use it to validate a set of ref updates before a push occurs. This hook has following rules:

- [shell-script](./rules#shell-script)

## pre-rebase

This hook has following rules:

- [shell-script](./rules#shell-script)

## pre-receive

This hook has following rules:

- [shell-script](./rules#shell-script)

## prepare-commit-msg

This hook is invoked by git-commit right after preparing the default log message, and before the editor is started. This hook has following rules:

- [prepare-message](./rules#prepare-message)

## update

This hook has following rules:

- [shell-script](./rules#shell-script)
