use std::path::PathBuf;
use std::{env, fs, io};

const HOOK_NAMES: [&str; 14] = [
    "applypatch-msg",
    "commit-msg",
    "fsmonitor-watchman",
    "post-update",
    "pre-applypatch",
    "pre-commit",
    "pre-merge-commit",
    "pre-push",
    "pre-rebase",
    "pre-receive",
    "prepare-commit-msg",
    "push-to-checkout",
    "sendemail-validate",
    "update",
];

pub(crate) fn read_hooks() -> Vec<(&'static str, PathBuf)> {
    let cwd = env::current_dir().expect("Failed to get current working directory");

    let hooks_directory = cwd.join(".git/hooks");

    if !hooks_directory.exists() {
        fs::create_dir(&hooks_directory).expect("Failed to create hooks directory");
    }

    return HOOK_NAMES
        .map(|hook_name| (hook_name, hooks_directory.join(hook_name)))
        .to_vec();
}

pub(crate) fn write_hook(path: &PathBuf, content: String) -> io::Result<()> {
    fs::write(path, content)
}

pub(crate) fn backup_hook(path: &PathBuf) -> io::Result<()> {
    match fs::read_to_string(&path) {
        Ok(content) => {
            let backup_path = path.with_extension("bak");
            fs::write(backup_path, content)
        }
        Err(e) => Err(e),
    }
}

pub(crate) fn build_hook_content(bin: &PathBuf, hook_name: &'static str) -> String {
    format!(
        "#!/bin/sh\n{} handle {}\n",
        bin.display(),
        hook_name
    )
}
