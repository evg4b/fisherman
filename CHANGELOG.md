# CHANGELOG

## Unrealized

## Implemented :
  - Added `add-to-index` section for include files in index after scripts
  - Added property to configure working directory for shell scripts

## [0.0.1-alpha.2 - [13 Nov 2020]](https://github.com/evg4b/fisherman/releases/tag/0.0.1-alpha.2)

## Implemented :
- Pre-push and `pre-commit` hooks
- Parallel shell script for `pre-push` and `pre-commit` hooks
- Version command
- Build for MacOS

## Fixed :
 - Problem to run outside the git repository.

## [0.0.1-alpha.1 - [3 Oct 2020]](https://github.com/evg4b/fisherman/releases/tag/0.0.1-alpha.1)

Implemented init, remove, handler commands.

### Implemented :
- **commit-msg** hook handling with rules:
  - MessageRegexp
  - MessagePrefix
  - MessageSuffix
  - StaticMessage
  - NotEmpty
- **commit-msg** hook handling with rule:
  - Message

Also added custom loading variables from global and hooks sections.
