#![allow(dead_code)]

use std::path::{Path, PathBuf};
use std::process::{Command, Output};
use std::sync::Once;

static COMPILE: Once = Once::new();

pub struct FishermanBinary {
    path: PathBuf,
}

impl FishermanBinary {
    pub fn build() -> Self {
        COMPILE.call_once(|| {
            let output = Command::new("cargo")
                .args(["build", "--release"])
                .output()
                .expect("Failed to build fisherman binary");

            if !output.status.success() {
                panic!(
                    "Failed to compile fisherman:\n{}",
                    String::from_utf8_lossy(&output.stderr)
                );
            }

            let binary_path = Self::binary_path();
            assert!(binary_path.exists(), "Binary not found after compilation");
        });

        Self {
            path: Self::binary_path(),
        }
    }

    pub fn path(&self) -> &Path {
        &self.path
    }

    pub fn run(&self, args: &[&str], working_dir: &Path) -> Output {
        Command::new(&self.path)
            .args(args)
            .current_dir(working_dir)
            .output()
            .expect("Failed to execute fisherman command")
    }

    pub fn install(&self, working_dir: &Path, force: bool) -> Output {
        let mut args = vec!["install"];
        if force {
            args.push("--force");
        }
        self.run(&args, working_dir)
    }

    pub fn handle(&self, hook: &str, working_dir: &Path, extra_args: &[&str]) -> Output {
        let mut args = vec!["handle", hook];
        args.extend_from_slice(extra_args);
        self.run(&args, working_dir)
    }

    pub fn explain(&self, hook: &str, working_dir: &Path) -> Output {
        self.run(&["explain", hook], working_dir)
    }

    fn binary_path() -> PathBuf {
        let manifest_dir = env!("CARGO_MANIFEST_DIR");
        let mut path = PathBuf::from(manifest_dir);
        path.push("target");
        path.push("release");
        path.push(Self::binary_name());
        path
    }

    #[cfg(windows)]
    fn binary_name() -> &'static str {
        "fisherman.exe"
    }

    #[cfg(not(windows))]
    fn binary_name() -> &'static str {
        "fisherman"
    }
}
