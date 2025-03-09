use crate::common::BError;
use crate::context::Context;
use crate::err;
use crate::hooks::errors::HookError;
use crate::hooks::GitHook;
use std::os::unix::fs::PermissionsExt;
use std::path::Path;
use std::{fs, io};

pub fn write_hook(context: &impl Context, hook: GitHook, content: String) -> Result<(), BError> {
    let hook_path = &context.hooks_dir().join(hook.as_str());

    if hook_path.exists() {
        err!(HookError::AlreadyExists {
            hook: hook_path.clone()
        });
    }

    fs::write(hook_path, content)?;
    fs::set_permissions(hook_path, fs::Permissions::from_mode(0o700))?;

    Ok(())
}

pub fn override_hook(context: &impl Context, hook: GitHook, content: String) -> io::Result<()> {
    let hook_path = &context.hooks_dir().join(hook.as_str());
    fs::write(hook_path, content)?;
    fs::set_permissions(hook_path, fs::Permissions::from_mode(0o700))
}

pub fn build_hook_content(bin: &Path, hook_name: GitHook) -> String {
    format!("#!/bin/sh\n{} handle {}\n", bin.display(), hook_name)
}
