use crate::common::ConfigFormat;
use core::configuration::Configuration;

#[macro_export] macro_rules! rule {
    ($params:expr) => {
        Rule {
            when: None,
            extract: None,
            params: $params,
        }
    };
}

#[macro_export] macro_rules! config {
    ($hook:expr => [ $( $rule:expr ),* $(,)? ]) => {{
        Configuration {
            hooks: std::collections::HashMap::from([
                ($hook, vec![$($rule),*])
            ]),
            extract: vec![],
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
