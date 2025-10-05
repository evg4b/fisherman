use crate::hooks::GitHook;
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
