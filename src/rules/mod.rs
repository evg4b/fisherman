use serde::Deserialize;

#[derive(Debug, Deserialize)]
#[serde(tag = "type")]
pub(crate) enum RuleRef {
    #[serde(rename = "exec")]
    ExecRule {
        command: String,
        args: Option<Vec<String>>,
    },
}
