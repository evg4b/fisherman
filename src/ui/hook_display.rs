use core::GitHook;
use std::path::PathBuf;

pub fn hook_display(hook: &GitHook, files: Vec<PathBuf>) -> String {
    format!(
        r#"
           .d8b.  |
          d88888b | hook: {}
          "Y888Y" |
 .          888   | {}
 8b.        888   | {}
 888b.      888   | {}
 888       .88P   |
 "Y8b.....d88P"   | fisherman: {}
  "Y8988888P"     |
"#,
        hook,
        get(&files, 0),
        get(&files, 1),
        get(&files, 2),
        env!("CARGO_PKG_VERSION"),
    )
}

fn get(files: &[PathBuf], index: usize) -> &str {
    files.get(index).and_then(|p| p.to_str()).unwrap_or("")
}

#[cfg(test)]
mod tests {
    use super::*;
    use core::GitHook;

    #[test]
    fn test_hook_display_no_files() {
        let result = hook_display(&GitHook::PreCommit, vec![]);
        assert!(result.contains("pre-commit"));
    }

    #[test]
    fn test_hook_display_with_files() {
        let files = vec![
            PathBuf::from("/home/.fisherman.toml"),
            PathBuf::from("/repo/.fisherman.toml"),
        ];
        let result = hook_display(&GitHook::CommitMsg, files);
        assert!(result.contains("commit-msg"));
        assert!(result.contains("/home/.fisherman.toml"));
        assert!(result.contains("/repo/.fisherman.toml"));
    }

    #[test]
    fn test_hook_display_contains_fisherman_version() {
        let result = hook_display(&GitHook::PreCommit, vec![]);
        assert!(result.contains("fisherman:"));
    }

    #[test]
    fn test_get_returns_empty_for_out_of_bounds() {
        let files: Vec<PathBuf> = vec![];
        assert_eq!(get(&files, 0), "");
        assert_eq!(get(&files, 5), "");
    }

    #[test]
    fn test_get_returns_path_at_index() {
        let files = vec![
            PathBuf::from("/path/to/first.toml"),
            PathBuf::from("/path/to/second.toml"),
        ];
        assert_eq!(get(&files, 0), "/path/to/first.toml");
        assert_eq!(get(&files, 1), "/path/to/second.toml");
        assert_eq!(get(&files, 2), "");
    }
}
