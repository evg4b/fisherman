use crate::hooks::GitHook;
use std::path::PathBuf;

pub(crate) fn logo() -> String {
    format!(
        r#"
 .d888  d8b          888
 d88P"  Y8P          888                        {:>30}
 888                 888
 888888 888 .d8888b  88888b.   .d88b.  888d888 88888b.d88b.   8888b.  88888b.
 888    888 88K      888 "88b d8P  Y8b 888P"   888 "888 "88b     "88b 888 "88b
 888    888 "Y8888b. 888  888 88888888 888     888  888  888 .d888888 888  888
 888    888      X88 888  888 Y8b.     888     888  888  888 888  888 888  888
 888    888  88888P' 888  888  "Y8888  888     888  888  888 "Y888888 888  888
"#,
        format!("Version: {}", env!("CARGO_PKG_VERSION"))
    )
}

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
