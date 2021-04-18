---
id: getting-started
title: Getting started
slug: /getting-started
---
## Install

Fisherman currently offers only two installation methods. Check out the available methods below.

### Get the binary

You can download the binary from [the releases page on GitHub](https://github.com/evg4b/fisherman/releases)
and add to your `PATH`. The fisherman_checksums.txt file contains the SHA-256 checksum for each file.

### Build from source

First, make sure you have Go properly installed and setup. Fisherman requires Go Modules.
Installing in another directory:

```bash
git clone git@github.com:evg4b/fisherman.git
cd fisherman

# Compiling binary to $GOPATH/bin
go install -v .

# Compiling it to another location.
# Replace <version> to fisherman version
go build -v -ldflags="-s -w -X fisherman/constants.Version=<version>"

./fisherman
```

## Initialize fisherman

Go to your repository's directory and run `init` command

```bash
cd ~/my-repository
fisherman init
```

When you already have installed hooks, fisherman will warn you about this.
To override your hooks use the `--force` flag. **Note:** in this case,
the contents of your hook files will be lost.

```bash
fisherman init --force
```

## Create configuration

Fisherman created a file called `.fisherman.yml` in the root of your repository.
The hooks selection should contain the configuration of hooks actions.
The example below checks that the commit message starts with `[fisherman]`.

```yaml
hooks:
  commit-msg:
    commit-prefix: '[fisherman]'
```

To quickly find common configuration solutions visit the [FAQ page](./faq.md).

:::caution
At the moment, the application works only with YAML files with `.yml` extension.
The config file must be named `.fisherman.yml` (NOT `.fisherman.yaml`).
:::

## Commit changes

Make your first commit under the supervision of a fisherman.

```bash
git add .
git commit -m '[fisherman] My first commit'
```
