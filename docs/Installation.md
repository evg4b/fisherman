This guide covers how to install Fisherman on your system and get started with Git hook management.

# Requirements

- **Rust** 1.70 or higher (for building from source)
- **Git** 2.0 or higher
- **Operating System**: Linux, macOS, or Windows (with WSL)

# Installation Methods

## From Source (Recommended)

Install Fisherman using Cargo, Rust's package manager:

```bash
cargo install --git https://github.com/evg4b/fisherman.git
```

This will:

1. Clone the repository
2. Build the binary
3. Install it to `~/.cargo/bin/fisherman`

Make sure `~/.cargo/bin` is in your PATH.

## Verify Installation

After installation, verify that Fisherman is available:

```bash
fisherman --version
```

You should see output like:

```
fisherman 0.0.1
```

---

# Quick Start

Once installed, follow these steps to start using Fisherman:

## 1. Create a Configuration File

Create a `.fisherman.toml` file in your Git repository:

```bash
cd /path/to/your/repo
touch .fisherman.toml
```

## 2. Configure Your Hooks

Edit `.fisherman.toml` with your preferred text editor:

```toml
# Example: Run tests before committing
[[hooks.pre-commit]]
type = "exec"
command = "cargo"
args = ["test"]

# Example: Validate commit message format
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore):\\s.+"
```

## 3. Install Hooks

Install the configured hooks into your repository:

```bash
fisherman install
```

This creates executable hook scripts in `.git/hooks/` that call Fisherman with the appropriate hook name.

## 4. Test Your Hooks

Try making a commit to test your hooks:

```bash
git add .
git commit -m "feat: test fisherman hooks"
```

If configured correctly, your hooks will execute and validate the commit.

---

# Configuration Locations

Fisherman supports three configuration scopes. See [Configuration](./Configuration) for details.

## Global Configuration

Create a global configuration that applies to all repositories:

```bash
touch ~/.fisherman.toml
```

Edit the file with rules you want to apply everywhere:

```toml
# ~/.fisherman.toml

# Enforce conventional commits globally
[[hooks.commit-msg]]
type = "message-regex"
regex = "^(feat|fix|docs|style|refactor|test|chore)(\\(.+\\))?:\\s.+"
```

## Repository Configuration

Create a repository-specific configuration:

```bash
cd /path/to/your/repo
touch .fisherman.toml
```

This configuration is shared with all developers via Git.

## Local Configuration

Create a local configuration (not shared via Git):

```bash
cd /path/to/your/repo
touch .git/.fisherman.toml
```

This is useful for personal preferences you don't want to commit to the repository.

---

# Installing Hooks

## Install All Configured Hooks

```bash
fisherman install
```

This installs all hooks that have rules defined in your configuration files.

## Install Specific Hooks

```bash
fisherman install pre-commit commit-msg
```

This installs only the specified hooks.

## Force Install (Override Existing Hooks)

```bash
fisherman install --force
```

This overwrites existing hook scripts. The original scripts are backed up with a `.bkp` extension.

---

# Uninstalling Hooks

To remove Fisherman hooks, simply delete the hook scripts from `.git/hooks/`:

```bash
cd /path/to/your/repo
rm .git/hooks/pre-commit
rm .git/hooks/commit-msg
# ... etc
```

If you forced installation and want to restore original hooks:

```bash
cd .git/hooks
mv pre-commit.bkp pre-commit
mv commit-msg.bkp commit-msg
# ... etc
```

---

# Updating Fisherman

To update Fisherman to the latest version:

```bash
cargo install --git https://github.com/evg4b/fisherman.git --force
```

After updating, reinstall hooks in your repositories:

```bash
cd /path/to/your/repo
fisherman install --force
```

This ensures the hook scripts point to the latest Fisherman binary.

---

# Platform-Specific Notes

## Linux

No special considerations. Fisherman should work out of the box.

## macOS

No special considerations. Fisherman should work out of the box.

## Windows (WSL)

Fisherman should work in WSL (Windows Subsystem for Linux). Make sure you're using a Unix-style Git installation within
WSL.

**Note:** Fisherman may not work correctly with native Windows Git due to differences in how hooks are executed.

---

# Troubleshooting

## `fisherman: command not found`

**Problem:** Cargo's binary directory is not in your PATH.

**Solution:** Add `~/.cargo/bin` to your PATH:

```bash
# For bash
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# For zsh
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

---

## Permission Denied

**Problem:** Hook scripts don't have execute permissions.

**Solution:** Fisherman automatically sets execute permissions, but if you encounter this issue:

```bash
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/commit-msg
# ... etc
```

---

## Hooks Not Running

**Problem:** Hooks are installed but don't execute.

**Possible causes:**

1. **Hooks disabled in Git configuration:**
   ```bash
   git config --get core.hooksPath
   ```
   If this returns a value, Git is using a custom hooks directory. Either:
    - Install Fisherman hooks in that directory, or
    - Unset the custom hooks path: `git config --unset core.hooksPath`

2. **Hook script not executable:**
   ```bash
   ls -la .git/hooks/pre-commit
   ```
   Should show `-rwxr-xr-x` permissions.

3. **Configuration errors:**
   ```bash
   fisherman explain pre-commit
   ```
   This shows what rules are configured and helps identify configuration issues.

---

## Cargo Installation Fails

**Problem:** `cargo install` fails with compilation errors.

**Solutions:**

1. **Update Rust:**
   ```bash
   rustup update
   ```

2. **Check Rust version:**
   ```bash
   rustc --version
   ```
   Make sure you have Rust 1.70 or higher.

3. **Clean cargo cache and retry:**
   ```bash
   cargo clean
   cargo install --git https://github.com/evg4b/fisherman.git --force
   ```

---

# Next Steps

Now that Fisherman is installed:

1. **Learn about configuration** - Read [Configuration](./Configuration) to understand how to structure your
   `.fisherman.toml` files
2. **Explore available rules** - Check [Rules Reference](./Rules-reference) for all available rule types
3. **Use variables and templates** - Learn about [Variables and Templates](./Variables-and-templates) for dynamic
   configurations
4. **See examples** - Browse [Examples](./Examples-of-usage) for real-world use cases
5. **Understand Git hooks** - Read [Git Hooks](./Git-hooks) for information about available hooks

---

# Getting Help

- **Documentation** - Check this wiki for comprehensive guides
- **Issues** - Report bugs or request features at [GitHub Issues](https://github.com/evg4b/fisherman/issues)
- **CLI Help** - Run `fisherman --help` for command-line usage information
