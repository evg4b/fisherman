---
id: cli
title: CLI
slug: /cli
---

## Init

This command initialize fisherman in git repository. Is can  be called only in root on git repository.

```bash
fisherman init [--mode local] [--force] [--abs]
```

This one has following flags:

- **force** - forces overwrites existing hooks (configuration file will not be overwritten).
- **mode** - Set configuration storage mode (`repo` by default, `local` or `global`. See more information about it [here](/docs/configuration/configuration-files#configuration-file-inheritance).
- **absolute** - Writes in hook absolute path to fisherman binary. This flag can be user for setup concrete version of fisherman in repository. Or it can be used for handling hooks without registration fisherman in the PATH variable.

## Remove

This command remove hooks and local and repository configs from git repository.

```bash
fisherman remove
```

## Version

This command print fisherman version on screen. Is has no params

```bash
fisherman version
```
