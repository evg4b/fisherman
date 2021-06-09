---
id: configuration-files
title: Configuration files
---

Fisherman supports very flexible configuration. All its configuration is contained in `.fisherman.yml` file.

## Configuration file inheritance

Fisherman supports inheritance of configuration files with different visibility levels.
At the moment, there are 3 places where fisherman is looking for configuration files.

- **Global configuration for user (Global config)** - Global configuration for the user.
  This config file is located in the user directory (`~/.fisherman.yml` for linux and `%USERPROFILE%/.fisherman.yml` for windows)
  The settings specified in this file will apply to which repository in which fisherman was initialized.
  The configuration has the lowest priority. All rules will be overridden in other configuration files.

- **Repository configuration (Repo config)** - This configuration is located directly in the repository
  and can be shared with your team. The rules from this file will overwrite the rules from the global configuration.

- **Local configuration for repository** - This configuration is located in `.git` directory in your repository.
  Rules in this file have the highest priority and well override all other rules.
  This was done primarily to add your own rules without changing repository configuration
  file or without configuring files in the repository.

The loaded configuration files and the path to them you can find in the header of hook handler.

``` text
           .d8d.  |
          d88888b | Hook: commit-msg-hook
          "Y888Y" |
 .          888   | Global config: /home/user/.fisherman.yml
 8b.        888   | Repo config: /home/user/documents/my-repo/.fisherman.yml
 888b.      888   | Local config: /home/user/documents/my-repo/.git/.fisherman.yml
 888       .88P   |
 "Y8b.....d88P"   | Fisherman: 0.0.1-alpha.1
  "Y8988888P"     |
```
