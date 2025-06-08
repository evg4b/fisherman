use crate::context::Context;
use anyhow::Result;

#[derive(Debug)]
pub enum RuleResult {
    Success { name: String, output: String, },
    Failure { name: String, message: String },
}

pub trait CompiledRule {
    fn sync(&self) -> bool;
    fn check(&self, ctx: &dyn Context) -> Result<RuleResult>;
}