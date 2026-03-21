use crate::context::Context;
use anyhow::Result;

#[derive(Debug)]
pub enum RuleResult {
    Success {
        name: String,
        output: Option<String>,
    },
    Failure {
        name: String,
        message: String,
    },
    Skipped {
        name: String,
    },
}

#[typetag::serde(tag = "type")]
pub trait Rule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult>;
}
