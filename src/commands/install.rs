use crate::commands::command::CliCommand;
use crate::ui::logo;
use anyhow::Result;
use clap::{Parser, ValueEnum};
use core::Context;
use core::GitHook;
use std::fs;

#[derive(Debug, Parser)]
pub struct InstallCommand {
    /// List of hooks to install (if not provided, only the configured
    /// hooks will be installed or all hooks if no configuration is found)
    hooks: Option<Vec<GitHook>>,
    /// Force the initialization of the hooks (override existing hooks)
    #[arg(short, long)]
    force: bool,
}

impl CliCommand for InstallCommand {
    fn exec(&self, context: &mut impl Context) -> Result<()> {
        println!("{}", logo());

        let selected_hooks = match self.hooks.as_ref() {
            Some(hooks) => hooks.clone(),
            None => context
                .configuration()?
                .get_configured_hooks()
                .unwrap_or_else(|| GitHook::value_variants().into()),
        };

        fs::create_dir_all(context.hooks_dir())?;
        for hook in selected_hooks {
            hook.install(context, self.force)?;
            println!("Hook {} installed", hook);
        }

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use core::Configuration;
    use core::MockContext;
    use std::collections::HashMap;
    use std::path::PathBuf;
    use tempdir::TempDir;

    #[test]
    fn test_exec_with_explicit_hooks() -> Result<()> {
        let dir = TempDir::new("fisherman_test")?;

        let cmd = InstallCommand {
            hooks: Some(vec![GitHook::PreCommit]),
            force: false,
        };

        let mut ctx = MockContext::new();
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());
        assert!(dir.path().join("pre-commit").exists());

        Ok(())
    }

    #[test]
    fn test_exec_with_configured_hooks() -> Result<()> {
        let dir = TempDir::new("fisherman_test")?;

        let cmd = InstallCommand {
            hooks: None,
            force: false,
        };

        let config = Configuration {
            hooks: HashMap::from([(GitHook::PreCommit, vec![])]),
            extract: vec![],
            files: vec![],
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration().return_once(move || Ok(config));
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());
        assert!(dir.path().join("pre-commit").exists());

        Ok(())
    }

    #[test]
    fn test_exec_with_no_configured_hooks_installs_all() -> Result<()> {
        let dir = TempDir::new("fisherman_test")?;

        let cmd = InstallCommand {
            hooks: None,
            force: false,
        };

        let config = Configuration {
            hooks: HashMap::new(),
            extract: vec![],
            files: vec![],
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration().return_once(move || Ok(config));
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());
        ctx.expect_bin()
            .return_const(PathBuf::from("/usr/bin/fisherman"));

        let result = cmd.exec(&mut ctx);
        assert!(result.is_ok());

        Ok(())
    }

    #[test]
    fn test_exec_configuration_error() {
        let dir = TempDir::new("fisherman_test").unwrap();

        let cmd = InstallCommand {
            hooks: None,
            force: false,
        };

        let mut ctx = MockContext::new();
        ctx.expect_configuration()
            .return_once(|| Err(anyhow::anyhow!("Config error")));
        ctx.expect_hooks_dir()
            .return_const(dir.path().to_path_buf());

        let result = cmd.exec(&mut ctx);
        assert!(result.is_err());
    }
}
