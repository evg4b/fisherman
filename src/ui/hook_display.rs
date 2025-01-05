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
        get_at_index(&files, 0),
        get_at_index(&files, 1),
        get_at_index(&files, 2),
        env!("CARGO_PKG_VERSION"),
    )
}

fn get_at_index(files: &[PathBuf], index: usize) -> String {
    files.get(index)
        .and_then(|p| p.to_str())
        .unwrap_or("")
        .to_string()
}
