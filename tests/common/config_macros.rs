//! Macro-based config generation for tests
//!
//! Provides a clean, declarative way to create fisherman configurations
//! in tests without verbose string concatenation.

/// Generate a complete fisherman configuration
///
/// # Examples
/// ```ignore
/// let config = config! {
///     extract: ["branch:^(?P<Type>feature|bugfix)"],
///     hooks: {
///         "pre-commit" => [
///             write_file!("output.txt", "content"),
///             branch_prefix!("feature/"),
///         ]
///     }
/// };
/// ```
#[macro_export]
macro_rules! config {
    // With extract patterns
    (
        extract: [$($extract:expr),* $(,)?],
        hooks: {
            $($hook:expr => [$($rules:expr),* $(,)?]),* $(,)?
        }
    ) => {{
        let mut result = String::new();
        result.push_str("extract = [");
        $(
            result.push_str(&format!("\"{}\"", $extract));
            result.push_str(", ");
        )*
        result.push_str("]\n\n");

        $(
            $(
                result.push_str(&format!("[[hooks.{}]]\n", $hook));
                result.push_str(&$rules);
                result.push('\n');
            )*
        )*

        result
    }};

    // Without extract patterns
    (
        hooks: {
            $($hook:expr => [$($rules:expr),* $(,)?]),* $(,)?
        }
    ) => {{
        let mut result = String::new();

        $(
            $(
                result.push_str(&format!("[[hooks.{}]]\n", $hook));
                result.push_str(&$rules);
                result.push('\n');
            )*
        )*

        result
    }};
}

/// Generate a write-file rule
#[macro_export]
macro_rules! write_file {
    ($path:expr, $content:expr) => {
        format!(r#"type = "write-file"
path = "{}"
content = "{}""#, $path, $content)
    };
    ($path:expr, $content:expr, append: $append:expr) => {
        format!(r#"type = "write-file"
path = "{}"
content = "{}"
append = {}"#, $path, $content, $append)
    };
    ($path:expr, $content:expr, when: $when:expr) => {
        format!(r#"type = "write-file"
path = "{}"
content = "{}"
when = "{}""#, $path, $content, $when)
    };
}

/// Generate a branch-name-regex rule
#[macro_export]
macro_rules! branch_regex {
    ($regex:expr) => {
        format!(r#"type = "branch-name-regex"
regex = "{}""#, $regex)
    };
    ($regex:expr, when: $when:expr) => {
        format!(r#"type = "branch-name-regex"
regex = "{}"
when = "{}""#, $regex, $when)
    };
}

/// Generate a branch-name-prefix rule
#[macro_export]
macro_rules! branch_prefix {
    ($prefix:expr) => {
        format!(r#"type = "branch-name-prefix"
prefix = "{}""#, $prefix)
    };
    ($prefix:expr, when: $when:expr) => {
        format!(r#"type = "branch-name-prefix"
prefix = "{}"
when = "{}""#, $prefix, $when)
    };
}

/// Generate a branch-name-suffix rule
#[macro_export]
macro_rules! branch_suffix {
    ($suffix:expr) => {
        format!(r#"type = "branch-name-suffix"
suffix = "{}""#, $suffix)
    };
    ($suffix:expr, when: $when:expr) => {
        format!(r#"type = "branch-name-suffix"
suffix = "{}"
when = "{}""#, $suffix, $when)
    };
}

/// Generate a message-regex rule
#[macro_export]
macro_rules! message_regex {
    ($regex:expr) => {
        format!(r#"type = "message-regex"
regex = "{}""#, $regex)
    };
    ($regex:expr, when: $when:expr) => {
        format!(r#"type = "message-regex"
regex = "{}"
when = "{}""#, $regex, $when)
    };
}

/// Generate a message-prefix rule
#[macro_export]
macro_rules! message_prefix {
    ($prefix:expr) => {
        format!(r#"type = "message-prefix"
prefix = "{}""#, $prefix)
    };
    ($prefix:expr, when: $when:expr) => {
        format!(r#"type = "message-prefix"
prefix = "{}"
when = "{}""#, $prefix, $when)
    };
}

/// Generate a message-suffix rule
#[macro_export]
macro_rules! message_suffix {
    ($suffix:expr) => {
        format!(r#"type = "message-suffix"
suffix = "{}""#, $suffix)
    };
    ($suffix:expr, when: $when:expr) => {
        format!(r#"type = "message-suffix"
suffix = "{}"
when = "{}""#, $suffix, $when)
    };
}

/// Generate an exec rule
#[macro_export]
macro_rules! exec {
    ($command:expr) => {
        format!(r#"type = "exec"
command = "{}""#, $command)
    };
    ($command:expr, args: [$($arg:expr),* $(,)?]) => {{
        let mut result = format!(r#"type = "exec"
command = "{}""#, $command);
        result.push_str("\nargs = [");
        $(
            result.push_str(&format!("\"{}\"", $arg));
            result.push_str(", ");
        )*
        result.push(']');
        result
    }};
    ($command:expr, args: [$($arg:expr),* $(,)?], env: {$($key:expr => $val:expr),* $(,)?}) => {{
        let mut result = format!(r#"type = "exec"
command = "{}""#, $command);
        result.push_str("\nargs = [");
        $(
            result.push_str(&format!("\"{}\"", $arg));
            result.push_str(", ");
        )*
        result.push(']');
        result.push_str("\n[env]");
        $(
            result.push_str(&format!("\n{} = \"{}\"", $key, $val));
        )*
        result
    }};
}

/// Generate a shell script rule
#[macro_export]
macro_rules! shell {
    ($script:expr) => {
        format!(r#"type = "shell"
script = """
{}
""""#, $script)
    };
    ($script:expr, env: {$($key:expr => $val:expr),* $(,)?}) => {{
        let mut result = format!(r#"type = "shell"
script = """
{}
""""#, $script);
        result.push_str("\n[env]");
        $(
            result.push_str(&format!("\n{} = \"{}\"", $key, $val));
        )*
        result
    }};
}
