package constants

// Version is constant to define app version
var Version = "x.x.x"

// Logo is string width logo
var Logo = `
 .d888  d8b          888
 d88P"  Y8P          888
 888                 888
 888888 888 .d8888b  88888b.   .d88b.  888d888 88888b.d88b.   8888b.  88888b.
 888    888 88K      888 "88b d8P  Y8b 888P"   888 "888 "88b     "88b 888 "88b
 888    888 "Y8888b. 888  888 88888888 888     888  888  888 .d888888 888  888
 888    888      X88 888  888 Y8b.     888     888  888  888 888  888 888  888
 888    888  88888P' 888  888  "Y8888  888     888  888  888 "Y888888 888  888

                         Fisherman version: {{.}}

`

// HookHeader is string width header for hook
var HookHeader = `
           .d8d.  |
          d88888b | Hook: {{.Hook}}
          "Y888Y" | 
 .          888   | Global config: {{.GlobalConfigPath}}
 8b.        888   | Repo config: {{.RepoConfigPath}}
 888b.      888   | Local config: {{.LocalConfigPath}}
 888       .88P   |
 "Y8b.....d88P"   | Fisherman: {{.Version}}
  "Y8988888P"     |

`
