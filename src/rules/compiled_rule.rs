use anyhow::Result;
use crate::context::Context;

#[derive(Debug)]
pub enum RuleResult {
    Success { name: String, output: String, },
    Failure { name: String, message: String },
}

pub trait CompiledRule {
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult>;
}