use anyhow::Result;
use crate::context::Context;

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
    Skipped,
}

#[typetag::serde(tag = "type")]
pub trait Rule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult>;
}