use crate::common::ConfigFormat;
use core::configuration::Configuration;

#[macro_export]
macro_rules! rule {
    ($params:expr $(, extract = $extract:expr)? $(, when = $when:expr)? ) => {
        core::rules::Rule {
            when: None $(.or(Some(core::scripting::Expression { condition: $when })))?,
            extract: None $(.or(Some($extract)))?,
            params: $params,
        }
    };
}

#[macro_export]
macro_rules! config {
    // Config with hooks only
    ($hook:expr => [ $( $rule:expr ),* $(,)? ]) => {{
        Configuration {
            hooks: std::collections::HashMap::from([
                ($hook, vec![$($rule),*])
            ]),
            extract: vec![],
            files: vec![],
        }
    }};
    // Config with hooks and extract
    ($hook:expr => [ $( $rule:expr ),* $(,)? ], extract = $extract:expr) => {{
        Configuration {
            hooks: std::collections::HashMap::from([
                ($hook, vec![$($rule),*])
            ]),
            extract: $extract,
            files: vec![],
        }
    }};
}

pub fn serialize_configuration(configuration: &Configuration, format: ConfigFormat) -> String {
    match format {
        ConfigFormat::Json => serde_json::to_string(configuration).unwrap(),
        ConfigFormat::Yaml => serde_yaml::to_string(configuration).unwrap(),
        ConfigFormat::Toml => toml::to_string(configuration).unwrap(),
    }
}
