---
id: installation
title: Installation
slug: /installation
---

## Get the binary
You can download the binary from [the releases page on GitHub](https://github.com/evg4b/fisherman/releases) and add to your $PATH. The fisherman_checksums.txt file contains the SHA-256 checksum for each file.

## Build from source

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
